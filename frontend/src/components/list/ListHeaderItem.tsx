import './list.css';

export type ListSortKey = 'code' | 'title' | 'lecturer' | 'timetable' | 'department';

type ListHeaderItemProps = {
  onSort?: (key: ListSortKey) => void;
};

const ListHeaderItem = ({ onSort }: ListHeaderItemProps) => {
  const handleSort = (key: ListSortKey) => {
    if (onSort) {
      onSort(key);
    }
  };

  return (
    <div className="list-header">
      <span className="list-cell university">大学名</span>
      <button type="button" className="list-cell code" onClick={() => handleSort('code')}>
        コード
      </button>
      <button type="button" className="list-cell name" onClick={() => handleSort('title')}>
        講義名
      </button>
      <button type="button" className="list-cell lecturer" onClick={() => handleSort('lecturer')}>
        担当
      </button>
      <button type="button" className="list-cell timetable" onClick={() => handleSort('timetable')}>
        時間割
      </button>
      <span className="list-cell semester">開講時期</span>
      <button type="button" className="list-cell department" onClick={() => handleSort('department')}>
        開講元
      </button>
      <span className="list-cell credit">単位数</span>
    </div>
  );
};

export default ListHeaderItem;
