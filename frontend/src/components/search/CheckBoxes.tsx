import { useEffect, useMemo, useState } from "react";
import { CheckBox as CommonCheckBox } from "../common";
import "./search.css";

type CheckBoxesProps = {
  idPrefix: string;
  contents: string[];
  selectedItems?: string[];
  onChange?: (items: string[]) => void;
};

type CheckedMap = Record<string, boolean>;

const buildCheckedMap = (contents: string[], selected: string[]): CheckedMap => {
  return contents.reduce<CheckedMap>((acc, item) => {
    acc[item] = selected.includes(item);
    return acc;
  }, {});
};

const CheckBoxes = ({ idPrefix, contents, selectedItems = [], onChange }: CheckBoxesProps) => {
  const [checked, setChecked] = useState<CheckedMap>(() => buildCheckedMap(contents, selectedItems));

  useEffect(() => {
    setChecked(buildCheckedMap(contents, selectedItems));
  }, [contents, selectedItems]);

  const selected = useMemo(() => contents.filter((item) => checked[item]), [checked, contents]);

  useEffect(() => {
    onChange?.(selected);
  }, [onChange, selected]);

  const toggle = (item: string, value?: boolean) => {
    setChecked((prev) => {
      const nextValue = typeof value === "boolean" ? value : !prev[item];
      return {
        ...prev,
        [item]: nextValue,
      };
    });
  };

  return (
    <div className="check-boxes-container">
      {contents.map((item, index) => {
        const elementId = `${idPrefix}-${index}`;
        return (
          <CommonCheckBox
            key={elementId}
            checkboxId={elementId}
            content={item}
            checked={Boolean(checked[item])}
            onChange={(value) => toggle(item, value)}
          />
        );
      })}
    </div>
  );
};

export default CheckBoxes;
