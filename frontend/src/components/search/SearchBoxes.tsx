import { useCallback, useEffect, useRef, useState } from 'react';
import ComboBox from '../common/ComboBox';
import SearchBox from '../common/SearchBox';
import {
  DEPARTMENTS_MENU,
  MOBILE_DEPARTMENTS_MENU,
  UNIVERSITIES_MENU,
  YEARS_MENU
} from '../../constants';
import type { SearchComboBox } from './types';
import './search.css';

type SearchBoxesProps = {
  onClickMenuItem?: (key: SearchComboBox, items: string[]) => void;
  onChangeSearchBox?: (title: string, lecturer: string, room: string) => void;
  resetSignal?: number;
};

const SearchBoxes = ({ onClickMenuItem, onChangeSearchBox, resetSignal }: SearchBoxesProps) => {
  const [title, setTitle] = useState('');
  const [lecturer, setLecturer] = useState('');
  const [room, setRoom] = useState('');
  const isFirstRender = useRef(true);

  const handleTitleChange = useCallback(
    (value: string) => {
      setTitle(value);
      onChangeSearchBox?.(value, lecturer, room);
    },
    [lecturer, onChangeSearchBox, room]
  );

  const handleLecturerChange = useCallback(
    (value: string) => {
      setLecturer(value);
      onChangeSearchBox?.(title, value, room);
    },
    [onChangeSearchBox, room, title]
  );

  const handleRoomChange = useCallback(
    (value: string) => {
      setRoom(value);
      onChangeSearchBox?.(title, lecturer, value);
    },
    [lecturer, onChangeSearchBox, title]
  );

  useEffect(() => {
    if (resetSignal === undefined) {
      return;
    }
    if (isFirstRender.current) {
      isFirstRender.current = false;
      return;
    }
    setTitle('');
    setLecturer('');
    setRoom('');
    onChangeSearchBox?.('', '', '');
  }, [onChangeSearchBox, resetSignal]);

  const handleSelect = (key: SearchComboBox) => {
    return (items: string[]) => {
      onClickMenuItem?.(key, items);
    };
  };

  return (
    <div className="search-box-wrapper">
      <ComboBox
        key={`university-${resetSignal}`}
        items={UNIVERSITIES_MENU}
        onSelectItem={handleSelect('university')}
      />
      <ComboBox
        key={`department-desktop-${resetSignal}`}
        items={DEPARTMENTS_MENU}
        className="desktop"
        onSelectItem={handleSelect('department')}
      />
      <ComboBox
        key={`department-mobile-${resetSignal}`}
        items={MOBILE_DEPARTMENTS_MENU}
        className="mobile"
        onSelectItem={handleSelect('department')}
      />
      <ComboBox
        key={`year-${resetSignal}`}
        items={YEARS_MENU}
        onSelectItem={handleSelect('year')}
      />
      <SearchBox placeholder="講義名" value={title} onChange={handleTitleChange} />
      <SearchBox placeholder="教員名" value={lecturer} onChange={handleLecturerChange} />
      <SearchBox placeholder="講義室名" value={room} onChange={handleRoomChange} />
    </div>
  );
};

export default SearchBoxes;
