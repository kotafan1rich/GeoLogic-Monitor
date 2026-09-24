import type { RatingHistoryEntry } from './types';

export function ratingChange(history: RatingHistoryEntry[]): string | null {
  if (history.length < 2) return null;
  const delta = Math.round((history[history.length - 1].value - history[history.length - 2].value) * 10) / 10;
  return `${delta > 0 ? '+' : ''}${delta.toFixed(1)}`;
}

export function formatCalculationDate(value: string): string {
  return new Date(value).toLocaleString('ru-RU', {
    timeZone: 'Europe/Moscow', day: '2-digit', month: '2-digit', year: 'numeric',
    hour: '2-digit', minute: '2-digit', second: '2-digit',
  });
}
