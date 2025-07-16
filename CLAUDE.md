# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Kleio is a self-hosted full-stack web application for managing vinyl record collections. It integrates with the Discogs API to automatically sync collections and provides comprehensive tracking of records, plays, equipment, and maintenance.

## Technology Stack

- **Backend**: Go 1.24.1 with SQLite database
- **Frontend**: SolidJS with TypeScript, built with Vite
- **Styling**: SCSS with CSS Modules
- **Database**: SQLite with manual SQL migrations
- **Deployment**: Docker with multi-stage builds

## Common Development Commands

### Backend Development
```bash
# Run backend server (localhost:8080)
make run

# Build backend binary
make build

# Run with Docker
make docker-run

# Database migrations run automatically on startup
```

### Frontend Development
```bash
# Install dependencies
npm install

# Run development server (localhost:3000)
npm run dev

# Build for production
npm run build

# Type checking
npm run typecheck

# Linting
npm run lint
```

### Docker Development
```bash
# Build and run full stack
make docker-run

# Build Docker image
docker build -t kleio .

# Production deployment with volumes
docker run -p 8080:8080 -v ./data:/app/data kleio
```

## Architecture

### Backend Structure (`internal/`)
- **Layered architecture**: Controller → Database → SQLite
- **Key packages**:
  - `api/`: HTTP handlers and routing
  - `database/`: SQLite operations and models
  - `discogs/`: Discogs API integration
  - `migrations/`: SQL migration files (numbered 001-011)

### Frontend Structure (`clio/`)
- **SolidJS reactive architecture** with context providers
- **Key directories**:
  - `src/components/`: Reusable UI components
  - `src/pages/`: Route-based page components
  - `src/lib/`: API client and utilities
  - `src/styles/`: Global SCSS and CSS modules

### Database Schema
- **Core entities**: Albums, Artists, Plays, Styluses, MaintenanceRecords
- **Complex relationships**: Many-to-many for album genres, artists, labels
- **Discogs integration**: Stores Discogs IDs for synchronization

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
- **Styling**: CSS Modules with SCSS preprocessing
- **Type safety**: Full TypeScript with strict checking

### Backend
- **Error handling**: Consistent JSON error responses
- **Database**: Direct SQL with prepared statements
- **Migrations**: Numbered files in `internal/migrations/`
- **Configuration**: Environment variables for Discogs integration

## Testing

The project currently uses manual testing. When adding tests:
- Backend: Use Go's standard testing package
- Frontend: Consider Vitest for SolidJS testing
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