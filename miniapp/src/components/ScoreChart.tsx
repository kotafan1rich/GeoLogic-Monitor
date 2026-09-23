import type { ScoreEntry } from '../mockData';

export default function ScoreChart({ history }: { history: ScoreEntry[] }) {
  const x = (index: number) => 40 + index * 120;
  const y = (score: number) => 150 - score * 12;

  return (
    <svg className="score-chart" viewBox="0 0 320 180" role="img" aria-label={`Smart Score от 0 до 10: ${history.map((entry) => `${entry.month} — ${entry.score.toFixed(1)}`).join(', ')}`}>
      {[0, 5, 10].map((score) => (
        <g key={score}>
          <line x1="32" y1={y(score)} x2="295" y2={y(score)} className="chart-grid" />
          <text x="24" y={y(score) + 4} textAnchor="end" className="chart-label">{score}</text>
        </g>
      ))}
      <polyline points={history.map((entry, index) => `${x(index)},${y(entry.score)}`).join(' ')} className="chart-line" />
      {history.map((entry, index) => (
        <g key={entry.month}>
          <circle cx={x(index)} cy={y(entry.score)} r="5" className="chart-dot" />
          <text x={x(index)} y={y(entry.score) - 12} textAnchor="middle" className="chart-value">{entry.score.toFixed(1)}</text>
          <text x={x(index)} y="174" textAnchor="middle" className="chart-label">{entry.month.split(' ')[0].slice(0, 3)}</text>
        </g>
      ))}
    </svg>
  );
}
