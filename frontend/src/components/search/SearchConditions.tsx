import { useEffect, useState } from "react";
import "./search.css";
import CheckBoxes from "./CheckBoxes";
import Timetable from "./Timetable";
import {
  GRADE_LABELS,
  QUARTER_OPTIONS,
  SearchCheckBox,
  SearchConditionsType,
  SearchTimetableCell,
} from "./types";

type SearchConditionsProps = {
  onCheckItem?: (type: SearchCheckBox, items: string[]) => void;
  onTimetableChange?: (items: SearchTimetableCell[]) => void;
};

const SearchConditions = ({ onCheckItem, onTimetableChange }: SearchConditionsProps) => {
  const [grades, setGrades] = useState<string[]>([]);
  const [quarters, setQuarters] = useState<string[]>([]);
  const [timetable, setTimetable] = useState<SearchTimetableCell[]>([]);

  useEffect(() => {
    if (!onTimetableChange) {
      return;
    }
    onTimetableChange(timetable);
  }, [timetable, onTimetableChange]);

  const handleCheckBoxesChange = (type: SearchConditionsType, items: string[]) => {
    if (type === "grade") {
      setGrades(items);
    }
    if (type === "quarter") {
      setQuarters(items);
    }
    onCheckItem?.(type as SearchCheckBox, items);
  };

  const handleTimetableChange = (items: SearchTimetableCell[]) => {
    setTimetable(items);
    onTimetableChange?.(items);
  };

  return (
    <div className="search-conditions">
      <div className="search-conditions-group">
        <h3 className="search-condition-title">対象学年</h3>
        <CheckBoxes idPrefix="grade" contents={GRADE_LABELS} selectedItems={grades} onChange={(items) => handleCheckBoxesChange("grade", items)} />
      </div>
      <div className="search-conditions-group">
        <h3 className="search-condition-title">開講クォーター</h3>
        <CheckBoxes idPrefix="quarter" contents={QUARTER_OPTIONS} selectedItems={quarters} onChange={(items) => handleCheckBoxesChange("quarter", items)} />
      </div>
      <div className="search-conditions-group">
        <h3 className="search-condition-title">時間割</h3>
        <Timetable value={timetable} onChange={handleTimetableChange} />
      </div>
    </div>
  );
};

export default SearchConditions;
