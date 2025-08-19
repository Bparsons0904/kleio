# Archive System Implementation Plan

## Problem Statement

Currently, Kleio has two critical issues with collection management:

1. **Sync Issue**: When releases are removed from your Discogs collection, they persist in the local database indefinitely
2. **Manual Management**: No way to manually archive/remove releases added by mistake (like mis-clicks)
3. **Data Loss Risk**: Any deletion approach risks losing valuable play history data

## Solution: Global Archive System

### Core Concept

Implement a **soft archive system** that preserves all data while allowing releases to be removed from the active collection. This solves both automatic sync cleanup and manual management needs.

### Key Features

- **Soft Archive**: Mark releases as archived instead of hard deletion
- **Play History Preservation**: All listening data remains intact for archived releases
- **Dual Triggers**: Both automatic (sync-based) and manual (user-initiated) archiving
- **Restore Capability**: Allow un-archiving of archived releases
- **Archive Reasons**: Track why releases were archived for better management

## Database Schema Changes

### New Columns for `releases` Table

```sql
-- Migration: Add archive support to releases table
ALTER TABLE releases ADD COLUMN archived BOOLEAN DEFAULT FALSE;
ALTER TABLE releases ADD COLUMN archived_at TIMESTAMP NULL;
ALTER TABLE releases ADD COLUMN archive_reason TEXT NULL;
ALTER TABLE releases ADD COLUMN sync_session_id TEXT NULL;
```

### Archive Reasons Enum

```go
type ArchiveReason string

const (
    ArchiveReasonSyncRemoved  ArchiveReason = "sync_removed"   // Auto-archived during sync
    ArchiveReasonUserAction   ArchiveReason = "user_action"    // Manual user archive
    ArchiveReasonDuplicate    ArchiveReason = "duplicate"      // Duplicate detection
    ArchiveReasonMistake      ArchiveReason = "mistake"        // User mistake/mis-click
)
```

## Implementation Plan

### Phase 1: Database Foundation

**Files to Create/Modify:**
- `internal/database/migrations/012_archive_system.sql`
- `internal/database/models.go` (update Release struct)

**Changes:**
1. Create migration to add archive columns
2. Update Release model with new fields
3. Add archive-related database methods

### Phase 2: Sync Enhancement (Two-Phase Sync)

**Files to Modify:**
- `internal/controller/collection.controller.go`
- `internal/database/release.database.go`

**Sync Flow:**
```
1. Start Sync → Generate unique sync_session_id
2. Mark Phase → Set all active releases with current sync_session_id
3. Update Phase → Update sync_session_id for releases found in Discogs
4. Cleanup Phase → Archive releases with old sync_session_id
```

**Pseudo-code:**
```go
func (c *Controller) SyncReleases() error {
    sessionID := generateSyncSessionID()
    
    // Mark all active releases with current session
    err := c.DB.MarkReleasesForSync(sessionID)
    if err != nil {
        return err
    }
    
    // Fetch and update releases from Discogs
    err = c.fetchAndUpdateReleases(sessionID)
    if err != nil {
        return err
    }
    
    // Archive releases not found in this sync
    err = c.DB.ArchiveStalereleases(sessionID, ArchiveReasonSyncRemoved)
    if err != nil {
        return err
    }
    
    return nil
}
```

### Phase 3: Manual Archive API

**Files to Create/Modify:**
- `internal/server/routes.go`
- `internal/server/release.server.go` (new endpoints)
- `internal/controller/release.controller.go`

**New API Endpoints:**
```
POST   /api/releases/{id}/archive    - Archive a release
POST   /api/releases/{id}/restore    - Restore archived release
GET    /api/releases/archived        - List archived releases
DELETE /api/releases/{id}            - Hard delete (admin only)
```

### Phase 4: Database Methods

**New Methods in `release.database.go`:**

```go
// Archive operations
func (db *Database) ArchiveRelease(releaseID int, reason ArchiveReason) error
func (db *Database) RestoreRelease(releaseID int) error
func (db *Database) GetArchivedReleases() ([]Release, error)

// Sync operations  
func (db *Database) MarkReleasesForSync(sessionID string) error
func (db *Database) UpdateReleaseWithSyncSession(releaseID int, sessionID string) error
func (db *Database) ArchiveStaleReleases(sessionID string, reason ArchiveReason) error

// Query modifications
func (db *Database) GetAllReleases() ([]Release, error) // Only active releases
func (db *Database) GetAllReleasesWithArchived() ([]Release, error) // Include archived
```

### Phase 5: Frontend Integration

**Files to Modify:**
- `clio/src/components/` (archive management components)
- `clio/src/pages/` (archived releases page)
- `clio/src/lib/api.ts` (new API calls)

**UI Features:**
- Archive button on release cards
- "Archived Releases" page
- Restore functionality
- Archive reason display
- Filter toggle (active/archived/all)

## Data Flow Examples

### Automatic Archive (Sync)

```
1. User removes album from Discogs collection
2. Next sync runs:
   - Session ID: "sync_2025_08_19_123456"
   - Mark all active releases with session ID
   - Fetch from Discogs (removed album not returned)
   - Archive releases still marked with old session ID
3. Album marked: archived=true, archive_reason="sync_removed"
4. Play history preserved, album hidden from active collection
```

### Manual Archive (User Action)

```
1. User clicks "Archive" on mistakenly added release
2. API call: POST /api/releases/1234/archive
3. Release marked: archived=true, archive_reason="user_action"
4. Release removed from active collection view
5. Play history preserved for analytics
```

## Benefits

1. **Data Preservation**: No play history lost, ever
2. **Clean Collection**: Only owned releases visible in main collection
3. **Mistake Recovery**: Easy to restore accidentally archived releases
4. **Audit Trail**: Track why releases were archived
5. **Analytics Intact**: Archived releases still contribute to listening statistics
6. **Flexible Management**: Both automatic and manual control

## Migration Strategy

1. **Backward Compatible**: New columns default to active state
2. **Gradual Rollout**: Can implement archive without affecting existing functionality
3. **Safe Testing**: Archive operations are reversible
4. **Data Integrity**: Foreign key relationships preserved

## Future Enhancements

- **Bulk Operations**: Archive multiple releases at once
- **Smart Suggestions**: Detect potential duplicates for archiving
- **Archive Analytics**: Reports on archived releases and reasons
- **Export Options**: Include/exclude archived releases in data exports
- **Archive Scheduling**: Auto-archive releases not played in X months

## Risk Mitigation

- **Reversible Operations**: All archives can be restored
- **Data Validation**: Verify play history preservation
- **Backup Strategy**: Ensure database backups before major archive operations
- **User Education**: Clear UI indicators for archive vs delete actions

This archive system provides a robust foundation for collection management while preserving the valuable play history data that makes Kleio valuable for music analytics.