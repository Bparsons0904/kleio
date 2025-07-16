/**
 * Safely converts a date string or Date object to a Date instance in local timezone
 */
function parseLocalDate(date: string | Date | null | undefined): Date | null {
  if (!date) return null;
  
  try {
    const dateObj = date instanceof Date ? date : new Date(date);
    if (isNaN(dateObj.getTime())) {
      return null;
    }
    return dateObj;
  } catch (error) {
    console.error("Error parsing date:", error);
    return null;
  }
}

/**
 * Formats a date for display with time (e.g., "Dec 25, 2023, 3:45 PM")
 * Always uses local timezone and user's locale
 */
export function useFormattedMediumDate(date: string | Date | null | undefined): string {
  if (!date) return "Never synced";

  const dateObj = parseLocalDate(date);
  if (!dateObj) return "Invalid date";

  return dateObj.toLocaleString(undefined, {
    month: "short",
    day: "numeric",
    year: "numeric",
    hour: "numeric",
    minute: "numeric",
    hour12: true,
  });
}

/**
 * Formats a date for display without time (e.g., "Dec 25, 2023")
 * Always uses local timezone and user's locale
 */
export function useFormattedShortDate(date: string | Date | null | undefined): string {
  if (!date) return "";
  
  const dateObj = parseLocalDate(date);
  if (!dateObj) return "Invalid date";

  return dateObj.toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

/**
 * Formats a date for HTML date inputs (YYYY-MM-DD format)
 * Always uses local timezone to prevent date shifting
 */
export function formatDateForInput(date: Date | string | null | undefined): string {
  const dateObj = parseLocalDate(date);
  if (!dateObj) return "";

  const year = dateObj.getFullYear();
  const month = String(dateObj.getMonth() + 1).padStart(2, "0");
  const day = String(dateObj.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
}

/**
 * Formats a date for display with fallback text for null dates
 * Uses local timezone and user's locale
 */
export function formatLocalDate(date: Date | string | null | undefined, fallback: string = "Never"): string {
  if (!date) return fallback;
  
  const dateObj = parseLocalDate(date);
  if (!dateObj) return "Invalid date";

  return dateObj.toLocaleDateString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

/**
 * Creates a date group key for local date grouping (YYYY-MM-DD in local timezone)
 * Prevents issues with UTC date shifting
 */
export function getLocalDateGroupKey(date: string | Date): string {
  const dateObj = parseLocalDate(date);
  if (!dateObj) return "";

  const year = dateObj.getFullYear();
  const month = String(dateObj.getMonth() + 1).padStart(2, "0");
  const day = String(dateObj.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
}

/**
 * Checks if two dates are on the same day in local timezone
 */
export function isSameLocalDay(date1: string | Date | null, date2: string | Date | null): boolean {
  if (!date1 || !date2) return false;
  
  const dateObj1 = parseLocalDate(date1);
  const dateObj2 = parseLocalDate(date2);
  
  if (!dateObj1 || !dateObj2) return false;
  
  return getLocalDateGroupKey(dateObj1) === getLocalDateGroupKey(dateObj2);
}
