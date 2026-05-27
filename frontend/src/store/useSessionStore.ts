import { create } from "zustand";
import { fetchProfile, login, register } from "../api/client";
import type { User } from "../types";

type AuthMode = "login" | "register";

interface SessionState {
  user: User | null;
  mode: AuthMode;
  loading: boolean;
  setMode: (mode: AuthMode) => void;
  hydrate: () => Promise<void>;
  signIn: (email: string, password: string) => Promise<void>;
  signUp: (name: string, email: string, password: string, code: string) => Promise<void>;
  signOut: () => void;
}

export const useSessionStore = create<SessionState>((set) => ({
  user: null,
  mode: "login",
  loading: false,
  setMode: (mode) => set({ mode }),
  hydrate: async () => {
    const token = localStorage.getItem("gobox_token");
    if (!token) {
      return;
    }
    set({ loading: true });
    try {
      const user = await fetchProfile();
      set({ user });
    } catch {
      localStorage.removeItem("gobox_token");
    } finally {
      set({ loading: false });
    }
  },
  signIn: async (email, password) => {
    set({ loading: true });
    try {
      const data = await login(email, password);
      localStorage.setItem("gobox_token", data.token);
      set({ user: data.user });
    } finally {
      set({ loading: false });
    }
  },
  signUp: async (name, email, password, code) => {
    set({ loading: true });
    try {
      const data = await register(name, email, password, code);
      localStorage.setItem("gobox_token", data.token);
      set({ user: data.user });
    } finally {
      set({ loading: false });
    }
  },
  signOut: () => {
    localStorage.removeItem("gobox_token");
    set({ user: null });
  }
}));
