import { formatCalculationDate } from '../rating';
import type { RatingHistoryEntry } from '../types';

export default function ScoreChart({ history }: { history: RatingHistoryEntry[] }) {
  if (history.length === 0) return <p className="muted">За последние 3 месяца расчётов нет.</p>;
  const first = Date.parse(history[0].calculated_at);
  const last = Date.parse(history[history.length - 1].calculated_at);
  const x = (entry: RatingHistoryEntry) => last === first ? 160
    : 40 + (Date.parse(entry.calculated_at) - first) / (last - first) * 240;
  const y = (score: number) => 150 - score * 12;
  const dateLabel = (value: string) => new Date(value).toLocaleDateString('ru-RU', {
    timeZone: 'Europe/Moscow', day: '2-digit', month: '2-digit',
  });

  return (
    <svg className="score-chart" viewBox="0 0 320 180" role="img"
      aria-label={`Smart Score от 0 до 10: ${history.map((entry) => `${formatCalculationDate(entry.calculated_at)} МСК — ${entry.value.toFixed(1)}`).join(', ')}`}>
      {[0, 5, 10].map((score) => <g key={score}>
        <line x1="32" y1={y(score)} x2="295" y2={y(score)} className="chart-grid" />
        <text x="24" y={y(score) + 4} textAnchor="end" className="chart-label">{score}</text>
      </g>)}
      {history.length > 1 && <polyline points={history.map((entry) => `${x(entry)},${y(entry.value)}`).join(' ')} className="chart-line" />}
      {history.map((entry, index) => <g key={index}>
        <circle cx={x(entry)} cy={y(entry.value)} r="5" className="chart-dot">
          <title>{formatCalculationDate(entry.calculated_at)} МСК: {entry.value.toFixed(1)}</title>
        </circle>
        {history.length <= 3 && <text x={x(entry)} y={y(entry.value) - 12} textAnchor="middle" className="chart-value">{entry.value.toFixed(1)}</text>}
      </g>)}
      <text x={x(history[0])} y="174" textAnchor="middle" className="chart-label">{dateLabel(history[0].calculated_at)}</text>
      {last !== first && <text x="280" y="174" textAnchor="middle" className="chart-label">{dateLabel(history[history.length - 1].calculated_at)}</text>}
    </svg>
  );
}
