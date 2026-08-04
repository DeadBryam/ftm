import ky from 'ky';
import type { HTTPError } from 'ky';

export function extractErrorMessage(body: string): string {
  const text = body.trim();
  if (!text) return '';

  try {
    const parsed = JSON.parse(text);
    if (typeof parsed === 'string') return parsed.trim();
    if (parsed && typeof parsed === 'object') {
      const value = (parsed as Record<string, unknown>).error ?? (parsed as Record<string, unknown>).message;
      if (typeof value === 'string' && value.trim()) return value.trim();
    }
    return text;
  } catch {
    return text;
  }
}

export const api = ky.create({
  prefix: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
  retry: {
    limit: 2,
    methods: ['get'],
    statusCodes: [408, 413, 429, 500, 502, 503, 504],
    backoffLimit: 10000,
  },
  hooks: {
    beforeError: [
      async ({ error }) => {
        const response = (error as HTTPError).response;
        if (!response) return error;

        try {
          const message = extractErrorMessage(await response.clone().text());
          if (message) error.message = message;
        } catch {
          return error;
        }

        return error;
      },
    ],
  },
});
