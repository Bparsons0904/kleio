/**
 * Calculate the cleanliness score for a record based ONLY on plays since cleaning
 * @param lastCleanedDate The date the record was last cleaned
 * @param playsSinceCleaning Number of plays since last cleaning
 * @returns Cleanliness score (0-100) based purely on play count
 */
export function getCleanlinessScore(
  lastCleanedDate: Date | null,
  playsSinceCleaning: number,
): number {
  // If never cleaned, assume the disc is dirty and needs cleaning
  if (!lastCleanedDate) {
    return 100; // Maximum "needs cleaning" score
  }

  // Play-based calculation (percentage of 5 plays)
  // Modified to ensure that exactly 4 plays is less than 80%
  const playScore = Math.min(100, (playsSinceCleaning / 5.01) * 100);

  // Return only the play-based score (independent of time)
  return playScore;
}

/**
 * Get color based on cleanliness score
 * @param score Cleanliness score (0-100)
 * @returns Hex color code
 */
export function getCleanlinessColor(score: number): string {
  if (score < 20) return "#35a173"; // Dark green
  if (score < 40) return "#59c48c"; // Medium green
  if (score < 60) return "#80d6aa"; // Light green
  if (score < 80) return "#f59e0b"; // Amber/yellow warning color
  return "#e9493e"; // Red for danger
}

/**
 * Calculate play recency score
 * @param lastPlayedDate The date the record was last played
 * @returns A number between 0-100 representing how recently played
 */
export function getPlayRecencyScore(lastPlayedDate: Date | null): number {
  // If never played, return 0
  if (!lastPlayedDate) return 0;

  const now = new Date();
  const daysElapsed =
    (now.getTime() - lastPlayedDate.getTime()) / (24 * 60 * 60 * 1000);

  // Convert days to score (0-100)
  if (daysElapsed <= 7) return 100; // Very recent: Last 7 days
  if (daysElapsed <= 30) return 80; // Recent: Last 30 days
  if (daysElapsed <= 90) return 60; // Moderately recent: Last 90 days
  if (daysElapsed <= 180) return 40; // Not recent: Last 180 days (6 months)
  if (daysElapsed <= 365) return 20; // Long ago: Last 365 days (1 year)
  return 0; // Very long ago: Over 1 year
}

/**
 * Get play recency color
 * @param score Play recency score (0-100)
 * @returns Hex color code
 */
export function getPlayRecencyColor(score: number): string {
  // Mid-way between original bright colors and muted colors
  if (score >= 80) return "#35a173"; // Green (between #2f855a and #3cb371)
  if (score >= 60) return "#59c48c"; // Light green (between #48bb78 and #66cdaa)
  if (score >= 40) return "#80d6aa"; // Yellow (between #ecc94b and #f0e68c)
  if (score >= 20) return "#f59e0b"; // Orange (between #ed8936 and #ffa07a)
  return "#e9493e"; // Red (between #e53e3e and #e9967a)
}

/**
 * Get text description for recency
 * @param lastPlayedDate The date the record was last played
 * @returns Text description of recency
 */
export function getPlayRecencyText(lastPlayedDate: Date | null): string {
  if (!lastPlayedDate) return "Never played";

  const now = new Date();
  const daysElapsed =
    (now.getTime() - lastPlayedDate.getTime()) / (24 * 60 * 60 * 1000);

  if (daysElapsed <= 7) return "Played this week";
  if (daysElapsed <= 30) return "Played this month";
  if (daysElapsed <= 90) return "Played in the last 3 months";
  if (daysElapsed <= 180) return "Played in the last 6 months";
  if (daysElapsed <= 365) return "Played in the last year";
  return "Not played recently";
}

/**
 * Get text description for cleanliness
 * @param score Cleanliness score (0-100)
 * @returns Text description of cleanliness
 */
export function getCleanlinessText(score: number): string {
  if (score < 20) return "Recently cleaned";
  if (score < 40) return "Clean";
  if (score < 60) return "May need cleaning soon";
  if (score < 80) return "Due for cleaning";
  return "Needs cleaning";
}

/**
 * Count plays since last cleaning
 * @param playHistory Array of play history items
 * @param lastCleanedDate Date of last cleaning
 * @returns Number of plays since last cleaning
 */
export function countPlaysSinceCleaning(
  playHistory: { playedAt: string }[],
  lastCleanedDate: Date | null,
): number {
  if (!lastCleanedDate) return playHistory.length; // If never cleaned, count all plays

  // Get the exact timestamp of the last cleaning
  const lastCleanedTime = lastCleanedDate.getTime();

  return playHistory.filter((play) => {
    const playDate = new Date(play.playedAt);
    const playTime = playDate.getTime();

    // Add a small epsilon (1 millisecond) to the cleaning time
    // This ensures that plays logged at exactly the same time as cleaning
    // are not counted toward the "plays since cleaning" total
    return playTime > lastCleanedTime + 1;
  }).length;
}

/**
 * Get the last cleaning date from cleaning history
 * @param cleaningHistory Array of cleaning history items
 * @returns Last cleaning date or null if never cleaned
 */
export function getLastCleaningDate(
  cleaningHistory: { cleanedAt: string }[] | undefined,
): Date | null {
  if (!cleaningHistory || cleaningHistory.length === 0) return null;

  // Sort by date descending
  const sortedHistory = [...cleaningHistory].sort((a, b) => {
    return new Date(b.cleanedAt).getTime() - new Date(a.cleanedAt).getTime();
  });

  return new Date(sortedHistory[0].cleanedAt);
}

/**
 * Get the last play date from play history
 * @param playHistory Array of play history items
 * @returns Last play date or null if never played
 */
export function getLastPlayDate(
  playHistory: { playedAt: string }[] | undefined,
): Date | null {
  if (!playHistory || playHistory.length === 0) return null;

  // Sort by date descending
  const sortedHistory = [...playHistory].sort((a, b) => {
    return new Date(b.playedAt).getTime() - new Date(a.playedAt).getTime();
  });

  return new Date(sortedHistory[0].playedAt);
}
