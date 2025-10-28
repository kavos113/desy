import CheckBoxes from "./CheckBoxes";
import Timetable from "./Timetable";
import { GRADE_LABELS, QUARTER_LABELS } from "../../constants";
import type { SearchCheckBox, SearchTimetableSelection } from "./types";
import "./search.css";

type SearchConditionsProps = {
  onCheckItem?: (type: SearchCheckBox, items: string[]) => void;
  onTimetableChange?: (items: SearchTimetableSelection[]) => void;
};

const SearchConditions = ({ onCheckItem, onTimetableChange }: SearchConditionsProps) => {
  const handleCheck = (type: SearchCheckBox) => {
    return (items: string[]) => {
      onCheckItem?.(type, items);
    };
  };

  return (
    <div className="search-conditions-container">
      <CheckBoxes checkboxId="grade" contents={GRADE_LABELS} onCheckItem={handleCheck("grade")} />
      <CheckBoxes
        checkboxId="quarter"
        contents={QUARTER_LABELS}
        onCheckItem={handleCheck("quarter")}
      />
      <Timetable onCheckItem={onTimetableChange} />
    </div>
  );
};

export default SearchConditions;
