# Deep Cleaning Feature Implementation Plan

## Overview
Add functionality to track and display "deep cleaning" as a special type of cleaning with a checkbox indicator and visual distinction.

## Database Changes

### 1. Create Migration (012_add_deep_clean_flag.sql)
- Add `is_deep_clean BOOLEAN DEFAULT 0` column to `cleaning_history` table
- Create query to retroactively mark existing records where notes contain "deep" (case insensitive)
- Migration will search notes for variations like "Deep Clean", "deep clean", "DEEP", etc.

## Backend Changes

### 2. Update Models (models.go)
- Add `IsDeepClean bool` field to `CleaningHistory` struct with proper JSON/DB tags

### 3. Update Database Layer (cleaningHistory.database.go)
- Update `CreateCleaningHistory` to include `is_deep_clean` in INSERT statement
- Update `UpdateCleaningHistory` to include `is_deep_clean` in UPDATE statement
- Update all SELECT statements to include `is_deep_clean` field
- Add `Scan()` calls for the new field in all query methods

### 4. Update Server/Controller (no changes needed)
- The JSON decode/encode will automatically handle the new field

## Frontend Changes

### 5. Update TypeScript Types (types/index.ts)
- Add `isDeepClean?: boolean` to `CleaningHistory` interface

### 6. Update Cleaning Form (RecordActionModal.tsx)
- Add "Deep Clean" checkbox above the notes field
- Update `handleLogCleaning()` and `handleLogBoth()` to include `isDeepClean` in API calls
- Add state management for the deep clean checkbox

### 7. Update Display Component (RecordHistoryItem.tsx)
- Add deep clean indicator (e.g., "Deep Cleaned" text or special icon) when `isDeepClean` is true
- Update styling to distinguish deep cleaning entries

### 8. Update Edit Panel (EditHistoryPanel.tsx)
- Add deep clean checkbox to the editing interface
- Update edit form submission to include `isDeepClean` field

## Implementation Details

**Migration Logic:**
```sql
ALTER TABLE cleaning_history ADD COLUMN is_deep_clean BOOLEAN DEFAULT 0;
UPDATE cleaning_history SET is_deep_clean = 1 WHERE LOWER(notes) LIKE '%deep%';
```

**Frontend UI:**
- Checkbox labeled "Deep Clean" in cleaning form
- Visual indicator (icon or badge) next to "Cleaned" text for deep cleanings
- Maintain existing styling but add visual distinction for deep cleans

**Backward Compatibility:**
- Existing cleaning records default to `is_deep_clean = false`
- Retroactive update based on notes content containing "deep" (case insensitive)
- No breaking changes to existing API or UI functionality

This approach maintains all existing functionality while cleanly adding the deep cleaning feature with both form input and visual feedback.