import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { domain } from "../../../wailsjs/go/models";
import { GetLectureDetails } from "../../../wailsjs/go/main/App";
import SimpleButton from "../common/SimpleButton";
import CourseDetail, { type RelatedCourseEntry } from "./CourseDetail";
import ListHeaderItem from "./ListHeaderItem";
import ListItem from "./ListItem";
import "./list.css";

type ListTableProps = {
  items: domain.LectureSummary[];
  onSort?: (key: string) => void;
  className?: string;
};

const CLOSE_DELAY_MS = 250;

const ListTable = ({ items, onSort, className }: ListTableProps) => {
  const [selectedLecture, setSelectedLecture] = useState<domain.Lecture | null>(null);
  const [isDetailOpen, setIsDetailOpen] = useState(false);
  const [isOverlayActive, setIsOverlayActive] = useState(false);
  const [isLoadingDetail, setIsLoadingDetail] = useState(false);
  const [lectureCache, setLectureCache] = useState<Record<number, domain.Lecture>>({});
  const closeTimerRef = useRef<NodeJS.Timeout | null>(null);

  const handleSort = useCallback(
    (key: string) => {
      onSort?.(key);
    },
    [onSort]
  );

  const storeLectureInCache = useCallback((lecture: domain.Lecture | null | undefined) => {
    if (!lecture || lecture.ID === 0) {
      return;
    }
    setLectureCache((prev) => {
      if (prev[lecture.ID]) {
        return prev;
      }
      return { ...prev, [lecture.ID]: lecture };
    });
  }, []);

  const prefetchRelatedLectures = useCallback(
    async (lecture: domain.Lecture) => {
      if (!lecture) {
        return;
      }
      const relatedIds = lecture.RelatedCourses ?? [];
      if (relatedIds.length === 0) {
        return;
      }

      const baseCache: Record<number, domain.Lecture> = { ...lectureCache, [lecture.ID]: lecture };
      const targets = relatedIds.filter((relatedId) => !baseCache[relatedId] && relatedId !== lecture.ID);
      if (targets.length === 0) {
        return;
      }

      const results = await Promise.allSettled(
        targets.map(async (relatedId) => {
          try {
            const detail = await GetLectureDetails(relatedId);
            return detail;
          } catch (error) {
            console.error("GetLectureDetails failed for related lecture", relatedId, error);
            throw error;
          }
        })
      );

      const updates: Record<number, domain.Lecture> = {};
      results.forEach((result) => {
        if (result.status === "fulfilled" && result.value && result.value.ID !== 0) {
          updates[result.value.ID] = result.value;
        }
      });

      if (Object.keys(updates).length > 0) {
        setLectureCache((prev) => ({ ...prev, ...updates }));
      }
    },
    [lectureCache]
  );

  const handleListItemClick = useCallback(async (id: number) => {
    setIsLoadingDetail(true);
    try {
      if (selectedLecture?.ID === id) {
        setIsDetailOpen(true);
        setIsOverlayActive(true);
        void prefetchRelatedLectures(selectedLecture);
        return;
      }

      let lecture = lectureCache[id];
      if (!lecture) {
        lecture = await GetLectureDetails(id);
      }

      if (!lecture) {
        return;
      }

      setSelectedLecture(lecture);
      setIsDetailOpen(true);
      setIsOverlayActive(true);
    } catch (error) {
      console.error("GetLectureDetails failed", error);
    } finally {
      setIsLoadingDetail(false);
    }
  }, [lectureCache, prefetchRelatedLectures, selectedLecture]);

  useEffect(() => {
    if (selectedLecture) {
      storeLectureInCache(selectedLecture);
      void prefetchRelatedLectures(selectedLecture);
    }
  }, [selectedLecture, prefetchRelatedLectures, storeLectureInCache]);

  const closeDetail = useCallback(() => {
    setIsDetailOpen(false);
    if (closeTimerRef.current) {
      clearTimeout(closeTimerRef.current);
    }
    closeTimerRef.current = setTimeout(() => {
      setIsOverlayActive(false);
    }, CLOSE_DELAY_MS);
  }, []);

  useEffect(() => {
    return () => {
      if (closeTimerRef.current) {
        clearTimeout(closeTimerRef.current);
      }
    };
  }, []);

  const handleRelatedCourseClick = useCallback(
    (id: number) => {
      void handleListItemClick(id);
    },
    [handleListItemClick]
  );

  const wrapperClassName = useMemo(() => {
    return ["list-wrapper", className].filter(Boolean).join(" ");
  }, [className]);

  const detailPanelClassName = useMemo(() => {
    return ["detail-panel", isDetailOpen ? "active" : ""].filter(Boolean).join(" ");
  }, [isDetailOpen]);

  const overlayClassName = useMemo(() => {
    return ["overlay", isOverlayActive ? "active" : ""].filter(Boolean).join(" ");
  }, [isOverlayActive]);

  const backButtonClassName = useMemo(() => {
    return ["back-button", isDetailOpen ? "visible" : ""].filter(Boolean).join(" ");
  }, [isDetailOpen]);

  const relatedCourseEntries = useMemo<RelatedCourseEntry[]>(() => {
    if (!selectedLecture) {
      return [];
    }

    const codes = selectedLecture.RelatedCourseCodes ?? [];
    const relatedIds = selectedLecture.RelatedCourses ?? [];

    if (codes.length === 0 && relatedIds.length === 0) {
      return [];
    }

    const normalize = (value?: string | null) => (value ? value.toUpperCase() : "");

    const cachedByCode = new Map<string, domain.Lecture>();
    Object.values(lectureCache).forEach((lecture) => {
      if (lecture?.Code) {
        cachedByCode.set(normalize(lecture.Code), lecture);
      }
    });

    const summaryByCode = new Map<string, domain.LectureSummary>();
    items.forEach((item) => {
      if (item?.Code) {
        const key = normalize(item.Code);
        if (!summaryByCode.has(key)) {
          summaryByCode.set(key, item);
        }
      }
    });

    const entries: RelatedCourseEntry[] = [];
    const entryByCode = new Map<string, RelatedCourseEntry>();

    codes.forEach((code) => {
      const key = normalize(code);
      if (!key) {
        return;
      }
      const cached = cachedByCode.get(key);
      const summary = summaryByCode.get(key);
      const entry: RelatedCourseEntry = {
        code,
        id: cached?.ID,
        title: cached?.Title ?? summary?.Title,
      };
      entries.push(entry);
      entryByCode.set(key, entry);
    });

    relatedIds.forEach((relatedId) => {
      const cached = lectureCache[relatedId];
      if (!cached) {
        return;
      }
      const key = normalize(cached.Code);
      if (key && entryByCode.has(key)) {
        const existing = entryByCode.get(key)!;
        existing.id = cached.ID;
        if (!existing.title) {
          existing.title = cached.Title;
        }
        return;
      }

      entries.push({
        code: cached.Code ?? `ID:${relatedId}`,
        id: cached.ID,
        title: cached.Title,
      });
    });

    return entries;
  }, [items, lectureCache, selectedLecture]);

  return (
    <div className={wrapperClassName}>
      <div className="list-table">
        <ListHeaderItem onSort={handleSort} />
        {items.length === 0 ? (
          <p>講義が見つかりませんでした。</p>
        ) : (
          items.map((item) => (
            <ListItem key={item.ID} item={item} onClick={handleListItemClick} />
          ))
        )}

        <div className={detailPanelClassName}>
          {isLoadingDetail && <p className="course-detail-text">読み込み中...</p>}
          <CourseDetail
            lecture={selectedLecture}
            relatedCourses={relatedCourseEntries}
            onSelectRelatedCourse={handleRelatedCourseClick}
          />
        </div>
        <SimpleButton text="戻る" className={backButtonClassName} onClick={closeDetail} />
        <div className={overlayClassName} onClick={closeDetail}></div>
      </div>
    </div>
  );
};

export default ListTable;
