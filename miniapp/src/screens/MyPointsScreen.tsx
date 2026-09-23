import { scoreChange } from '../mockData';
import type { Point } from '../mockData';

type Props = { points: Point[]; onSelect: (id: string) => void; onAdd: () => void };

export default function MyPointsScreen({ points, onSelect, onAdd }: Props) {
  return (
    <section className="screen" aria-labelledby="points-title">
      <h1 id="points-title">Мои точки</h1>
      <p className="muted small">Демо-данные. Изменения сохраняются до перезагрузки страницы.</p>
      {points.length === 0 && <p className="card">У вас пока нет точек. Добавьте первую.</p>}
      <ul className="point-list">
        {points.map((point) => (
          <li key={point.id}>
            <button className="card point-card" type="button" onClick={() => onSelect(point.id)}>
              <span className="point-name">{point.name}</span>
              <span className="muted">{point.address}</span>
              <span className="score-row">
                <span>Smart Score <strong>{point.history[2].score.toFixed(1)}</strong><span className="muted"> / 10</span></span>
                <span className="change">{scoreChange(point)} за месяц</span>
              </span>
            </button>
          </li>
        ))}
      </ul>
      <button className="button button-primary" type="button" onClick={onAdd}>📍 Добавить точку</button>
    </section>
  );
}
