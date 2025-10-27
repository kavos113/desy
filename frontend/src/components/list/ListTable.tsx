import { useEffect, useState } from 'react';
import { GetLectureDetails } from '../../../wailsjs/go/main/App';
import { domain } from '../../../wailsjs/go/models';
import CourseDetailPanel from './CourseDetailPanel';
import ListHeaderItem, { ListSortKey } from './ListHeaderItem';
import ListItemRow from './ListItemRow';
import './list.css';

export type ListTableProps = {
  items: domain.LectureSummary[];
  loading?: boolean;
  statusMessage?: string | null;
  errorMessage?: string | null;
  onSort?: (key: ListSortKey) => void;
};

const ListTable = ({ items, loading = false, statusMessage, errorMessage, onSort }: ListTableProps) => {
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [detail, setDetail] = useState<domain.Lecture | null>(null);
  const [detailLoading, setDetailLoading] = useState(false);
  const [detailError, setDetailError] = useState<string | null>(null);
  const [isDetailOpen, setDetailOpen] = useState(false);
  const [isOverlayInteractive, setOverlayInteractive] = useState(false);

  useEffect(() => {
    setSelectedId(null);
    setDetail(null);
    setDetailError(null);
    setDetailOpen(false);
    setOverlayInteractive(false);
  }, [items]);

  const handleSort = (key: ListSortKey) => {
    if (onSort) {
      onSort(key);
    }
  };

  const handleSelect = async (lectureId: number) => {
    setSelectedId(lectureId);
    setDetailLoading(true);
    setDetailError(null);
    setDetailOpen(true);
    setOverlayInteractive(true);

    try {
      const data = await GetLectureDetails(lectureId);
      setDetail(data ?? null);
    } catch (error) {
      console.error(error);
      setDetailError('講義詳細の取得に失敗しました。');
    } finally {
      setDetailLoading(false);
    }
  };

  const handleCloseDetail = () => {
    setDetailOpen(false);
    setTimeout(() => {
      setOverlayInteractive(false);
    }, 250);
  };

  const overlayClassNames = [
    'overlay',
    isDetailOpen ? 'overlay--visible' : '',
    isOverlayInteractive ? 'overlay--interactive' : ''
  ]
    .filter(Boolean)
    .join(' ');

  const hasMessage = Boolean(statusMessage || errorMessage);

  return (
    <div className="list-table">
      <div className="list-table-header">
        <ListHeaderItem onSort={handleSort} />
        {hasMessage ? (
          <div className="list-table-messages">
            {statusMessage && <span className="list-table-status">{statusMessage}</span>}
            {errorMessage && <span className="list-table-error">{errorMessage}</span>}
          </div>
        ) : null}
      </div>

      {loading ? (
        <p className="list-table-placeholder">検索結果を読み込み中です...</p>
      ) : items.length === 0 ? (
        <p className="list-table-placeholder">検索条件を入力して講義を探してください。</p>
      ) : (
        <div className="list-table-scroll">
          {items.map((item) => (
            <ListItemRow key={item.ID} item={item} onSelect={handleSelect} selected={item.ID === selectedId} />
          ))}
        </div>
      )}

      <CourseDetailPanel lecture={detail} open={isDetailOpen} loading={detailLoading} errorMessage={detailError} />

      <button
        type="button"
        className={`back-button${isDetailOpen ? ' back-button--visible' : ''}`}
        onClick={handleCloseDetail}
      >
        戻る
      </button>

      <div className={overlayClassNames} onClick={handleCloseDetail} role="presentation" />
    </div>
  );
};

export default ListTable;
