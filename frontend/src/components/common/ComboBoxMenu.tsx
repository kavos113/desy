import { useMemo, useState } from "react";
import "./common.css";
import { Menu, isMenuLeaf, isMenuNode } from "./Menu";

type ComboBoxMenuProps = {
  items: Menu;
  selectedKeys: Set<string>;
  onSelect: (key: string) => void;
  className?: string;
};

const ComboBoxMenu = ({
  items,
  selectedKeys,
  onSelect,
  className,
}: ComboBoxMenuProps) => {
  const [openItems, setOpenItems] = useState<Record<string, boolean>>({});

  const containerClassName = useMemo(() => {
    return ["menu", className].filter(Boolean).join(" ");
  }, [className]);

  const showSubMenu = (key: string) => {
    setOpenItems((previous) => {
      if (previous[key]) {
        return previous;
      }
      return { ...previous, [key]: true };
    });
  };

  const hideSubMenu = (key: string) => {
    setOpenItems((previous) => {
      if (!previous[key]) {
        return previous;
      }
      return { ...previous, [key]: false };
    });
  };

  const isOpen = (key: string) => {
    return Boolean(openItems[key]);
  };

  if (isMenuLeaf(items)) {
    return (
      <ul className={containerClassName}>
        {items.map((item) => {
          const checked = selectedKeys.has(item);
          return (
            <li
              key={item}
              className={`menuItem ${checked ? "check" : "notCheck"}`}
              onClick={() => onSelect(item)}
            >
              <span className="menuText">{item}</span>
              {checked ? "\u2713" : ""}
            </li>
          );
        })}
      </ul>
    );
  }

  if (!isMenuNode(items)) {
    return null;
  }

  return (
    <ul className={containerClassName}>
      {Object.entries(items).map(([key, value]) => {
        const hasSubMenu = isMenuNode(value) || isMenuLeaf(value);
        const checked = selectedKeys.has(key);
        const open = hasSubMenu && isOpen(key);

        return (
          <li
            key={key}
            className="menuItem"
            onMouseEnter={() => (hasSubMenu ? showSubMenu(key) : undefined)}
            onMouseLeave={() => (hasSubMenu ? hideSubMenu(key) : undefined)}
          >
            <button
              type="button"
              className={`menuTitle ${checked ? "check" : "notCheck"}`}
              onClick={() => onSelect(key)}
            >
              <span className="menuText">{key}</span>
              {checked ? "\u2713" : ""}
              {isMenuNode(value) ? <span className="arrow">\u25B6</span> : null}
            </button>
            {isMenuNode(value) ? (
              <ComboBoxMenu
                items={value}
                selectedKeys={selectedKeys}
                onSelect={onSelect}
                className={`subMenu ${open ? "open" : "close"}`}
              />
            ) : null}
            {isMenuLeaf(value) ? (
              <ComboBoxMenu
                items={value}
                selectedKeys={selectedKeys}
                onSelect={onSelect}
                className={`subMenu ${open ? "open" : "close"}`}
              />
            ) : null}
          </li>
        );
      })}
    </ul>
  );
};

export default ComboBoxMenu;
