import { useRef } from 'react';
import ScoreChart from '../components/ScoreChart';
import { formatCalculationDate, ratingChange } from '../rating';
import { useRatingHistory } from '../useRatingHistory';
import type { TrackedLocation } from '../types';

type Props = {
  point: TrackedLocation; categoryName: string; busy: boolean; error: string;
  onDelete: () => void; onUnauthorized: () => void;
};

export default function PointDetailsScreen({ point, categoryName, busy, error, onDelete, onUnauthorized }: Props) {
  const dialogRef = useRef<HTMLDialogElement>(null);
  const { state, retry } = useRatingHistory(point.id, onUnauthorized);
  const change = state.status === 'success' ? ratingChange(state.data) : null;

  return (
    <section className="screen" aria-labelledby="details-title">
      <h1 id="details-title">{point.name}</h1>
      <p className="muted">{point.address}</p>
      <p className="category-tag">{categoryName}</p>
      <div className="card score-summary">
        <h2>Текущий Smart Score</h2>
        {point.rating === null ? <p className="muted">Рейтинг пока недоступен</p>
          : <p className="score-number">{point.rating.toFixed(1)}<span> / 10</span></p>}
        {point.rating_calculated_at && <p className="muted small">Последнее обновление: <time dateTime={point.rating_calculated_at}>{formatCalculationDate(point.rating_calculated_at)}</time> (МСК)</p>}
        {change !== null && <p className="small">{change} с предыдущего расчёта</p>}
      </div>
      <section className="card" aria-labelledby="history-title">
        <h2 id="history-title">История за 3 месяца</h2>
        <p className="muted small">История Smart Score доступна за последние 3 месяца.</p>
        {state.status === 'loading' && <p role="status">Загрузка истории…</p>}
        {state.status === 'error' && <div role="alert"><p>{state.message}</p><button className="button" type="button" onClick={retry}>Повторить загрузку</button></div>}
        {state.status === 'success' && <>
          <ScoreChart history={state.data} />
          {state.data.length > 0 && <dl className="history-list">
            {state.data.map((entry, index) => <div key={index}>
              <dt><time dateTime={entry.calculated_at}>{formatCalculationDate(entry.calculated_at)}</time> (МСК)</dt>
              <dd>{entry.value.toFixed(1)}</dd>
            </div>)}
          </dl>}
          <p className="muted small">История будет пополняться после следующих пересчётов.</p>
        </>}
      </section>
      <section className="card" aria-labelledby="factors-title">
        <h2 id="factors-title">Что учитывает Smart Score</h2>
        <dl className="data-list">
          <div><dt>Конкуренция</dt><dd>Количество и близость похожих заведений рядом</dd></div>
          <div><dt>Инфраструктура</dt><dd>Объекты и сервисы вокруг точки</dd></div>
          <div><dt>Транспортная доступность</dt><dd>Насколько удобно добраться до локации</dd></div>
        </dl>
      </section>
      <button className="button button-danger" type="button" disabled={busy} onClick={() => dialogRef.current?.showModal()}>Удалить точку</button>
      <dialog ref={dialogRef} className="delete-dialog" aria-labelledby="delete-title" aria-describedby="delete-description"
        onCancel={(event) => { if (busy) event.preventDefault(); }}>
        <h2 id="delete-title">Удалить точку?</h2>
        <p id="delete-description">{point.name}</p>
        {error && <p className="error-message" role="alert">{error}</p>}
        {busy && <p role="status">Удаление точки…</p>}
        <div className="actions">
          <button className="button button-danger" type="button" disabled={busy} onClick={onDelete}>Удалить</button>
          <button className="button" type="button" autoFocus disabled={busy} onClick={() => dialogRef.current?.close()}>Отмена</button>
        </div>
      </dialog>
    </section>
  );
}
