export function useFormattedMediumDate(date: string) {
  if (!date) return "Never synced";

  try {
    const dateObj = new Date(date);

    if (isNaN(dateObj.getTime())) {
      console.log("Invalid date created from:", date);
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
  // });
}
