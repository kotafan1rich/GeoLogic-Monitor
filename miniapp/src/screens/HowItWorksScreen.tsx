export default function HowItWorksScreen() {
  return (
    <section className="screen" aria-labelledby="help-title">
      <h1 id="help-title">Как это работает</h1>
      <ol className="steps">
        <li className="card">
          <h2>Добавьте точку</h2>
          <p>Укажите название, категорию и адрес или поделитесь геолокацией.</p>
        </li>
        <li className="card">
          <h2>GeoLogic анализирует окружение</h2>
          <p>При добавлении точки сервис рассчитывает Smart Score по данным об окружении бизнеса.</p>
        </li>
        <li className="card">
          <h2>Уведомления в MAX — в планах</h2>
          <p>Доставка уведомлений о важных изменениях пока не подключена.</p>
        </li>
      </ol>
      <p className="card info-card">Smart Score пересчитывается раз в месяц. История Smart Score доступна за последние 3 месяца.</p>
    </section>
  );
}
