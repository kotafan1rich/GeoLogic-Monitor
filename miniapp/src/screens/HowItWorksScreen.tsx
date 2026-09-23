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
          <p>Сервис оценивает ситуацию вокруг бизнеса и отслеживает новые события.</p>
        </li>
        <li className="card">
          <h2>Получайте уведомления в MAX</h2>
          <p>О важных изменениях вы узнаёте сразу.</p>
        </li>
      </ol>
      <p className="card info-card">Smart Score пересчитывается раз в месяц. История доступна за последние 3 месяца.</p>
    </section>
  );
}
