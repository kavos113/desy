import { useCallback } from 'react';
import CheckBoxes from './CheckBoxes';
import Timetable from './Timetable';
import CheckBox from '../common/CheckBox';
import { GRADE_LABELS, QUARTER_LABELS } from '../../constants';
import type { SearchCheckBox, SearchTimetableSelection } from './types';
import './search.css';

type SearchConditionsProps = {
  onCheckItem?: (type: SearchCheckBox, items: string[]) => void;
  onTimetableChange?: (items: SearchTimetableSelection[]) => void;
  onToggleFilterNotResearch?: (value: boolean) => void;
  resetSignal?: number;
  filterNotResearch?: boolean;
};

const SearchConditions = ({
  onCheckItem,
  onTimetableChange,
  onToggleFilterNotResearch,
  resetSignal,
  filterNotResearch
}: SearchConditionsProps) => {
  const handleGradeCheck = useCallback(
    (items: string[]) => {
      onCheckItem?.('grade', items);
    },
    [onCheckItem]
  );

  const handleQuarterCheck = useCallback(
    (items: string[]) => {
      onCheckItem?.('quarter', items);
    },
    [onCheckItem]
  );

  const handleFilterToggle = useCallback(
    (value: boolean) => {
      onToggleFilterNotResearch?.(value);
    },
    [onToggleFilterNotResearch]
  );

  return (
    <div className="search-conditions-container">
      <CheckBoxes
        checkboxId="grade"
        contents={GRADE_LABELS}
        onCheckItem={handleGradeCheck}
        resetSignal={resetSignal}
      />
      <CheckBoxes
        checkboxId="quarter"
        contents={QUARTER_LABELS}
        onCheckItem={handleQuarterCheck}
        resetSignal={resetSignal}
      />
      <Timetable onCheckItem={onTimetableChange} resetSignal={resetSignal} />
      <CheckBox
        checkboxId="filterNotResearch"
        content="講究除外"
        onChange={handleFilterToggle}
        checked={filterNotResearch}
      />
    </div>
  );
};

export default SearchConditions;
