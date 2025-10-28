import { ChangeEvent, useEffect, useState } from "react";
import "./common.css";

type SearchBoxProps = {
  placeholder: string;
  value?: string;
  defaultValue?: string;
  onChange?: (value: string) => void;
};

const SearchBox = ({
  placeholder,
  value,
  defaultValue = "",
  onChange,
}: SearchBoxProps) => {
  const [internalValue, setInternalValue] = useState(defaultValue);

  useEffect(() => {
    setInternalValue(defaultValue);
  }, [defaultValue]);

  const inputValue = value ?? internalValue;

  const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
    const nextValue = event.target.value;

    if (value === undefined) {
      setInternalValue(nextValue);
    }

    onChange?.(nextValue);
  };

  return (
    <div className="search-box-container">
      <input
        className="search-box-input"
        type="text"
        placeholder={placeholder}
        value={inputValue}
        onChange={handleChange}
      />
    </div>
  );
};

export default SearchBox;
