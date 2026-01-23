import { z } from 'zod';
import { ApiErrorSchema, createApiResponseSchema } from './types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/todos';

export class ApiException extends Error {
  constructor(
    public statusCode: number,
    message: string,
    public validationErrors?: z.ZodError
  ) {
    super(message);
    this.name = 'ApiException';
  }
}

export class ValidationException extends ApiException {
  constructor(zodError: z.ZodError) {
    const message = formatZodError(zodError);
    super(0, message, zodError);
    this.name = 'ValidationException';
  }
}

function formatZodError(error: z.ZodError): string {
  return error.issues
    .map((e) => `${e.path.join('.')}: ${e.message}`)
    .join('; ');
}

async function request<T>(
  endpoint: string,
  schema: z.ZodType<T>,
  options?: RequestInit
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;

  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    let json: unknown;
    try {
      json = await response.json();
    } catch {
      throw new ApiException(
        response.status,
        response.ok ? 'Invalid JSON response' : `HTTP ${response.status}`
      );
    }

    if (!response.ok) {
      const errorResult = ApiErrorSchema.safeParse(json);
      const errorMessage = errorResult.success
        ? errorResult.data.error
        : 'Unknown error';
      throw new ApiException(response.status, errorMessage);
    }

    const responseSchema = createApiResponseSchema(schema);
    const result = responseSchema.safeParse(json);

    if (!result.success) {
      console.error('API Response validation failed:', result.error);
      throw new ValidationException(result.error);
    }

    return result.data.data;
  } catch (error) {
    if (error instanceof ApiException) throw error;
    const wrappedError = new ApiException(0, error instanceof Error ? error.message : 'Network error');
    if (error instanceof Error) {
      wrappedError.cause = error;
    }
    throw wrappedError;
  }
}

export const api = {
  get: <T>(endpoint: string, schema: z.ZodType<T>) =>
    request(endpoint, schema, { method: 'GET' }),

  post: <T>(endpoint: string, schema: z.ZodType<T>, data?: unknown) =>
    request(endpoint, schema, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    }),

  patch: <T>(endpoint: string, schema: z.ZodType<T>) =>
    request(endpoint, schema, { method: 'PATCH' }),

  delete: <T>(endpoint: string, schema: z.ZodType<T>) =>
    request(endpoint, schema, { method: 'DELETE' }),
};
