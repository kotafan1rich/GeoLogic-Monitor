import type { PointDraft } from '../mockData';

type Props = { draft: PointDraft; onConfirm: () => void; onEdit: () => void };

export default function ConfirmPointScreen({ draft, onConfirm, onEdit }: Props) {
  return (
    <section className="screen" aria-labelledby="confirm-title">
      <h1 id="confirm-title">Подтверждение точки</h1>
      <p className="muted">Проверьте данные перед добавлением.</p>
      <dl className="card data-list">
        <div><dt>Название</dt><dd>{draft.name}</dd></div>
        <div><dt>Категория</dt><dd>{draft.category}</dd></div>
        <div><dt>Адрес</dt><dd>{draft.address}</dd></div>
      </dl>
      <div className="actions">
        <button className="button button-primary" type="button" onClick={onConfirm}>Подтвердить</button>
        <button className="button" type="button" onClick={onEdit}>Изменить</button>
      </div>
    </section>
  );
}
