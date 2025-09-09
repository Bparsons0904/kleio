# Step 1: Multi-User Database Architecture Migration

## Overview

This document details the first step in migrating Kleio from SQLite to PostgreSQL 18 with multi-user support. The key architectural changes include shared data models, user-specific collections via join tables, GORM integration, and UUID7 primary keys.

## Architecture Decision: Shared vs User-Specific Data

### Shared Data (No user ownership)
- **Releases, Artists, Labels, Genres, Styles, Tracks**: Same vinyl data serves all users
- **Single source of truth**: One release record can be owned by multiple users
- **Discogs integration**: Release data synced once, shared across users

### User-Specific Data  
- **Collection ownership**: `user_releases` join table for who owns what releases
- **Play history, Styluses, Cleaning history**: Personal user data and activities
- **User preferences**: Personal ratings, notes, custom folder organization

## What We Need to Create

### 1. GORM Models with UUID7 PKs (`internal/models/models.go`)

#### Shared Models
```go
type Release struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
    DiscogsID   int       `gorm:"uniqueIndex;not null" json:"discogsId"`
    Title       string    `gorm:"not null" json:"title"`
    Year        *int      `json:"year"`
    ResourceURL string    `json:"resourceUrl"`
    Thumb       string    `json:"thumb"`
    CoverImage  string    `json:"coverImage"`
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
    
    // Relationships
    Artists     []Artist  `gorm:"many2many:release_artists" json:"artists,omitempty"`
    Labels      []Label   `gorm:"many2many:release_labels" json:"labels,omitempty"`
    Genres      []Genre   `gorm:"many2many:release_genres" json:"genres,omitempty"`
    Styles      []Style   `gorm:"many2many:release_styles" json:"styles,omitempty"`
    Tracks      []Track   `gorm:"foreignKey:ReleaseID" json:"tracks,omitempty"`
}

type Artist struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
    DiscogsID   int       `gorm:"uniqueIndex;not null" json:"discogsId"`
    Name        string    `gorm:"not null" json:"name"`
    ResourceURL string    `json:"resourceUrl"`
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type Label struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
    DiscogsID   int       `gorm:"uniqueIndex;not null" json:"discogsId"`
    Name        string    `gorm:"not null" json:"name"`
    ResourceURL string    `json:"resourceUrl"`
    EntityType  string    `json:"entityType"`
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type Genre struct {
    ID   uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
    Name string    `gorm:"uniqueIndex;not null" json:"name"`
}

type Style struct {
    ID   uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
    Name string    `gorm:"uniqueIndex;not null" json:"name"`
}

type Track struct {
    ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
    ReleaseID       uuid.UUID `gorm:"type:uuid;not null;index" json:"releaseId"`
    Position        string    `gorm:"not null" json:"position"` // e.g., "A1", "B2"
    Title           string    `gorm:"not null" json:"title"`
    DurationText    string    `json:"durationText"`    // Original format (e.g., "3:45")
    DurationSeconds int       `json:"durationSeconds"` // Normalized duration
    CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
    UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
    
    Release         Release   `gorm:"foreignKey:ReleaseID" json:"release,omitempty"`
}
```

#### User Models
```go
type User struct {
    ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
    Username        string    `gorm:"uniqueIndex;not null" json:"username"`
    Email           string    `gorm:"uniqueIndex;not null" json:"email"`
    DiscogsUsername string    `json:"discogsUsername"`
    CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
    UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
    
    // Relationships
    AuthTokens      []AuthToken    `gorm:"foreignKey:UserID" json:"-"`
    UserReleases    []UserRelease  `gorm:"foreignKey:UserID" json:"-"`
    PlayHistory     []PlayHistory  `gorm:"foreignKey:UserID" json:"-"`
    Styluses        []Stylus       `gorm:"foreignKey:UserID" json:"-"`
}

type AuthToken struct {
    Token        string    `gorm:"primaryKey" json:"token"`
    UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
    DiscogsToken string    `json:"discogsToken"` // For Discogs API
    CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
    ExpiresAt    *time.Time `json:"expiresAt"`
    
    User         User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
```

