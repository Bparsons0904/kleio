// src/pages/LogPlay/LogPlay.tsx (updated)
import { Component, createSignal, createEffect, For, Show } from "solid-js";
import { useAppContext } from "../../provider/Provider";
import styles from "./LogPlay.module.scss";
import { Release } from "../../types";
import RecordActionModal from "../../components/RecordActionModal/RecordActionModal";
import { exportHistory } from "../../utils/mutations/export";
import { RecordStatusIndicator } from "../../components/StatusIndicators/StatusIndicators";
import { getLastPlayDate, getLastCleaningDate } from "../../utils/playStatus";

const LogPlay: Component = () => {
  const { releases, showError } = useAppContext();

  const [filteredReleases, setFilteredReleases] = createSignal<Release[]>([]);
  const [searchTerm, setSearchTerm] = createSignal("");
  const [selectedRelease, setSelectedRelease] = createSignal<Release | null>(
    null,
  );
  const [isModalOpen, setIsModalOpen] = createSignal(false);
  const [showStatusDetails, setShowStatusDetails] = createSignal(false);

  // Filter releases based on search term
  createEffect(() => {
    const term = searchTerm().toLowerCase();
    if (!term) {
      setFilteredReleases(releases?.());
      return;
    }

    const filtered = releases().filter(
      (release) =>
        release.title.toLowerCase().includes(term) ||
        release.artists.some((artist) =>
          artist.artist?.name.toLowerCase().includes(term),
        ),
    );

    setFilteredReleases(filtered);
  });

  // Make sure selected release is always up to date with the latest data
  createEffect(() => {
    if (selectedRelease()) {
      for (const release of releases()) {
        if (release.id === selectedRelease().id) {
          setSelectedRelease(release);
          break;
        }
      }
    }
  });

  const handleReleaseClick = (release: Release) => {
    setSelectedRelease(release);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    // We keep the selected release so clicking the card again reopens the modal
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
      <h1 class={styles.title}>Log Play & Cleaning</h1>
      <p class={styles.subtitle}>
        Record when you play or clean records from your collection.
      </p>

      <div class={styles.logForm}>
        {/* Search input */}
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
              No releases found. Try a different search term.
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
      <div class={styles.exportSection}>
        <button class={styles.exportButton} onClick={handleExport}>
          Export Play & Cleaning History
        </button>
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
