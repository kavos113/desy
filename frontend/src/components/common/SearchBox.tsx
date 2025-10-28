import { ChangeEvent, InputHTMLAttributes } from "react";
import "./common.css";

type SearchBoxProps = {
  placeholder?: string;
  value?: string;
  defaultValue?: string;
  className?: string;
  onChange?: (value: string) => void;
} & Omit<InputHTMLAttributes<HTMLInputElement>, "type" | "value" | "defaultValue" | "placeholder" | "onChange" | "className">;

const SearchBox = ({
  placeholder,
  value,
  defaultValue,
  className,
  onChange,
  ...rest
}: SearchBoxProps) => {
  const containerClass = ["search-box-container", className].filter(Boolean).join(" ");

  const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange?.(event.target.value);
  };

  return (
    <div className={containerClass}>
      <input
        {...rest}
        type="text"
        className="search-box-input"
        value={value}
        defaultValue={defaultValue}
        placeholder={placeholder}
        onChange={handleChange}
      />
    </div>
  );
};

export default SearchBox;
