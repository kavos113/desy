import { useEffect, useState } from "react";
import "./common.css";
import ComboBoxMenu from "./ComboBoxMenu";
import { Menu } from "./types";

type ComboBoxProps = {
  items: Menu;
  className?: string;
  initialSelectedItems?: string[];
  onSelectItem?: (items: string[]) => void;
};

const toUniqueSelection = (items: string[]) => Array.from(new Set(items));

const ComboBox = ({ items, className, initialSelectedItems = [], onSelectItem }: ComboBoxProps) => {
  const [selectedItems, setSelectedItems] = useState<string[]>(() => toUniqueSelection(initialSelectedItems));

  useEffect(() => {
    setSelectedItems(toUniqueSelection(initialSelectedItems));
  }, [initialSelectedItems]);

  const handleSelect = (key: string) => {
    setSelectedItems((prev) => {
      const exists = prev.includes(key);
      const next = exists ? prev.filter((item) => item !== key) : [...prev, key];
      onSelectItem?.(next);
      return next;
    });
  };

  const containerClass = ["combobox-box-container", className].filter(Boolean).join(" ");

  return (
    <div className={containerClass}>
      <div className="comboMenu">
        <ComboBoxMenu items={items} onClickMenuItem={handleSelect} />
      </div>
    </div>
  );
};

export default ComboBox;
