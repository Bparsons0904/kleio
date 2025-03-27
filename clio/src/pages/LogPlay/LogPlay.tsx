import { Component, createSignal, createEffect, For } from "solid-js";
import { useAppContext } from "../../provider/Provider";
import styles from "./LogPlay.module.scss";
import { Release } from "../../types";

const LogPlay: Component = () => {
  console.log("Rendering LogPlay");
  const { releases } = useAppContext();
  const [filteredReleases, setFilteredReleases] = createSignal<Release[]>([]);
  const [searchTerm, setSearchTerm] = createSignal("");
  const [selectedDate, setSelectedDate] = createSignal(
    new Date().toISOString().split("T")[0],
  );
  const [selectedRelease, setSelectedRelease] = createSignal<Release | null>(
    null,
  );

  // Filter releases based on search term
  createEffect(() => {
    const term = searchTerm().toLowerCase();
    if (!term) {
      console.log("Rendering LogPlay term ");
      setFilteredReleases(releases?.());
      return;
    }

    console.log("Rendering LogPlay before ");
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
  };

  const handleLogPlay = () => {
    // This will be implemented later
    console.log("Logging play:", {
      releaseId: selectedRelease()?.id,
      date: selectedDate(),
    });
  };

  return (
    <div class={styles.container}>
      <h1 class={styles.title}>Log Play</h1>
      <p class={styles.subtitle}>
        Record when you play a record from your collection.
      </p>

      <div class={styles.logPlayForm}>
        <div class={styles.datePickerSection}>
          <label class={styles.label} for="playDate">
            Date Played
          </label>
          <input
            type="date"
            id="playDate"
            class={styles.datePicker}
            value={selectedDate()}
            onInput={(e) => setSelectedDate(e.target.value)}
          />
        </div>

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

        <div class={styles.releasesSection}>
          <h2 class={styles.sectionTitle}>Your Collection</h2>
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
                    </div>
                    <div class={styles.releaseInfo}>
                      <h3 class={styles.releaseTitle}>{release.title}</h3>
                      <p class={styles.releaseArtist}>
                        {release.artists
                          .filter((artist) => artist.role !== "Producer")
                          .map((artist) => artist.artist?.name)
                          .join(", ")}
                      </p>
                      {release.year && (
                        <p class={styles.releaseYear}>{release.year}</p>
                      )}
                    </div>
                  </div>
                )}
              </For>
            </div>
          )}
        </div>

        <div class={styles.actionSection}>
          <button
            class={styles.logButton}
            disabled={!selectedRelease()}
            onClick={handleLogPlay}
          >
            Log Play
          </button>
          {!selectedRelease() && (
            <p class={styles.selectionHint}>
              Please select a release to log a play
            </p>
          )}
        </div>
      </div>
    </div>
  );
};

export default LogPlay;
