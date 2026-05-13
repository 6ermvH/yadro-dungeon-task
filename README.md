# YADRO Dungeon Task

Реализация обработчика событий подземелья на Go для тестового задания YADRO Telecom.

## Запуск

```bash
go run ./cmd/dungeon task/config.json task/events
```

Также поддерживаются произвольные пути:

```bash
go run ./cmd/dungeon <путь_к_конфигу> <путь_к_событиям>
```

## Юнит-тесты

```bash
go test ./internal/...
```

## E2E-тесты

```bash
go test ./tests/...
```

Запуск всех тестов:

```bash
go test ./...
```

## Линтер

В проекте настроен `golangci-lint` v2:

```bash
golangci-lint run ./...
```

## Архитектура

```text
cmd -> app
app -> config, parser, simulator, output
config -> domain
parser -> domain
simulator -> domain
output -> domain
domain -> stdlib
```

Ответственности пакетов:

- `cmd/dungeon` — точка входа.
- `internal/app` — оркестрация CLI.
- `internal/config` — загрузка и валидация JSON-конфигурации.
- `internal/parser` — парсинг файла событий.
- `internal/domain` — основные типы и методы состояния.
- `internal/simulator` — правила игры, обработчики событий, переходы состояний, расчёт отчёта.
- `internal/output` — форматирование вывода.