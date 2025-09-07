# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Kleio is a self-hosted full-stack web application for managing vinyl record collections. It integrates with the Discogs API to automatically sync collections and provides comprehensive tracking of records, plays, equipment, and maintenance.

## Technology Stack

- **Backend**: Go 1.25 with Fiber web framework, SQLite database
- **Frontend**: SolidJS with TypeScript, Vite build system
- **Styling**: SCSS with CSS Modules, comprehensive design system
- **Database**: SQLite with numbered migration files (001-011)
- **Development Tools**: Air for live reload, Docker for containerization
- **Monitoring**: Prometheus metrics integration

## Common Development Commands

### Backend Development
```bash
# Run backend server (localhost:8080) 
make run

# Build backend binary
make build

# Live reload development (requires air)
make watch

# Run tests
make test

# Clean built binary
make clean

# Database migrations run automatically on startup
```

### Frontend Development (in clio/ directory)
```bash
# Install dependencies
npm install

# Run development server (localhost:3000)
npm run dev
# Alternative: npm start

# Build for production
npm run build

# Preview production build
npm run serve

# NOTE: No separate typecheck/lint commands in package.json - handled by editor/IDE
```

### Docker Development
```bash
# Build and run full stack with Docker Compose
make docker-run

# Stop Docker containers
make docker-down

# Build Docker image manually
docker build -t kleio .

# Production deployment with volumes
docker run -p 8080:8080 -v ./data:/app/data kleio
```

### Development Workflow
```bash
# Terminal 1: Backend with live reload
make watch

# Terminal 2: Frontend development server
cd clio && npm run dev

# Production build process
cd clio && npm run build && cd .. && make build
```

## Architecture

### Backend Structure (`internal/`)
- **Clean architecture**: Controller → Database → SQLite
- **Key packages**:
  - `controller/`: Fiber HTTP handlers (auth, collection, playHistory, etc.)
  - `database/`: SQLite operations and models
  - `server/`: Fiber app configuration and routing setup
  - `database/migrations/`: SQL migration files (numbered 001-011, auto-applied on startup)

### Frontend Structure (`clio/`)
- **SolidJS reactive architecture** with context providers
- **Key directories**:
  - `src/components/`: Reusable UI components (modals, dropdowns, panels)
  - `src/pages/`: Route-based page components
  - `src/lib/`: API client and utilities
  - `src/styles/`: Global SCSS and CSS modules
- **Component Architecture**: Each component has own directory with `.tsx` and `.module.scss` files
- **Styling System**: Comprehensive design system defined in `STYLE_GUIDELINES.md`

### Database Schema
- **Core entities**: Albums, Artists, Plays, Styluses, MaintenanceRecords
- **Complex relationships**: Many-to-many for album genres, artists, labels
- **Discogs integration**: Stores Discogs IDs for synchronization
- **Migration system**: Numbered files applied automatically on startup

### Build System
- **Backend**: Standard Go build with CGO enabled for SQLite
- **Frontend**: Vite build system with SolidJS plugin
- **Development**: Air for backend live reload, Vite dev server for frontend
- **Production**: Multi-stage Docker build serving both frontend and backend from single Go binary

## Key Features

1. **Discogs Integration**: Automatic collection sync with token-based authentication
2. **Play Tracking**: Log plays with stylus selection, duration, and notes
3. **Equipment Management**: Track stylus wear and lifespan calculations
4. **Analytics**: Charts for play frequency, duration, and genre analysis
5. **Collection Management**: Browse, search, and organize vinyl records
6. **Maintenance Tracking**: Record cleaning history and schedules

## API Design

RESTful endpoints with proper HTTP methods:
- `GET /api/albums` - List albums with pagination
- `POST /api/plays` - Log new play session
- `PUT /api/styluses/{id}` - Update stylus information
- `DELETE /api/albums/{id}` - Archive album from collection

Authentication uses Discogs personal access tokens.

## Development Patterns

### Frontend
- **Reactive data**: Use SolidJS resources for API calls
- **State management**: Context providers for global state
- **Styling**: CSS Modules with SCSS preprocessing, follow `STYLE_GUIDELINES.md`
- **Type safety**: Full TypeScript with JSX preserve mode
- **Component structure**: Each component in own directory with `.tsx` and `.module.scss`
- **Dependencies**: Axios for HTTP client, Chart.js for analytics, Fuse.js for search

### Backend
- **Web framework**: Fiber (Express-like framework for Go)
- **Error handling**: Consistent JSON error responses
- **Database**: Direct SQL with prepared statements, sqlite3 driver
- **Migrations**: Numbered files in `internal/database/migrations/`
- **Configuration**: Environment variables for Discogs integration
- **Architecture**: Controller pattern with clean separation of concerns

### Development Environment
- **Live reload**: Air configured with `.air.toml` (excludes frontend dirs)
- **Port separation**: Backend on :8080, Frontend dev server on :3000
- **Database**: SQLite file in root directory (`sqlite.db`)
- **Assets**: Stored in `assets/` directory

## Testing

The project currently uses manual testing. When adding tests:
- Backend: Use Go's standard testing package (`make test`)
- Frontend: Consider Vitest for SolidJS testing (not currently configured)
- Integration: Test API endpoints with real database

## Deployment

Production deployment uses Docker with:
- Multi-stage build for optimized image size
- Volume mounting for SQLite database persistence
- Environment variables for configuration
- Single container serving both frontend and backend

## Discogs Integration

- Requires personal access token from Discogs
- Syncs collection data including release information
- Handles rate limiting and API pagination
- Stores Discogs IDs for data consistency

## File Structure Summary

```
kleio/
├── cmd/api/main.go           # Application entry point
├── internal/                 # Private application code
│   ├── controller/           # HTTP request handlers
│   ├── database/             # Database models and operations
│   │   └── migrations/       # SQL migration files (001-011)
│   └── server/               # Fiber server setup
├── clio/                     # Frontend SolidJS application
│   ├── src/
│   │   ├── components/       # Reusable UI components
│   │   ├── pages/            # Route-based page components
│   │   ├── lib/              # API client and utilities
│   │   └── styles/           # Global SCSS and CSS modules
│   ├── package.json          # Frontend dependencies
│   └── tsconfig.json         # TypeScript configuration
├── assets/                   # Static assets (images, etc.)
├── Makefile                  # Build and development commands
├── go.mod                    # Go module dependencies
├── .air.toml                 # Live reload configuration
├── docker-compose.yml        # Docker development setup
└── Dockerfile                # Production Docker build
```

## Important Notes

- **Monorepo structure**: Backend and frontend in same repository
- **Production serving**: Single Go binary serves both API and static files
- **Database location**: SQLite file in project root (`sqlite.db`)
- **Migration system**: Automatic application on startup, numbered files
- **Docker ports**: Production uses 38080, development uses 8080 for backend