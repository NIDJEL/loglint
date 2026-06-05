# loglint

`loglint` — это линтер для Go, который проверяет сообщения в логах.

Линтер написан через `golang.org/x/tools/go/analysis`. Его можно запускать напрямую через `go run` или как кастомный линтер для `golangci-lint`.

## Что проверяет

Линтер проверяет, что лог-сообщения:

- начинаются со строчной буквы;
- написаны на английском языке;
- не содержат спецсимволы и эмодзи;
- не содержат чувствительные данные: `password`, `token`, `api_key`, `secret`.

## Примеры

Неправильно:

```go
slog.Info("Starting server")
slog.Info("запуск сервера")
slog.Info("warning: something went wrong...")
slog.Info("server started 🔥")
slog.Info("token: " + token)
```

Правильно:

```go
slog.Info("starting server")
slog.Error("failed to connect")
slog.Info("user authenticated successfully")
slog.Info("api request completed")
slog.Info("token validated")
```

## Запуск тестов

```bash
go test ./...
```

## Запуск напрямую

Проверить примеры:

```bash
go run ./cmd/loglint ./examples
```

На Windows:

```powershell
go run .\cmd\loglint\ ./examples
```

Проверить весь проект:

```bash
go run ./cmd/loglint ./...
```

## Запуск через golangci-lint

Сначала нужно установить `golangci-lint`:

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
```

Проверить установку:

```bash
golangci-lint --version
```

Собрать кастомный `golangci-lint` с этим линтером:

```bash
golangci-lint custom -v
```

После сборки появится файл `custom-gcl` или `custom-gcl.exe`.

Запуск на Windows:

```powershell
.\custom-gcl.exe run ./examples
```

Запуск на Linux/macOS:

```bash
./custom-gcl run ./examples
```

## Пример вывода

```text
examples/bad.go:8:2: log message should start with lowercase letter
examples/bad.go:10:2: log message should not contain special characters or emoji
examples/bad.go:12:2: log message should be written in English
examples/bad.go:13:2: log message should not contain sensitive data
```

## Что реализовано

- analyzer на базе `golang.org/x/tools/go/analysis`;
- проверка лог-сообщений по требованиям;
- тесты через `analysistest`;
- тестовые файлы в `testdata`;
- запуск через `go run`;
- интеграция с `golangci-lint`.