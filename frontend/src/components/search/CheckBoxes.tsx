import { useEffect, useState } from 'react';
import CheckBox from '../common/CheckBox';
import './search.css';

type CheckBoxesProps = {
  checkboxId: string;
  contents: string[];
  onCheckItem?: (items: string[]) => void;
  resetSignal?: number;
};

const buildInitialState = (length: number): boolean[] => {
  return Array.from({ length }, () => false);
};

const CheckBoxes = ({ checkboxId, contents, onCheckItem, resetSignal }: CheckBoxesProps) => {
  const [checked, setChecked] = useState<boolean[]>(() => {
    return buildInitialState(contents.length);
  });

  useEffect(() => {
    setChecked(buildInitialState(contents.length));
  }, [contents]);

  useEffect(() => {
    if (resetSignal === undefined) {
      return;
    }
    setChecked(buildInitialState(contents.length));
  }, [contents.length, resetSignal]);

  useEffect(() => {
    const selected = contents.filter((_, index) => checked[index]);
    onCheckItem?.(selected);
  }, [checked, contents, onCheckItem]);

  const handleCheck = (index: number, value: boolean) => {
    const zeroBasedIndex = index - 1;
    if (zeroBasedIndex < 0 || zeroBasedIndex >= contents.length) {
      return;
    }

    setChecked((previous) => {
      const next = [...previous];
      next[zeroBasedIndex] = value;
      return next;
    });
  };

  return (
    <div className="check-boxes-container">
      {contents.map((content, index) => {
        const id = `${checkboxId}${index + 1}`;
        return (
          <CheckBox
            key={id}
            checkboxId={id}
            content={content}
            onCheckItem={handleCheck}
            checked={checked[index]}
          />
        );
      })}
    </div>
  );
};

export default CheckBoxes;
