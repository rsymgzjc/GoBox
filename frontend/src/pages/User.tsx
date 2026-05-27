import { useEffect, useState } from "react";
import { savePreferences, sendRegisterCode } from "../api/client";
import { useSessionStore } from "../store/useSessionStore";

export function UserPage() {
  const { user, mode, setMode, signIn, signUp, hydrate, loading } = useSessionStore();
  const [form, setForm] = useState({ name: "", email: "", password: "", code: "" });
  const [theme, setTheme] = useState("aurora");
  const [message, setMessage] = useState("");
  const [countdown, setCountdown] = useState(0);

  useEffect(() => {
    hydrate();
  }, [hydrate]);

  useEffect(() => {
    if (countdown <= 0) {
      return;
    }
    const timer = window.setTimeout(() => setCountdown((current) => current - 1), 1000);
    return () => window.clearTimeout(timer);
  }, [countdown]);

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setMessage("");
    try {
      if (mode === "login") {
        await signIn(form.email, form.password);
        setMessage("登录成功");
      } else {
        await signUp(form.name, form.email, form.password, form.code);
        setMessage("注册成功");
      }
    } catch (err) {
      setMessage(err instanceof Error ? err.message : "认证失败");
    }
  }

  async function handleSendCode() {
    setMessage("");
    try {
      const data = await sendRegisterCode(form.email);
      setCountdown(data.cooldownSeconds);
      if (data.previewCode) {
        setForm((current) => ({ ...current, code: data.previewCode || current.code }));
        setMessage(`开发模式验证码：${data.previewCode}`);
        return;
      }
      setMessage(`验证码已发送，有效期 ${data.expiresInMinutes} 分钟`);
    } catch (err) {
      setMessage(err instanceof Error ? err.message : "发送验证码失败");
    }
  }

  async function handleSavePreference() {
    try {
      await savePreferences({ dashboard_theme: theme });
      setMessage("偏好已保存");
    } catch (err) {
      setMessage(err instanceof Error ? err.message : "保存失败");
    }
  }

  if (user) {
    return (
      <section className="panel user-panel">
        <div className="panel-header">
          <div>
            <h2>欢迎回来，{user.name}</h2>
            <p>{user.email}</p>
          </div>
          <span className="badge badge-strong">{user.role}</span>
        </div>
        <div className="settings-card">
          <label className="field">
            仪表盘主题
            <select value={theme} onChange={(event) => setTheme(event.target.value)}>
              <option value="aurora">Aurora</option>
              <option value="ember">Ember</option>
              <option value="graphite">Graphite</option>
            </select>
          </label>
          <button className="primary-button" onClick={handleSavePreference}>
            保存偏好
          </button>
        </div>
        {message ? <p className="success-text">{message}</p> : null}
      </section>
    );
  }

  return (
    <section className="auth-layout">
      <div className="panel auth-panel">
        <div className="panel-header">
          <div>
            <h2>{mode === "login" ? "登录 GoBox" : "创建账户"}</h2>
            <p>注册时需要先发送邮箱验证码，再输入验证码完成注册。</p>
          </div>
          <button className="ghost-button" onClick={() => setMode(mode === "login" ? "register" : "login")}>
            切换到{mode === "login" ? "注册" : "登录"}
          </button>
        </div>
        <form className="auth-form" onSubmit={handleSubmit}>
          {mode === "register" ? (
            <label className="field">
              昵称
              <input value={form.name} onChange={(event) => setForm({ ...form, name: event.target.value })} />
            </label>
          ) : null}
          <label className="field">
            邮箱
            <input value={form.email} onChange={(event) => setForm({ ...form, email: event.target.value })} />
          </label>
          <label className="field">
            密码
            <input
              type="password"
              value={form.password}
              onChange={(event) => setForm({ ...form, password: event.target.value })}
            />
          </label>
          {mode === "register" ? (
            <>
              <label className="field">
                邮箱验证码
                <input value={form.code} onChange={(event) => setForm({ ...form, code: event.target.value })} />
              </label>
              <button
                className="ghost-button"
                disabled={!form.email || countdown > 0}
                onClick={handleSendCode}
                type="button"
              >
                {countdown > 0 ? `${countdown}s 后重发` : "发送验证码"}
              </button>
            </>
          ) : null}
          <button className="primary-button" disabled={loading} type="submit">
            {loading ? "提交中..." : mode === "login" ? "登录" : "注册"}
          </button>
        </form>
        {message ? <p className={message.includes("成功") || message.includes("验证码") ? "success-text" : "error-text"}>{message}</p> : null}
      </div>
    </section>
  );
}
