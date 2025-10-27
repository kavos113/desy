import { useMemo, useState } from "react";
import "./search.css";
import FetchButton from "./FetchButton";
import SearchField from "./SearchField";
import {
  Day,
  Period,
  SearchConditionsType,
  SearchTimetableCell,
} from "./types";
import { SearchLectures as defaultSearchLectures } from "../../../wailsjs/go/main/App";
import { domain } from "../../../wailsjs/go/models";

type SearchProps = {
  className?: string;
  onSearch?: (results: domain.LectureSummary[]) => void;
  onBack?: () => void;
  onSearchStart?: () => void;
  onSearchError?: (message: string) => void;
  searchLecturesFn?: typeof defaultSearchLectures;
};

type SearchCondition = {
  university: string[];
  department: string[];
  year: string[];
  title: string[];
  lecturer: string[];
  grade: string[];
  quarter: string[];
  timetable: SearchTimetableCell[];
};

const initialCondition: SearchCondition = {
  university: [],
  department: [],
  year: [],
  title: [],
  lecturer: [],
  grade: [],
  quarter: [],
  timetable: [],
};

const GRADE_TO_LEVEL: Record<string, number> = {
  学士1年: 1,
  学士2年: 2,
  学士3年: 3,
  修士1年: 4,
  修士2年: 5,
  博士課程: 6,
};

const DAY_TO_QUERY: Record<Day, string> = {
  月: "Monday",
  火: "Tuesday",
  水: "Wednesday",
  木: "Thursday",
  金: "Friday",
};

const QUARTER_TO_SEMESTER: Record<string, string> = {
  "1Q": "First",
  "2Q": "Second",
  "3Q": "Third",
  "4Q": "Fourth",
};

const toPeriodNumber = (period: Period): number => Number(period);

const buildTimeTables = (cells: SearchTimetableCell[], quarters: string[]): domain.TimeTable[] => {
  const quarterList = quarters.length ? quarters : [""];
  const entries: domain.TimeTable[] = [];

  quarterList.forEach((quarter) => {
    const semester = QUARTER_TO_SEMESTER[quarter] ?? "";
    cells.forEach((cell) => {
      entries.push(
        domain.TimeTable.createFrom({
          Semester: semester,
          DayOfWeek: DAY_TO_QUERY[cell.day],
          Period: toPeriodNumber(cell.period),
        })
      );
    });
  });

  return entries;
};

const buildSearchQuery = (condition: SearchCondition): domain.SearchQuery => {
  const title = condition.title[0] ?? "";
  const lecturer = condition.lecturer[0] ?? "";
  const yearValue = condition.year[0] ? Number(condition.year[0]) : 0;
  const levels = condition.grade
    .map((grade) => GRADE_TO_LEVEL[grade])
    .filter((level): level is number => typeof level === "number" && !Number.isNaN(level));

  const timeTables = buildTimeTables(condition.timetable, condition.quarter);

  return domain.SearchQuery.createFrom({
    Title: title,
    Keywords: [],
    Departments: condition.department,
    Year: yearValue,
    TeacherName: lecturer,
    TimeTables: timeTables,
    Levels: levels,
  });
};

const Search = ({ className, onSearch, onBack, onSearchStart, onSearchError, searchLecturesFn }: SearchProps) => {
  const [condition, setCondition] = useState<SearchCondition>(initialCondition);
  const [searching, setSearching] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const searchLectures = searchLecturesFn ?? defaultSearchLectures;

  const universitySummary = useMemo(() => condition.university.join(", "), [condition.university]);
  const departmentSummary = useMemo(() => condition.department.join(", "), [condition.department]);
  const yearSummary = useMemo(() => condition.year.join(", "), [condition.year]);

  const handleConditionChange = (key: SearchConditionsType, items: string[]) => {
    setCondition((prev) => ({
      ...prev,
      [key]: items,
    }) as SearchCondition);
  };

  const handleTimetableChange = (items: SearchTimetableCell[]) => {
    setCondition((prev) => ({
      ...prev,
      timetable: items,
    }));
  };

  const handleSearch = async () => {
    if (searching) {
      return;
    }

    onSearchStart?.();
    setSearching(true);
    setError(null);

    try {
      const query = buildSearchQuery(condition);
      const results = await searchLectures(query);
      onSearch?.(results);
      onBack?.();
    } catch (err) {
      console.error(err);
      setError("検索に失敗しました。");
      onSearchError?.("検索に失敗しました。");
    } finally {
      setSearching(false);
    }
  };

  const wrapperClassName = ["search-wrapper", className].filter(Boolean).join(" ");

  return (
    <div className={wrapperClassName}>
      <FetchButton className="fetch" />
      <SearchField onConditionChange={handleConditionChange} onTimetableChange={handleTimetableChange} />
      <div className="search-actions">
        <button type="button" className="button primary" onClick={handleSearch} disabled={searching}>
          {searching ? "Searching..." : "Search"}
        </button>
        {error && <p className="search-summary" role="alert">{error}</p>}
        <p className="search-summary">大学: {universitySummary || "未選択"}</p>
        <p className="search-summary">開講: {departmentSummary || "未選択"}</p>
        <p className="search-summary">年度: {yearSummary || "未選択"}</p>
        <button type="button" className="button ghost back" onClick={() => onBack?.()}>
          戻る
        </button>
      </div>
    </div>
  );
};

export default Search;
