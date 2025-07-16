import { Component, createSignal, createEffect, Show } from "solid-js";
import {
  TimeFrame,
  GroupFrequency,
  useDateRange,
} from "../../provider/DateRangeContext";
import { formatDateForInput } from "../../utils/dates";
import styles from "./ChartControls.module.scss";
import SearchableDropdown, {
  DropdownOption,
} from "../SearchableDropdown/SearchableDropdown";

interface ChartControlsProps {
  showFrequencyControls?: boolean;
  showFilters?: boolean;
  onFilterChange?: (filter: string) => void;
  filterOptions?: DropdownOption[];
  filterLabel?: string;
  filterValue?: string;
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
  const [internalSelectedFilter, setInternalSelectedFilter] = createSignal("");


  // Initialize custom date fields when the component mounts
  createEffect(() => {
    setCustomStartDate(formatDateForInput(dateRange().start));
    setCustomEndDate(formatDateForInput(dateRange().end));
  });

  const handleTimeFrameChange = (value: string) => {
    if (value === "custom") {
      setShowCustomDate(true);
    } else {
      setShowCustomDate(false);
      setTimeFrame(value as TimeFrame);
    }
  };

  const handleFrequencyChange = (value: string) => {
    setGroupFrequency(value as GroupFrequency);
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

  const handleFilterChange = (value: string) => {
    setInternalSelectedFilter(value);
    if (props.onFilterChange) {
      props.onFilterChange(value);
    }
  };

  const selectedFilter = () =>
    props.filterValue !== undefined
      ? props.filterValue
      : internalSelectedFilter();

  // Time frame options for dropdown
  const timeFrameOptions = [
    { value: "7d", label: "Last 7 Days" },
    { value: "30d", label: "Last 30 Days" },
    { value: "90d", label: "Last 90 Days" },
    { value: "1y", label: "Last Year" },
    { value: "all", label: "All Time" },
    { value: "custom", label: "Custom Date Range" },
  ];

  // Frequency options for dropdown
  const frequencyOptions = [
    { value: "daily", label: "Daily" },
    { value: "weekly", label: "Weekly" },
    { value: "monthly", label: "Monthly" },
  ];

  return (
    <div class={styles.controlsContainer}>
      <div class={styles.controlGroup}>
        <SearchableDropdown
          label="Time Period:"
          options={timeFrameOptions}
          value={timeFrame()}
          onChange={handleTimeFrameChange}
          placeholder="Select time period"
        />
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
          <SearchableDropdown
            label="Group By:"
            options={frequencyOptions}
            value={groupFrequency()}
            onChange={handleFrequencyChange}
            placeholder="Select grouping"
          />
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
          <SearchableDropdown
            label={props.filterLabel || "Filter By:"}
            options={props.filterOptions}
            value={selectedFilter()}
            onChange={handleFilterChange}
            placeholder="Select filter"
            isSearchable
          />
        </div>
      </Show>
    </div>
  );
};

export default ChartControls;
