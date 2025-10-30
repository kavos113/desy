import { ChangeEvent, useEffect, useState } from 'react';
import './common.css';

type CheckBoxProps = {
  checkboxId: string;
  content: string;
  onCheckItem?: (index: number, value: boolean) => void;
  onChange?: (value: boolean) => void;
  checked?: boolean;
};

const CheckBox = ({ checkboxId, content, onCheckItem, onChange, checked }: CheckBoxProps) => {
  const [internalChecked, setInternalChecked] = useState(checked ?? false);

  useEffect(() => {
    if (typeof checked === 'boolean') {
      setInternalChecked(checked);
    }
  }, [checked]);

  const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
    const nextChecked = event.target.checked;

    setInternalChecked(nextChecked);
    onChange?.(nextChecked);

    const index = Number.parseInt(checkboxId.slice(-1), 10);
    onCheckItem?.(index, nextChecked);
  };

  return (
    <div className="check-box-container">
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
