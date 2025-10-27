import { ChangeEvent, useMemo, useState } from "react";
import "./search.css";
import { DEPARTMENT_OPTIONS, YEAR_OPTIONS } from "../../constants";
import { SearchComboBox } from "./types";

type SearchBoxesProps = {
  onSelectMenuItem?: (key: SearchComboBox, items: string[]) => void;
  onChangeSearchBox?: (title: string, lecturer: string) => void;
};

const UNIVERSITY_OPTIONS = ["東京工業大学", "一橋大学"];

const toStringOptions = (options: Array<string | number>): string[] => options.map((option) => String(option));

const SearchBoxes = ({ onSelectMenuItem, onChangeSearchBox }: SearchBoxesProps) => {
  const yearOptions = useMemo(() => toStringOptions(YEAR_OPTIONS), []);

  const [universities, setUniversities] = useState<string[]>([]);
  const [departments, setDepartments] = useState<string[]>([]);
  const [years, setYears] = useState<string[]>([]);
  const [title, setTitle] = useState("");
  const [lecturer, setLecturer] = useState("");

  const notifyTextChange = (nextTitle: string, nextLecturer: string) => {
    onChangeSearchBox?.(nextTitle, nextLecturer);
  };

  const handleSelect = (event: ChangeEvent<HTMLSelectElement>, type: SearchComboBox) => {
    const selected = Array.from(event.target.selectedOptions).map((option) => option.value);

    switch (type) {
      case "university":
        setUniversities(selected);
        onSelectMenuItem?.("university", selected);
        break;
      case "department":
        setDepartments(selected);
        onSelectMenuItem?.("department", selected);
        break;
      case "year":
        setYears(selected);
        onSelectMenuItem?.("year", selected);
        break;
      default:
        break;
    }
  };

  const handleTitleChange = (value: string) => {
    setTitle(value);
    notifyTextChange(value, lecturer);
  };

  const handleLecturerChange = (value: string) => {
    setLecturer(value);
    notifyTextChange(title, value);
  };

  return (
    <div className="search-box-wrapper">
      <div className="search-box-row">
        <select
          className="search-box-select"
          multiple
          value={universities}
          onChange={(event) => handleSelect(event, "university")}
        >
          {UNIVERSITY_OPTIONS.map((option) => (
            <option key={option} value={option}>
              {option}
            </option>
          ))}
        </select>
        <select
          className="search-box-select"
          multiple
          value={departments}
          onChange={(event) => handleSelect(event, "department")}
        >
          {DEPARTMENT_OPTIONS.map((option) => (
            <option key={option} value={option}>
              {option}
            </option>
          ))}
        </select>
        <select
          className="search-box-select mobile"
          multiple
          value={departments}
          onChange={(event) => handleSelect(event, "department")}
        >
          {DEPARTMENT_OPTIONS.map((option) => (
            <option key={option} value={option}>
              {option}
            </option>
          ))}
        </select>
        <select
          className="search-box-select desktop"
          multiple
          value={years}
          onChange={(event) => handleSelect(event, "year")}
        >
          {yearOptions.map((option) => (
            <option key={option} value={option}>
              {option}
            </option>
          ))}
        </select>
      </div>
      <div className="search-box-row">
        <input
          value={title}
          placeholder="講義名"
          onChange={(event) => handleTitleChange(event.target.value)}
        />
        <input
          value={lecturer}
          placeholder="教員名"
          onChange={(event) => handleLecturerChange(event.target.value)}
        />
      </div>
    </div>
  );
};

export default SearchBoxes;
