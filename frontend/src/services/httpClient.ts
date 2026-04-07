import axios, { type AxiosInstance } from 'axios';
import type { ApiResponse } from '@/types';
import { toast } from 'sonner';

// ─── Abstract HttpClient Interface ───────────────────────────────────────────

export interface HttpClient {
  get<T = unknown>(url: string, params?: Record<string, unknown>): Promise<T>;
  post<T = unknown>(url: string, data?: unknown): Promise<T>;
  put<T = unknown>(url: string, data?: unknown): Promise<T>;
  delete<T = unknown>(url: string): Promise<T>;
}

// ─── Auth token helpers (lazily imported to avoid circular deps) ─────────────

function getToken(): string | null {
  try {
    const raw = localStorage.getItem('auth-storage');
    if (!raw) return null;
    const parsed = JSON.parse(raw);
    return parsed?.state?.token ?? null;
  } catch {
    return null;
  }
}

function clearAuth(): void {
  try {
    localStorage.removeItem('auth-storage');
  } catch {
    // ignore
  }
}

// ─── Error class ─────────────────────────────────────────────────────────────

export class ApiError extends Error {
  constructor(
    public code: number,
    message: string,
    public data?: unknown,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

// ─── Axios Implementation ────────────────────────────────────────────────────

export class AxiosHttpClient implements HttpClient {
  private instance: AxiosInstance;

  constructor(baseURL = '/api') {
    this.instance = axios.create({ baseURL, timeout: 15000 });

    // Request interceptor — attach Authorization header
    this.instance.interceptors.request.use((config) => {
      const token = getToken();
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    // Response interceptor — unwrap envelope { code, message, data }
    this.instance.interceptors.response.use(
      (response) => {
        const body = response.data as ApiResponse<unknown>;
        if (body.code !== 200) {
          if (body.code === 401) {
            clearAuth();
            window.location.href = '/login';
          }
          toast.error(body.message || '请求失败');
          throw new ApiError(body.code, body.message, body.data);
        }
        return body.data as never;
      },
      (error) => {
        if (axios.isAxiosError(error)) {
          const status = error.response?.status;
          if (status === 401) {
            clearAuth();
            window.location.href = '/login';
          }
          const body = error.response?.data as ApiResponse<unknown> | undefined;
          const message = body?.message ?? error.message;
          toast.error(message || '网络请求失败');
          throw new ApiError(
            status ?? 500,
            message,
            body?.data,
          );
        }
        toast.error('网络连接异常');
        throw error;
      },
    );
  }

  async get<T = unknown>(url: string, params?: Record<string, unknown>): Promise<T> {
    return this.instance.get(url, { params }) as Promise<T>;
  }

  async post<T = unknown>(url: string, data?: unknown): Promise<T> {
    return this.instance.post(url, data) as Promise<T>;
  }

  async put<T = unknown>(url: string, data?: unknown): Promise<T> {
    return this.instance.put(url, data) as Promise<T>;
  }

  async delete<T = unknown>(url: string): Promise<T> {
    return this.instance.delete(url) as Promise<T>;
  }
}

// ─── Native Fetch Implementation ─────────────────────────────────────────────

export class FetchHttpClient implements HttpClient {
  constructor(private baseURL = '/api') {}

  private async request<T>(
    method: string,
    url: string,
    data?: unknown,
    params?: Record<string, unknown>,
  ): Promise<T> {
    let fullUrl = `${this.baseURL}${url}`;

    // Append query params
    if (params) {
      const searchParams = new URLSearchParams();
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          searchParams.append(key, String(value));
        }
      });
      const qs = searchParams.toString();
      if (qs) fullUrl += `?${qs}`;
    }

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    const token = getToken();
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const init: RequestInit = { method, headers };
    if (data !== undefined && method !== 'GET') {
      init.body = JSON.stringify(data);
    }

    const response = await fetch(fullUrl, init);
    const body = (await response.json()) as ApiResponse<T>;

    if (body.code === 401) {
      clearAuth();
      window.location.href = '/login';
      throw new ApiError(401, body.message);
    }

    if (body.code !== 200) {
      throw new ApiError(body.code, body.message, body.data);
    }

    return body.data;
  }

  async get<T = unknown>(url: string, params?: Record<string, unknown>): Promise<T> {
    return this.request<T>('GET', url, undefined, params);
  }

  async post<T = unknown>(url: string, data?: unknown): Promise<T> {
    return this.request<T>('POST', url, data);
  }

  async put<T = unknown>(url: string, data?: unknown): Promise<T> {
    return this.request<T>('PUT', url, data);
  }

  async delete<T = unknown>(url: string): Promise<T> {
    return this.request<T>('DELETE', url);
  }
}

// ─── Factory ─────────────────────────────────────────────────────────────────

export type HttpClientType = 'axios' | 'fetch';

export function createHttpClient(type: HttpClientType = 'axios'): HttpClient {
  switch (type) {
    case 'fetch':
      return new FetchHttpClient();
    case 'axios':
    default:
      return new AxiosHttpClient();
  }
}
