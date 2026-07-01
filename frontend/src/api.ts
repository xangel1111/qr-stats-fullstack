import type { ComputationResponse, ComputationSummary } from './types';

const BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080';

/** ApiError carries the API's error code so the UI can react to it. */
export class ApiError extends Error {
  constructor(
    public readonly status: number,
    public readonly code: string,
    message: string,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

async function throwApiError(res: Response): Promise<never> {
  let code = 'ERROR';
  let message = res.statusText || `HTTP ${res.status}`;
  try {
    const body = await res.json();
    if (body?.error) {
      code = body.error.code ?? code;
      message = body.error.message ?? message;
    }
  } catch {
    // response had no JSON body; keep defaults
  }
  throw new ApiError(res.status, code, message);
}

export async function login(username: string, password: string): Promise<string> {
  const res = await fetch(`${BASE_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });
  if (!res.ok) await throwApiError(res);
  const data = await res.json();
  return data.token as string;
}

export async function computeQR(token: string, matrix: number[][]): Promise<ComputationResponse> {
  const res = await fetch(`${BASE_URL}/api/v1/qr`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ matrix }),
  });
  if (!res.ok) await throwApiError(res);
  return (await res.json()) as ComputationResponse;
}

export async function listComputations(token: string): Promise<ComputationSummary[]> {
  const res = await fetch(`${BASE_URL}/api/v1/computations?limit=10`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) await throwApiError(res);
  const data = await res.json();
  return (data.items ?? []) as ComputationSummary[];
}
