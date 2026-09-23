import { ratingChange } from '../rating';
import { useRatingHistory } from '../useRatingHistory';
import type { LoadState, TrackedLocation } from '../types';

type Props = {
  points: LoadState<TrackedLocation[]>; notice: string;
  onSelect: (point: TrackedLocation) => void; onAdd: () => void;
  onRetry: () => void; onUnauthorized: () => void;
};

function PointCard({ point, onSelect, onUnauthorized }: {
  point: TrackedLocation; onSelect: Props['onSelect']; onUnauthorized: () => void;
}) {
  const { state, retry } = useRatingHistory(point.id, onUnauthorized);
  const change = state.status === 'success' ? ratingChange(state.data) : null;
  return (
    <li>
      <button className="card point-card" type="button" onClick={() => onSelect(point)}>
        <span className="point-name">{point.name}</span>
        <span className="muted">{point.address}</span>
        <span className="score-row">
          <span>Smart Score {point.rating === null ? <span className="muted">пока недоступен</span>
            : <><strong>{point.rating.toFixed(1)}</strong><span className="muted"> / 10</span></>}</span>
          {change !== null && <span className="change">{change} с предыдущего расчёта</span>}
          {state.status === 'loading' && <span className="muted small" role="status">Загрузка истории…</span>}
        </span>
      </button>
      {state.status === 'error' && <div className="small error-message" role="alert">
        <p>История недоступна. {state.message}</p>
        <button className="back-button" type="button" onClick={retry}>Повторить загрузку истории</button>
      </div>}
    </li>
  );
}

export default function MyPointsScreen({ points, notice, onSelect, onAdd, onRetry, onUnauthorized }: Props) {
  return (
    <section className="screen" aria-labelledby="points-title">
      <h1 id="points-title">Мои точки</h1>
      {notice && <p role="status">{notice}</p>}
      {points.status === 'loading' && <p role="status">Загрузка точек…</p>}
      {points.status === 'error' && <div className="card" role="alert"><p>{points.message}</p><button className="button" type="button" onClick={onRetry}>Повторить загрузку</button></div>}
      {points.status === 'success' && (points.data.length === 0
        ? <p className="card">У вас пока нет точек. Добавьте первую.</p>
        : <ul className="point-list">{points.data.map((point) =>
          <PointCard key={point.id} point={point} onSelect={onSelect} onUnauthorized={onUnauthorized} />)}</ul>)}
      <button className="button button-primary" type="button" onClick={onAdd}>📍 Добавить точку</button>
    </section>
  );
}
