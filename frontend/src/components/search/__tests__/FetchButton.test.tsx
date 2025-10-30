import { act, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

const fetchStatusListeners: Array<(payload: unknown) => void> = [];

const { mockGreet, mockScrape, mockScrapeTest } = vi.hoisted(() => ({
  mockGreet: vi.fn(),
  mockScrape: vi.fn(),
  mockScrapeTest: vi.fn()
}));

vi.mock('../../../../wailsjs/runtime/runtime', () => {
  return {
    EventsOn: vi.fn((event: string, callback: (payload: unknown) => void) => {
      if (event === 'fetch_status') {
        fetchStatusListeners.push(callback);
      }
      return () => {
        if (event !== 'fetch_status') {
          return;
        }
        const index = fetchStatusListeners.indexOf(callback);
        if (index !== -1) {
          fetchStatusListeners.splice(index, 1);
        }
      };
    })
  };
});

vi.mock('../../../../wailsjs/go/main/App', () => ({
  Greet: mockGreet,
  Scrape: mockScrape,
  ScrapeTest: mockScrapeTest
}));

import FetchButton from '../FetchButton';

describe('FetchButton', () => {
  beforeEach(() => {
    fetchStatusListeners.length = 0;
    mockGreet.mockReset();
    mockScrape.mockReset();
    mockScrapeTest.mockReset();
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it('進捗イベントが完了メッセージを表示する', async () => {
    render(<FetchButton />);

    await waitFor(() => {
      expect(fetchStatusListeners.length).toBeGreaterThan(0);
    });
    const listener = fetchStatusListeners[0];

    act(() => {
      listener({ Total: 120, Current: 12, Code: 'EEE.A123', Title: '講義名' });
    });

    expect(screen.getByText('12 / 120 EEE.A123 講義名')).toBeInTheDocument();

    act(() => {
      listener({ Total: 120, Current: 120, Code: 'EEE.A123', Title: '講義名' });
    });

    expect(screen.getByText('完了しました')).toBeInTheDocument();
  });

  it('Fetchボタン成功後に完了メッセージを表示する', async () => {
    mockScrape.mockResolvedValue(undefined);
    const user = userEvent.setup();

    render(<FetchButton />);

    const fetchButton = screen.getByRole('button', { name: 'Fetch' });
    await act(async () => {
      await user.click(fetchButton);
    });

    expect(mockScrape).toHaveBeenCalled();

    await screen.findByText('完了しました');
  });
});
