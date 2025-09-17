# Repository Guidelines

## Project Structure & Module Organization
Backend Go services live at the repository root: `main.go` wires Echo routes, `monitor.go` schedules checks, `notifications.go` wraps Pushover, and persistence helpers sit in `storage.go`. Frontend React code resides in `src/` with `components/`, `context/`, and `services/` modules; compiled assets land in `dist/`. Automation scripts live in `scripts/`, and container assets are `Dockerfile` plus `docker-compose.yml`. Keep environment samples in `env.example`; never commit real secrets.

## Build, Test, and Development Commands
- `./scripts/dev.sh` (or `scripts\\dev.bat` on Windows) boots Go API and Vite dev server together.
- `go run .` serves the API; pair with `npm run dev` when debugging frontend only.
- `npm run build` compiles the React UI into `dist/`; `go build -o gjallarhorn .` produces a standalone binary.
- `./scripts/build.sh` bundles backend and embedded frontend; `docker-compose up -d` runs the full stack with Docker.

## Coding Style & Naming Conventions
Run `go fmt ./...` before committing; prefer idiomatic Go naming (CamelCase types, mixedCase functions). Frontend code follows the shared ESLint config: two-space indentation, single quotes, no trailing semicolons, PascalCase components in `components/` and camelCase utilities. Tailwind classes stay inline; extract shared styles into `src/index.css` tokens. Use descriptive file names (e.g., `ServiceHealthCard.jsx`).

## Testing Guidelines
Add backend tests alongside Go sources using `_test.go` files and run `go test ./...`; target coverage for new logic and include table-driven cases for monitor workflows. Frontend tests are currently absent—introduce Vitest/React Testing Library under `src/__tests__/` when contributing UI logic, and document manual smoke steps in the PR until automated coverage exists.

## Commit & Pull Request Guidelines
Follow the lightweight history pattern: concise, capitalized commit subject in the imperative mood (≤72 chars) with optional wrapped body describing motivation and effects; reference issues using `Refs #123` when applicable. For PRs, provide a problem statement, summary of changes, testing evidence (`go test ./...`, `npm run build`), and screenshots for UI tweaks. Request review from a backend and frontend peer when impacts cross the boundary.

## Configuration & Security Notes
Copy `env.example` to `.env` for local secrets; never commit `.env`. Pushover tokens grant production access—scope them to test applications and rotate after leaks. Docker deployments mount data at `/data`; ensure volumes are access-controlled when sharing hosts.
