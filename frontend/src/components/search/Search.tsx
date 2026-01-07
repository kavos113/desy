import type { KeyboardEvent } from 'react';
import { useCallback, useMemo, useState } from 'react';
import SimpleButton from '../common/SimpleButton';
import FetchButton from './FetchButton';
import SearchField from './SearchField';
import {
  dayToDomain,
  gradeLabelToLevel,
  parseKeywordInput,
  parseYearLabel,
  periodForQueriesFromPeriod,
  quarterToSemester
} from '../../constants';
import type { SearchConditionKey, SearchState, SearchTimetableSelection } from './types';
import { SearchLectures } from '../../../wailsjs/go/main/App';
import { domain } from '../../../wailsjs/go/models';
import './search.css';

type SearchProps = {
  className?: string;
  onSearch?: (results: domain.LectureSummary[]) => void;
  onBack?: () => void;
};

const createInitialState = (): SearchState => ({
  university: [],
  department: [],
  year: [],
  title: [],
  lecturer: [],
  room: [],
  grade: [],
  quarter: [],
  timetable: [],
  filterNotResearch: false
});

const buildSearchQuery = (state: SearchState) => {
  const title = state.title[0] ?? '';
  const teacherName = state.lecturer[0] ?? '';
  const room = state.room[0] ?? '';
  const keywords = parseKeywordInput(title);
  const levels = state.grade
    .map(gradeLabelToLevel)
    .filter((value): value is number => typeof value === 'number');
  const yearValue =
    state.year
      .map(parseYearLabel)
      .find((value): value is number => typeof value === 'number') ?? 0;

  const timetables = state.timetable.map((item) => {
    const periods = periodForQueriesFromPeriod(item.period);
    return periods.map((period) =>
      domain.TimeTable.createFrom({
        DayOfWeek: dayToDomain(item.day),
        Period: period
      })
    );
  }).flat();

  const semesters = state.quarter
    .map((label) => quarterToSemester(label))
    .filter((value): value is string => typeof value === 'string');

  return domain.SearchQuery.createFrom({
    Title: title,
    Keywords: keywords,
    Departments: state.department,
    Year: yearValue,
    TeacherName: teacherName,
    Room: room,
    Semester: semesters,
    TimeTables: timetables,
    Levels: levels,
    FilterNotResearch: state.filterNotResearch
  });
};

const Search = ({ className, onSearch, onBack }: SearchProps) => {
  const [condition, setCondition] = useState<SearchState>(() => createInitialState());
  const [isSearching, setIsSearching] = useState(false);
  const [resetCounter, setResetCounter] = useState(0);

  const wrapperClassName = useMemo(() => {
    return ['search-wrapper', className].filter(Boolean).join(' ');
  }, [className]);

  const handleConditionChange = useCallback((key: SearchConditionKey, items: string[]) => {
    setCondition((previous) => ({ ...previous, [key]: items }));
  }, []);

  const handleTimetableChange = useCallback((items: SearchTimetableSelection[]) => {
    setCondition((previous) => ({ ...previous, timetable: items }));
  }, []);

  const handleFilterNotResearchChange = useCallback((value: boolean) => {
    setCondition((previous) => ({ ...previous, filterNotResearch: value }));
  }, []);

  const handleSearch = useCallback(async () => {
    if (isSearching) {
      return;
    }
    setIsSearching(true);
    try {
      const query = buildSearchQuery(condition);
      console.log(query)
      const results = await SearchLectures(query);
      onSearch?.(results);
      onBack?.();
    } catch (error) {
      console.error('SearchLectures failed', error);
    } finally {
      setIsSearching(false);
    }
  }, [condition, isSearching, onBack, onSearch]);

  const handleBack = useCallback(() => {
    onBack?.();
  }, [onBack]);

  const handleReset = useCallback(() => {
    setCondition(createInitialState());
    setResetCounter((previous) => previous + 1);
  }, []);

  const handleKeyDown = useCallback(
    (event: KeyboardEvent<HTMLDivElement>) => {
      if (event.key !== 'Enter' || event.isDefaultPrevented()) {
        return;
      }
      event.preventDefault();
      handleSearch();
    },
    [handleSearch]
  );

  const selectedUniversity = useMemo(() => condition.university.join(', '), [condition.university]);
  const selectedDepartment = useMemo(() => condition.department.join(', '), [condition.department]);
  const selectedYear = useMemo(() => condition.year.join(', '), [condition.year]);

  return (
    <div className={wrapperClassName} onKeyDown={handleKeyDown}>
      <div className="fetch">
        <FetchButton />
      </div>
      <SearchField
        onClickMenuItem={handleConditionChange}
        onTimetableChange={handleTimetableChange}
        onToggleFilterNotResearch={handleFilterNotResearchChange}
        resetSignal={resetCounter}
        filterNotResearch={condition.filterNotResearch}
      />
      <div>
        <div className="search-actions">
          <SimpleButton
            text="Search"
            className="button"
            onClick={handleSearch}
            disabled={isSearching}
          />
          <SimpleButton
            text="リセット"
            className="button"
            onClick={handleReset}
            disabled={isSearching}
          />
        </div>
        <p>大学: {selectedUniversity}</p>
        <p>開講: {selectedDepartment}</p>
        <p>年度: {selectedYear}</p>
        <SimpleButton text="戻る" className="back" onClick={handleBack} />
      </div>
    </div>
  );
};

export default Search;
