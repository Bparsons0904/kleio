// src/components/charts/PlayFrequencyChart.tsx
import { Component, createEffect, createSignal, onMount, Show } from "solid-js";
import {
  Chart,
  Title,
  Tooltip,
  Legend,
  Colors,
  BarController,
  BarElement,
  CategoryScale,
  LinearScale,
  TimeScale,
} from "chart.js";
import { Line } from "solid-chartjs";
import { useAppContext } from "../../provider/Provider";
import { useDateRange, GroupFrequency } from "../../provider/DateRangeContext";
import ChartControls from "./ChartControls";
import styles from "./PlayFrequencyChart.module.scss";
import { PlayHistory } from "../../types";

// Register Chart.js components
Chart.register(
  Title,
  Tooltip,
  Legend,
  Colors,
  BarController,
  BarElement,
  CategoryScale,
  LinearScale,
  TimeScale,
);

interface PlayHistoryDataPoint {
  date: Date;
  count: number;
}

interface ChartData {
  labels: string[];
  datasets: {
    label: string;
    data: number[];
    backgroundColor?: string;
    borderColor?: string;
    borderWidth?: number;
  }[];
}

interface PlayFrequencyChartProps {
  filter?: string;
}

const PlayFrequencyChart: Component<PlayFrequencyChartProps> = (props) => {
  const filter = () => props.filter || "";

  // Handle filter change
  const handleFilterChange = (newFilter: string) => {
    if (newFilter.startsWith("HEADER:")) return;
  };

  const { playHistory, releases } = useAppContext();
  const { dateRange, groupFrequency } = useDateRange();

  const [chartData, setChartData] = createSignal<ChartData>({
    labels: [],
    datasets: [
      {
        label: "Records Played",
        data: [],
        backgroundColor: "rgba(54, 162, 235, 0.5)",
        borderColor: "rgba(54, 162, 235, 1)",
        borderWidth: 1,
      },
    ],
  });

  const [filterOptions, setFilterOptions] = createSignal<
    { value: string; label: string }[]
  >([]);
  const [isLoading, setIsLoading] = createSignal(true);

  // Chart options
  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: "top" as const,
      },
      title: {
        display: true,
        text: "Records Played Over Time",
      },
    },
    scales: {
      y: {
        beginAtZero: true,
        title: {
          display: true,
          text: "Number of Records",
        },
        ticks: {
          stepSize: 1,
          precision: 0,
        },
      },
      x: {
        title: {
          display: true,
          text: "Date",
        },
      },
    },
  };

  // Generate filter options (artists and genres)
  onMount(() => {
    // Generate artist options
    const artistSet = new Set<string>();
    const genreSet = new Set<string>();

    releases().forEach((release) => {
      // Add main artists (skip producers etc.)
      release.artists
        .filter((a) => a.role !== "Producer")
        .forEach((a) => {
          if (a.artist?.name) {
            artistSet.add(a.artist.name);
          }
        });

      // Add genres
      release.genres.forEach((genre) => {
        genreSet.add(genre.name);
      });
    });

    // Convert to array and sort alphabetically
    const artistOptions = Array.from(artistSet)
      .map((name) => ({
        value: `artist:${name}`,
        label: name,
      }))
      .sort((a, b) => a.label.localeCompare(b.label));

    const genreOptions = Array.from(genreSet)
      .map((name) => ({
        value: `genre:${name}`,
        label: name,
      }))
      .sort((a, b) => a.label.localeCompare(b.label));

    // Combine options with separators
    setFilterOptions([
      { value: "", label: "All Records" },
      { value: "HEADER:ARTISTS", label: "--- Artists ---" },
      ...artistOptions,
      { value: "HEADER:GENRES", label: "--- Genres ---" },
      ...genreOptions,
    ]);

    setIsLoading(false);
  });

  // Process data when dependencies change
  createEffect(() => {
    const range = dateRange();
    const frequency = groupFrequency();
    const currentFilter = filter();

    if (isLoading() || !playHistory().length) return;

    // Filter play history data by date range
    let filteredHistory = playHistory().filter((item) => {
      const playDate = new Date(item.playedAt);
      return playDate >= range.start && playDate <= range.end;
    });

    // Apply additional filters (artist or genre)
    if (currentFilter && !currentFilter.startsWith("HEADER:")) {
      const [type, value] = currentFilter.split(":");

      if (type === "artist") {
        filteredHistory = filteredHistory.filter((item) => {
          return item.release.artists.some(
            (a) => a.artist?.name === value && a.role !== "Producer",
          );
        });
      } else if (type === "genre") {
        filteredHistory = filteredHistory.filter((item) => {
          return item.release.genres.some((g) => g.name === value);
        });
      }
    }

    // Group by frequency
    const groupedData = groupDataByFrequency(filteredHistory, frequency);

    // Update chart data
    setChartData({
      labels: groupedData.map((d) => formatDateLabel(d.date, frequency)),
      datasets: [
        {
          label: "Records Played",
          data: groupedData.map((d) => d.count),
          backgroundColor: "rgba(54, 162, 235, 0.5)",
          borderColor: "rgba(54, 162, 235, 1)",
          borderWidth: 1,
        },
      ],
    });
  });

  // Group play history data by the selected frequency
  const groupDataByFrequency = (
    history: PlayHistory[],
    frequency: GroupFrequency,
  ): PlayHistoryDataPoint[] => {
    // Create a map to store grouped data
    const groupedMap = new Map<string, number>();

    // Sort history by date (oldest first) to ensure chronological order
    const sortedHistory = [...history].sort(
      (a, b) => new Date(a.playedAt).getTime() - new Date(b.playedAt).getTime(),
    );

    // Generate all dates in the range to ensure no gaps
    const allDates = generateDateRange(
      dateRange().start,
      dateRange().end,
      frequency,
    );

    // Initialize all dates with 0 count
    allDates.forEach((date) => {
      const key = getDateKey(date, frequency);
      groupedMap.set(key, 0);
    });

    // Count plays for each date group
    sortedHistory.forEach((play) => {
      const playDate = new Date(play.playedAt);
      const key = getDateKey(playDate, frequency);

      groupedMap.set(key, (groupedMap.get(key) || 0) + 1);
    });

    // Convert map to array of data points
    return Array.from(groupedMap)
      .map(([key, count]) => ({
        date: parseDateKey(key, frequency),
        count,
      }))
      .sort((a, b) => a.date.getTime() - b.date.getTime());
  };

  // Get a consistent string key for a date based on frequency
  const getDateKey = (date: Date, frequency: GroupFrequency): string => {
    const year = date.getFullYear();
    const month = date.getMonth();
    const day = date.getDate();

    switch (frequency) {
      case "monthly":
        return `${year}-${month + 1}`;
      case "weekly":
        // Get the Monday of the week
        const d = new Date(date);
        const day1 = d.getDate() - d.getDay() + (d.getDay() === 0 ? -6 : 1);
        d.setDate(day1);
        return `${d.getFullYear()}-${d.getMonth() + 1}-${d.getDate()}`;
      case "daily":
      default:
        return `${year}-${month + 1}-${day}`;
    }
  };

  // Parse a date key back to a Date object
  const parseDateKey = (key: string, frequency: GroupFrequency): Date => {
    const parts = key.split("-").map(Number);

    switch (frequency) {
      case "monthly":
        return new Date(parts[0], parts[1] - 1, 1);
      case "weekly":
      case "daily":
      default:
        return new Date(parts[0], parts[1] - 1, parts[2] || 1);
    }
  };

  // Generate all dates in a range based on frequency
  const generateDateRange = (
    start: Date,
    end: Date,
    frequency: GroupFrequency,
  ): Date[] => {
    const dates: Date[] = [];
    const current = new Date(start);

    while (current <= end) {
      dates.push(new Date(current));

      switch (frequency) {
        case "monthly":
          current.setMonth(current.getMonth() + 1);
          break;
        case "weekly":
          current.setDate(current.getDate() + 7);
          break;
        case "daily":
        default:
          current.setDate(current.getDate() + 1);
          break;
      }
    }

    return dates;
  };

  // Format date for display on chart
  const formatDateLabel = (date: Date, frequency: GroupFrequency): string => {
    const options: Intl.DateTimeFormatOptions = {};

    switch (frequency) {
      case "monthly":
        options.year = "numeric";
        options.month = "short";
        break;
      case "weekly":
        // Start of week format
        return `Week of ${date.toLocaleDateString(undefined, {
          month: "short",
          day: "numeric",
        })}`;
      case "daily":
      default:
        options.month = "short";
        options.day = "numeric";
        // Only show year if it's not the current year
        if (date.getFullYear() !== new Date().getFullYear()) {
          options.year = "numeric";
        }
        break;
    }

    return date.toLocaleDateString(undefined, options);
  };

  return (
    <div class={styles.chartWrapper}>
      <Show
        when={!isLoading()}
        fallback={<div class={styles.loading}>Loading chart data...</div>}
      >
        <Line
          data={chartData()}
          options={chartOptions}
          width={800}
          height={400}
        />
      </Show>
    </div>
  );
};

export default PlayFrequencyChart;
