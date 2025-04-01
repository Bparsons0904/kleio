import { createContext, createSignal, ParentProps, useContext } from "solid-js";

export type DateRange = {
  start: Date;
  end: Date;
};

export type TimeFrame = "7d" | "30d" | "90d" | "1y" | "all" | "custom";
export type GroupFrequency = "daily" | "weekly" | "monthly";

type DateRangeContextType = {
  dateRange: () => DateRange;
  timeFrame: () => TimeFrame;
  groupFrequency: () => GroupFrequency;
  setTimeFrame: (timeFrame: TimeFrame) => void;
  setGroupFrequency: (frequency: GroupFrequency) => void;
  setCustomDateRange: (range: DateRange) => void;
};

const DateRangeContext = createContext<DateRangeContextType>();

export function DateRangeProvider(props: ParentProps) {
  const [dateRange, setDateRange] = createSignal<DateRange>(
    getDateRangeForTimeFrame("30d"),
  );
  const [timeFrame, setTimeFrame] = createSignal<TimeFrame>("30d");
  const [groupFrequency, setGroupFrequency] =
    createSignal<GroupFrequency>("daily");

  // Update date range when time frame changes
  const updateTimeFrame = (newTimeFrame: TimeFrame) => {
    setTimeFrame(newTimeFrame);
    setDateRange(getDateRangeForTimeFrame(newTimeFrame));
  };

  // Allow setting a custom date range
  const setCustomDateRange = (range: DateRange) => {
    setTimeFrame("custom" as TimeFrame);
    setDateRange(range);
  };

  const store = {
    dateRange,
    timeFrame,
    groupFrequency,
    setTimeFrame: updateTimeFrame,
    setGroupFrequency,
    setCustomDateRange,
  };

  return (
    <DateRangeContext.Provider value={store}>
      {props.children}
    </DateRangeContext.Provider>
  );
}

export function useDateRange() {
  const context = useContext(DateRangeContext);
  if (!context) {
    throw new Error("useDateRange must be used within a DateRangeProvider");
  }
  return context;
}

// Helper function to calculate date range based on time frame
export function getDateRangeForTimeFrame(timeFrame: TimeFrame): DateRange {
  const end = new Date();
  const start = new Date();

  switch (timeFrame) {
    case "7d":
      start.setDate(end.getDate() - 7);
      break;
    case "30d":
      start.setDate(end.getDate() - 30);
      break;
    case "90d":
      start.setDate(end.getDate() - 90);
      break;
    case "1y":
      start.setFullYear(end.getFullYear() - 1);
      break;
    case "all":
      // Go back 10 years as a default for "all time"
      start.setFullYear(end.getFullYear() - 10);
      break;
  }

  // Set time to beginning/end of day
  start.setHours(0, 0, 0, 0);
  end.setHours(23, 59, 59, 999);

  return { start, end };
}
