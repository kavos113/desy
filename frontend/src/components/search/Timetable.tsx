import { useEffect, useMemo, useState } from "react";
import "./search.css";
import { DAY_OPTIONS, PERIOD_OPTIONS, Day, Period, SearchTimetableCell } from "./types";

type TimetableProps = {
  value?: SearchTimetableCell[];
  onChange?: (items: SearchTimetableCell[]) => void;
};

type GridKey = `${Day}-${Period}`;

const toKey = (cell: SearchTimetableCell): GridKey => `${cell.day}-${cell.period}`;
const fromKey = (key: GridKey): SearchTimetableCell => {
  const [day, period] = key.split("-") as [Day, Period];
  return { day, period };
};

const Timetable = ({ value = [], onChange }: TimetableProps) => {
  const [selected, setSelected] = useState<Set<GridKey>>(() => new Set(value.map(toKey)));
  const [hoverRow, setHoverRow] = useState<Day | null>(null);
  const [hoverColumn, setHoverColumn] = useState<Period | null>(null);

  useEffect(() => {
    setSelected(new Set(value.map(toKey)));
  }, [value]);

  const selectedCells = useMemo(() => Array.from(selected).map(fromKey), [selected]);

  useEffect(() => {
    onChange?.(selectedCells);
  }, [onChange, selectedCells]);

  const toggleCell = (day: Day, period: Period) => {
    setSelected((prev) => {
      const next = new Set(prev);
      const key = `${day}-${period}` as GridKey;
      if (next.has(key)) {
        next.delete(key);
      } else {
        next.add(key);
      }
      return next;
    });
  };

  const toggleDay = (day: Day) => {
    setSelected((prev) => {
      const next = new Set(prev);
      const allKeys = PERIOD_OPTIONS.map((period) => `${day}-${period}` as GridKey);
      const shouldAdd = !allKeys.every((key) => next.has(key));
      allKeys.forEach((key) => {
        if (shouldAdd) {
          next.add(key);
        } else {
          next.delete(key);
        }
      });
      return next;
    });
  };

  const togglePeriod = (period: Period) => {
    setSelected((prev) => {
      const next = new Set(prev);
      const allKeys = DAY_OPTIONS.map((day) => `${day}-${period}` as GridKey);
      const shouldAdd = !allKeys.every((key) => next.has(key));
      allKeys.forEach((key) => {
        if (shouldAdd) {
          next.add(key);
        } else {
          next.delete(key);
        }
      });
      return next;
    });
  };

  return (
    <div className="timetable-container">
      <table className="timetable-table">
        <thead>
          <tr>
            <th></th>
            {DAY_OPTIONS.map((day) => (
              <th
                key={day}
                onMouseEnter={() => setHoverRow(day)}
                onMouseLeave={() => setHoverRow(null)}
                onClick={() => toggleDay(day)}
              >
                {day}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {PERIOD_OPTIONS.map((period) => (
            <tr key={period}>
              <th
                onMouseEnter={() => setHoverColumn(period)}
                onMouseLeave={() => setHoverColumn(null)}
                onClick={() => togglePeriod(period)}
              >
                {period}
              </th>
              {DAY_OPTIONS.map((day) => {
                const key = `${day}-${period}` as GridKey;
                const isActive = selected.has(key);
                const isHover = hoverRow === day || hoverColumn === period;
                return (
                  <td
                    key={key}
                    className={[isActive ? "is-active" : "", isHover ? "is-hover" : ""].filter(Boolean).join(" ")}
                    onClick={() => toggleCell(day, period)}
                  />
                );
              })}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default Timetable;
