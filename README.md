![CI](https://github.com/T-AKHMAD/gnotes-cli/actions/workflows/ci.yml/badge.svg)

# gnotes-cli

CLI-клиент для проекта **gopher-notes**: логин, logout, профиль, заметки (list/create/get/delete).
Токен сохраняется локально в `~/.gnotes/token`.

Server: https://github.com/T-AKHMAD/gopher-notes

## Requirements
- Go

## Run

> Убедись, что сервер **gopher-notes** запущен на `http://localhost:8080`.

### Login (сохраняет токен)
```bash
go run ./cmd/gnotes login --email a@b.com --password 123
```

### Me
```bash
go run ./cmd/gnotes me
```

### Notes list
```bash
go run ./cmd/gnotes notes list
```

### Notes create
```bash
go run ./cmd/gnotes notes create --title "t1" --body "b1"
```

### Notes get
```bash
go run ./cmd/gnotes notes get 1
```

### Notes delete
```bash
go run ./cmd/gnotes notes delete 1
```

### Logout (удаляет сессию на сервере + локальный токен)
```bash
go run ./cmd/gnotes logout
```

## Options
У всех команд есть флаг:
- `--base-url` (по умолчанию `http://localhost:8080`)

