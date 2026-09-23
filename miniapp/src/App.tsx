import { useEffect, useRef, useState } from 'react';
import { createMockPoint, initialPoints } from './mockData';
import type { PointDraft } from './mockData';
import HomeScreen from './screens/HomeScreen';
import AddPointScreen from './screens/AddPointScreen';
import ConfirmPointScreen from './screens/ConfirmPointScreen';
import MyPointsScreen from './screens/MyPointsScreen';
import PointDetailsScreen from './screens/PointDetailsScreen';
import HowItWorksScreen from './screens/HowItWorksScreen';

type Screen = 'home' | 'add' | 'confirm' | 'points' | 'details' | 'help';
const emptyDraft: PointDraft = { name: '', category: 'Кофейня', address: '' };

export default function App() {
  const [screen, setScreen] = useState<Screen>('home');
  const [points, setPoints] = useState(initialPoints);
  const [draft, setDraft] = useState<PointDraft>(emptyDraft);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const mainRef = useRef<HTMLElement>(null);
  const selectedPoint = points.find((point) => point.id === selectedId);

  useEffect(() => {
    mainRef.current?.focus({ preventScroll: true });
    window.scrollTo(0, 0);
  }, [screen]);

  function startAdding() {
    setDraft({ ...emptyDraft });
    setScreen('add');
  }

  function confirmPoint() {
    const point = createMockPoint(draft);
    setPoints((current) => [...current, point]);
    setSelectedId(point.id);
    setDraft({ ...emptyDraft });
    setScreen('details');
  }

  function goBack() {
    setScreen(screen === 'confirm' ? 'add' : screen === 'details' ? 'points' : 'home');
  }
  return (
    <div className="app">
      <header className="header">
        <span className="brand">GeoLogic</span>
      </header>

      {screen !== 'home' && (
        <nav className="navigation" aria-label="Навигация">
          <button type="button" className="back-button" onClick={goBack}>← Назад</button>
          <button type="button" className="back-button" onClick={() => setScreen('home')}>Главная</button>
        </nav>
      )}
      <main ref={mainRef} tabIndex={-1} className={`main${screen === 'home' ? '' : ' main-screen'}`}>
        {screen === 'home' && <HomeScreen onAdd={startAdding} onPoints={() => setScreen('points')} onHelp={() => setScreen('help')} />}
        {screen === 'add' && <AddPointScreen draft={draft} onChange={setDraft} onSubmit={(value) => { setDraft(value); setScreen('confirm'); }} />}
        {screen === 'confirm' && <ConfirmPointScreen draft={draft} onConfirm={confirmPoint} onEdit={() => setScreen('add')} />}
        {screen === 'points' && <MyPointsScreen points={points} onAdd={startAdding} onSelect={(id) => { setSelectedId(id); setScreen('details'); }} />}
        {screen === 'details' && selectedPoint && (
          <PointDetailsScreen point={selectedPoint} onDelete={() => {
            setPoints((current) => current.filter((point) => point.id !== selectedPoint.id));
            setSelectedId(null);
            setScreen('points');
          }} />
        )}
        {screen === 'help' && <HowItWorksScreen />}
      </main>

      <footer className="footer">Уведомления приходят в MAX</footer>
    </div>
  );
}
