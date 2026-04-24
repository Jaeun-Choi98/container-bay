# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Backend (Go)
```bash
cd backend
go run ./cmd/main.go          # run the server
go build -o server ./cmd/main.go  # build binary
go test ./...                 # run all tests
go test ./internal/...        # run specific package tests
```
The server reads `backend/env.ini` at startup — it must exist before running.

### Frontend (React/TypeScript)
```bash
cd front
npm start        # dev server on :3000
npm run build    # production build → front/build/
npm test         # run tests
```

### Full deployment
The backend serves the React SPA via `SpaHandlerRoot("build", "index.html")`, looking for a `build/` folder relative to the working directory. Copy `front/build/` into `backend/build/` after building the frontend.

## Configuration

The backend uses an INI config file at `backend/env.ini`. Copy `env.exampleini` as a starting point. Required sections and keys:

| Section | Keys |
|---|---|
| `[LOG]` | `MAX_AGE` (days, negative = delete old) |
| `[REST]` | `IP`, `PORT` |
| `[DIR]` | `REPO_DIR`, `SHELL_DIR`, `VOLUME_DIR` |
| `[GITAUTH]` | `ID`, `PASSWD` (for cloning private repos) |
| `[DOCKER]` | `DOCKER_REPO_IP`, `DOCKER_REPO_PORT`, `DOCKER_REPO_ID`, `DOCKER_REPO_PASSWD` |
| `[BUILD_SERVER]` | `BUILD_SERVER_PASSWD` (sudo password on the build host) |
| `[SHELL]` | `SHELL_NAME` (e.g. `bash`) |
| `[REDIS]` | `IP`, `PORT`, `PASSWD`, `DB`, `PROTOCOL`, `TIMEOUT` |

Frontend API base URL is set via `REACT_APP_API_BASE_URL` in `front/.env`.

## Architecture

### Backend
`cmd/main.go` → `container.NewContainer()` wires all dependencies as a singleton DI container → `app.NewApplication()` starts the HTTP server with graceful shutdown on SIGINT/SIGTERM.

Layers (inner to outer):
- **Config** (`internal/config/`) — reads `env.ini` at startup; all other layers receive `*Config`
- **Redis** (`internal/redis/`) — generic repository abstraction over go-redis; used to track Docker registry login sessions
- **Service** (`internal/service/`) — `ApiServiceInterface` defines all business operations; implemented by `ApiService` in `api-service/`
- **Transport** (`internal/transport/http/rest/`) — Gin router, middleware (CORS, JWT, logging, error, SPA), controllers, request/response types

**How Docker operations work:** All Docker commands are executed as shell scripts using the `github.com/Jaeun-Choi98/modules/shell` package. Scripts are written to `cfg.ShellDir`, run with `cfg.ShellName` (bash), and prefix every command with `echo <BuildSvrPasswd> | sudo -S` to gain root. Remote daemons are targeted with `docker -H tcp://HOST:PORT`. Output (`stdout`/`stderr`) is parsed by splitting on double-spaces and rejoining fields with `;` as a column separator — this is what the frontend receives and renders.

**Build flow:** `/build` endpoint clones a Git repo into `cfg.RepoDir/<pjtName>`, runs `docker build`, tags, pushes to the private registry, then prunes and removes the local clone.

### Frontend
Standard CRA app. Key patterns:
- `api/client.ts` — singleton `ApiClient` class; attaches `Bearer <token>` from `localStorage` on every request
- `api/hooks/useApi.ts` — `useApi<T>(apiFunc)` hook encapsulates loading/error/result state for any async call; call `.execute(...args)` to trigger, `.reset()` to clear
- `services/` — thin wrappers that call `apiClient` methods with typed request/response shapes from `api/types/`
- All API responses follow `BaseResponse { result: number, data: any }` where `result === 0` means success; `data["stdout"]` and `data["stderr"]` are `string[]` with `;`-separated columns (first element is always a header label, slice from index 1 to get data rows)
