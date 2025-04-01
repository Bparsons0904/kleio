import { Component, createSignal, createEffect, Show } from "solid-js";
import {
  TimeFrame,
  GroupFrequency,
  useDateRange,
} from "../../provider/DateRangeContext";
import styles from "./ChartControls.module.scss";

interface ChartControlsProps {
  showFrequencyControls?: boolean;
  showFilters?: boolean;
  onFilterChange?: (filter: string) => void;
  filterOptions?: { value: string; label: string }[];
  filterLabel?: string;
}

const ChartControls: Component<ChartControlsProps> = (props) => {
  const {
    timeFrame,
    setTimeFrame,
    groupFrequency,
    setGroupFrequency,
    dateRange,
    setCustomDateRange,
  } = useDateRange();

  const [customStartDate, setCustomStartDate] = createSignal("");
  const [customEndDate, setCustomEndDate] = createSignal("");
  const [showCustomDate, setShowCustomDate] = createSignal(false);
  const [selectedFilter, setSelectedFilter] = createSignal("");

  // Format date to YYYY-MM-DD for input type date
  const formatDateForInput = (date: Date) => {
    return date.toISOString().split("T")[0];
  };

  // Initialize custom date fields when the component mounts
  createEffect(() => {
    setCustomStartDate(formatDateForInput(dateRange().start));
    setCustomEndDate(formatDateForInput(dateRange().end));
  });

  const handleTimeFrameChange = (event: Event) => {
    const target = event.target as HTMLSelectElement;
    const value = target.value as TimeFrame;

    if (value === "custom") {
      setShowCustomDate(true);
    } else {
      setShowCustomDate(false);
      setTimeFrame(value);
    }
  };

  const handleFrequencyChange = (event: Event) => {
    const target = event.target as HTMLSelectElement;
    setGroupFrequency(target.value as GroupFrequency);
  };

  const handleCustomDateSubmit = () => {
    const start = new Date(customStartDate());
    const end = new Date(customEndDate());

    // Ensure end is at the end of the day
    end.setHours(23, 59, 59, 999);

    // Only update if valid dates
    if (!isNaN(start.getTime()) && !isNaN(end.getTime()) && start <= end) {
      setCustomDateRange({ start, end });
    }
  };

  const handleFilterChange = (event: Event) => {
    const target = event.target as HTMLSelectElement;
    setSelectedFilter(target.value);
    if (props.onFilterChange) {
      props.onFilterChange(target.value);
    }
  };

  return (
    <div class={styles.controlsContainer}>
      <div class={styles.controlGroup}>
        <label class={styles.label}>Time Period:</label>
        <select
          class={styles.select}
          value={timeFrame()}
          onChange={handleTimeFrameChange}
        >
          <option value="7d">Last 7 Days</option>
          <option value="30d">Last 30 Days</option>
          <option value="90d">Last 90 Days</option>
          <option value="1y">Last Year</option>
          <option value="all">All Time</option>
          <option value="custom">Custom Date Range</option>
        </select>
      </div>

      <Show when={showCustomDate()}>
        <div class={styles.customDateContainer}>
          <div class={styles.dateInputGroup}>
            <label class={styles.label}>Start:</label>
            <input
              type="date"
              class={styles.dateInput}
              value={customStartDate()}
              onChange={(e) => setCustomStartDate(e.target.value)}
            />
          </div>
          <div class={styles.dateInputGroup}>
            <label class={styles.label}>End:</label>
            <input
              type="date"
              class={styles.dateInput}
              value={customEndDate()}
              onChange={(e) => setCustomEndDate(e.target.value)}
            />
          </div>
          <button class={styles.applyButton} onClick={handleCustomDateSubmit}>
            Apply
          </button>
        </div>
      </Show>

      <Show when={props.showFrequencyControls}>
        <div class={styles.controlGroup}>
          <label class={styles.label}>Group By:</label>
          <select
            class={styles.select}
            value={groupFrequency()}
            onChange={handleFrequencyChange}
          >
            <option value="daily">Daily</option>
            <option value="weekly">Weekly</option>
            <option value="monthly">Monthly</option>
          </select>
        </div>
      </Show>

      <Show
        when={
          props.showFilters &&
          props.filterOptions &&
          props.filterOptions.length > 0
        }
      >
        <div class={styles.controlGroup}>
          <label class={styles.label}>
            {props.filterLabel || "Filter By:"}
          </label>
          <select
            class={styles.select}
            value={selectedFilter()}
            onChange={handleFilterChange}
          >
            <option value="">All</option>
            {props.filterOptions?.map((option) => (
              <option value={option.value}>{option.label}</option>
            ))}
          </select>
        </div>
      </Show>
    </div>
  );
};

export default ChartControls;
