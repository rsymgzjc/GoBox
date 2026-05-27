import { Routes, Route } from "react-router-dom";
import { AppShell } from "./components/AppShell";
import { AdminPage } from "./pages/Admin";
import { HomePage } from "./pages/Home";
import { ToolsPage } from "./pages/Tools";
import { UserPage } from "./pages/User";

export default function App() {
  return (
    <AppShell>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/tools" element={<ToolsPage />} />
        <Route path="/user" element={<UserPage />} />
        <Route path="/admin" element={<AdminPage />} />
      </Routes>
    </AppShell>
  );
}
