import type { PointDraft } from '../types';

type Props = {
  draft: PointDraft; categoryName: string; busy: boolean; error: string;
  onConfirm: () => void; onEdit: () => void;
};

export default function ConfirmPointScreen({ draft, categoryName, busy, error, onConfirm, onEdit }: Props) {
  return (
    <section className="screen" aria-labelledby="confirm-title">
      <h1 id="confirm-title">Подтверждение точки</h1>
      <p className="muted">Проверьте данные перед добавлением.</p>
      <dl className="card data-list">
        <div><dt>Название</dt><dd>{draft.name.trim()}</dd></div>
        <div><dt>Категория</dt><dd>{categoryName}</dd></div>
        <div><dt>Адрес</dt><dd>{draft.selectedAddress?.address}</dd></div>
      </dl>
      {error && <p className="error-message" role="alert">{error}</p>}
      {busy && <p role="status">Сохраняем точку и рассчитываем Smart Score…</p>}
      <div className="actions">
        <button className="button button-primary" type="button" disabled={busy} onClick={onConfirm}>Подтвердить</button>
        <button className="button" type="button" disabled={busy} onClick={onEdit}>Изменить</button>
      </div>
    </section>
  );
}
