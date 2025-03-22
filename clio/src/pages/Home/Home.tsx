import { Component } from "solid-js";
import styles from "./Home.module.scss";

const Home: Component = () => {
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
            <button class={styles.button}>Log Now</button>
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
            <button class={styles.button}>View Stats</button>
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
            <button class={styles.button}>View Collection</button>
          </div>
        </div>

        <div class={styles.card}>
          <div class={styles.cardHeader}>
            <h2>Go to Discogs</h2>
          </div>
          <div class={styles.cardBody}>
            <p>Visit your Discogs profile to manage your collection.</p>
          </div>
          <div class={styles.cardFooter}>
            <a
              href="https://www.discogs.com/user/collection"
              target="_blank"
              rel="noopener noreferrer"
              class={styles.button}
            >
              Open Discogs
            </a>
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
            <button class={styles.button}>Sync Now</button>
          </div>
        </div>

        <div class={styles.card}>
          <div class={styles.cardHeader}>
            <h2>View Stats</h2>
          </div>
          <div class={styles.cardBody}>
            <p>Explore insights about your collection and listening habits.</p>
          </div>
          <div class={styles.cardFooter}>
            <button class={styles.button}>View Insights</button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Home;
