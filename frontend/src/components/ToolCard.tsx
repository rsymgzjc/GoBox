import type { Tool } from "../types";

export function ToolCard({
  tool,
  active,
  onSelect
}: {
  tool: Tool;
  active: boolean;
  onSelect: (tool: Tool) => void;
}) {
  return (
    <button className={`tool-card ${active ? "active" : ""}`} onClick={() => onSelect(tool)}>
      <div className="tool-card-top">
        <span className="badge">{tool.category}</span>
        {tool.isFeatured ? <span className="badge badge-strong">精选</span> : null}
      </div>
      <strong>{tool.name}</strong>
      <p>{tool.description}</p>
    </button>
  );
}
