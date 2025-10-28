import { useEffect, useMemo, useState } from "react";
import { DAYS, Day, PERIODS, Period } from "../../constants";
import type { SearchTimetableSelection } from "./types";
import "./search.css";

type TimetableState = Record<Day, Record<Period, boolean>>;

type TimetableProps = {
  onCheckItem?: (items: SearchTimetableSelection[]) => void;
};

const createInitialState = (): TimetableState => {
  return DAYS.reduce<TimetableState>((acc, day) => {
    acc[day] = PERIODS.reduce<Record<Period, boolean>>((rowAcc, period) => {
      rowAcc[period] = false;
      return rowAcc;
    }, {} as Record<Period, boolean>);
    return acc;
  }, {} as TimetableState);
};

const Timetable = ({ onCheckItem }: TimetableProps) => {
  const [checked, setChecked] = useState<TimetableState>(() => createInitialState());

  useEffect(() => {
    const selected: SearchTimetableSelection[] = [];
    for (const day of DAYS) {
      for (const period of PERIODS) {
        if (checked[day][period]) {
          selected.push({ day, period });
        }
      }
    }
    onCheckItem?.(selected);
  }, [checked, onCheckItem]);

  const toggleCell = (day: Day, period: Period) => {
    setChecked((previous) => {
      const next: TimetableState = { ...previous, [day]: { ...previous[day] } };
      next[day][period] = !next[day][period];
      return next;
    });
  };

  const toggleDay = (day: Day) => {
    setChecked((previous) => {
      const row = { ...previous[day] };
      for (const period of PERIODS) {
        row[period] = !row[period];
      }
      return { ...previous, [day]: row };
    });
  };

  const togglePeriod = (period: Period) => {
    setChecked((previous) => {
      const next: TimetableState = { ...previous };
      for (const day of DAYS) {
        const row = { ...next[day] };
        row[period] = !row[period];
        next[day] = row;
      }
      return next;
    });
  };

  const dayHeaders = useMemo(() => DAYS, []);
  const periodRows = useMemo(() => PERIODS, []);

  return (
    <div className="timetable-container">
      <table className="table">
        <thead>
          <tr className="days">
            <th className="left"></th>
            {dayHeaders.map((day) => (
              <th
                key={day}
                className={`mainContent day day-${day}`}
                onClick={() => toggleDay(day)}
              >
                {day}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {periodRows.map((period) => (
            <tr key={period} className="row">
              <th className="left period" onClick={() => togglePeriod(period)}>
                {period}
              </th>
              {dayHeaders.map((day) => {
                const isChecked = checked[day][period];
                const classNames = [
                  "mainContent",
                  "box",
                  `box-${day}`,
                  isChecked ? "checked" : "",
                ]
                  .filter(Boolean)
                  .join(" ");

                return (
                  <td
                    key={`${day}-${period}`}
                    className={classNames}
                    onClick={() => toggleCell(day, period)}
                  ></td>
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
