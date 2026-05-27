import { useEffect, useState } from "react";
import { fetchSummary } from "../api/client";
import { StatPill } from "../components/StatPill";
import type { Summary } from "../types";

export function AdminPage() {
  const [summary, setSummary] = useState<Summary | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    fetchSummary()
      .then(setSummary)
      .catch((err) => setError(err.message));
  }, []);

  return (
    <section className="panel">
      <div className="panel-header">
        <div>
          <h2>平台总览</h2>
          <p>适合作为后续接入真实监控与 BI 面板前的管理后台起点。</p>
        </div>
      </div>
      {summary ? (
        <>
          <div className="stats-row">
            <StatPill label="工具总数" value={summary.toolCount} />
            <StatPill label="运行次数" value={summary.usageCount} />
            <StatPill label="用户数" value={summary.userCount} />
          </div>
          <div className="leaderboard">
            {summary.topTools.map((item) => (
              <div className="leaderboard-item" key={item.toolSlug}>
                <strong>{item.toolSlug}</strong>
                <span>{item.count} 次</span>
              </div>
            ))}
          </div>
        </>
      ) : null}
      {error ? <p className="error-text">{error}</p> : null}
    </section>
  );
}
