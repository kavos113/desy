import { render, screen, waitFor } from '@testing-library/react';
import { act } from 'react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import Search from '../Search';

vi.mock('../FetchButton', () => ({
  __esModule: true,
  default: () => <div data-testid="fetch-button" />
}));

const searchLecturesMock = vi.fn<(query: any) => Promise<any>>(() => Promise.resolve([]));

vi.mock('../../../../wailsjs/go/main/App', () => ({
  SearchLectures: (query: any) => searchLecturesMock(query)
}));

describe('Search', () => {
  beforeEach(() => {
    searchLecturesMock.mockClear();
    searchLecturesMock.mockResolvedValue([]);
  });

  it('検索時に研究系科目除外の指定を送信できる', async () => {
    const user = userEvent.setup();
    render(<Search />);

    const checkbox = screen.getByLabelText('研究系科目を除外');
    await act(async () => {
      await user.click(checkbox);
      await user.click(screen.getByRole('button', { name: 'Search' }));
    });

    await waitFor(() => {
      expect(searchLecturesMock).toHaveBeenCalledTimes(1);
    });

    const [[query]] = searchLecturesMock.mock.calls;
    expect(query.FilterNotResearch ?? false).toBe(true);
  });

  it('講義室名を検索条件として送信できる', async () => {
    const user = userEvent.setup();
    render(<Search />);

    const roomInput = screen.getByPlaceholderText('講義室名');
    await act(async () => {
      await user.type(roomInput, '本館');
      await user.click(screen.getByRole('button', { name: 'Search' }));
    });

    await waitFor(() => {
      expect(searchLecturesMock).toHaveBeenCalledTimes(1);
    });

    const [[query]] = searchLecturesMock.mock.calls;
    expect(query.Room ?? '').toBe('本館');
  });
});
