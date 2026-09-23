export const categories = [
  'Кофейня', 'Кафе', 'Ресторан', 'Пекарня', 'Фастфуд', 'Другое заведение общепита',
] as const;

export type Category = typeof categories[number];
export type PointDraft = { name: string; category: Category; address: string };
export type ScoreEntry = { month: string; score: number };
export type Point = PointDraft & {
  id: string;
  updatedAt: string;
  history: [ScoreEntry, ScoreEntry, ScoreEntry];
  factors: { name: string; description: string }[];
};

export const testAddress = 'Санкт-Петербург, Невский проспект, 28';

export function createMockPoint(
  draft: PointDraft,
  scores: [number, number, number] = [6.8, 7.1, 7.4],
): Point {
  const now = new Date();
  const history = scores.map((score, index) => ({
    month: new Date(now.getFullYear(), now.getMonth() - 2 + index, 1)
      .toLocaleDateString('ru-RU', { month: 'long', year: 'numeric' }),
    score,
  })) as Point['history'];

  return {
    ...draft,
    id: crypto.randomUUID(),
    updatedAt: new Date(now.getFullYear(), now.getMonth(), 1).toLocaleDateString('ru-RU'),
    history,
    factors: [
      { name: 'Конкуренция', description: 'Количество и близость похожих заведений рядом' },
      { name: 'Инфраструктура', description: 'Объекты и сервисы вокруг точки' },
      { name: 'Транспортная доступность', description: 'Насколько удобно добраться до локации' },
    ],
  };
}

export const initialPoints: Point[] = [
  createMockPoint({ name: 'Кофейня на Невском', category: 'Кофейня', address: testAddress }, [7.2, 7.6, 8.1]),
  createMockPoint({ name: 'Пекарня у Фонтанки', category: 'Пекарня', address: 'Санкт-Петербург, набережная реки Фонтанки, 54' }, [7.9, 7.7, 7.5]),
  createMockPoint({ name: 'Кафе на Петроградской', category: 'Кафе', address: 'Санкт-Петербург, Большой проспект П.С., 45' }, [6.7, 7, 7]),
];

export function scoreChange(point: Point) {
  const change = point.history[2].score - point.history[1].score;
  return `${change > 0 ? '+' : ''}${change.toFixed(1)}`;
}
