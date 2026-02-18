# Trello Clone

A Kanban board web application. Create boards with columns, add cards, and drag them between columns.

## Stack

- **Backend**: Go, PostgreSQL 16, server-side sessions
- **Frontend**: SvelteKit (Svelte 5), TypeScript
- **Auth**: Email/password (bcrypt) + OAuth2 (Google, Microsoft)

## Architecture

The app is split into two independent processes that talk over HTTP.

```
Browser
  │
  ├── GET /api/*  ──────────────────────────────► Go backend (:8080)
  │                                                    │
  │   (Vite proxy in dev; same origin in prod)         └── PostgreSQL
  │
  └── Everything else ──► SvelteKit (:5173 dev)
```

**Request flow**

1. SvelteKit's `+layout.server.ts` runs on every navigation, calls `GET /api/auth/me`, and either returns the user or redirects to `/login`. Auth state is never stored client-side.
2. Page components call the Go API directly from the browser via `$lib/api.ts`. All API calls include the session cookie automatically.
3. The Go backend validates the session cookie on every protected route via the `RequireAuth` middleware, which looks up the token in the DB and injects the `*User` into the request context.

**Session management**

Sessions are stored as rows in PostgreSQL. The cookie holds a random token; there is no JWT. On login/signup the server sets an HTTP-only, `SameSite=Lax` cookie. On logout the token row is deleted and the cookie is cleared.

**Data model**

```
users
  └── boards (user_id FK)
        └── board_columns (board_id FK, ordered by position)
              └── cards (column_id FK, ordered by position)
```

Card positions are integers. Moving a card is a single transaction: close the gap in the source column, open a gap in the target column, then update the card's `column_id` and `position`.

## Prerequisites

- Go 1.24+
- Node.js 25+ and pnpm
- Docker (for PostgreSQL)

## Getting Started

```bash
# Start PostgreSQL
docker compose up -d

# Start the Go backend (runs migrations automatically)
make dev-backend

# In a separate terminal, start the frontend
make dev-frontend
```

Open http://localhost:5173.

## Configuration

Copy `.env.example` to `.env` and adjust as needed:

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | `postgres://trello:trello@localhost:5432/trello?sslmode=disable` | PostgreSQL connection string |
| `PORT` | `8080` | Backend port |
| `COOKIE_DOMAIN` | `localhost` | Session cookie domain |
| `BASE_URL` | `http://localhost:5173` | Frontend origin (for CORS) |
| `GOOGLE_CLIENT_ID` | | Google OAuth2 client ID |
| `GOOGLE_CLIENT_SECRET` | | Google OAuth2 client secret |
| `MICROSOFT_CLIENT_ID` | | Microsoft OAuth2 client ID |
| `MICROSOFT_CLIENT_SECRET` | | Microsoft OAuth2 client secret |

OAuth is optional. Leave the client ID/secret empty to skip.

## Project Structure

```
backend/
  cmd/server/main.go          # Entry point
  internal/
    auth/                      # Signup, login, logout, OAuth2, sessions
    board/                     # Boards, columns, cards CRUD + card moves
    database/                  # DB connection, migrations
    httputil/                  # JSON response helpers
    server/                    # HTTP mux, middleware

frontend/
  src/
    lib/
      api.ts                   # Typed fetch wrapper
      types.ts                 # TypeScript interfaces
      components/              # Navbar, BoardView, BoardColumn, BoardCard
    routes/
      +page.svelte             # Dashboard (list/create boards)
      login/+page.svelte       # Login page
      signup/+page.svelte      # Signup page
      boards/[id]/+page.svelte # Kanban board view
```

## API Endpoints

### Auth

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/api/auth/signup` | No | Create account |
| POST | `/api/auth/login` | No | Log in |
| POST | `/api/auth/logout` | Yes | Log out |
| GET | `/api/auth/me` | Yes | Current user |
| GET | `/api/auth/oauth/{provider}` | No | Start OAuth flow |
| GET | `/api/auth/oauth/{provider}/callback` | No | OAuth callback |

### Boards

| Method | Path | Description |
|---|---|---|
| GET | `/api/boards` | List boards |
| POST | `/api/boards` | Create board (auto-creates Todo/Doing/Done columns) |
| GET | `/api/boards/{id}` | Get board with columns and cards |
| DELETE | `/api/boards/{id}` | Delete board |

### Columns

| Method | Path | Description |
|---|---|---|
| POST | `/api/boards/{boardID}/columns` | Add column |
| PATCH | `/api/columns/{id}` | Rename/reorder column |
| DELETE | `/api/columns/{id}` | Delete column |

### Cards

| Method | Path | Description |
|---|---|---|
| POST | `/api/columns/{columnID}/cards` | Create card |
| PATCH | `/api/cards/{id}` | Update title/description |
| DELETE | `/api/cards/{id}` | Delete card |
| POST | `/api/cards/{id}/move` | Move card to column + position |

## Make Targets

| Target | Description |
|---|---|
| `make dev-backend` | Run Go backend with hot reload |
| `make dev-frontend` | Run SvelteKit dev server |
| `make dev` | Start DB + both servers |
| `make build` | Build backend binary + frontend |
| `make test` | Run Go tests |
| `make db-up` | Start PostgreSQL |
| `make db-down` | Stop PostgreSQL |
