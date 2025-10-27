import './list.css';
import { domain } from '../../../wailsjs/go/models';
import { LEVEL_LABELS } from '../../constants';

type ListItemRowProps = {
  item: domain.LectureSummary;
  onSelect: (id: number) => void;
  selected?: boolean;
};

const toTimetableText = (timetables: domain.TimeTable[] | undefined): string => {
  if (!timetables?.length) {
    return '';
  }

  return timetables
    .map((timetable) => {
      const semester = timetable.Semester ? `${timetable.Semester}` : '';
      const day = timetable.DayOfWeek ? `${timetable.DayOfWeek}` : '';
      const period = timetable.Period ? `${timetable.Period}限` : '';
      const room = timetable.Room?.Name ? `(${timetable.Room.Name})` : '';
      return [semester, day, period].filter(Boolean).join('') + room;
    })
    .join(', ');
};

const ListItemRow = ({ item, onSelect, selected }: ListItemRowProps) => {
  return (
    <div
      className={`list-item-row${selected ? ' selected' : ''}`}
      onClick={() => onSelect(item.ID)}
      role="button"
      tabIndex={0}
      onKeyDown={(event) => {
        if (event.key === 'Enter' || event.key === ' ') {
          event.preventDefault();
          onSelect(item.ID);
        }
      }}
    >
      <span className="list-cell university">{item.University || '-'}</span>
      <span className="list-cell code">{item.Code || '-'}</span>
      <span className="list-cell name">{item.Title || '-'}</span>
      <span className="list-cell lecturer">
        {(item.Teachers ?? []).map((teacher: domain.Teacher) => teacher.Name).join(', ') || '-'}
      </span>
      <span className="list-cell timetable">{toTimetableText(item.Timetables)}</span>
      <span className="list-cell semester">{item.Year ? `${item.Year}年度` : '-'}</span>
    <span className="list-cell department">{item.Department || '-'}</span>
    <span className="list-cell credit">{item.Level ? LEVEL_LABELS[item.Level] ?? item.Level : '-'}</span>
    </div>
  );
};

export default ListItemRow;
