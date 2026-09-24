import { useEffect, useRef, useState } from 'react';
import type { Dispatch, SetStateAction } from 'react';
import { api, ApiError, errorMessage } from '../api';
import type { Address, BusinessType, LoadState, PointDraft } from '../types';

type Props = {
  draft: PointDraft;
  onChange: Dispatch<SetStateAction<PointDraft>>;
  businessTypes: LoadState<BusinessType[]>;
  onRetryTypes: () => void;
  onUnauthorized: () => void;
  onSubmit: () => void;
};

export default function AddPointScreen({ draft, onChange, businessTypes, onRetryTypes, onUnauthorized, onSubmit }: Props) {
  const [suggestions, setSuggestions] = useState<LoadState<Address[]> | null>(null);
  const [searchRevision, setSearchRevision] = useState(0);
  const [locating, setLocating] = useState(false);
  const [locationError, setLocationError] = useState('');
  const locationVersion = useRef(0);
  const reverseController = useRef<AbortController | null>(null);

  useEffect(() => {
    if (draft.selectedAddress || !draft.addressQuery.trim() || locating) {
      setSuggestions(null);
      return;
    }
    const controller = new AbortController();
    setSuggestions({ status: 'loading' });
    const timeout = window.setTimeout(() => {
      api.suggestions(draft.addressQuery.trim(), controller.signal).then((data) => {
        if (!controller.signal.aborted) setSuggestions({ status: 'success', data });
      }).catch((error: unknown) => {
        if (controller.signal.aborted) return;
        if (error instanceof ApiError && error.status === 401) onUnauthorized();
        setSuggestions({ status: 'error', message: errorMessage(error) });
      });
    }, 350);
    return () => { window.clearTimeout(timeout); controller.abort(); };
  }, [draft.addressQuery, draft.selectedAddress, locating, searchRevision, onUnauthorized]);

  useEffect(() => () => {
    locationVersion.current += 1;
    reverseController.current?.abort();
  }, []);

  function cancelLocation() {
    locationVersion.current += 1;
    reverseController.current?.abort();
    setLocating(false);
    setLocationError('');
  }

  function chooseAddress(address: Address) {
    cancelLocation();
    setSuggestions(null);
    onChange((current) => ({ ...current, addressQuery: address.address, selectedAddress: address }));
  }

  function locate() {
    cancelLocation();
    if (!navigator.geolocation) {
      setLocationError('Геолокация недоступна. Найдите адрес через поиск.');
      return;
    }
    const version = locationVersion.current;
    const controller = new AbortController();
    reverseController.current = controller;
    setLocating(true);
    setSuggestions(null);
    navigator.geolocation.getCurrentPosition(async (position) => {
      if (version !== locationVersion.current) return;
      try {
        const address = await api.address(position.coords.latitude, position.coords.longitude, controller.signal);
        if (version !== locationVersion.current) return;
        onChange((current) => ({ ...current, addressQuery: address.address, selectedAddress: address }));
      } catch (error) {
        if (version !== locationVersion.current) return;
        if (error instanceof ApiError && error.status === 401) onUnauthorized();
        setLocationError(errorMessage(error) + ' Можно найти адрес через поиск.');
      } finally {
        if (version === locationVersion.current) setLocating(false);
      }
    }, (error) => {
      if (version !== locationVersion.current) return;
      const message = error.code === 1 ? 'Доступ к геолокации запрещён.'
        : error.code === 3 ? 'Время ожидания геолокации истекло.' : 'Не удалось определить местоположение.';
      setLocationError(message + ' Найдите адрес через поиск.');
      setLocating(false);
    }, { timeout: 10000, maximumAge: 0 });
  }

  const validCategory = businessTypes.status === 'success'
    && businessTypes.data.some((type) => type.id === draft.business_type_id);
  const canSubmit = Boolean(draft.name.trim() && validCategory && draft.selectedAddress && !locating);

  return (
    <section className="screen" aria-labelledby="add-title">
      <h1 id="add-title">Добавить точку</h1>
      <form className="point-form" onSubmit={(event) => { event.preventDefault(); if (canSubmit) onSubmit(); }}>
        <label htmlFor="point-name">Название</label>
        <input id="point-name" value={draft.name} required pattern={'.*\\S.*'} maxLength={100}
          placeholder="Например, «Уют»"
          onChange={(event) => onChange((current) => ({ ...current, name: event.target.value }))} />
        <label htmlFor="point-category">Категория</label>
        <select id="point-category" value={draft.business_type_id} required
          disabled={businessTypes.status !== 'success' || businessTypes.data.length === 0}
          onChange={(event) => onChange((current) => ({ ...current, business_type_id: event.target.value }))}>
          <option value="">Выберите категорию</option>
          {businessTypes.status === 'success' && businessTypes.data.map((type) =>
            <option key={type.id} value={type.id}>{type.infra_type.name}</option>)}
        </select>
        {businessTypes.status === 'loading' && <p role="status">Загрузка категорий…</p>}
        {businessTypes.status === 'error' && <div role="alert"><p>{businessTypes.message}</p><button className="button" type="button" onClick={onRetryTypes}>Повторить загрузку категорий</button></div>}
        {businessTypes.status === 'success' && businessTypes.data.length === 0 && <div className="card" role="status"><p>Категории пока недоступны. Добавление точки невозможно.</p><button className="button" type="button" onClick={onRetryTypes}>Обновить категории</button></div>}
        <label htmlFor="point-address">Адрес</label>
        <input id="point-address" value={draft.addressQuery} required maxLength={500}
          placeholder="Город, улица, дом" autoComplete="off" aria-describedby="location-hint"
          onChange={(event) => {
            cancelLocation();
            setSuggestions(null);
            onChange((current) => ({ ...current, addressQuery: event.target.value, selectedAddress: null }));
          }} />
        <p className="muted small" id="location-hint">Выберите адрес из подсказок или поделитесь геолокацией. Поиск доступен для Санкт-Петербурга.</p>
        {suggestions?.status === 'loading' && <p role="status">Поиск адресов…</p>}
        {suggestions?.status === 'error' && <div role="alert"><p>{suggestions.message}</p><button className="button" type="button" onClick={() => setSearchRevision((value) => value + 1)}>Повторить поиск</button></div>}
        {suggestions?.status === 'success' && (suggestions.data.length === 0
          ? <p role="status">Адреса не найдены. Уточните запрос.</p>
          : <ul className="address-suggestions" aria-label="Варианты адреса">
            {suggestions.data.map((address, index) => <li key={index}><button type="button" className="button" onClick={() => chooseAddress(address)}>{address.address}</button></li>)}
          </ul>)}
        {draft.selectedAddress && <p className="muted small" role="status">Адрес выбран.</p>}
        {locationError && <p className="error-message" role="alert">{locationError}</p>}
        <div className="actions">
          <button className="button" type="button" disabled={locating} onClick={locate}>{locating ? 'Определение адреса…' : 'Поделиться геолокацией'}</button>
          <button className="button button-primary" type="submit" disabled={!canSubmit}>Добавить</button>
        </div>
      </form>
    </section>
  );
}
