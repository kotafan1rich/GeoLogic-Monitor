import { categories, testAddress } from '../mockData';
import type { Category, PointDraft } from '../mockData';

type Props = {
  draft: PointDraft;
  onChange: (draft: PointDraft) => void;
  onSubmit: (draft: PointDraft) => void;
};

export default function AddPointScreen({ draft, onChange, onSubmit }: Props) {
  return (
    <section className="screen" aria-labelledby="add-title">
      <h1 id="add-title">Добавить точку</h1>
      <form className="point-form" onSubmit={(event) => {
        event.preventDefault();
        onSubmit({ ...draft, name: draft.name.trim(), address: draft.address.trim() });
      }}>
        <label htmlFor="point-name">Название</label>
        <input id="point-name" value={draft.name} required pattern={'.*\\S.*'} maxLength={100}
          placeholder="Например, «Уют»"
          onChange={(event) => onChange({ ...draft, name: event.target.value })} />
        <label htmlFor="point-category">Категория</label>
        <select id="point-category" value={draft.category}
          onChange={(event) => onChange({ ...draft, category: event.target.value as Category })}>
          {categories.map((category) => <option key={category}>{category}</option>)}
        </select>
        <label htmlFor="point-address">Адрес</label>
        <input id="point-address" value={draft.address} required pattern={'.*\\S.*'} maxLength={250}
          placeholder="Город, улица, дом" aria-describedby="location-hint"
          onChange={(event) => onChange({ ...draft, address: event.target.value })} />
        <p className="muted small" id="location-hint">В демо геолокация подставляет тестовый адрес Санкт-Петербурга.</p>
        <div className="actions">
          <button className="button" type="button" onClick={() => onChange({ ...draft, address: testAddress })}>Поделиться геолокацией</button>
          <button className="button button-primary" type="submit">Добавить</button>
        </div>
      </form>
    </section>
  );
}
