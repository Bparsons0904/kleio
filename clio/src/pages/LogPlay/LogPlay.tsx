import { Component, createSignal, createEffect, For, Show } from "solid-js";
import { useAppContext } from "../../provider/Provider";
import styles from "./LogPlay.module.scss";
import { Release, Stylus } from "../../types";
import { createPlayHistory } from "../../utils/api";

const LogPlay: Component = () => {
  console.log("Rendering LogPlay");
  const { releases, styluses } = useAppContext();
  const [filteredReleases, setFilteredReleases] = createSignal<Release[]>([]);
  const [searchTerm, setSearchTerm] = createSignal("");
  const [selectedDate, setSelectedDate] = createSignal(
    new Date().toISOString().split("T")[0],
  );
  const [selectedRelease, setSelectedRelease] = createSignal<Release | null>(
    null,
  );
  const [selectedStylus, setSelectedStylus] = createSignal<Stylus | null>(null);
  const [isSubmitting, setIsSubmitting] = createSignal(false);

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
  };

  const handleLogPlay = async () => {
    if (!selectedRelease()) {
      return;
    }

    try {
      setIsSubmitting(true);
      const payload = {
        releaseId: selectedRelease()!.id,
        playedAt: new Date(selectedDate()).toISOString(),
        stylusId: selectedStylus()?.id || null,
      };

      console.log("Logging play:", payload);

      const response = await createPlayHistory(payload);

      if (response.status === 201) {
        // Success feedback
        alert("Play logged successfully!");
        // Optionally reset form
        setSelectedRelease(null);
      } else {
        throw new Error("Failed to log play");
      }
    } catch (error) {
      console.error("Error logging play:", error);
      alert("Failed to log play. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  // Get only active styluses for the dropdown
  const activeStyluses = () => {
    return (
      styluses()?.filter((stylus) => stylus.active || stylus.primary) || []
    );
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
      <h1 class={styles.title}>Log Play</h1>
      <p class={styles.subtitle}>
        Record when you play a record from your collection.
      </p>

      <div class={styles.logPlayForm}>
        {/* ROW 1: Date and stylus selection */}
        <div class={styles.controlsRow}>
          <div class={styles.formControl}>
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

          <div class={styles.formControl}>
            <label class={styles.label} for="stylusSelect">
              Stylus Used
            </label>
            <select
              id="stylusSelect"
              class={styles.stylusSelect}
              value={selectedStylus()?.id || ""}
              onChange={(e) => {
                const id = parseInt(e.target.value);
                const stylus = styluses().find((s) => s.id === id);
                setSelectedStylus(stylus || null);
              }}
            >
              <option value="">None</option>
              <For each={activeStyluses()}>
                {(stylus) => (
                  <option value={stylus.id}>
                    {stylus.name} {stylus.primary ? "(Primary)" : ""}
                  </option>
                )}
              </For>
            </select>
          </div>
        </div>

        {/* ROW 2: Search input */}
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

        {/* ROW 3: Collection title and log button */}
        <div class={styles.collectionHeader}>
          <h2 class={styles.sectionTitle}>Your Collection</h2>
          <button
            class={styles.logButton}
            disabled={!selectedRelease() || isSubmitting()}
            onClick={handleLogPlay}
          >
            {isSubmitting() ? "Logging..." : "Log Play"}
          </button>
        </div>

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
                    </div>
                  </div>
                )}
              </For>
            </div>
          )}
        </div>

        {/* Removed the bottom action section since we moved the button up */}
        {!selectedRelease() && (
          <p class={styles.selectionHint}>
            Please select a release to log a play
          </p>
        )}
      </div>
    </div>
  );
};

export default LogPlay;
