import { useCallback, useEffect, useState } from "react";
import SimpleButton from "../../common/SimpleButton";
import "./search.css";
import { Greet, Scrape } from "../../../wailsjs/go/main/App";
import { EventsOn } from "../../../wailsjs/runtime/runtime";

const DEFAULT_STATUS = "Not fetched";

type FetchStatusEventPayload = unknown;

type Unsubscribe = () => void;

const FetchButton = () => {
  const [status, setStatus] = useState(DEFAULT_STATUS);
  const [isFetching, setIsFetching] = useState(false);

  useEffect(() => {
    let unsubscribe: Unsubscribe | undefined;

    try {
      const result = EventsOn("fetch_status", (payload: FetchStatusEventPayload) => {
        setStatus(String(payload));
      });

      if (typeof result === "function") {
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
    setStatus("Fetching...");
    try {
      await Scrape();
      setStatus("Fetched");
    } catch (error) {
      console.error("Scrape failed", error);
      setStatus("Fetch failed");
    } finally {
      setIsFetching(false);
    }
  }, []);

  const handleGreeting = useCallback(async () => {
    try {
      const message = await Greet("Fetch-Test");
      setStatus(message);
    } catch (error) {
      console.error("Greet failed", error);
      setStatus("Greeting failed");
    }
  }, []);

  return (
    <div>
      <div className="button-wrapper">
        <SimpleButton text="Fetch" onClick={handleFetch} disabled={isFetching} />
        <SimpleButton text="Fetch-Test" onClick={handleGreeting} />
      </div>
      <div>
        <p className="fetch-status">{status}</p>
      </div>
    </div>
  );
};

export default FetchButton;
