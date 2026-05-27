import { useEffect, useMemo, useState } from "react";
import { fetchTools, runTool } from "../api/client";
import { ToolCard } from "../components/ToolCard";
import { StatPill } from "../components/StatPill";
import type { Tool, ToolResult } from "../types";

export function ToolsPage() {
  const [tools, setTools] = useState<Tool[]>([]);
  const [activeTool, setActiveTool] = useState<Tool | null>(null);
  const [input, setInput] = useState("");
  const [output, setOutput] = useState("");
  const [meta, setMeta] = useState<Record<string, unknown> | undefined>();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    fetchTools()
      .then((data) => {
        setTools(data);
        setActiveTool(data[0] ?? null);
      })
      .catch((err) => setError(err.message));
  }, []);

  const categories = useMemo(() => Array.from(new Set(tools.map((tool) => tool.category))), [tools]);

  async function handleRun() {
    if (!activeTool) {
      return;
    }
    setLoading(true);
    setError("");
    try {
      const options = activeTool.slug === "uuid-generate" ? { count: "3" } : undefined;
      const result: ToolResult = await runTool(activeTool.slug, input, options);
      setOutput(result.output);
      setMeta(result.meta);
    } catch (err) {
      setError(err instanceof Error ? err.message : "执行失败");
      setOutput("");
      setMeta(undefined);
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="tools-layout">
      <aside className="panel">
        <div className="panel-header">
          <h2>工具目录</h2>
          <p>{tools.length} 个内置工具</p>
        </div>
        <div className="category-row">
          {categories.map((category) => (
            <span key={category} className="badge">
              {category}
            </span>
          ))}
        </div>
        <div className="tool-grid">
          {tools.map((tool) => (
            <ToolCard key={tool.slug} tool={tool} active={activeTool?.slug === tool.slug} onSelect={setActiveTool} />
          ))}
        </div>
      </aside>

      <div className="panel panel-workbench">
        <div className="panel-header">
          <div>
            <h2>{activeTool?.name ?? "选择一个工具"}</h2>
            <p>{activeTool?.description ?? "从左侧选择工具后开始使用"}</p>
          </div>
          <button className="primary-button" disabled={loading || !activeTool} onClick={handleRun}>
            {loading ? "执行中..." : "运行工具"}
          </button>
        </div>
        <label className="field">
          输入内容
          <textarea
            placeholder={activeTool?.inputHint || "请输入内容"}
            value={input}
            onChange={(event) => setInput(event.target.value)}
          />
        </label>
        <label className="field">
          输出结果
          <textarea
            placeholder={activeTool?.outputHint || "结果会展示在这里"}
            value={output}
            readOnly
          />
        </label>
        <div className="stats-row">
          {meta
            ? Object.entries(meta).map(([key, value]) => <StatPill key={key} label={key} value={String(value)} />)
            : <StatPill label="状态" value={error ? "失败" : "待执行"} />}
        </div>
        {error ? <p className="error-text">{error}</p> : null}
      </div>
    </section>
  );
}
