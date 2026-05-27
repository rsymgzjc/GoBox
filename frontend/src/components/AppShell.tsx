import type { ReactNode } from "react";
import { Link, NavLink } from "react-router-dom";
import { useSessionStore } from "../store/useSessionStore";

const navItems = [
  { to: "/", label: "首页" },
  { to: "/tools", label: "工具台" },
  { to: "/user", label: "用户中心" },
  { to: "/admin", label: "管理视图" }
];

export function AppShell({ children }: { children: ReactNode }) {
  const user = useSessionStore((state) => state.user);
  const signOut = useSessionStore((state) => state.signOut);

  return (
    <div className="app-shell">
      <header className="topbar">
        <Link className="brand" to="/">
          <span className="brand-mark">GB</span>
          <div>
            <strong>GoBox</strong>
            <p>在线工具箱平台</p>
          </div>
        </Link>
        <nav className="nav">
          {navItems.map((item) => (
            <NavLink key={item.to} className="nav-link" to={item.to}>
              {item.label}
            </NavLink>
          ))}
        </nav>
        <div className="topbar-user">
          {user ? (
            <>
              <span>{user.name}</span>
              <button className="ghost-button" onClick={signOut}>
                退出
              </button>
            </>
          ) : (
            <span>未登录</span>
          )}
        </div>
      </header>
      <main className="page">{children}</main>
    </div>
  );
}
