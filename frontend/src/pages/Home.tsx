import { Link } from "react-router-dom";

export function HomePage() {
  return (
    <section className="hero-grid">
      <div className="hero-card hero-card-large">
        <span className="eyebrow">Go + React</span>
        <h1>把常用开发小工具放进一个反应很快的工作台。</h1>
        <p>
          GoBox 提供认证、工具运行、统计分析和响应式界面。它已经具备可继续扩展为 20+
          工具平台的完整骨架。
        </p>
        <div className="hero-actions">
          <Link className="primary-button" to="/tools">
            进入工具台
          </Link>
          <Link className="ghost-button" to="/user">
            登录体验
          </Link>
        </div>
      </div>
      <div className="hero-card">
        <h2>当前内置能力</h2>
        <ul className="feature-list">
          <li>JSON / Base64 / URL / Hash / 时间戳 / UUID / Slug</li>
          <li>JWT 登录态、用户偏好保存、管理员统计视图</li>
          <li>Docker Compose、Nginx 反向代理、GitHub Actions</li>
        </ul>
      </div>
      <div className="hero-card">
        <h2>默认管理员</h2>
        <p>邮箱：admin@gobox.local</p>
        <p>密码：admin123456</p>
      </div>
    </section>
  );
}
