import { useCallback, useEffect, useRef, useState } from 'react';
import { api, ApiError, errorMessage } from './api';
import { getInitData } from './max';
import type { BusinessType, LoadState, PointDraft, TrackedLocation } from './types';
import HomeScreen from './screens/HomeScreen';
import AddPointScreen from './screens/AddPointScreen';
import ConfirmPointScreen from './screens/ConfirmPointScreen';
import MyPointsScreen from './screens/MyPointsScreen';
import PointDetailsScreen from './screens/PointDetailsScreen';
import HowItWorksScreen from './screens/HowItWorksScreen';

type Screen = 'home' | 'add' | 'confirm' | 'points' | 'details' | 'help';
const emptyDraft: PointDraft = { name: '', business_type_id: '', addressQuery: '', selectedAddress: null };

export default function App() {
  const [screen, setScreen] = useState<Screen>('home');
  const [points, setPoints] = useState<LoadState<TrackedLocation[]>>({ status: 'loading' });
  const [businessTypes, setBusinessTypes] = useState<LoadState<BusinessType[]>>({ status: 'loading' });
  const [pointsRevision, setPointsRevision] = useState(0);
  const [typesRevision, setTypesRevision] = useState(0);
  const [draft, setDraft] = useState<PointDraft>(emptyDraft);
  const [selectedPoint, setSelectedPoint] = useState<TrackedLocation | null>(null);
  const [expired, setExpired] = useState(false);
  const [busy, setBusy] = useState(false);
  const [mutationError, setMutationError] = useState('');
  const [notice, setNotice] = useState('');
  const mutationPending = useRef(false);
  const mainRef = useRef<HTMLElement>(null);
  const hasInitData = Boolean(getInitData());
  const onUnauthorized = useCallback(() => setExpired(true), []);

  useEffect(() => {
    mainRef.current?.focus({ preventScroll: true });
    window.scrollTo(0, 0);
  }, [screen, expired]);

  useEffect(() => {
    if (!hasInitData || expired) return;
    const controller = new AbortController();
    setBusinessTypes({ status: 'loading' });
    api.businessTypes(controller.signal).then((data) => {
      if (!controller.signal.aborted) setBusinessTypes({ status: 'success', data });
    }).catch((error: unknown) => {
      if (controller.signal.aborted) return;
      if (error instanceof ApiError && error.status === 401) onUnauthorized();
      setBusinessTypes({ status: 'error', message: errorMessage(error) });
    });
    return () => controller.abort();
  }, [typesRevision, hasInitData, expired, onUnauthorized]);

  useEffect(() => {
    if (!hasInitData || expired) return;
    const controller = new AbortController();
    setPoints({ status: 'loading' });
    api.points(controller.signal).then((data) => {
      if (!controller.signal.aborted) setPoints({ status: 'success', data });
    }).catch((error: unknown) => {
      if (controller.signal.aborted) return;
      if (error instanceof ApiError && error.status === 401) onUnauthorized();
      const message = error instanceof ApiError && error.status === 404
        ? 'Пользователь ещё не зарегистрирован. Регистрация через /start в боте должна быть настроена на стороне сервиса.'
        : errorMessage(error);
      setPoints({ status: 'error', message });
    });
    return () => controller.abort();
  }, [pointsRevision, hasInitData, expired, onUnauthorized]);

  function startAdding() {
    setDraft({ ...emptyDraft });
    setMutationError('');
    setNotice('');
    setScreen('add');
  }

  function showPoints() {
    setMutationError('');
    setPointsRevision((value) => value + 1);
    setScreen('points');
  }

  async function confirmPoint() {
    if (mutationPending.current || !draft.selectedAddress || !draft.name.trim() || !draft.business_type_id) return;
    mutationPending.current = true;
    setBusy(true);
    setMutationError('');
    try {
      const point = await api.createPoint({
        name: draft.name.trim(), business_type_id: draft.business_type_id, ...draft.selectedAddress,
      });
      setPoints((current) => ({
        status: 'success', data: [...(current.status === 'success' ? current.data : []), point],
      }));
      setSelectedPoint(point);
      setDraft({ ...emptyDraft });
      setScreen('details');
    } catch (error) {
      if (error instanceof ApiError && error.status === 401) onUnauthorized();
      let message = errorMessage(error);
      if (error instanceof ApiError && error.status === 404) {
        message = 'Пользователь или категория не найдены. Проверьте регистрацию в боте и обновите категории.';
      }
      if (error instanceof ApiError && (error.status === 0 || error.status >= 500 && error.status !== 503)) {
        message += ' Результат создания неизвестен. Проверьте «Мои точки» перед повторным добавлением.';
      }
      setMutationError(message);
    } finally {
      mutationPending.current = false;
      setBusy(false);
    }
  }

  async function deletePoint() {
    if (!selectedPoint || mutationPending.current) return;
    mutationPending.current = true;
    setBusy(true);
    setMutationError('');
    try {
      await api.deletePoint(selectedPoint.id);
      setPoints((current) => current.status === 'success'
        ? { status: 'success', data: current.data.filter((point) => point.id !== selectedPoint.id) }
        : current);
      setSelectedPoint(null);
      showPoints();
    } catch (error) {
      if (error instanceof ApiError && error.status === 401) onUnauthorized();
      if (error instanceof ApiError && error.status === 404) {
        setNotice('Точка больше недоступна.');
        showPoints();
      } else {
        setMutationError(errorMessage(error));
      }
    } finally {
      mutationPending.current = false;
      setBusy(false);
    }
  }

  function goBack() {
    setMutationError('');
    if (screen === 'details') showPoints();
    else setScreen(screen === 'confirm' ? 'add' : 'home');
  }

  const categoryName = (id: string) => businessTypes.status === 'success'
    ? businessTypes.data.find((type) => type.id === id)?.infra_type.name ?? 'Категория недоступна'
    : 'Категория недоступна';

  return (
    <div className="app">
      <header className="header"><span className="brand">GeoLogic</span></header>
      {hasInitData && !expired && screen !== 'home' && (
        <nav className="navigation" aria-label="Навигация">
          <button type="button" className="back-button" disabled={busy} onClick={goBack}>← Назад</button>
          <button type="button" className="back-button" disabled={busy} onClick={() => { setMutationError(''); setScreen('home'); }}>Главная</button>
        </nav>
      )}
      <main ref={mainRef} tabIndex={-1} className={`main${screen === 'home' && hasInitData && !expired ? '' : ' main-screen'}`}>
        {!hasInitData || expired ? (
          <section className="screen" aria-labelledby="auth-title">
            <h1 id="auth-title">Откройте GeoLogic через MAX</h1>
            <p className="card" role="status">{expired
              ? 'Сессия недействительна или истекла. Заново откройте Mini App через MAX.'
              : 'Данные запуска MAX недоступны. Откройте мини-приложение из MAX.'}</p>
          </section>
        ) : (
          <>
            {screen === 'home' && <HomeScreen onAdd={startAdding} onPoints={showPoints} onHelp={() => setScreen('help')} />}
            {screen === 'add' && <AddPointScreen draft={draft} onChange={setDraft} businessTypes={businessTypes}
              onRetryTypes={() => setTypesRevision((value) => value + 1)} onUnauthorized={onUnauthorized}
              onSubmit={() => { setMutationError(''); setScreen('confirm'); }} />}
            {screen === 'confirm' && <ConfirmPointScreen draft={draft} categoryName={categoryName(draft.business_type_id)}
              busy={busy} error={mutationError} onConfirm={confirmPoint} onEdit={() => { setMutationError(''); setScreen('add'); }} />}
            {screen === 'points' && <MyPointsScreen points={points} notice={notice} onRetry={showPoints}
              onAdd={startAdding} onUnauthorized={onUnauthorized}
              onSelect={(point) => { setSelectedPoint(point); setMutationError(''); setScreen('details'); }} />}
            {screen === 'details' && selectedPoint && <PointDetailsScreen key={selectedPoint.id}
              point={selectedPoint} categoryName={categoryName(selectedPoint.business_type_id)}
              busy={busy} error={mutationError} onDelete={deletePoint} onUnauthorized={onUnauthorized} />}
            {screen === 'help' && <HowItWorksScreen />}
          </>
        )}
      </main>
      <footer className="footer">Уведомления в MAX — в планах</footer>
    </div>
  );
}
