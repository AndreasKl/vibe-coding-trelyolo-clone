# FlowBoard

> A fast, focused Kanban board — built with Go and SvelteKit.

![FlowBoard demo](docs/screenshots/demo.gif)

---

## Screenshots

| Dashboard | Board |
|-----------|-------|
| ![Dashboard](docs/screenshots/dashboard.png) | ![Board](docs/screenshots/board.png) |

| Drag &amp; drop | Sign in |
|------------|---------|
| ![Drag and drop](docs/screenshots/drag.png) | ![Login](docs/screenshots/login.png) |

> **Note**: Drop your screenshots into `docs/screenshots/` and remove this note.

---

## Features

- **Boards** — create as many boards as you need, each pre-loaded with *Todo / Doing / Done* columns
- **Columns** — add, rename, delete, and drag to reorder
- **Cards** — create inline, edit in place, drag between columns
- **Auth** — email/password sign-up or one-click sign-in via Google / Microsoft OAuth2
- **Sessions** — HTTP-only cookies backed by the database; no JWT, no localStorage
- **Real-time feel** — optimistic UI with instant drag feedback

---

## Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.24, `net/http` (Go 1.22 mux) |
| Database | PostgreSQL 16, `pgx/v5` |
| Frontend | SvelteKit 2, Svelte 5 (runes), TypeScript |
| Auth | bcrypt + OAuth2 (Google, Microsoft) |
| Dev DB | Docker Compose |

---

## Getting Started

**Prerequisites**: Go 1.24+, Node.js 20+, pnpm, Docker

```bash
# 1. Clone and enter the repo
git clone https://github.com/your-org/flowboard.git
cd flowboard

# 2. Configure environment
cp .env.example .env          # edit if needed

# 3. Start PostgreSQL
docker compose up -d

# 4. Start the Go backend (runs DB migrations automatically)
make dev-backend              # http://localhost:8080

# 5. In a second terminal, start the frontend
make dev-frontend             # http://localhost:5173
```

Open **http://localhost:5173** and sign up.

---

## Configuration

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | `postgres://trello:trello@localhost:5432/trello?sslmode=disable` | PostgreSQL connection string |
| `PORT` | `8080` | Backend port |
| `COOKIE_DOMAIN` | `localhost` | Session cookie domain |
| `BASE_URL` | `http://localhost:5173` | Frontend origin (for CORS) |
| `GOOGLE_CLIENT_ID` | *(empty)* | Google OAuth2 — leave blank to disable |
| `GOOGLE_CLIENT_SECRET` | *(empty)* | |
| `MICROSOFT_CLIENT_ID` | *(empty)* | Microsoft OAuth2 — leave blank to disable |
| `MICROSOFT_CLIENT_SECRET` | *(empty)* | |

---

## Architecture

```
Browser
  │
  ├── /api/*  ──────────────────────► Go backend (:8080)
  │         (Vite proxy in dev;              │
  │          same origin in prod)       PostgreSQL
  │
  └── everything else ──► SvelteKit (:5173 in dev)
```

**Auth flow**: every navigation triggers `+layout.server.ts`, which calls `GET /api/auth/me`. If that fails, the user is redirected to `/login`. Auth state never lives in the browser — no JWT, no localStorage.

**Data model**:

```
users
  └── boards
        └── board_columns  (ordered by position)
              └── cards    (ordered by position)
```

Moving a card or column is a single transaction: close the gap at the source, open a gap at the target, update the row.

---

## Project Structure

```
backend/
  cmd/server/main.go
  internal/
    auth/          # signup, login, logout, OAuth2, sessions
    board/         # boards, columns, cards CRUD + move operations
    database/      # connection + embedded migrations
    httputil/      # JSON/error response helpers
    server/        # HTTP mux + middleware chain

frontend/src/
  lib/
    api.ts         # typed fetch wrapper
    types.ts       # TypeScript interfaces
    components/    # Navbar, BoardView, BoardColumn, BoardCard
  routes/
    +page.svelte               # dashboard
    login/, signup/            # auth pages
    boards/[id]/+page.svelte   # kanban board
```

---

## API Reference

### Auth

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/auth/signup` | Create account |
| POST | `/api/auth/login` | Log in |
| POST | `/api/auth/logout` | Log out |
| GET | `/api/auth/me` | Current user |
| GET | `/api/auth/oauth/{provider}` | Start OAuth flow |
| GET | `/api/auth/oauth/{provider}/callback` | OAuth callback |

### Boards, Columns, Cards

| Method | Path | Description |
|--------|------|-------------|
| GET/POST | `/api/boards` | List / create boards |
| GET/DELETE | `/api/boards/{id}` | Get / delete board |
| POST | `/api/boards/{boardID}/columns` | Add column |
| PATCH/DELETE | `/api/columns/{id}` | Rename or delete column |
| POST | `/api/columns/{id}/move` | Reorder column `{ position }` |
| POST | `/api/columns/{columnID}/cards` | Create card |
| PATCH/DELETE | `/api/cards/{id}` | Update or delete card |
| POST | `/api/cards/{id}/move` | Move card `{ column_id, position }` |

---

## Development

```bash
make test            # run all tests
make test-backend    # Go tests (requires PostgreSQL)
make test-frontend   # Vitest only

# Run a single Go test
cd backend && go test ./internal/board/... -run TestMoveCard

# Run a single frontend test
cd frontend && pnpm test --run -- BoardCard

# Type-check + lint frontend
cd frontend && pnpm check && pnpm lint
```

---

## Make Targets

| Target | Description |
|--------|-------------|
| `make dev-backend` | Go backend with hot reload |
| `make dev-frontend` | SvelteKit dev server |
| `make dev` | DB + both servers |
| `make build` | Production binary + frontend bundle |
| `make test` | All tests |
| `make db-up` | Start PostgreSQL |
| `make db-down` | Stop PostgreSQL |
