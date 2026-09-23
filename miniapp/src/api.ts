import { getInitData } from './max';
import type {
  Address, ApiErrorCode, BusinessType, CreatedTrackedLocation,
  CreatePointRequest, RatingHistoryEntry, TrackedLocation,
} from './types';

export class ApiError extends Error {
  readonly status: number;
  readonly code?: ApiErrorCode;

  constructor(status: number, code?: ApiErrorCode) {
    super('GeoLogic API request failed');
    this.status = status;
    this.code = code;
  }
}

export function errorMessage(error: unknown): string {
  if (error instanceof ApiError) {
    if (error.status === 401) return 'Заново откройте Mini App через MAX.';
    if (error.status === 400) return 'Проверьте введённые данные.';
    if (error.status === 404) return 'Запрошенные данные не найдены.';
    if (error.status === 503) return 'Сервис временно недоступен. Попробуйте позже.';
    if (error.status === 0) return 'Не удалось связаться с API. Проверьте соединение.';
    return 'Не удалось получить ответ API. Попробуйте позже.';
  }
  return 'Не удалось выполнить запрос. Попробуйте позже.';
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const initData = getInitData();
  if (!initData) throw new ApiError(401, 'unauthorized');

  const controller = new AbortController();
  const abort = () => controller.abort();
  const signal = options.signal;
  if (signal?.aborted) controller.abort();
  signal?.addEventListener('abort', abort, { once: true });
  const timeout = window.setTimeout(abort, 30000);

  try {
    const headers = new Headers(options.headers);
    headers.set('X-Max-Init-Data', initData);
    if (options.body) headers.set('Content-Type', 'application/json');
    const response = await fetch('/api/v1' + path, {
      ...options, headers, signal: controller.signal, credentials: 'omit', cache: 'no-store',
    });
    if (!response.ok) {
      const body = await response.json().catch(() => null) as { code?: ApiErrorCode } | null;
      throw new ApiError(response.status, body?.code);
    }
    if (response.status === 204) return undefined as T;
    return await response.json() as T;
  } catch (error) {
    if (signal?.aborted) throw error;
    if (error instanceof ApiError) throw error;
    throw new ApiError(0);
  } finally {
    window.clearTimeout(timeout);
    signal?.removeEventListener('abort', abort);
  }
}

export const api = {
  businessTypes: (signal?: AbortSignal) =>
    request<BusinessType[]>('/business-types', { signal }),
  points: (signal?: AbortSignal) =>
    request<TrackedLocation[]>('/tracked-locations', { signal }),
  createPoint: (point: CreatePointRequest) =>
    request<CreatedTrackedLocation>('/tracked-locations', {
      method: 'POST',
      body: JSON.stringify({
        name: point.name, business_type_id: point.business_type_id,
        address: point.address, lat: point.lat, lon: point.lon,
      }),
    }),
  deletePoint: (id: string) =>
    request<void>('/tracked-locations/' + encodeURIComponent(id), { method: 'DELETE' }),
  history: (id: string, signal?: AbortSignal) =>
    request<RatingHistoryEntry[]>('/tracked-locations/' + encodeURIComponent(id) + '/rating-history?months=3', { signal }),
  suggestions: (query: string, signal?: AbortSignal) =>
    request<Address[]>('/geocoding/suggestions?' + new URLSearchParams({ query, limit: '5' }), { signal }),
  address: (lat: number, lon: number, signal?: AbortSignal) =>
    request<Address>('/geocoding/address?' + new URLSearchParams({ lat: String(lat), lon: String(lon) }), { signal }),
};
