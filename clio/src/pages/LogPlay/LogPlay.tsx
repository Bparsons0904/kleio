// src/pages/LogPlay/LogPlay.tsx (modified)
import { Component, createSignal, createEffect, For, Show } from "solid-js";
import { useAppContext } from "../../provider/Provider";
import styles from "./LogPlay.module.scss";
import { Release } from "../../types";
import RecordActionModal from "../../components/RecordActionModal/RecordActionModal";
import { RecordStatusIndicator } from "../../components/StatusIndicators/StatusIndicators";
import { getLastPlayDate } from "../../utils/playStatus";
import {
  createPlayHistory,
  createCleaningHistory,
} from "../../utils/mutations/post";

const LogPlay: Component = () => {
  const { releases, styluses, showSuccess, showError, setKleioStore } =
    useAppContext();

  const [filteredReleases, setFilteredReleases] = createSignal<Release[]>([]);
  const [searchTerm, setSearchTerm] = createSignal("");
  const [selectedRelease, setSelectedRelease] = createSignal<Release | null>(
    null,
  );
  const [isModalOpen, setIsModalOpen] = createSignal(false);
  const [showStatusDetails, setShowStatusDetails] = createSignal(false);
  const [sortBy, setSortBy] = createSignal("artist"); // Default sort by album title

  // Filter releases based on search term and sort them
  createEffect(() => {
    const term = searchTerm().toLowerCase();
    let filtered = [...releases()];

    // Apply search filter
    if (term) {
      filtered = filtered.filter(
        (release) =>
          release.title.toLowerCase().includes(term) ||
          release.artists.some((artist) =>
            artist.artist?.name.toLowerCase().includes(term),
          ),
      );
    }

    // Apply sorting
    filtered = sortReleases(filtered, sortBy());

    setFilteredReleases(filtered);
  });

  // Sort releases based on selected sort option
  const sortReleases = (releases: Release[], sortOption: string): Release[] => {
    switch (sortOption) {
      case "lastPlayed":
        // Sort by last played date (most recent first)
        return [...releases].sort((a, b) => {
          const dateA = getLastPlayDate(a.playHistory);
          const dateB = getLastPlayDate(b.playHistory);

          // If both have play history
          if (dateA && dateB) return dateB.getTime() - dateA.getTime();
          // If only A has play history
          if (dateA && !dateB) return -1;
          // If only B has play history
          if (!dateA && dateB) return 1;
          // If neither has play history, sort by title
          return a.title.localeCompare(b.title);
        });

      case "recentlyPlayed":
        // Show only recently played records (last 30 days) sorted by date
        return [...releases]
          .filter((release) => {
            const lastPlayed = getLastPlayDate(release.playHistory);
            if (!lastPlayed) return false;

            const thirtyDaysAgo = new Date();
            thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);

            return lastPlayed >= thirtyDaysAgo;
          })
          .sort((a, b) => {
            const dateA = getLastPlayDate(a.playHistory)!;
            const dateB = getLastPlayDate(b.playHistory)!;
            return dateB.getTime() - dateA.getTime();
          });

      case "genre":
        // Sort by primary genre
        return [...releases].sort((a, b) => {
          const genreA = a.genres[0]?.name || "Unknown";
          const genreB = b.genres[0]?.name || "Unknown";
          return genreA.localeCompare(genreB);
        });

      case "artist":
        // Sort by primary artist
        return [...releases].sort((a, b) => {
          const artistA =
            a.artists.find((a) => a.role !== "Producer")?.artist?.name ||
            "Unknown";
          const artistB =
            b.artists.find((a) => a.role !== "Producer")?.artist?.name ||
            "Unknown";
          return artistA.localeCompare(artistB);
        });

      case "album":
        // Sort by album title
        return [...releases].sort((a, b) => a.title.localeCompare(b.title));

      case "year":
        // Sort by release year (newest first)
        return [...releases].sort((a, b) => {
          const yearA = a.year || 0;
          const yearB = b.year || 0;
          return yearB - yearA;
        });

      case "recentlyAdded":
        // Sort by date added to collection (using createdAt date)
        return [...releases].sort((a, b) => {
          const dateA = new Date(a.createdAt).getTime();
          const dateB = new Date(b.createdAt).getTime();
          return dateB - dateA;
        });

      case "playCount":
        // Sort by number of plays (most played first)
        return [...releases].sort((a, b) => {
          const countA = a.playHistory?.length || 0;
          const countB = b.playHistory?.length || 0;
          return countB - countA;
        });

      default:
        // Default sort by title
        return [...releases].sort((a, b) => a.title.localeCompare(b.title));
    }
  };

  const handleReleaseClick = (release: Release) => {
    setSelectedRelease(release);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    // We keep the selected release so clicking the card again reopens the modal
  };

  const handleQuickPlay = async (release: Release) => {
    const primaryStylus = styluses().find((stylus) => stylus.primary);

    try {
      const result = await createPlayHistory({
        releaseId: release.id,
        playedAt: new Date().toISOString(),
        stylusId: primaryStylus?.id,
      });

      if (result.success) {
        showSuccess(`Logged play for "${release.title}"`);
        setKleioStore(result.data);
      } else {
        throw new Error("Failed to log play");
      }
    } catch (error) {
      console.error("Error quick logging play:", error);
      showError("Failed to log play. Please try again.");
    }
  };

  const handleQuickCleaning = async (release: Release) => {
    try {
      const result = await createCleaningHistory({
        releaseId: release.id,
        cleanedAt: new Date().toISOString(),
      });

      if (result.success) {
        showSuccess(`Logged cleaning for "${release.title}"`);
        setKleioStore(result.data);
      } else {
        throw new Error("Failed to log cleaning");
      }
    } catch (error) {
      console.error("Error quick logging cleaning:", error);
      showError("Failed to log cleaning. Please try again.");
    }
  };

  return (
    <div class={styles.container}>
      <h1 class={styles.title}>Log Play & Cleaning</h1>
      <p class={styles.subtitle}>
        Record when you play or clean records from your collection.
      </p>

      <div class={styles.logForm}>
        {/* Search and sorting controls */}
        <div class={styles.controlsRow}>
          <div class={styles.searchSection}>
            <label class={styles.label} for="releaseSearch">
              Search Your Collection
            </label>
            <input
              type="text"
              id="releaseSearch"
              class={styles.searchInput}
              value={searchTerm()}
              onInput={(e) => setSearchTerm(e.target.value)}
              placeholder="Search by title or artist..."
            />
          </div>

          <div class={styles.sortSection}>
            <label class={styles.label} for="sortOptions">
              Sort By
            </label>
            <select
              id="sortOptions"
              class={styles.select}
              value={sortBy()}
              onInput={(e) => setSortBy(e.target.value)}
            >
              <option value="album">Album (A-Z)</option>
              <option value="artist">Artist (A-Z)</option>
              <option value="genre">Genre (A-Z)</option>
              <option value="lastPlayed">Last Played</option>
              <option value="recentlyPlayed">Recently Played (30 days)</option>
              <option value="year">Release Year (newest first)</option>
              <option value="recentlyAdded">Recently Added</option>
              <option value="playCount">Most Played</option>
            </select>
          </div>
        </div>

        {/* Status details toggle */}
        <div class={styles.optionsSection}>
          <label class={styles.toggleLabel}>
            <input
              type="checkbox"
              checked={showStatusDetails()}
              onChange={(e) => setShowStatusDetails(e.target.checked)}
              class={styles.toggleInput}
            />
            <span class={styles.toggleSwitch}></span>
            Show status details
          </label>
        </div>

        {/* Collection heading */}
        <h2 class={styles.sectionTitle}>Your Collection</h2>

        {/* Release list */}
        <div class={styles.releasesSection}>
          {filteredReleases().length === 0 ? (
            <p class={styles.noResults}>
              No releases found. Try a different search term or sort option.
            </p>
          ) : (
            <div class={styles.releasesList}>
              <For each={filteredReleases()}>
                {(release) => (
                  <div
                    class={`${styles.releaseCard} ${selectedRelease()?.id === release.id ? styles.selected : ""}`}
                    onClick={() => handleReleaseClick(release)}
                  >
                    <div class={styles.releaseCardContainer}>
                      <div class={styles.releaseImageContainer}>
                        {release.thumb ? (
                          <img
                            src={release.thumb}
                            alt={release.title}
                            class={styles.releaseImage}
                          />
                        ) : (
                          <div class={styles.noImage}>No Image</div>
                        )}
                        {release.year && (
                          <div class={styles.releaseYear}>{release.year}</div>
                        )}
                      </div>
                      <div class={styles.releaseInfo}>
                        <h3 class={styles.releaseTitle}>{release.title}</h3>
                        <p class={styles.releaseArtist}>
                          {release.artists
                            .filter((artist) => artist.role !== "Producer")
                            .map((artist) => artist.artist?.name)
                            .join(", ")}
                        </p>

                        <div class={styles.statusSection}>
                          <RecordStatusIndicator
                            playHistory={release.playHistory}
                            cleaningHistory={release.cleaningHistory}
                            showDetails={false}
                            onPlayClick={() => handleQuickPlay(release)}
                            onCleanClick={() => handleQuickCleaning(release)}
                          />
                        </div>
                      </div>
                    </div>

                    {/* Details section that spans full width */}
                    <Show when={showStatusDetails()}>
                      <div class={styles.fullWidthDetails}>
                        <RecordStatusIndicator
                          playHistory={release.playHistory}
                          cleaningHistory={release.cleaningHistory}
                          showDetails={true}
                        />
                      </div>
                    </Show>
                  </div>
                )}
              </For>
            </div>
          )}
        </div>
      </div>

      {/* The modal */}
      <Show when={selectedRelease()}>
        <RecordActionModal
          isOpen={isModalOpen()}
          onClose={handleCloseModal}
          release={selectedRelease()!}
        />
      </Show>
    </div>
  );
};

export default LogPlay;
