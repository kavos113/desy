import { useCallback } from 'react';
import SearchBoxes from './SearchBoxes';
import SearchConditions from './SearchConditions';
import type { SearchComboBox, SearchConditionKey, SearchTimetableSelection } from './types';
import './search.css';

type SearchFieldProps = {
  onClickMenuItem?: (key: SearchConditionKey, items: string[]) => void;
  onTimetableChange?: (items: SearchTimetableSelection[]) => void;
  onToggleFilterNotResearch?: (value: boolean) => void;
};

const SearchField = ({
  onClickMenuItem,
  onTimetableChange,
  onToggleFilterNotResearch
}: SearchFieldProps) => {
  const handleSearchBoxChange = useCallback(
    (title: string, lecturer: string, room: string) => {
      onClickMenuItem?.('title', title ? [title] : []);
      onClickMenuItem?.('lecturer', lecturer ? [lecturer] : []);
      onClickMenuItem?.('room', room ? [room] : []);
    },
    [onClickMenuItem]
  );

  const handleClickMenuItem = useCallback(
    (key: SearchComboBox, items: string[]) => {
      onClickMenuItem?.(key, items);
    },
    [onClickMenuItem]
  );

  const handleConditionChange = useCallback(
    (key: SearchConditionKey, items: string[]) => {
      onClickMenuItem?.(key, items);
    },
    [onClickMenuItem]
  );

  return (
    <div className="search-container">
      <SearchBoxes
        onClickMenuItem={handleClickMenuItem}
        onChangeSearchBox={handleSearchBoxChange}
      />
      <SearchConditions
        onCheckItem={handleConditionChange}
        onTimetableChange={onTimetableChange}
        onToggleFilterNotResearch={onToggleFilterNotResearch}
      />
    </div>
  );
};

export default SearchField;
