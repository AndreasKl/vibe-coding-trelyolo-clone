# Backend

Go HTTP server with PostgreSQL. Runs on port 8080.

## Structure

```
cmd/server/
  main.go               # Entry point: loads config, connects DB, runs migrations, starts server

internal/
  auth/
    handler.go          # HTTP handlers: signup, login, logout, /me
    oauth.go            # OAuth2 flow (Google, Microsoft)
    middleware.go       # RequireAuth middleware
    store.go            # DB queries: users, sessions
    session.go          # Session token generation
    password.go         # bcrypt helpers

  board/
    handler.go          # HTTP handlers: boards, columns, cards CRUD + card move
    store.go            # DB queries for boards/columns/cards
    model.go            # Domain types

  database/
    # DB connection and embedded SQL migrations (auto-run on startup)

  httputil/
    # JSON response helpers (WriteJSON, WriteError)

  server/
    # http.ServeMux wiring, middleware chain

  testutil/
    # Shared test helpers (test DB setup)
```

## Development

```bash
# From repo root
make db-up          # Start PostgreSQL in Docker
make dev-backend    # Run server (go run ./cmd/server)
```

Migrations run automatically on startup. No manual migration step needed.

## Testing

```bash
cd backend && go test ./...
```

Tests use a real PostgreSQL instance. Set `DATABASE_URL` (or use the default from `.env.example`) before running.

## Configuration

Set via environment variables or a `.env` file in the repo root:

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | `postgres://trello:trello@localhost:5432/trello?sslmode=disable` | PostgreSQL DSN |
| `PORT` | `8080` | Listen port |
| `COOKIE_DOMAIN` | `localhost` | Session cookie domain |
| `BASE_URL` | `http://localhost:5173` | Frontend origin (CORS allowed origin) |
| `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET` | | Google OAuth2 (optional) |
| `MICROSOFT_CLIENT_ID` / `MICROSOFT_CLIENT_SECRET` | | Microsoft OAuth2 (optional) |

## API

See the [root README](../README.md#api-endpoints) for the full endpoint reference.

## Conventions

**Package structure**
- Each domain area has a `Handler` (HTTP) and a `Store` (DB queries) — these are the only two types per package
- Handlers are methods on `Handler`; stores are methods on `Store`. No functions that take a `*sql.DB` directly
- `server.go` is the only place that wires routes — handlers never know about each other

**Handlers**
- Inline anonymous structs for request bodies (`var req struct{ ... }`), not named request types
- Validate input in the handler before calling the store; return early on error
- Ownership is verified in the handler via a store query before any mutation (e.g. `CardOwner`, `ColumnBoardOwner`)
- Responses: `httputil.JSON` for success, `httputil.Error` for errors, `w.WriteHeader(204)` for no-body deletes
- Path values extracted with `r.PathValue("id")` (Go 1.22+)

**Store (DB layer)**
- All public methods take `context.Context` as the first argument
- SQL is written inline as raw string literals — no query builder or ORM
- `RETURNING` clauses used on INSERT/UPDATE to avoid a second SELECT
- `COALESCE($n, column)` used for partial updates (PATCH semantics) instead of building dynamic queries
- Mutations that touch multiple rows use explicit transactions with `defer tx.Rollback()`
- Unexported helpers (`listColumns`, `listCards`) for sub-queries used only within the package

**Auth**
- Sessions are random tokens stored in the DB; the cookie holds the token, not a JWT
- `RequireAuth` middleware reads the cookie, looks up the session, and stores the `*User` in `context`
- Handlers retrieve the user with `auth.UserFromContext(r.Context())`
- Password comparison is constant-time via `bcrypt.CompareHashAndPassword`
