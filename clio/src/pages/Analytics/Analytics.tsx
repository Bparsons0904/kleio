// In Analytics.tsx
import { Component, createSignal, createEffect } from "solid-js";
import PlayFrequencyChart from "../../components/charts/PlayFrequencyChart";
import {
  DateRangeProvider,
  useDateRange,
} from "../../provider/DateRangeContext";
import styles from "./Analytics.module.scss";
import PlayDurationChart from "../../components/charts/PlayDurationChart";
import DistributionCharts from "../../components/charts/DistributionChart";
import ChartControls from "../../components/charts/ChartControls";
import { useAppContext } from "../../provider/Provider";

const AnalyticsContent: Component = () => {
  const { releases } = useAppContext();
  const [filter, setFilter] = createSignal("");
  const [filterOptions, setFilterOptions] = createSignal<
    { value: string; label: string }[]
  >([]);

  // Generate filter options (artists and genres)
  createEffect(() => {
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
  });

  // Handle filter change
  const handleFilterChange = (newFilter: string) => {
    if (newFilter.startsWith("HEADER:")) return;
    setFilter(newFilter);
  };

  return (
    <div class={styles.dashboard}>
      <h2 class={styles.dashboardTitle}>Listening Analytics</h2>
      <p class={styles.dashboardDescription}>
        Visualize your vinyl listening habits with interactive charts. Use the
        filters to explore your collection by time period, artist, or genre.
      </p>

      {/* Shared controls for both time-based charts */}
      <div class={styles.controlsSection}>
        <ChartControls
          showFrequencyControls={true}
          showFilters={true}
          filterOptions={filterOptions()}
          filterValue={filter()}
          onFilterChange={handleFilterChange}
          filterLabel="Filter by Artist/Genre:"
        />
      </div>

      <div class={styles.chartContainer}>
        <h3 class={styles.chartTitle}>Records Played Over Time</h3>
        <PlayFrequencyChart filter={filter()} />
      </div>

      <div class={styles.chartContainer}>
        <h3 class={styles.chartTitle}>Listening Time Over Time</h3>
        <PlayDurationChart filter={filter()} />
      </div>

      <div class={styles.chartContainer}>
        <DistributionCharts />
      </div>
    </div>
  );
};

// Wrapper component that provides the DateRangeProvider
const Analytics: Component = () => {
  return (
    <DateRangeProvider>
      <AnalyticsContent />
    </DateRangeProvider>
  );
};

export default Analytics;
