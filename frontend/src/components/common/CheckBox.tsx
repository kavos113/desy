import { ChangeEvent, useEffect, useState } from "react";
import "./common.css";

type CheckBoxProps = {
  checkboxId: string;
  content: string;
  checked?: boolean;
  defaultChecked?: boolean;
  className?: string;
  onCheckItem?: (index: number, value: boolean) => void;
  onChange?: (value: boolean) => void;
};

const CheckBox = ({
  checkboxId,
  content,
  checked,
  defaultChecked = false,
  className,
  onCheckItem,
  onChange,
}: CheckBoxProps) => {
  const [internalChecked, setInternalChecked] = useState(defaultChecked);
  const isControlled = typeof checked === "boolean";

  useEffect(() => {
    if (isControlled) {
      setInternalChecked(checked);
    }
  }, [checked, isControlled]);

  useEffect(() => {
    if (!isControlled) {
      setInternalChecked(defaultChecked);
    }
  }, [defaultChecked, isControlled]);

  const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
    const nextChecked = event.target.checked;

    if (!isControlled) {
      setInternalChecked(nextChecked);
    }

    onChange?.(nextChecked);

    const index = Number.parseInt(checkboxId.slice(-1), 10);
    onCheckItem?.(index, nextChecked);
  };

  const containerClass = ["check-box-container", className].filter(Boolean).join(" ");

  return (
    <div className={containerClass}>
      <input
        type="checkbox"
        id={checkboxId}
        className="check-box"
        checked={internalChecked}
        onChange={handleChange}
      />
      <label htmlFor={checkboxId} className="check-box-label">
        {content}
      </label>
    </div>
  );
};

export default CheckBox;
