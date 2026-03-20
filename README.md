# URL Shortener

Сервис, предоставляющий простой API для создания сокращенных ссылок по формату:

- Ссылка уникальна — на один оригинальный URL ссылается только одна сокращенная ссылка
- Длина ссылки — 10 символов
- Символы: латинский алфавит (нижний и верхний регистр), цифры и `_`
---
- В качестве хранилища используется postgres и отдельный пакет.
- Пакет представляет собой LRU-cache. По моему мнению самое оптимальное решение, чтоб не сожрало всю оперативку.
- Базово покрыл unit тестами.

Сервис на Docker Hub: `gamexost/url-shortener`

## Запуск

**In-memory:**
```
docker run -p 8080:8080 gamexost/url-shortener:latest
```

**Postgres:**
```
docker run -p 8080:8080 \
  -e STORAGE_TYPE=postgres \
  -e DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=disable \
  gamexost/url-shortener:latest
```
`postgres://user:password@host:port/bd_name?sslmode=disable`

