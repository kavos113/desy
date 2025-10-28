import SearchBoxes from "./SearchBoxes";
import SearchConditions from "./SearchConditions";
import type {
  SearchComboBox,
  SearchConditionKey,
  SearchTimetableSelection,
} from "./types";
import "./search.css";

type SearchFieldProps = {
  onClickMenuItem?: (key: SearchConditionKey, items: string[]) => void;
  onTimetableChange?: (items: SearchTimetableSelection[]) => void;
};

const SearchField = ({ onClickMenuItem, onTimetableChange }: SearchFieldProps) => {
  const handleSearchBoxChange = (title: string, lecturer: string) => {
    onClickMenuItem?.("title", title ? [title] : []);
    onClickMenuItem?.("lecturer", lecturer ? [lecturer] : []);
  };

  const handleClickMenuItem = (key: SearchComboBox, items: string[]) => {
    onClickMenuItem?.(key, items);
  };

  const handleConditionChange = (key: SearchConditionKey, items: string[]) => {
    onClickMenuItem?.(key, items);
  };

  return (
    <div className="search-container">
      <SearchBoxes
        onClickMenuItem={handleClickMenuItem}
        onChangeSearchBox={handleSearchBoxChange}
      />
      <SearchConditions
        onCheckItem={handleConditionChange}
        onTimetableChange={onTimetableChange}
      />
    </div>
  );
};

export default SearchField;
