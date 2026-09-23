import { useEffect, useState } from 'react';
import { api, ApiError, errorMessage } from './api';
import type { LoadState, RatingHistoryEntry } from './types';

export function useRatingHistory(id: string, onUnauthorized: () => void) {
  const [state, setState] = useState<LoadState<RatingHistoryEntry[]>>({ status: 'loading' });
  const [revision, setRevision] = useState(0);

  useEffect(() => {
    const controller = new AbortController();
    setState({ status: 'loading' });
    api.history(id, controller.signal).then((data) => {
      if (!controller.signal.aborted) setState({ status: 'success', data });
    }).catch((error: unknown) => {
      if (controller.signal.aborted) return;
      if (error instanceof ApiError && error.status === 401) onUnauthorized();
      setState({ status: 'error', message: errorMessage(error) });
    });
    return () => controller.abort();
  }, [id, revision, onUnauthorized]);

  return { state, retry: () => setRevision((value) => value + 1) };
}
