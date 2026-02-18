# Frontend

SvelteKit application built with Svelte 5 runes and TypeScript.

## Structure

```
src/
  lib/
    api.ts              # Typed fetch wrapper for all backend API calls
    types.ts            # TypeScript interfaces (Board, Column, Card, User)
    index.ts            # Barrel export
    components/
      Navbar.svelte     # Top navigation bar
      BoardView.svelte  # Full board layout, manages drag state
      BoardColumn.svelte# Column with card list and add-card form
      BoardCard.svelte  # Individual draggable card
  routes/
    +layout.svelte         # Root layout with Navbar
    +layout.server.ts      # Checks /api/auth/me, redirects if unauthenticated
    +page.svelte           # Dashboard: list and create boards
    login/+page.svelte     # Login form (email/password + OAuth buttons)
    signup/+page.svelte    # Signup form
    boards/[id]/+page.svelte # Kanban board view
```

## Development

```bash
pnpm install
pnpm dev        # Dev server at http://localhost:5173
```

The dev server proxies `/api/*` to the Go backend at `http://localhost:8080`. Start the backend first.

## Testing

```bash
pnpm test       # Run Vitest unit tests
pnpm check      # TypeScript + Svelte type checking
pnpm lint       # ESLint
```

## Building

```bash
pnpm build      # Production build
pnpm preview    # Preview production build locally
```

## Conventions

**Svelte 5 runes**
- Always use `$props()`, `$state`, `$derived`, `$effect` — no legacy `export let` or stores
- To pass server-loaded data into a component without triggering the `state_referenced_locally` warning, hold an override in a nullable `$state` and resolve with `$derived`: `let override = $state<T | null>(null); let value = $derived(override ?? initialValue)`
- Auth state flows from `+layout.server.ts` → `page.data` (via `$app/state`), not from a store

**API calls**
- All backend communication goes through `$lib/api.ts`, which is a thin typed wrapper around `fetch`
- Every function throws `ApiError` (with `.status`) on non-2xx responses — callers catch and handle errors locally
- Cookies are sent automatically (`credentials: 'include'`); no manual token handling in components

**Components**
- Each component owns its own `<style>` block; no global stylesheet
- Components accept callbacks (e.g. `onrefresh`) as props rather than emitting custom events
- State that needs to be refreshed from the server is re-fetched via `api.getBoard()` rather than manually patching local state

**TypeScript**
- All API shapes are defined in `$lib/types.ts` as interfaces — no inline type definitions in components
- Partial updates use pointer-style optional fields (`{ title?: string; description?: string }`) to match the backend's PATCH semantics
