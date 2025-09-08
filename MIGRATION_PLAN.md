# Kleio Production Migration Plan

## Overview
This document outlines the migration strategy for converting Kleio from a local SQLite-based application to a production-ready system with PostgreSQL, Fiber framework, Zitadel authentication, and enhanced architecture.

## Current Architecture Analysis

### Current State
- **Backend**: Go 1.24.1 with standard HTTP server
- **Database**: SQLite with manual migrations (001-011)
- **Architecture**: Simple controller→database layered pattern
- **Frontend**: SolidJS (clio/) served separately
- **Authentication**: Basic Discogs token storage
- **API Integration**: Server-side Discogs API calls

### Target Architecture
- **Backend**: Fiber web framework with clean architecture
- **Database**: PostgreSQL with connection pooling
- **Architecture**: Repository pattern with dependency injection
- **Frontend**: Embedded SolidJS served by Go application
- **Authentication**: Zitadel OIDC with JWT validation
- **API Integration**: Client-side Discogs calls with server coordination
- **Real-time**: WebSocket support for live updates

## Migration Strategy: **Incremental Conversion**

**Recommendation**: Convert the existing codebase rather than starting fresh to preserve domain logic, business rules, and existing database schema.

## Implementation Phases

### Phase 1: Infrastructure Foundation

#### 1. Database Migration to PostgreSQL
- **Current**: SQLite with 11 migration files
- **Target**: PostgreSQL with equivalent schema
- **Tasks**:
  - Convert SQLite-specific syntax to PostgreSQL
  - Update data types (TEXT → VARCHAR, INTEGER → SERIAL, etc.)
  - Implement connection pooling with pgx or similar
  - Update database.go connection logic
  - Test migration scripts against PostgreSQL instance

#### 2. Framework Migration to Fiber
- **Current**: Standard net/http with manual routing
- **Target**: Fiber framework with middleware pipeline
- **Tasks**:
  - Replace http.Server with fiber.New()
  - Convert all handler functions to fiber.Ctx
  - Migrate routing from manual setup to Fiber routing
  - Add Fiber middleware: CORS, compression, helmet, logging
  - Update response/request handling patterns

#### 3. Project Structure Reorganization
- **Current**: Simple internal/ structure
- **Target**: Clean architecture following vim/main pattern
- **New Structure**:
  ```
  cmd/
    api/main.go
    migration/
  config/
    config.go
  internal/
    app/app.go           # Dependency injection container
    controllers/         # HTTP handlers
    repositories/        # Data access layer
    services/           # Business logic layer
    models/             # Data models
    handlers/           # Route definitions
    websockets/         # WebSocket management
    events/             # Event bus system
    logger/             # Structured logging
    middleware/         # Authentication & validation
  ```

### Phase 2: Authentication & Security

#### 4. Zitadel Integration
- **Current**: Simple token table in SQLite
- **Target**: Full OIDC authentication with Zitadel
- **Tasks**:
  - Set up Zitadel instance configuration
  - Implement OIDC client in Go
  - Create JWT validation middleware
  - Replace auth table with proper user management
  - Update frontend to handle OIDC flow
  - Implement user session management

### Phase 3: Real-time Features

#### 5. WebSocket Implementation
- **Current**: HTTP-only communication
- **Target**: WebSocket support for real-time updates
- **Tasks**:
  - Implement WebSocket hub manager (based on vim/main pattern)
  - Create connection management system
  - Add real-time collection sync notifications
  - Implement client connection handling
  - Add event broadcasting for play logging, sync status

### Phase 4: API Architecture Redesign

#### 6. Client-Side Discogs API Strategy
- **Current**: Server makes all Discogs API calls
- **Target**: Client-side API calls with server coordination
- **Rationale**: Discogs rate limits per user+IP, so client calls avoid server bottleneck
- **Tasks**:
  - Create API proxy endpoints that return instructions to client
  - Implement client-side rate limiting coordination
  - Add backend validation of client API responses
  - Maintain server-side fallback for critical operations
  - Update frontend to handle direct Discogs communication

### Phase 5: Frontend Integration

#### 7. Frontend Serving Strategy
- **Current**: Separate SolidJS dev server
- **Target**: Embedded frontend served by Go application
- **Options**:
  - **Option A**: Embed using Go embed directive (recommended)
  - **Option B**: Keep separate with nginx proxy
- **Tasks**:
  - Build SolidJS app for production
  - Embed static files in Go binary using `//go:embed`
  - Add static file serving through Fiber
  - Implement proper asset caching headers
  - Update build process and Docker configuration

### Phase 6: Production Readiness

#### 8. Observability & Reliability
- **Current**: Basic logging
- **Target**: Comprehensive observability
- **Tasks**:
  - Implement structured logging (based on vim/main logger pattern)
  - Add health check endpoints (`/health`, `/ready`)
  - Implement metrics collection and exposure
  - Add graceful shutdown handling
  - Create proper configuration management with environment variables

#### 9. Deployment & Scaling
- **Current**: Single Docker container with SQLite
- **Target**: Production-ready deployment
- **Tasks**:
  - Update Docker configuration for new architecture
  - Add environment-specific configurations
  - Implement database connection management
  - Add container health checks
  - Create production deployment documentation

## Migration Sequence

### Recommended Order
1. **Start with Phase 1**: Infrastructure changes provide foundation
2. **Parallel Phase 2 & 3**: Auth and WebSockets can be developed concurrently
3. **Phase 4**: API redesign after core architecture is solid
4. **Phase 5 & 6**: Frontend integration and production readiness

### Key Checkpoints
- [ ] Database migration completed and tested
- [ ] Fiber framework integrated with all existing endpoints
- [ ] Project structure reorganized with dependency injection
- [ ] Zitadel authentication working
- [ ] WebSocket functionality implemented
- [ ] Client-side Discogs API strategy working
- [ ] Frontend embedded and serving correctly
- [ ] Production observability in place
- [ ] Deployment pipeline updated

## Risk Mitigation

### Preserve Business Logic
- Keep all existing controller logic during migration
- Maintain existing API contracts initially
- Test each phase thoroughly before proceeding

### Database Migration Safety
- Create PostgreSQL migration scripts from existing SQLite schema
- Test migrations on copy of production data
- Maintain backward compatibility during transition

### Client-Side API Risks
- Implement server-side fallback for Discogs API calls
- Add proper error handling for client-side failures
- Monitor rate limiting effectiveness

## Benefits of Incremental Approach

1. **Lower Risk**: Migrate piece-by-piece while maintaining functionality
2. **Preserve Investment**: Reuse existing business logic and database schema
3. **Continuous Operation**: Application remains functional during migration
4. **Learning Opportunity**: Understand production requirements during migration
5. **Flexibility**: Can adjust approach based on learnings from each phase

## File Structure Reference

Based on vim/main project structure, the final organization should follow:
- Clear separation of concerns
- Dependency injection through app container
- Repository pattern for data access
- Service layer for business logic
- Proper middleware and handler organization