import { useCallback, useEffect, useState } from 'react';
import SimpleButton from '../common/SimpleButton';
import './search.css';
import { Greet, Scrape, ScrapeAll, ScrapeTest } from '../../../wailsjs/go/main/App';
import { EventsOn } from '../../../wailsjs/runtime/runtime';

const DEFAULT_STATUS = 'Not fetched';
const COMPLETE_STATUS = '完了しました';

type ScrapeProgressPayload = {
  total?: number;
  current?: number;
  code?: string;
  title?: string;
  Total?: number;
  Current?: number;
  Code?: string;
  Title?: string;
};

type FetchStatusEventPayload = ScrapeProgressPayload | string | number | null | undefined;

type Unsubscribe = () => void;

const isObject = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null;

const normalizeProgress = (
  payload: ScrapeProgressPayload
): { total: number; current: number; code?: string; title?: string } | null => {
  if (!isObject(payload)) {
    return null;
  }

  const totalValue = payload.total ?? payload.Total;
  if (typeof totalValue !== 'number' || Number.isNaN(totalValue)) {
    return null;
  }

  const currentValue = payload.current ?? payload.Current;
  const current =
    typeof currentValue === 'number' && !Number.isNaN(currentValue) ? currentValue : 0;
  const code = (payload.code ?? payload.Code) || undefined;
  const title = (payload.title ?? payload.Title) || undefined;

  return {
    total: Math.max(0, totalValue),
    current: Math.max(0, current),
    code: typeof code === 'string' && code.trim().length > 0 ? code.trim() : undefined,
    title: typeof title === 'string' && title.trim().length > 0 ? title.trim() : undefined
  };
};

const formatProgressStatus = (progress: {
  total: number;
  current: number;
  code?: string;
  title?: string;
}): string => {
  const { total, current, code, title } = progress;
  const safeTotal = total > 0 ? total : 0;
  const safeCurrent = current > 0 ? current : 0;
  const parts = [`${safeCurrent} / ${safeTotal}`];
  if (code) {
    parts.push(code);
  }
  if (title) {
    parts.push(title);
  }
  return parts.join(' ');
};

const FetchButton = () => {
  const [status, setStatus] = useState(DEFAULT_STATUS);
  const [isFetching, setIsFetching] = useState(false);

  useEffect(() => {
    let unsubscribe: Unsubscribe | undefined;

    try {
      const result = EventsOn('fetch_status', (payload: FetchStatusEventPayload) => {
        if (isObject(payload)) {
          const progress = normalizeProgress(payload as ScrapeProgressPayload);
          if (progress) {
            if (progress.total <= 0 || progress.current >= progress.total) {
              setStatus(COMPLETE_STATUS);
            } else {
              setStatus(formatProgressStatus(progress));
            }
            return;
          }
        }
        if (payload !== undefined && payload !== null) {
          setStatus(String(payload));
        }
      });

      if (typeof result === 'function') {
        unsubscribe = result;
      }
    } catch (error) {
      // イベント未定義でも処理を続行する
    }

    return () => {
      unsubscribe?.();
    };
  }, []);

  const handleFetch = useCallback(async () => {
    setIsFetching(true);
    setStatus('Fetching...');
    try {
      await Scrape();
      setStatus(COMPLETE_STATUS);
    } catch (error) {
      console.error('Scrape failed', error);
      setStatus('Fetch failed');
    } finally {
      setIsFetching(false);
    }
  }, []);

  const handleFetchTest = useCallback(async () => {
    setIsFetching(true);
    setStatus('Fetch-Test...');
    try {
      await ScrapeTest();
      setStatus(COMPLETE_STATUS);
    } catch (error) {
      console.error('Scrape failed', error);
      setStatus('Fetch-Test failed');
    } finally {
      setIsFetching(false);
    }
  }, []);

  const handleGreeting = useCallback(async () => {
    try {
      const message = await Greet('Fetch-Test');
      setStatus(message);
    } catch (error) {
      console.error('Greet failed', error);
      setStatus('Greeting failed');
    }
  }, []);

  const handleFetchAll = useCallback(async () => {
    setIsFetching(true);
    setStatus('Fetch-All...');
    try {
      await ScrapeAll();
      setStatus(COMPLETE_STATUS);
    } catch (error) {
      console.error('ScrapeAll failed', error);
      setStatus('Fetch-All failed');
    } finally {
      setIsFetching(false);
    }
  }, []);

  return (
    <div>
      <div className="button-wrapper">
        <SimpleButton text="Fetch" onClick={handleFetch} disabled={isFetching} />
        <SimpleButton text="Fetch-All" onClick={handleFetchAll} disabled={isFetching} />
      </div>
      <div>
        <p className="fetch-status">{status}</p>
      </div>
    </div>
  );
};

export default FetchButton;
