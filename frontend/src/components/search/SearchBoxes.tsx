import { useEffect, useState } from "react";
import ComboBox from "../common/ComboBox";
import SearchBox from "../common/SearchBox";
import {
  DEPARTMENTS_MENU,
  MOBILE_DEPARTMENTS_MENU,
  UNIVERSITIES_MENU,
  YEARS_MENU,
} from "../../constants";
import type { SearchComboBox } from "./types";
import "./search.css";

type SearchBoxesProps = {
  onClickMenuItem?: (key: SearchComboBox, items: string[]) => void;
  onChangeSearchBox?: (title: string, lecturer: string) => void;
};

const SearchBoxes = ({ onClickMenuItem, onChangeSearchBox }: SearchBoxesProps) => {
  const [title, setTitle] = useState("");
  const [lecturer, setLecturer] = useState("");

  useEffect(() => {
    onChangeSearchBox?.(title, lecturer);
  }, [lecturer, onChangeSearchBox, title]);

  const handleSelect = (key: SearchComboBox) => {
    return (items: string[]) => {
      onClickMenuItem?.(key, items);
    };
  };

  return (
    <div className="search-box-wrapper">
      <ComboBox items={UNIVERSITIES_MENU} onSelectItem={handleSelect("university")}
      />
      <ComboBox
        items={DEPARTMENTS_MENU}
        className="desktop"
        onSelectItem={handleSelect("department")}
      />
      <ComboBox
        items={MOBILE_DEPARTMENTS_MENU}
        className="mobile"
        onSelectItem={handleSelect("department")}
      />
      <ComboBox items={YEARS_MENU} onSelectItem={handleSelect("year")} />
      <SearchBox
        placeholder="講義名"
        value={title}
        onChange={setTitle}
      />
      <SearchBox
        placeholder="教員名"
        value={lecturer}
        onChange={setLecturer}
      />
    </div>
  );
};

export default SearchBoxes;
