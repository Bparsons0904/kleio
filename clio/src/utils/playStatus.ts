/**
 * Calculate the cleanliness score for a record
 * @param lastCleanedDate The date the record was last cleaned
 * @param playsSinceCleaning Number of plays since last cleaning
 * @returns Cleanliness score (0-100)
 */
export function getCleanlinessScore(
  lastCleanedDate: Date | null,
  playsSinceCleaning: number,
): number {
  // If never cleaned, return 100 (needs cleaning)
  if (!lastCleanedDate) return 100;

  // Time-based calculation (percentage of 6 months elapsed)
  const sixMonthsInMs = 6 * 30 * 24 * 60 * 60 * 1000;
  const timeElapsed = Date.now() - lastCleanedDate.getTime();
  const timeScore = Math.min(100, (timeElapsed / sixMonthsInMs) * 100);

  // Play-based calculation (percentage of 5 plays)
  const playScore = Math.min(100, (playsSinceCleaning / 5) * 100);

  // Return the higher score (worse case)
  return Math.max(timeScore, playScore);
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
  if (daysElapsed <= 60) return 60; // Moderately recent: Last 60 days
  if (daysElapsed <= 90) return 40; // Not recent: Last 90 days
  return 20; // Long ago: Over 90 days
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
  if (daysElapsed <= 60) return "Played in the last 2 months";
  if (daysElapsed <= 90) return "Played in the last 3 months";
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

  return playHistory.filter((play) => {
    const playDate = new Date(play.playedAt);
    return playDate.getTime() > lastCleanedDate.getTime();
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
