// import { createMemo } from "solid-js";

export function useFormattedMediumDate(dateSignal) {
  // const date = dateSignal();

  // Handle loading state
  // if (!date || date === "Never synced") return "Never synced";

  try {
    const dateObj = new Date(dateSignal);

    // Check if we got a valid date
    if (isNaN(dateObj.getTime())) {
      console.log("Invalid date created from:", dateSignal);
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
