const API_BASE_URL = 'http://localhost:8080/api/todos';

interface ApiResponse<T> {
  data: T;
}

interface ApiError {
  error: string;
}

export class ApiException extends Error {
  constructor(public statusCode: number, message: string) {
    super(message);
    this.name = 'ApiException';
  }
}

async function request<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;

  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    const json = await response.json();

    if (!response.ok) {
      const error = json as ApiError;
      throw new ApiException(response.status, error.error || 'Unknown error');
    }

    const result = json as ApiResponse<T>;
    return result.data;
  } catch (error) {
    if (error instanceof ApiException) throw error;
    throw new ApiException(0, error instanceof Error ? error.message : 'Network error');
  }
}

export const api = {
  get: <T>(endpoint: string) => request<T>(endpoint, { method: 'GET' }),
  post: <T>(endpoint: string, data?: unknown) => request<T>(endpoint, {
    method: 'POST',
    body: data ? JSON.stringify(data) : undefined,
  }),
  patch: <T>(endpoint: string) => request<T>(endpoint, { method: 'PATCH' }),
  delete: <T>(endpoint: string) => request<T>(endpoint, { method: 'DELETE' }),
};
