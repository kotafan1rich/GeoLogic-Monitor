declare global {
  interface Window {
    WebApp?: { initData?: string };
  }
}

export function getInitData(): string {
  return typeof window === 'undefined' ? '' : window.WebApp?.initData ?? '';
}
