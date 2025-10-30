import { useCallback, useMemo, useState } from 'react';
import SimpleButton from '../common/SimpleButton';
import FetchButton from './FetchButton';
import SearchField from './SearchField';
import {
  dayToDomain,
  gradeLabelToLevel,
  parseKeywordInput,
  parseYearLabel,
  periodToNumber
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

const INITIAL_STATE: SearchState = {
  university: [],
  department: [],
  year: [],
  title: [],
  lecturer: [],
  grade: [],
  quarter: [],
  timetable: [],
  filterNotResearch: false
};

const buildSearchQuery = (state: SearchState) => {
  const title = state.title[0] ?? '';
  const teacherName = state.lecturer[0] ?? '';
  const keywords = parseKeywordInput(state.title.join(' '));
  const levels = state.grade
    .map(gradeLabelToLevel)
    .filter((value): value is number => typeof value === 'number');
  const yearValue =
    state.year
      .map(parseYearLabel)
      .find((value): value is number => typeof value === 'number' && !Number.isNaN(value)) ?? 0;

  const timetables = state.timetable.map((item) =>
    domain.TimeTable.createFrom({
      DayOfWeek: dayToDomain(item.day),
      Period: periodToNumber(item.period)
    })
  );

  return domain.SearchQuery.createFrom({
    Title: title,
    Keywords: keywords,
    Departments: state.department,
    Year: yearValue,
    TeacherName: teacherName,
    TimeTables: timetables,
    Levels: levels,
    FilterNotResearch: state.filterNotResearch
  });
};

const Search = ({ className, onSearch, onBack }: SearchProps) => {
  const [condition, setCondition] = useState<SearchState>(INITIAL_STATE);
  const [isSearching, setIsSearching] = useState(false);

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
    setIsSearching(true);
    try {
      const query = buildSearchQuery(condition);
      const results = await SearchLectures(query);
      onSearch?.(results);
      onBack?.();
    } catch (error) {
      console.error('SearchLectures failed', error);
    } finally {
      setIsSearching(false);
    }
  }, [condition, onBack, onSearch]);

  const handleBack = useCallback(() => {
    onBack?.();
  }, [onBack]);

  const selectedUniversity = useMemo(() => condition.university.join(', '), [condition.university]);
  const selectedDepartment = useMemo(() => condition.department.join(', '), [condition.department]);
  const selectedYear = useMemo(() => condition.year.join(', '), [condition.year]);

  return (
    <div className={wrapperClassName}>
      <div className="fetch">
        <FetchButton />
      </div>
      <SearchField
        onClickMenuItem={handleConditionChange}
        onTimetableChange={handleTimetableChange}
        onToggleFilterNotResearch={handleFilterNotResearchChange}
      />
      <div>
        <SimpleButton
          text="Search"
          className="button"
          onClick={handleSearch}
          disabled={isSearching}
        />
        <p>大学: {selectedUniversity}</p>
        <p>開講: {selectedDepartment}</p>
        <p>年度: {selectedYear}</p>
        <SimpleButton text="戻る" className="back" onClick={handleBack} />
      </div>
    </div>
  );
};

export default Search;
