type Props = { onAdd: () => void; onPoints: () => void; onHelp: () => void };

export default function HomeScreen({ onAdd, onPoints, onHelp }: Props) {
  return (
    <section className="welcome" aria-labelledby="welcome-title">
      <div className="location-mark" aria-hidden="true">📍</div>
      <h1 id="welcome-title">GeoLogic</h1>
      <p className="description">Управление геолокациями бизнеса</p>
      <div className="actions">
        <button className="button button-primary" type="button" onClick={onAdd}>📍 Добавить точку</button>
        <button className="button" type="button" onClick={onPoints}>📌 Мои точки</button>
        <button className="button" type="button" onClick={onHelp}>ℹ️ Как это работает</button>
      </div>
    </section>
  );
}
