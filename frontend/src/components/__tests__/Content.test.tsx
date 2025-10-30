import { act, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { domain } from '../../../wailsjs/go/models';

const mockLectures = [
  domain.LectureSummary.createFrom({
    ID: 1,
    University: '東京大学',
    Department: '工学部',
    Title: 'B講義',
    Code: 'B001',
    Timetables: [],
    Teachers: []
  }),
  domain.LectureSummary.createFrom({
    ID: 2,
    University: '東京大学',
    Department: '工学部',
    Title: 'A講義',
    Code: 'A001',
    Timetables: [],
    Teachers: []
  })
];

vi.mock('../search/Search', () => {
  const MockSearch = ({ className, onSearch, onBack }: any) => (
    <div data-testid="search-panel" className={className}>
      <button
        type="button"
        onClick={() => {
          onSearch?.(mockLectures);
          onBack?.();
        }}
      >
        trigger-search
      </button>
    </div>
  );

  return {
    __esModule: true,
    default: MockSearch
  };
});

vi.mock('../list/ListTable', () => {
  const MockListTable = ({ items, className, onSort }: any) => (
    <div>
      <div data-testid="list-titles" className={className}>
        {items.map((item: domain.LectureSummary) => item.Title).join(',')}
      </div>
      <button type="button" onClick={() => onSort?.('title')}>
        sort-title
      </button>
    </div>
  );

  return {
    __esModule: true,
    default: MockListTable
  };
});

import Content from '../Content';

describe('Content', () => {
  it('検索結果を受け取り並び替えができる', async () => {
    const user = userEvent.setup();
    await act(async () => {
      render(<Content />);
    });

    await act(async () => {
      await user.click(screen.getByRole('button', { name: 'trigger-search' }));
    });

    const list = await screen.findByTestId('list-titles');
    expect(list.textContent).toBe('B講義,A講義');

    await act(async () => {
      await user.click(screen.getByRole('button', { name: 'sort-title' }));
    });
    expect(list.textContent).toBe('A講義,B講義');

    await act(async () => {
      await user.click(screen.getByRole('button', { name: 'sort-title' }));
    });
    expect(list.textContent).toBe('B講義,A講義');
  });
});
