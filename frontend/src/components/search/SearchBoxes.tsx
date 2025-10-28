import { useState } from "react";
import "./search.css";
import { ComboBox, SearchBox } from "../common";
import { SearchComboBox } from "./types";
import {
  DEPARTMENT_MENU,
  MOBILE_DEPARTMENT_MENU,
  UNIVERSITY_MENU,
  YEAR_MENU,
} from "./menus";

type SearchBoxesProps = {
  onSelectMenuItem?: (key: SearchComboBox, items: string[]) => void;
  onChangeSearchBox?: (title: string, lecturer: string) => void;
};

const SearchBoxes = ({ onSelectMenuItem, onChangeSearchBox }: SearchBoxesProps) => {
  const [title, setTitle] = useState("");
  const [lecturer, setLecturer] = useState("");

  const handleSelect = (type: SearchComboBox) => (items: string[]) => {
    onSelectMenuItem?.(type, items);
  };

  const handleTitleChange = (value: string) => {
    setTitle(value);
    onChangeSearchBox?.(value, lecturer);
  };

  const handleLecturerChange = (value: string) => {
    setLecturer(value);
    onChangeSearchBox?.(title, value);
  };

  return (
    <div className="search-box-wrapper">
      <div className="search-box-row">
        <ComboBox items={UNIVERSITY_MENU} onSelectItem={handleSelect("university")} />
        <ComboBox
          className="desktop"
          items={DEPARTMENT_MENU}
          onSelectItem={handleSelect("department")}
        />
        <ComboBox
          className="mobile"
          items={MOBILE_DEPARTMENT_MENU}
          onSelectItem={handleSelect("department")}
        />
        <ComboBox items={YEAR_MENU} onSelectItem={handleSelect("year")} />
      </div>
      <div className="search-box-row search-box-row--inputs">
        <SearchBox
          placeholder="講義名"
          value={title}
          onChange={handleTitleChange}
        />
        <SearchBox
          placeholder="教員名"
          value={lecturer}
          onChange={handleLecturerChange}
        />
      </div>
    </div>
  );
};

export default SearchBoxes;
