// src/pages/Home/Home.tsx (Modified)
import { Component } from "solid-js";
import styles from "./Home.module.scss";
import { refreshCollection } from "../../utils/api";
import { useNavigate } from "@solidjs/router";
import { useAppContext } from "../../provider/Provider";
import { exportHistory } from "../../utils/mutations/export";
import FolderSelector from "../../components/FolderSelector/FolderSelector";

const Home: Component = () => {
  const { setIsSyncing, showError } = useAppContext();
  const navigate = useNavigate();

  const handleLogPlay = () => {
    navigate("/log");
  };

  const handleManageStyluses = () => {
    navigate("/equipment");
  };

  const handleResync = async () => {
    try {
      const response = await refreshCollection();
      if (response.status === 200) {
        setIsSyncing(true);
      }
    } catch (error) {
      console.error("Error resyncing:", error);
    }
  };

  const handleExport = async () => {
    try {
      await exportHistory();
      // No need for success notification since it's a direct download
    } catch (error) {
      console.error("Error exporting history:", error);
      showError("Failed to export history. Please try again.");
    }
  };

  return (
    <div class={styles.container}>
      <h1 class={styles.title}>Welcome to Kleio</h1>
      <p class={styles.subtitle}>Your personal vinyl collection tracker</p>

      <div class={styles.cardGrid}>
        <div class={styles.card}>
          <div class={styles.cardHeader}>
            <h2>Log Play</h2>
          </div>
          <div class={styles.cardBody}>
            <p>Record when you play a record from your collection.</p>
          </div>
          <div class={styles.cardFooter}>
            <button class={styles.button} on:click={handleLogPlay}>
              Log Now
            </button>
          </div>
        </div>

        <div class={styles.card}>
          <div class={styles.cardHeader}>
            <h2>View Play Time</h2>
          </div>
          <div class={styles.cardBody}>
            <p>See statistics about your listening habits.</p>
          </div>
          <div class={styles.cardFooter}>
            <button
              class={styles.button}
              on:click={() => navigate("/playHistory")}
            >
              View Stats
            </button>
          </div>
        </div>

        <div class={styles.card}>
          <div class={styles.cardHeader}>
            <h2>View Collection</h2>
          </div>
          <div class={styles.cardBody}>
            <p>Browse and search through your vinyl collection.</p>
          </div>
          <div class={styles.cardFooter}>
            <button
              class={styles.button}
              on:click={() => navigate("/collection")}
            >
              View Collection
            </button>
          </div>
        </div>

        <div class={styles.card}>
          <div class={styles.cardHeader}>
            <h2>View Equipment</h2>
          </div>
          <div class={styles.cardBody}>
            <p>View, Edit and add equipment to your profile.</p>
          </div>
          <div class={styles.cardFooter}>
            <button class={styles.button} on:click={handleManageStyluses}>
              View Equipment
            </button>
          </div>
        </div>

        <div class={styles.card}>
          <div class={styles.cardHeader}>
            <h2>Refresh Collection</h2>
          </div>
          <div class={styles.cardBody}>
            <p>Sync your Kleio collection with your Discogs library.</p>
          </div>
          <div class={styles.cardFooter}>
            <button class={styles.button} on:click={handleResync}>
              Sync Now
            </button>
          </div>
        </div>

        <div class={styles.card}>
          <div class={styles.cardHeader}>
            <h2>View Analytics</h2>
          </div>
          <div class={styles.cardBody}>
            <p>Explore insights about your collection and listening habits.</p>
          </div>
          <div class={styles.cardFooter}>
            <button
              class={styles.button}
              on:click={() => navigate("/analytics")}
            >
              View Insights
            </button>
          </div>
        </div>
      </div>

      <div class={styles.folderSelectorSection}>
        <FolderSelector />
      </div>

      <div class={styles.exportSection}>
        <button class={styles.exportButton} onClick={handleExport}>
          Export Play & Cleaning History
        </button>
      </div>
    </div>
  );
};

export default Home;
