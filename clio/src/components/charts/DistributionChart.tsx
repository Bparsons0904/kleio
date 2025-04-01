// src/components/charts/DistributionCharts.tsx
import { Component, createEffect, createSignal, onMount, Show } from "solid-js";
import {
  Chart,
  Title,
  Tooltip,
  Legend,
  Colors,
  ArcElement,
  PieController,
  DoughnutController,
} from "chart.js";
import { Pie, Doughnut } from "solid-chartjs";
import { useAppContext } from "../../provider/Provider";
import { useDateRange } from "../../provider/DateRangeContext";
import ChartControls from "./ChartControls";
import { PlayHistory } from "../../types";
import styles from "./DistributionChart.module.scss";

// Register Chart.js components
Chart.register(
  Title,
  Tooltip,
  Legend,
  Colors,
  ArcElement,
  PieController,
  DoughnutController,
);

type DistributionType = "artist" | "genre";

interface DistributionDataItem {
  label: string;
  count: number;
  duration: number;
}

interface ChartData {
  labels: string[];
  datasets: {
    label: string;
    data: number[];
    backgroundColor: string[];
    borderColor?: string[];
    borderWidth?: number;
  }[];
}

const DistributionCharts: Component = () => {
  const { playHistory } = useAppContext();
  const { dateRange } = useDateRange();

  const [distributionType, setDistributionType] =
    createSignal<DistributionType>("genre");
  const [countChartData, setCountChartData] = createSignal<ChartData>({
    labels: [],
    datasets: [
      {
        label: "Play Count",
        data: [],
        backgroundColor: [],
      },
    ],
  });

  const [durationChartData, setDurationChartData] = createSignal<ChartData>({
    labels: [],
    datasets: [
      {
        label: "Play Duration",
        data: [],
        backgroundColor: [],
      },
    ],
  });

  const [showTopCount, setShowTopCount] = createSignal(10);
  const [isLoading, setIsLoading] = createSignal(true);

  // Chart options for count chart
  const countChartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: "right" as const,
      },
      title: {
        display: true,
        text: "Play Count Distribution",
      },
      tooltip: {
        callbacks: {
          label: function (context) {
            const label = context.label || "";
            const value = (context.raw as number) || 0;
            const total = context.dataset.data.reduce(
              (a, b) => a + (b as number),
              0,
            );
            const percentage = Math.round((value / total) * 100);
            return `${label}: ${value} plays (${percentage}%)`;
          },
        },
      },
    },
  };

  // Chart options for duration chart
  const durationChartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: "right" as const,
      },
      title: {
        display: true,
        text: "Listening Time Distribution",
      },
      tooltip: {
        callbacks: {
          label: function (context) {
            const label = context.label || "";
            const value = (context.raw as number) || 0;
            const total = context.dataset.data.reduce(
              (a, b) => a + (b as number),
              0,
            );
            const percentage = Math.round((value / total) * 100);

            // Format minutes into hours and minutes
            let timeLabel = `${value} minutes`;
            if (value >= 60) {
              const hours = Math.floor(value / 60);
              const minutes = Math.round(value % 60);
              timeLabel = `${hours}h ${minutes}m`;
            }

            return `${label}: ${timeLabel} (${percentage}%)`;
          },
        },
      },
    },
  };

  // Define reusable chart colors (will cycle through these)
  const chartColors = [
    "rgba(255, 99, 132, 0.8)", // Red
    "rgba(54, 162, 235, 0.8)", // Blue
    "rgba(255, 206, 86, 0.8)", // Yellow
    "rgba(75, 192, 192, 0.8)", // Teal
    "rgba(153, 102, 255, 0.8)", // Purple
    "rgba(255, 159, 64, 0.8)", // Orange
    "rgba(199, 199, 199, 0.8)", // Gray
    "rgba(83, 102, 255, 0.8)", // Indigo
    "rgba(78, 205, 196, 0.8)", // Turquoise
    "rgba(255, 99, 71, 0.8)", // Tomato
    "rgba(144, 238, 144, 0.8)", // Light green
    "rgba(255, 182, 193, 0.8)", // Light pink
  ];

  // Get additional colors by rotating hue
  const getExtendedColors = (count: number): string[] => {
    if (count <= chartColors.length) {
      return chartColors.slice(0, count);
    }

    // Need more colors, generate them
    const colors = [...chartColors];
    for (let i = chartColors.length; i < count; i++) {
      const hue = (i * 137.508) % 360; // Golden angle approximation for good distribution
      colors.push(`hsla(${hue}, 70%, 60%, 0.8)`);
    }

    return colors;
  };

  // Initialize on mount
  onMount(() => {
    setIsLoading(false);
  });

  // Update charts when dependencies change
  createEffect(() => {
    const range = dateRange();
    const type = distributionType();
    const topCount = showTopCount();

    if (isLoading() || !playHistory().length) return;

    // Filter play history data by date range
    const filteredHistory = playHistory().filter((item) => {
      const playDate = new Date(item.playedAt);
      return playDate >= range.start && playDate <= range.end;
    });

    if (filteredHistory.length === 0) {
      // No data in the selected range
      resetCharts();
      return;
    }

    // Calculate distribution based on selected type
    const distribution = calculateDistribution(filteredHistory, type);

    // Sort by count descending and take top N
    const sortedData = [...distribution]
      .sort((a, b) => b.count - a.count)
      .slice(0, topCount);

    // Sort by duration descending and take top N
    const sortedDataByDuration = [...distribution]
      .sort((a, b) => b.duration - a.duration)
      .slice(0, topCount);

    // Generate colors
    const countColors = getExtendedColors(sortedData.length);
    const durationColors = getExtendedColors(sortedDataByDuration.length);

    // Update count chart data
    setCountChartData({
      labels: sortedData.map((d) => d.label),
      datasets: [
        {
          label: "Play Count",
          data: sortedData.map((d) => d.count),
          backgroundColor: countColors,
        },
      ],
    });

    // Update duration chart data
    setDurationChartData({
      labels: sortedDataByDuration.map((d) => d.label),
      datasets: [
        {
          label: "Listening Time",
          data: sortedDataByDuration.map((d) => d.duration),
          backgroundColor: durationColors,
        },
      ],
    });
  });

  // Reset charts to empty state
  const resetCharts = () => {
    setCountChartData({
      labels: [],
      datasets: [
        {
          label: "Play Count",
          data: [],
          backgroundColor: [],
        },
      ],
    });

    setDurationChartData({
      labels: [],
      datasets: [
        {
          label: "Listening Time",
          data: [],
          backgroundColor: [],
        },
      ],
    });
  };

  // Calculate distribution data
  const calculateDistribution = (
    history: PlayHistory[],
    type: DistributionType,
  ): DistributionDataItem[] => {
    // Map to store distribution data
    const distMap = new Map<string, DistributionDataItem>();

    // Process each play history entry
    history.forEach((play) => {
      // Determine the labels based on type
      const labels = getLabelsForPlay(play, type);

      // Get play duration in minutes
      const durationMinutes = getPlayDurationMinutes(play);

      // Update the distribution map for each label
      labels.forEach((label) => {
        if (distMap.has(label)) {
          const item = distMap.get(label)!;
          item.count += 1;
          item.duration += durationMinutes;
        } else {
          distMap.set(label, {
            label,
            count: 1,
            duration: durationMinutes,
          });
        }
      });
    });

    return Array.from(distMap.values());
  };

  // Get appropriate labels for a play based on distribution type
  const getLabelsForPlay = (play, type: DistributionType): string[] => {
    if (type === "artist") {
      // Return primary artists (excluding producers, etc.)
      return play.release.artists
        .filter((a) => a.role !== "Producer")
        .map((a) => a.artist?.name || "Unknown")
        .filter(Boolean);
    } else if (type === "genre") {
      // Return genres
      return play.release.genres.map((g) => g.name);
    }

    return [];
  };

  // Helper to get play duration in minutes
  const getPlayDurationMinutes = (play) => {
    // If release has play_duration, use it (it's stored in seconds)
    if (play.release.playDuration) {
      return Math.round(play.release.playDuration / 60);
    }

    // Otherwise estimate based on format
    // For vinyl, typical LP is ~40 minutes
    return 40; // Default to 40 minutes if no duration info available
  };

  // Handle distribution type change
  const handleDistributionTypeChange = (event: Event) => {
    const target = event.target as HTMLSelectElement;
    setDistributionType(target.value as DistributionType);
  };

  // Handle top count change
  const handleTopCountChange = (event: Event) => {
    const target = event.target as HTMLSelectElement;
    setShowTopCount(parseInt(target.value, 10));
  };

  return (
    <div class={styles.chartsContainer}>
      <h3 class={styles.chartTitle}>
        {distributionType() === "artist" ? "Artists" : "Genres"} Distribution
      </h3>

      <div class={styles.controls}>
        <ChartControls showFrequencyControls={false} showFilters={false} />

        <div class={styles.typeControls}>
          <div class={styles.controlGroup}>
            <label class={styles.label}>Distribution Type:</label>
            <select
              class={styles.select}
              value={distributionType()}
              onChange={handleDistributionTypeChange}
            >
              <option value="genre">By Genre</option>
              <option value="artist">By Artist</option>
            </select>
          </div>

          <div class={styles.controlGroup}>
            <label class={styles.label}>Show Top:</label>
            <select
              class={styles.select}
              value={showTopCount()}
              onChange={handleTopCountChange}
            >
              <option value="5">5</option>
              <option value="10">10</option>
              <option value="15">15</option>
              <option value="20">20</option>
              <option value="50">50</option>
            </select>
          </div>
        </div>
      </div>

      <Show
        when={!isLoading()}
        fallback={<div class={styles.loading}>Loading chart data...</div>}
      >
        <div class={styles.chartsWrapper}>
          <div class={styles.chartWrapper}>
            <h4 class={styles.subTitle}>By Play Count</h4>
            <Show
              when={countChartData().labels.length > 0}
              fallback={
                <div class={styles.noData}>
                  No data available for the selected period
                </div>
              }
            >
              <div class={styles.pieChartContainer}>
                <Pie data={countChartData()} options={countChartOptions} />
              </div>
            </Show>
          </div>

          <div class={styles.chartWrapper}>
            <h4 class={styles.subTitle}>By Listening Time</h4>
            <Show
              when={durationChartData().labels.length > 0}
              fallback={
                <div class={styles.noData}>
                  No data available for the selected period
                </div>
              }
            >
              <div class={styles.pieChartContainer}>
                <Doughnut
                  data={durationChartData()}
                  options={durationChartOptions}
                />
              </div>
            </Show>
          </div>
        </div>
      </Show>
    </div>
  );
};

export default DistributionCharts;
