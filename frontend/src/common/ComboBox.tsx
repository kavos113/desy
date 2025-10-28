import { useEffect, useMemo, useState } from "react";
import ComboBoxMenu from "./ComboBoxMenu";
import { Menu } from "./Menu";
import "./common.css";

type ComboBoxProps = {
  items: Menu;
  defaultSelectedItems?: string[];
  onSelectItem?: (selected: string[]) => void;
};

const ComboBox = ({ items, defaultSelectedItems, onSelectItem }: ComboBoxProps) => {
  const [selectedItems, setSelectedItems] = useState<string[]>(() => {
    return defaultSelectedItems ? [...defaultSelectedItems] : [];
  });

  useEffect(() => {
    if (defaultSelectedItems) {
      setSelectedItems([...defaultSelectedItems]);
    }
  }, [defaultSelectedItems]);

  const selectedKeySet = useMemo(() => {
    return new Set(selectedItems);
  }, [selectedItems]);

  const handleSelect = (key: string) => {
    setSelectedItems((previous) => {
      const exists = previous.includes(key);
      const next = exists
        ? previous.filter((value) => value !== key)
        : [...previous, key];

      onSelectItem?.(next);
      return next;
    });
  };

  return (
    <div className="combobox-box-container">
      <div className="comboMenu">
        <ComboBoxMenu items={items} selectedKeys={selectedKeySet} onSelect={handleSelect} />
      </div>
    </div>
  );
};

export default ComboBox;
