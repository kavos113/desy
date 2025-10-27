import { useState } from "react";
import "./content.css";
import Search from "./search/Search";
import ListTable from "./list/ListTable";
import { ListSortKey } from "./list/ListHeaderItem";
import { domain } from "../../wailsjs/go/models";

const collator = new Intl.Collator("ja", { sensitivity: "base", numeric: true });

const timetableText = (timetables: domain.TimeTable[] | undefined): string => {
  if (!timetables?.length) {
    return "";
  }

  return timetables
    .map((item) => {
      const semester = item.Semester ? `${item.Semester}` : "";
      const day = item.DayOfWeek ? `${item.DayOfWeek}` : "";
      const period = item.Period ? `${item.Period}限` : "";
      const room = item.Room?.Name ? `(${item.Room.Name})` : "";
      return [semester, day, period].filter(Boolean).join("") + room;
    })
    .join(", ");
};

const initialSortState: Record<ListSortKey, boolean> = {
  title: true,
  code: true,
  lecturer: true,
  department: true,
  timetable: true,
};

const sortValue = (item: domain.LectureSummary, key: ListSortKey): string => {
  switch (key) {
    case "title":
      return item.Title ?? "";
    case "code":
      return item.Code ?? "";
    case "lecturer":
      return (item.Teachers ?? []).map((teacher: domain.Teacher) => teacher.Name).join(", ");
    case "department":
      return item.Department ?? "";
    case "timetable":
      return timetableText(item.Timetables);
    default:
      return "";
  }
};

const Content = () => {
  const [items, setItems] = useState<domain.LectureSummary[]>([]);
  const [loading, setLoading] = useState(false);
  const [statusMessage, setStatusMessage] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [searchOpen, setSearchOpen] = useState(false);
  const [sortState, setSortState] = useState(initialSortState);

  const handleSearchStart = () => {
    setLoading(true);
    setStatusMessage(null);
    setErrorMessage(null);
  };

  const handleSearch = (results: domain.LectureSummary[]) => {
  setItems(results.slice());
    setLoading(false);
    setSearchOpen(false);
    setStatusMessage(results.length ? `${results.length}件の講義が見つかりました。` : "講義は見つかりませんでした。");
    setErrorMessage(null);
  setSortState({ ...initialSortState });
  };

  const handleSearchError = (message: string) => {
    setLoading(false);
    setErrorMessage(message);
  };

  const handleSort = (key: ListSortKey) => {
    const ascending = sortState[key];
    setSortState((prev) => ({
      ...prev,
      [key]: !prev[key],
    }));

    setItems((prevItems) => {
      const next = [...prevItems];
      next.sort((a, b) => {
        const aValue = sortValue(a, key);
        const bValue = sortValue(b, key);
        const comparison = collator.compare(aValue, bValue);
        if (comparison === 0) {
          return 0;
        }
        return ascending ? comparison : -comparison;
      });
      return next;
    });
  };

  return (
    <div className="content-container">
      <Search
        className={`content-search${searchOpen ? " search-open" : ""}`}
        onSearchStart={handleSearchStart}
        onSearch={handleSearch}
        onSearchError={handleSearchError}
        onBack={() => setSearchOpen(false)}
      />
      <button type="button" className="search-menu-button" onClick={() => setSearchOpen(true)}>
        メニュー
      </button>
      <div className="content-table">
        <ListTable
          items={items}
          loading={loading}
          statusMessage={statusMessage}
          errorMessage={errorMessage}
          onSort={handleSort}
        />
      </div>
    </div>
  );
};

export default Content;
