import { useCallback } from 'react';
import CheckBoxes from './CheckBoxes';
import Timetable from './Timetable';
import { GRADE_LABELS, QUARTER_LABELS } from '../../constants';
import type { SearchCheckBox, SearchTimetableSelection } from './types';
import './search.css';

type SearchConditionsProps = {
  onCheckItem?: (type: SearchCheckBox, items: string[]) => void;
  onTimetableChange?: (items: SearchTimetableSelection[]) => void;
};

const SearchConditions = ({ onCheckItem, onTimetableChange }: SearchConditionsProps) => {
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

  return (
    <div className="search-conditions-container">
      <CheckBoxes checkboxId="grade" contents={GRADE_LABELS} onCheckItem={handleGradeCheck} />
      <CheckBoxes checkboxId="quarter" contents={QUARTER_LABELS} onCheckItem={handleQuarterCheck} />
      <Timetable onCheckItem={onTimetableChange} />
    </div>
  );
};

export default SearchConditions;
