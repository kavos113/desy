import { MouseEventHandler, useMemo } from "react";
import { domain } from "../../../wailsjs/go/models";
import { formatSemesters, formatTeachers, formatTimetables } from "./utils";
import "./list.css";

type ListItemProps = {
  item: domain.LectureSummary;
  onClick?: (id: number) => void;
  className?: string;
};

const ListItem = ({ item, onClick, className }: ListItemProps) => {
  const handleClick: MouseEventHandler<HTMLDivElement> = (event) => {
    onClick?.(item.ID);
    event.stopPropagation();
  };

  const teacherText = useMemo(() => formatTeachers(item.Teachers), [item.Teachers]);
  const timetableText = useMemo(
    () => formatTimetables(item.Timetables, { includeRoom: false }),
    [item.Timetables]
  );
  const semesterText = useMemo(() => formatSemesters(item.Timetables), [item.Timetables]);
  const creditText = item.Credit !== undefined && item.Credit !== null ? String(item.Credit) : "--";

  const containerClassName = useMemo(() => {
    const classList = ["item-wrapper", "list-item"];
    if (className) {
      classList.push(className);
    }
    return classList.join(" ");
  }, [className]);

  return (
    <div className={containerClassName} onClick={handleClick}>
      <div className="item university">
        <p className="text">{item.University}</p>
      </div>
      <div className="item code">
        <p className="text">{item.Code}</p>
      </div>
      <div className="item name">
        <p className="text">{item.Title}</p>
      </div>
      <div className="item lecturer">
        <p className="text">{teacherText}</p>
      </div>
      <div className="item timetable">
        <p className="text">{timetableText}</p>
      </div>
      <div className="item semester">
        <p className="text">{semesterText}</p>
      </div>
      <div className="item department">
        <p className="text">{item.Department}</p>
      </div>
      <div className="item credit">
        <p className="text">{creditText}</p>
      </div>
    </div>
  );
};

export default ListItem;
