export function useFormattedMediumDate(date: string) {
  if (!date) return "Never synced";

  try {
    const dateObj = new Date(date);

    if (isNaN(dateObj.getTime())) {
      return "Invalid date";
    }

    return dateObj.toLocaleString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
      hour: "numeric",
      minute: "numeric",
      hour12: true,
    });
  } catch (error) {
    console.error("Error formatting date:", error);
    return "Error formatting date";
  }
}

export function useFormattedShortDate(date: string) {
  if (!date) return "";
  try {
    const dateObj = new Date(date);

    if (isNaN(dateObj.getTime())) {
      return "Invalid date";
    }

    return dateObj.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  } catch (error) {
    console.error("Error formatting date:", error);
    return "Error formatting date";
  }
}

export const formatDateForInput = (date: Date): string => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
};
