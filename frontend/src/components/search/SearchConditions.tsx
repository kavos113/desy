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
};

const SearchConditions = ({
  onCheckItem,
  onTimetableChange,
  onToggleFilterNotResearch
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
      <CheckBoxes checkboxId="grade" contents={GRADE_LABELS} onCheckItem={handleGradeCheck} />
      <CheckBoxes checkboxId="quarter" contents={QUARTER_LABELS} onCheckItem={handleQuarterCheck} />
      <Timetable onCheckItem={onTimetableChange} />
      <CheckBox
        checkboxId="filterNotResearch"
        content="研究系科目を除外"
        onChange={handleFilterToggle}
      />
    </div>
  );
};

export default SearchConditions;
