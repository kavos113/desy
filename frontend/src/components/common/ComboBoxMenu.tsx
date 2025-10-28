import { useEffect, useState } from "react";
import "./common.css";
import { Menu, MenuRecord } from "./types";

type ComboBoxMenuProps = {
  items: Menu;
  onClickMenuItem?: (key: string) => void;
  className?: string;
  isSubMenu?: boolean;
};

const isMenuArray = (value: Menu): value is string[] => Array.isArray(value);

const isMenuRecord = (value: Menu): value is MenuRecord => !Array.isArray(value) && value !== null;

const buildClassName = (base: string, className?: string) =>
  [base, className].filter(Boolean).join(" ");

const ComboBoxMenu = ({ items, onClickMenuItem, className, isSubMenu = false }: ComboBoxMenuProps) => {
  const [openItems, setOpenItems] = useState<Record<string, boolean>>({});
  const [selectItems, setSelectItems] = useState<Record<string, boolean>>({});

  useEffect(() => {
    setOpenItems({});
    setSelectItems({});
  }, [items]);

  const handleMouseEnter = (key: string) => {
    setOpenItems((prev) => ({ ...prev, [key]: true }));
  };

  const handleMouseLeave = (key: string) => {
    setOpenItems((prev) => ({ ...prev, [key]: false }));
  };

  const handleClick = (key: string) => {
    onClickMenuItem?.(key);
    setSelectItems((prev) => ({ ...prev, [key]: !prev[key] }));
  };

  const renderArrayItems = (menuItems: string[], extraClassName?: string) => {
    const classes = ["subMenu", extraClassName].filter(Boolean).join(" ");

    return (
      <ul className={classes.trim()}>
        {menuItems.map((subItem) => (
          <li
            key={subItem}
            className={buildClassName(
              "menuItem",
              selectItems[subItem] ? "check" : "notCheck"
            )}
            onClick={() => handleClick(subItem)}
          >
            <span className="menuText">{subItem}</span>
            {selectItems[subItem] ? "✓" : ""}
          </li>
        ))}
      </ul>
    );
  };

  if (isMenuArray(items)) {
    const rootClass = ["menu", isSubMenu ? "subMenu" : "", className].filter(Boolean).join(" ");

    return (
      <ul className={rootClass.trim()}>
        {items.map((item) => (
          <li
            key={item}
            className={buildClassName("menuItem", selectItems[item] ? "check" : "notCheck")}
            onClick={() => handleClick(item)}
          >
            <span className="menuText">{item}</span>
            {selectItems[item] ? "✓" : ""}
          </li>
        ))}
      </ul>
    );
  }

  const rootClass = ["menu", isSubMenu ? "subMenu" : "", className].filter(Boolean).join(" ");

  return (
    <ul className={rootClass.trim()}>
      {Object.entries(items).map(([key, value]) => {
        const isArray = isMenuArray(value);
        const isRecord = isMenuRecord(value);
        const open = !!openItems[key];
        const visibilityClass = open ? "open" : "close";

        return (
          <li
            key={key}
            className="menuItem"
            onMouseEnter={() => handleMouseEnter(key)}
            onMouseLeave={() => handleMouseLeave(key)}
          >
            <div
              className={buildClassName("menuTitle", selectItems[key] ? "check" : "notCheck")}
              onClick={() => handleClick(key)}
            >
              <span className="menuText">{key}</span>
              {selectItems[key] ? "✓" : ""}
              {isRecord ? <span className="arrow">▶</span> : null}
            </div>
            {isRecord ? (
              <ComboBoxMenu
                items={value}
                onClickMenuItem={onClickMenuItem}
                className={buildClassName("subMenu", visibilityClass)}
                isSubMenu
              />
            ) : null}
            {isArray ? renderArrayItems(value, visibilityClass) : null}
          </li>
        );
      })}
    </ul>
  );
};

export default ComboBoxMenu;
