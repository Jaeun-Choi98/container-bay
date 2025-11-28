import { useState, useCallback } from 'react';

export interface UseApiResult<T> {
  result: T | null;
  loading: boolean;
  error: string | null;
  execute: (...args: any[]) => Promise<void>;
  reset: () => void;
}

export function useApi<T>(apiFunc: (...args: any[]) => Promise<T>)
  : UseApiResult<T> {
  const [result, setResult] = useState<T | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const execute = useCallback(async (...args: any[]): Promise<void> => {
    setLoading(true);
    setError(null);
    try {
      const result = await apiFunc(...args);
      setResult(result);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '[API] unknown error type';
      setError(errorMessage);
      console.error('[API] error:', errorMessage);
    } finally {
      setLoading(false);
    }
  }, [apiFunc]);
  const reset = useCallback((): void => {
    setResult(null);
    setError(null);
    setLoading(false);
  }, [])

  return { result, loading, error, execute, reset }
}