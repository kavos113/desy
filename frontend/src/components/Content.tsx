import { useCallback, useMemo, useState } from "react";
import { domain } from "../../wailsjs/go/models";
import SimpleButton from "./common/SimpleButton";
import ListTable from "./list/ListTable";
import Search from "./search/Search";
import { formatTeachers, formatTimetables } from "./list/utils";
import "./content.css";

type SortKey = "title" | "code" | "lecturer" | "department" | "timetable";

type Comparator = (
  left: domain.LectureSummary,
  right: domain.LectureSummary,
  ascending: boolean
) => number;

const toComparableString = (value: string | undefined | null) => {
  return (value ?? "").toString().trim().toLowerCase();
};

const compareText = (left: string, right: string, ascending: boolean) => {
  const normalizedLeft = toComparableString(left);
  const normalizedRight = toComparableString(right);

  if (normalizedLeft === normalizedRight) {
    return 0;
  }

  const result = normalizedLeft < normalizedRight ? -1 : 1;
  return ascending ? result : -result;
};

const buildComparators = (): Record<SortKey, Comparator> => ({
  title: (left, right, ascending) => compareText(left.Title ?? "", right.Title ?? "", ascending),
  code: (left, right, ascending) => compareText(left.Code ?? "", right.Code ?? "", ascending),
  lecturer: (left, right, ascending) =>
    compareText(formatTeachers(left.Teachers), formatTeachers(right.Teachers), ascending),
  department: (left, right, ascending) =>
    compareText(left.Department ?? "", right.Department ?? "", ascending),
  timetable: (left, right, ascending) =>
    compareText(
      formatTimetables(left.Timetables, { includeRoom: false }),
      formatTimetables(right.Timetables, { includeRoom: false }),
      ascending
    ),
});

const Content = () => {
  const [lectures, setLectures] = useState<domain.LectureSummary[]>([]);
  const [sortState, setSortState] = useState<Record<SortKey, boolean>>({
    title: true,
    code: true,
    lecturer: true,
    department: true,
    timetable: true,
  });
  const [isSearchVisible, setIsSearchVisible] = useState(false);

  const comparators = useMemo(buildComparators, []);

  const handleSearch = useCallback((results: domain.LectureSummary[]) => {
    setLectures(results);
    setIsSearchVisible(false);
  }, []);

  const handleSort = useCallback(
    (key: string) => {
      if (!isSortKey(key)) {
        return;
      }

      setLectures((previousLectures) => {
        const nextLectures = [...previousLectures];
        const comparator = comparators[key];
        const ascending = sortState[key];
        nextLectures.sort((left, right) => comparator(left, right, ascending));
        return nextLectures;
      });

      setSortState((previousSortState) => ({
        ...previousSortState,
        [key]: !previousSortState[key],
      }));
    },
    [comparators, sortState]
  );

  const handleMenuClick = useCallback(() => {
    setIsSearchVisible(true);
  }, []);

  const handleBack = useCallback(() => {
    setIsSearchVisible(false);
  }, []);

  const searchPanelClassName = useMemo(() => {
    return ["search-panel", isSearchVisible ? "search-visible" : ""].filter(Boolean).join(" ");
  }, [isSearchVisible]);

  return (
    <div className="content-container">
      <Search className={searchPanelClassName} onSearch={handleSearch} onBack={handleBack} />
      <SimpleButton text="メニュー" className="search-menu" onClick={handleMenuClick} />
      <ListTable className="table-panel" items={lectures} onSort={handleSort} />
    </div>
  );
};

function isSortKey(value: string): value is SortKey {
  return ["title", "code", "lecturer", "department", "timetable"].includes(value);
}

export default Content;
