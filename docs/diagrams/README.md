# Диаграммы

Как в референсном проекте, исходники PlantUML хранятся в src, изображения — в rendered. Диаграммы показывают проектируемый MVP.

| Диаграмма | Исходник | Изображение |
| --- | --- | --- |
| Контейнеры | [PlantUML](src/c4-containers.puml) | [SVG](rendered/c4-containers.svg) |
| Компоненты API | [PlantUML](src/c4-components.puml) | [SVG](rendered/c4-components.svg) |
| Путь предпринимателя | [PlantUML](src/cjm-owner.puml) | [SVG](rendered/cjm-owner.svg) |
| ER-модель | [PlantUML](src/er-model.puml) | [SVG](rendered/er-model.svg) |
| Мониторинг и доставка | [PlantUML](src/monitoring-sequence.puml) | [SVG](rendered/monitoring-sequence.svg) |

SVG генерируются из исходников PlantUML. Пример повторной генерации через Kroki из корня репозитория:

```bash
curl -sS --fail -X POST -H 'Content-Type: text/plain' \
  --data-binary @docs/diagrams/src/er-model.puml \
  https://kroki.io/plantuml/svg \
  -o docs/diagrams/rendered/er-model.svg
```

Команда передаёт содержимое исходника сервису Kroki. После изменения схемы обновлять и исходник, и SVG.
