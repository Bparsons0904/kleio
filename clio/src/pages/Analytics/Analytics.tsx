import { Component } from "solid-js";
import PlayFrequencyChart from "../../components/charts/PlayFrequencyChart";
import { DateRangeProvider } from "../../provider/DateRangeContext";
import styles from "./Analytics.module.scss";
import PlayDurationChart from "../../components/charts/PlayDurationChart";
import DistributionCharts from "../../components/charts/DistributionChart";

const Analytics: Component = () => {
  return (
    <DateRangeProvider>
      <div class={styles.dashboard}>
        <h2 class={styles.dashboardTitle}>Listening Analytics</h2>
        <p class={styles.dashboardDescription}>
          Visualize your vinyl listening habits with interactive charts. Use the
          filters to explore your collection by time period, artist, or genre.
        </p>

        <div class={styles.chartContainer}>
          <PlayFrequencyChart />
        </div>

        <div class={styles.chartContainer}>
          <PlayDurationChart />
        </div>

        <div class={styles.chartContainer}>
          <DistributionCharts />
        </div>
      </div>
    </DateRangeProvider>
  );
};

export default Analytics;
