import "./list.css";

type ListHeaderItemProps = {
  onSort?: (key: string) => void;
};

const ListHeaderItem = ({ onSort }: ListHeaderItemProps) => {
  const handleSort = (key: string) => () => {
    onSort?.(key);
  };

  return (
    <div className="item-wrapper header">
      <div className="item university">
        <p className="text">大学名</p>
      </div>
      <div className="item code sortable" onClick={handleSort("code")}>
        <p className="text">コード</p>
      </div>
      <div className="item name sortable" onClick={handleSort("title")}>
        <p className="text">講義名</p>
      </div>
      <div className="item lecturer sortable" onClick={handleSort("lecturer")}>
        <p className="text">担当</p>
      </div>
      <div className="item timetable sortable" onClick={handleSort("timetable")}>
        <p className="text">時間割</p>
      </div>
      <div className="item semester">
        <p className="text">開講時期</p>
      </div>
      <div className="item department sortable" onClick={handleSort("department")}>
        <p className="text">開講元</p>
      </div>
      <div className="item credit">
        <p className="text">単位数</p>
      </div>
    </div>
  );
};

export default ListHeaderItem;
