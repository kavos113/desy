import { useEffect, useState } from "react";
import "./search.css";
import { Greet, Scrape } from "../../../wailsjs/go/main/App";
import { EventsOn } from "../../../wailsjs/runtime/runtime";
import { SimpleButton } from "../common";

type FetchButtonProps = {
  className?: string;
};

const FetchButton = ({ className }: FetchButtonProps) => {
  const [status, setStatus] = useState("未取得");
  const [loading, setLoading] = useState(false);
  const [testing, setTesting] = useState(false);

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }

    try {
      const off = EventsOn("fetch_status", (message?: string) => {
        if (typeof message === "string" && message.length > 0) {
          setStatus(message);
        }
      });

      return () => {
        off?.();
      };
    } catch (error) {
      console.warn("Failed to subscribe fetch_status event", error);
      return undefined;
    }
  }, []);

  const handleScrape = async () => {
    if (loading) {
      return;
    }
    setLoading(true);
    setStatus("スクレイピングを実行中です...");
    try {
      await Scrape();
      setStatus("最新のシラバス情報を取得しました。");
    } catch (error) {
      console.error(error);
      setStatus("スクレイピングに失敗しました。");
    } finally {
      setLoading(false);
    }
  };

  const handleTest = async () => {
    if (testing) {
      return;
    }
    setTesting(true);
    setStatus("テスト呼び出しを実行中です...");
    try {
      const greeting = await Greet("Fetch");
      setStatus(greeting);
    } catch (error) {
      console.error(error);
      setStatus("テスト呼び出しに失敗しました。");
    } finally {
      setTesting(false);
    }
  };

  return (
    <div className={["fetch-panel", className].filter(Boolean).join(" ")}>
      <div className="fetch-buttons">
        <SimpleButton
          text={loading ? "Fetching..." : "Fetch"}
          className="secondary"
          onClick={handleScrape}
          disabled={loading}
        />
        <SimpleButton
          text={testing ? "Testing..." : "Fetch-Test"}
          className="ghost"
          onClick={handleTest}
          disabled={testing}
        />
      </div>
      <p className="search-summary">{status}</p>
    </div>
  );
};

export default FetchButton;
