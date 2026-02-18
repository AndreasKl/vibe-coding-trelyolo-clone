# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Start PostgreSQL (required before backend)
make db-up

# Run both servers (each in its own terminal)
make dev-backend     # Go on :8080
make dev-frontend    # SvelteKit on :5173

# Tests
make test            # both
make test-backend    # go test ./... (requires PostgreSQL)
make test-frontend   # vitest --run

# Run a single Go test
cd backend && go test ./internal/board/... -run TestMoveCard

# Run a single frontend test
cd frontend && pnpm test --run -- BoardCard

# Type-check + lint frontend
cd frontend && pnpm check && pnpm lint
```

Config is read from a `.env` file in the **repo root**. Copy `.env.example` to `.env` to get started. OAuth (Google, Microsoft) is optional — leave those vars empty to skip.

## Architecture

Two independent processes. In dev, SvelteKit proxies `/api/*` to the Go backend; in production they share the same origin.

```
Browser
  ├── /api/*  ──► Go backend (:8080) ──► PostgreSQL
  └── rest   ──► SvelteKit (:5173)
```

**Auth flow**: Every navigation triggers `+layout.server.ts`, which calls `GET /api/auth/me`. If it fails, the user is redirected to `/login`. Auth state never lives in the client — no JWT, no localStorage. Sessions are random tokens in a DB table; the cookie holds only the token.

**Data model**: `users → boards → board_columns → cards`. All ordered by integer `position`. Moving a card or column is a single transaction: close the gap at the source, open a gap at the target, update the item.

## Backend Conventions

Each domain (`auth`, `board`) exposes exactly two types: a `Handler` (HTTP) and a `Store` (DB). `server/server.go` is the only file that wires routes — handlers never reference each other.

**Handlers**
- Use inline anonymous structs for request bodies: `var req struct{ Title string \`json:"title"\` }`
- Verify ownership via a store query *before* any mutation (`CardOwner`, `ColumnBoardOwner`)
- Retrieve the authenticated user with `auth.UserFromContext(r.Context())`
- Success: `httputil.JSON(w, status, body)` — No-body deletes: `w.WriteHeader(http.StatusNoContent)`
- Path params: `r.PathValue("id")` (Go 1.22+)

**Store (DB layer)**
- All public methods take `context.Context` as the first argument
- Raw SQL inline — no ORM or query builder
- `RETURNING` on INSERT/UPDATE to avoid a follow-up SELECT
- `COALESCE($n, column)` for partial updates (PATCH semantics)
- Multi-row mutations use explicit transactions with `defer tx.Rollback()`

## Frontend Conventions

All API calls go through `$lib/api.ts`, a thin typed wrapper around `fetch`. It throws `ApiError` (with `.status`) on non-2xx responses — callers catch and handle errors locally. Cookies are sent automatically (`credentials: 'include'`).

All TypeScript shapes are in `$lib/types.ts`. No inline type definitions in components.

**Svelte 5 rules**
- Always use `$props()`, `$state`, `$derived`, `$effect` — never `export let` or legacy stores
- To avoid the `state_referenced_locally` warning when a prop needs a local override:
  ```ts
  let override = $state<T | null>(null);
  let value = $derived(override ?? props.initialValue);
  ```
- Auth state flows from `+layout.server.ts` → `page.data` (via `$app/state`), not a store
- After logout: `goto(url, { invalidateAll: true })` to force the layout server load to re-run
- Components accept callbacks as props (e.g. `onrefresh`) rather than emitting custom events
- Stale state is fixed by re-fetching via `api.getBoard()`, not by patching local state

## API Reference

| Method | Path | Notes |
|--------|------|-------|
| POST | `/api/auth/signup\|login\|logout` | |
| GET | `/api/auth/me` | Used by layout to gate all pages |
| GET | `/api/auth/oauth/{provider}` + `/callback` | Google, Microsoft |
| GET/POST | `/api/boards` | POST auto-creates Todo/Doing/Done columns |
| GET/DELETE | `/api/boards/{id}` | |
| POST | `/api/boards/{boardID}/columns` | |
| PATCH/DELETE | `/api/columns/{id}` | PATCH accepts `{ name?, position? }` |
| POST | `/api/columns/{id}/move` | `{ position }` — transactional reorder |
| POST | `/api/columns/{columnID}/cards` | |
| PATCH/DELETE | `/api/cards/{id}` | |
| POST | `/api/cards/{id}/move` | `{ column_id, position }` — transactional move |

## Testing Notes

- Backend tests hit a real PostgreSQL instance. Shared setup helpers live in `internal/testutil/`.
- Frontend test mocks: any new `api.*` function must be added to both `vi.hoisted` and `vi.mock` in the relevant test file.
- Card action buttons are icon-only with `aria-label="Edit"` / `aria-label="Delete"`; use `getByRole('button', { name: '...' })` in tests.
