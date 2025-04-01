// src/pages/Analytics/Analytics.tsx
import { Component, createSignal, createEffect } from "solid-js";
import PlayFrequencyChart from "../../components/charts/PlayFrequencyChart";
import { DateRangeProvider } from "../../provider/DateRangeContext";
import styles from "./Analytics.module.scss";
import PlayDurationChart from "../../components/charts/PlayDurationChart";
import DistributionCharts from "../../components/charts/DistributionChart";
import ChartControls from "../../components/charts/ChartControls";
import { useAppContext } from "../../provider/Provider";
import { DropdownOption } from "../../components/SearchableDropdown/SearchableDropdown";

const AnalyticsContent: Component = () => {
  const { releases } = useAppContext();
  const [filter, setFilter] = createSignal("");
  const [filterOptions, setFilterOptions] = createSignal<DropdownOption[]>([]);

  // Generate filter options (artists, genres, and records)
  createEffect(() => {
    // Using sets to avoid duplicates
    const artistSet = new Set<string>();
    const genreSet = new Set<string>();
    // We'll use the records directly to maintain title/ID pairs
    const recordOptions: DropdownOption[] = [];

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

      // Add record titles
      recordOptions.push({
        value: `record:${release.id}`,
        label: release.title,
      });
    });

    // Convert sets to sorted arrays
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

    // Sort record options by title
    recordOptions.sort((a, b) => a.label.localeCompare(b.label));

    // Combine options with headers
    setFilterOptions([
      { value: "", label: "All Records" },
      { value: "HEADER:RECORDS", label: "--- Records ---", disabled: true },
      ...recordOptions,
      { value: "HEADER:ARTISTS", label: "--- Artists ---", disabled: true },
      ...artistOptions,
      { value: "HEADER:GENRES", label: "--- Genres ---", disabled: true },
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
        filters to explore your collection by time period, artist, genre, or
        specific record.
      </p>

      {/* Shared controls for both time-based charts */}
      <div class={styles.controlsSection}>
        <ChartControls
          showFrequencyControls={true}
          showFilters={true}
          filterOptions={filterOptions()}
          filterValue={filter()}
          onFilterChange={handleFilterChange}
          filterLabel="Filter by Record/Artist/Genre:"
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