#### Collection Ownership Join Table
```go
type UserRelease struct {
    UserID     uuid.UUID `gorm:"type:uuid;primary_key" json:"userId"`
    ReleaseID  uuid.UUID `gorm:"type:uuid;primary_key" json:"releaseId"`  
    FolderID   int       `json:"folderId"`   // User's Discogs folder
    Rating     int       `json:"rating"`     // User's personal rating
    Notes      string    `json:"notes"`      // User's personal notes
    AddedAt    time.Time `gorm:"autoCreateTime" json:"addedAt"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
    
    // Relationships
    User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Release    Release   `gorm:"foreignKey:ReleaseID" json:"release,omitempty"`
}
```

#### User-Specific Activity Models
```go
type PlayHistory struct {
    ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
    UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"userId"`
    ReleaseID uuid.UUID  `gorm:"type:uuid;not null;index" json:"releaseId"`
    StylusID  *uuid.UUID `gorm:"type:uuid" json:"stylusId"`
    PlayedAt  time.Time  `gorm:"not null" json:"playedAt"`
    Notes     string     `json:"notes"`
    CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
    UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
    
    // Relationships
    User      User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Release   Release    `gorm:"foreignKey:ReleaseID" json:"release,omitempty"`
    Stylus    *Stylus    `gorm:"foreignKey:StylusID" json:"stylus,omitempty"`
}

type Stylus struct {
    ID               uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
    UserID           uuid.UUID  `gorm:"type:uuid;not null;index" json:"userId"`
    Name             string     `gorm:"not null" json:"name"`
    Manufacturer     string     `json:"manufacturer"`
    ExpectedLifespan int        `json:"expectedLifespan"`
    PurchaseDate     *time.Time `json:"purchaseDate"`
    Active           bool       `gorm:"default:false" json:"active"`
    Primary          bool       `gorm:"default:false" json:"primary"`
    Owned            bool       `gorm:"default:false" json:"owned"`
    BaseModel        bool       `gorm:"default:false" json:"baseModel"`
    CreatedAt        time.Time  `gorm:"autoCreateTime" json:"createdAt"`
    UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
    
    User             User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type CleaningHistory struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid7()" json:"id"`
    UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
    ReleaseID uuid.UUID `gorm:"type:uuid;not null;index" json:"releaseId"`
    CleanedAt time.Time `gorm:"not null" json:"cleanedAt"`
    Notes     string    `json:"notes"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
    
    // Relationships
    User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Release   Release   `gorm:"foreignKey:ReleaseID" json:"release,omitempty"`
}
```

### 2. GORM Database Connection (`internal/database/gorm.go`)
```go
package database

import (
    "fmt"
    "log/slog"
    "os"
    "time"
    
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    
    "kleio/internal/models"
)

var DB *gorm.DB

func InitializeGORM() error {
    dsn := buildPostgreSQLDSN()
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
    }
    
    // Configure connection pool
    sqlDB, err := db.DB()
    if err != nil {
        return fmt.Errorf("failed to get database instance: %w", err)
    }
    
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    // Enable UUID extension for gen_random_uuid7()
    if err := db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto").Error; err != nil {
        return fmt.Errorf("failed to create pgcrypto extension: %w", err)
    }
    
    // Auto-migrate schemas
    if err := autoMigrate(db); err != nil {
        return fmt.Errorf("failed to auto-migrate schemas: %w", err)
    }
    
    DB = db
    slog.Info("Successfully connected to PostgreSQL with GORM")
    return nil
}

func buildPostgreSQLDSN() string {
    host := getEnvOrDefault("DB_HOST", "localhost")
    port := getEnvOrDefault("DB_PORT", "5432")
    user := getEnvOrDefault("DB_USER", "kleio")
    password := getEnvOrDefault("DB_PASSWORD", "")
    dbname := getEnvOrDefault("DB_NAME", "kleio")
    sslmode := getEnvOrDefault("DB_SSL_MODE", "disable")
    
    return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        host, port, user, password, dbname, sslmode)
}

func autoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.User{},
        &models.AuthToken{},
        &models.Release{},
        &models.Artist{},
        &models.Label{},
        &models.Genre{},
        &models.Style{},
        &models.Track{},
        &models.UserRelease{},
        &models.PlayHistory{},
        &models.Stylus{},
        &models.CleaningHistory{},
    )
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

### 3. Repository Layer (`internal/repository/`)

#### User Repository (`internal/repository/user_repository.go`)
```go
package repository

import (
    "kleio/internal/models"
    "gorm.io/gorm"
    "github.com/google/uuid"
)

type UserRepository struct {
    DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
    return r.DB.Create(user).Error
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
    var user models.User
    err := r.DB.First(&user, "id = ?", id).Error
    return &user, err
}

func (r *UserRepository) GetUserByToken(token string) (*models.User, error) {
    var authToken models.AuthToken
    err := r.DB.Preload("User").First(&authToken, "token = ?", token).Error
    if err != nil {
        return nil, err
    }
    return &authToken.User, nil
}

func (r *UserRepository) CreateAuthToken(token *models.AuthToken) error {
    return r.DB.Create(token).Error
}
```

#### Collection Repository (`internal/repository/collection_repository.go`)
```go
package repository

import (
    "kleio/internal/models"
    "gorm.io/gorm"
    "github.com/google/uuid"
)

type CollectionRepository struct {
    DB *gorm.DB
}

func NewCollectionRepository(db *gorm.DB) *CollectionRepository {
    return &CollectionRepository{DB: db}
}

func (r *CollectionRepository) GetUserCollection(userID uuid.UUID, offset, limit int) ([]models.UserRelease, error) {
    var userReleases []models.UserRelease
    err := r.DB.Preload("Release").
        Preload("Release.Artists").
        Preload("Release.Labels").
        Preload("Release.Genres").
        Preload("Release.Styles").
        Where("user_id = ?", userID).
        Offset(offset).
        Limit(limit).
        Find(&userReleases).Error
    
    return userReleases, err
}

func (r *CollectionRepository) AddReleaseToCollection(userRelease *models.UserRelease) error {
    return r.DB.Create(userRelease).Error
}

func (r *CollectionRepository) RemoveReleaseFromCollection(userID, releaseID uuid.UUID) error {
    return r.DB.Where("user_id = ? AND release_id = ?", userID, releaseID).
        Delete(&models.UserRelease{}).Error
}

func (r *CollectionRepository) UpdateUserRelease(userRelease *models.UserRelease) error {
    return r.DB.Save(userRelease).Error
}
```

### 4. Updated Go Dependencies (`go.mod`)
Add these dependencies:
```go
require (
    gorm.io/gorm v1.25.5
    gorm.io/driver/postgres v1.5.4
    github.com/google/uuid v1.6.0
    // ... existing dependencies
)
```

## What We Need to Update

### 1. Database Initialization (`internal/database/database.go`)
Replace SQLite initialization with GORM:
```go
func New() Database {
    if dbInstance != nil {
        return *dbInstance
    }
    
    // Use GORM for PostgreSQL, fallback to SQLite for development
    if os.Getenv("DB_TYPE") == "postgres" {
        if err := InitializeGORM(); err != nil {
            log.Fatalf("Failed to initialize PostgreSQL with GORM: %v", err)
        }
        // Create wrapper for compatibility
        dbInstance = &Database{DB: nil} // GORM instance accessed via global DB variable
    } else {
        // Existing SQLite initialization for development/migration
        // ... keep existing SQLite code
    }
    
    return *dbInstance
}
```

### 2. All Controllers (`internal/controller/*.go`)  
- Replace direct SQL queries with GORM repository calls
- Add user context extraction from auth tokens  
- Update collection endpoints to use UserRelease join table
- Add user filtering to all personal data queries

Example updates:
```go
// Before (direct SQL)
func GetReleases(c *fiber.Ctx) error {
    rows, err := db.Query("SELECT * FROM releases WHERE folder_id = ?", folderID)
    // ... manual row scanning
}

// After (GORM with user context)
func GetUserCollection(c *fiber.Ctx) error {
    userID := extractUserIDFromContext(c)
    collectionRepo := repository.NewCollectionRepository(database.DB)
    
    userReleases, err := collectionRepo.GetUserCollection(userID, offset, limit)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(userReleases)
}
```

## Data Migration Strategy (Phase 1.5)

### Current Data → Multi-User Structure
1. **Create your user record**: Your existing data becomes the first user account
2. **Releases become shared**: Import releases once, multiple users can own the same release
3. **Create UserRelease records**: Your collection ownership via join table with your personal ratings/notes
4. **User-specific data**: PlayHistory, Styluses, CleaningHistory linked to your user ID
5. **Preserve relationships**: All foreign keys maintained with UUID mapping table

### Migration Script Structure
```go
type DataMigration struct {
    SQLiteDB     *sql.DB
    PostgresDB   *gorm.DB
    UUIDMapping  map[int]uuid.UUID // Maps old integer IDs to new UUIDs
}

func (m *DataMigration) MigrateData() error {
    // 1. Create first user (you)
    user := &models.User{
        Username: "your-username",
        Email:    "your-email",
        // ... other fields
    }
    
    // 2. Migrate shared data (releases, artists, labels, etc.)
    // 3. Create UserRelease records for your collection
    // 4. Migrate your personal data (play history, styluses, etc.)
    // 5. Validate data integrity
}
```

## Key Benefits

### Data Efficiency
- **No duplicate release data**: Same vinyl release shared across all users
- **Normalized relationships**: Artists, labels, genres shared efficiently
- **Storage optimization**: Only user-specific data duplicated per user

### GORM Advantages
- **Auto-migrations**: Schema changes handled automatically
- **Type safety**: Compile-time query validation
- **Relationship management**: Automatic foreign key handling and joins
- **Hooks and callbacks**: Before/after save logic for business rules
- **Connection pooling**: Built-in PostgreSQL optimizations

### PostgreSQL 18 Features
- **UUID7 primary keys**: Time-ordered, better for indexing and pagination
- **Built-in UUID generation**: Native `gen_random_uuid7()` function
- **Advanced indexing**: Partial indexes for user-specific queries
- **JSON support**: JSONB columns for flexible metadata storage
- **Performance**: Superior concurrent access and query optimization

### Multi-User Scalability
- **User isolation**: Complete separation of personal data
- **Shared knowledge base**: Efficient release data management
- **Flexible permissions**: Easy to add user roles and access controls
- **API consistency**: Same endpoints serve user-specific data based on auth context

This architecture provides a solid foundation for scaling Kleio to multiple users while preserving all existing functionality and data.