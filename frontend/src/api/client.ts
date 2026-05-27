import axios from "axios";
import type { Summary, Tool, ToolResult, User } from "../types";

const client = axios.create({
  baseURL: "/api/v1",
  timeout: 8000
});

client.interceptors.request.use((config) => {
  const token = localStorage.getItem("gobox_token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

type Envelope<T> = {
  success: boolean;
  message: string;
  data: T;
  error?: string;
};

export async function fetchTools(): Promise<Tool[]> {
  const { data } = await client.get<Envelope<Tool[]>>("/tools");
  return data.data;
}

export async function runTool(slug: string, input: string, options?: Record<string, string>): Promise<ToolResult> {
  const { data } = await client.post<Envelope<ToolResult>>(`/tools/${slug}/run`, { input, options });
  return data.data;
}

export async function login(email: string, password: string): Promise<{ user: User; token: string }> {
  const { data } = await client.post<Envelope<{ user: User; token: string }>>("/auth/login", { email, password });
  return data.data;
}

export async function sendRegisterCode(email: string): Promise<{ cooldownSeconds: number; expiresInMinutes: number; previewCode?: string }> {
  const { data } = await client.post<Envelope<{ cooldownSeconds: number; expiresInMinutes: number; previewCode?: string }>>(
    "/auth/register/send-code",
    { email }
  );
  return data.data;
}

export async function register(
  name: string,
  email: string,
  password: string,
  code: string
): Promise<{ user: User; token: string }> {
  const { data } = await client.post<Envelope<{ user: User; token: string }>>("/auth/register", {
    name,
    email,
    password,
    code
  });
  return data.data;
}

export async function fetchProfile(): Promise<User> {
  const { data } = await client.get<Envelope<User>>("/me");
  return data.data;
}

export async function savePreferences(preferences: Record<string, string>): Promise<Record<string, string>> {
  const { data } = await client.put<Envelope<Record<string, string>>>("/me/preferences", { preferences });
  return data.data;
}

export async function fetchSummary(): Promise<Summary> {
  const { data } = await client.get<Envelope<Summary>>("/stats/summary");
  return data.data;
}
