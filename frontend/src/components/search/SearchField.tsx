import "./search.css";
import SearchBoxes from "./SearchBoxes";
import SearchConditions from "./SearchConditions";
import { SearchConditionsType, SearchTimetableCell } from "./types";

type SearchFieldProps = {
  onConditionChange?: (key: SearchConditionsType, items: string[]) => void;
  onTimetableChange?: (items: SearchTimetableCell[]) => void;
};

const SearchField = ({ onConditionChange, onTimetableChange }: SearchFieldProps) => {
  const handleSelect = (key: SearchConditionsType, items: string[]) => {
    onConditionChange?.(key, items);
  };

  const handleSearchBoxChange = (title: string, lecturer: string) => {
    onConditionChange?.("title", title ? [title] : []);
    onConditionChange?.("lecturer", lecturer ? [lecturer] : []);
  };

  const handleTimetable = (items: SearchTimetableCell[]) => {
    onTimetableChange?.(items);
  };

  return (
    <div className="search-container">
      <SearchBoxes onSelectMenuItem={handleSelect} onChangeSearchBox={handleSearchBoxChange} />
      <SearchConditions onCheckItem={handleSelect} onTimetableChange={handleTimetable} />
    </div>
  );
};

export default SearchField;
