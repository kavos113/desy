import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { domain } from "../../../wailsjs/go/models";
import { GetLectureDetails } from "../../../wailsjs/go/main/App";
import SimpleButton from "../common/SimpleButton";
import CourseDetail from "./CourseDetail";
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
  const closeTimerRef = useRef<NodeJS.Timeout | null>(null);

  const handleSort = useCallback(
    (key: string) => {
      onSort?.(key);
    },
    [onSort]
  );

  const handleListItemClick = useCallback(async (id: number) => {
    setIsLoadingDetail(true);
    try {
      const lecture = await GetLectureDetails(id);
      setSelectedLecture(lecture);
      setIsDetailOpen(true);
      setIsOverlayActive(true);
    } catch (error) {
      console.error("GetLectureDetails failed", error);
    } finally {
      setIsLoadingDetail(false);
    }
  }, []);

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

  const tableClassName = useMemo(() => {
    return ["list-table", className].filter(Boolean).join(" ");
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

  return (
    <div className={tableClassName}>
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
        <CourseDetail lecture={selectedLecture} />
      </div>
      <SimpleButton text="戻る" className={backButtonClassName} onClick={closeDetail} />
      <div className={overlayClassName} onClick={closeDetail}></div>
    </div>
  );
};

export default ListTable;
