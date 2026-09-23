import { useRef } from 'react';
import ScoreChart from '../components/ScoreChart';
import type { Point } from '../mockData';

export default function PointDetailsScreen({ point, onDelete }: { point: Point; onDelete: () => void }) {
  const dialogRef = useRef<HTMLDialogElement>(null);

  return (
    <section className="screen" aria-labelledby="details-title">
      <h1 id="details-title">{point.name}</h1>
      <p className="muted">{point.address}</p>
      <p className="category-tag">{point.category}</p>
      <p className="muted small">Демонстрационные оценки и факторы, включая новые точки.</p>
      <div className="card score-summary">
        <h2>Текущий Smart Score</h2>
        <p className="score-number">{point.history[2].score.toFixed(1)}<span> / 10</span></p>
        <p className="muted small">Последнее обновление: {point.updatedAt}</p>
      </div>
      <section className="card" aria-labelledby="history-title">
        <h2 id="history-title">История за 3 месяца</h2>
        <ScoreChart history={point.history} />
        <dl className="history-list">
          {point.history.map((entry) => <div key={entry.month}><dt>{entry.month}</dt><dd>{entry.score.toFixed(1)}</dd></div>)}
        </dl>
      </section>
      <section className="card" aria-labelledby="factors-title">
        <h2 id="factors-title">Что учитывает Smart Score</h2>
        <dl className="data-list">
          {point.factors.map((factor) => <div key={factor.name}><dt>{factor.name}</dt><dd>{factor.description}</dd></div>)}
        </dl>
      </section>
      <button className="button button-danger" type="button" onClick={() => dialogRef.current?.showModal()}>Удалить точку</button>
      <dialog ref={dialogRef} className="delete-dialog" aria-labelledby="delete-title" aria-describedby="delete-description">
        <h2 id="delete-title">Удалить точку?</h2>
        <p id="delete-description">{point.name}</p>
        <div className="actions">
          <button className="button button-danger" type="button" onClick={onDelete}>Удалить</button>
          <button className="button" type="button" autoFocus onClick={() => dialogRef.current?.close()}>Отмена</button>
        </div>
      </dialog>
    </section>
  );
}
