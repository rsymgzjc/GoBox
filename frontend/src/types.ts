export interface Tool {
  id: number;
  name: string;
  slug: string;
  category: string;
  description: string;
  inputHint: string;
  outputHint: string;
  isFeatured: boolean;
}

export interface ToolResult {
  toolSlug: string;
  input: string;
  output: string;
  success: boolean;
  latencyMs: number;
  meta?: Record<string, unknown>;
}

export interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  lastLoginAt?: string;
  preferences?: Array<{ key: string; value: string }>;
}

export interface Summary {
  toolCount: number;
  usageCount: number;
  userCount: number;
  topTools: Array<{ toolSlug: string; count: number }>;
}
