// pages/LogPlay/LogPlay.tsx
import { Component, createSignal, createEffect, For, Show } from "solid-js";
import { useAppContext } from "../../provider/Provider";
import styles from "./LogPlay.module.scss";
import { Release, Stylus } from "../../types";
import { createPlayHistory, createCleaningHistory } from "../../utils/api";
import RecordActionModal from "../../components/RecordActionModal/RecordActionModal";

const LogPlay: Component = () => {
  const {
    releases,
    styluses,
    showSuccess,
    showError,
    setKleioStore: setAuthPayload,
  } = useAppContext();

  const [filteredReleases, setFilteredReleases] = createSignal<Release[]>([]);
  const [searchTerm, setSearchTerm] = createSignal("");
  const [selectedRelease, setSelectedRelease] = createSignal<Release | null>(
    null,
  );
  const [isModalOpen, setIsModalOpen] = createSignal(false);

  // Form state for the modal
  const [selectedDate, setSelectedDate] = createSignal(
    new Date().toISOString().split("T")[0],
  );
  const [selectedStylus, setSelectedStylus] = createSignal<Stylus | null>(null);
  const [notes, setNotes] = createSignal("");
  const [isSubmittingPlay, setIsSubmittingPlay] = createSignal(false);
  const [isSubmittingCleaning, setIsSubmittingCleaning] = createSignal(false);
  const [isSubmittingBoth, setIsSubmittingBoth] = createSignal(false);

  // Set primary stylus as the default selected stylus
  createEffect(() => {
    const primaryStylus = styluses()?.find(
      (stylus) => stylus.primary && stylus.active,
    );
    if (primaryStylus) {
      setSelectedStylus(primaryStylus);
    }
  });

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

  const handleReleaseClick = (release: Release) => {
    setSelectedRelease(release);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    // We keep the selected release so clicking the card again reopens the modal
  };

  const handleLogPlay = async () => {
    if (!selectedRelease()) {
      return;
    }

    try {
      setIsSubmittingPlay(true);
      const payload = {
        releaseId: selectedRelease()!.id,
        playedAt: new Date(selectedDate()).toISOString(),
        stylusId: selectedStylus()?.id || null,
        notes: notes(),
      };

      const response = await createPlayHistory(payload);

      if (response.status === 201) {
        showSuccess("Play logged successfully!");
        setIsModalOpen(false);
        setNotes("");
        setAuthPayload(response.data);
      } else {
        throw new Error("Failed to log play");
      }
    } catch (error) {
      console.error("Error logging play:", error);
      showError("Failed to log play. Please try again.");
    } finally {
      setIsSubmittingPlay(false);
    }
  };

  const handleLogCleaning = async () => {
    if (!selectedRelease()) {
      return;
    }

    try {
      setIsSubmittingCleaning(true);
      const payload = {
        releaseId: selectedRelease()!.id,
        cleanedAt: new Date(selectedDate()).toISOString(),
        notes: notes(),
      };

      const response = await createCleaningHistory(payload);

      if (response.status === 201) {
        showSuccess("Cleaning logged successfully!");
        setIsModalOpen(false);
        setNotes("");
        setAuthPayload(response.data);
      } else {
        throw new Error("Failed to log cleaning");
      }
    } catch (error) {
      console.error("Error logging cleaning:", error);
      showError("Failed to log cleaning. Please try again.");
    } finally {
      setIsSubmittingCleaning(false);
    }
  };

  const handleLogBoth = async () => {
    if (!selectedRelease()) {
      return;
    }

    try {
      setIsSubmittingBoth(true);

      // Create cleaning payload
      const cleaningPayload = {
        releaseId: selectedRelease()!.id,
        cleanedAt: new Date(selectedDate()).toISOString(),
        notes: notes(),
      };

      // Create play payload
      const playPayload = {
        releaseId: selectedRelease()!.id,
        playedAt: new Date(selectedDate()).toISOString(),
        stylusId: selectedStylus()?.id || null,
        notes: notes(),
      };

      // Log cleaning first
      const cleaningResponse = await createCleaningHistory(cleaningPayload);
      if (cleaningResponse.status !== 201) {
        throw new Error("Failed to log cleaning");
      }

      // Then log play
      const playResponse = await createPlayHistory(playPayload);
      if (playResponse.status !== 201) {
        throw new Error("Failed to log play");
      }

      showSuccess("Both cleaning and play logged successfully!");
      setIsModalOpen(false);
      setNotes("");
      setAuthPayload(playResponse.data); // Use the latest response data
    } catch (error) {
      console.error("Error logging both:", error);
      showError("Failed to log both activities. Please try again.");
    } finally {
      setIsSubmittingBoth(false);
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
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
                      {release.playHistory && release.playHistory.length > 0 ? (
                        <p class={styles.lastPlayed}>
                          Last played:{" "}
                          {formatDate(release.playHistory[0].playedAt)}
                        </p>
                      ) : (
                        <p class={styles.neverPlayed}>Never played</p>
                      )}
                      {/* Add last cleaned info */}
                      {release.cleaningHistory &&
                      release.cleaningHistory.length > 0 ? (
                        <p class={styles.lastCleaned}>
                          Last cleaned:{" "}
                          {formatDate(release.cleaningHistory[0].cleanedAt)}
                        </p>
                      ) : (
                        <p class={styles.neverCleaned}>Never cleaned</p>
                      )}
                    </div>
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
          date={selectedDate()}
          setDate={setSelectedDate}
          selectedStylus={selectedStylus()}
          setSelectedStylus={setSelectedStylus}
          styluses={styluses()}
          notes={notes()}
          setNotes={setNotes}
          onLogPlay={handleLogPlay}
          onLogCleaning={handleLogCleaning}
          onLogBoth={handleLogBoth}
          isSubmittingPlay={isSubmittingPlay()}
          isSubmittingCleaning={isSubmittingCleaning()}
          isSubmittingBoth={isSubmittingBoth()}
        />
      </Show>
    </div>
  );
};

export default LogPlay;
