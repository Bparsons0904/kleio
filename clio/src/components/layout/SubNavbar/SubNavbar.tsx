// src/components/SubNavbar/SubNavbar.tsx
import { Component, createEffect, createSignal, For, Show } from "solid-js";
import styles from "./SubNavbar.module.scss";
import { BsVinylFill } from "solid-icons/bs";
import { useAppContext } from "../../../provider/Provider";
import { Release } from "../../../types";
import { createPlayHistory } from "../../../utils/mutations/post";

const SubNavbar: Component = () => {
  const { releases, showSuccess, showError, setKleioStore, styluses } =
    useAppContext();
  const [suggestType, setSuggestType] = createSignal("one");
  const [suggestedReleases, setSuggestedReleases] = createSignal<Release[]>([]);
  const [showSuggestions, setShowSuggestions] = createSignal(false);
  const [recordOfTheDay, setRecordOfTheDay] = createSignal<Release | null>(
    null,
  );

  // Check for record of the day on component mount
  createEffect(() => {
    const today = new Date().toLocaleDateString();
    const storedRecord = localStorage.getItem("recordOfTheDay");
    const storedDate = localStorage.getItem("recordOfTheDayDate");

    // If we have a stored record and it's from today, use it
    if (storedRecord && storedDate === today) {
      const recordId = parseInt(storedRecord);
      const record = releases().find((r) => r.id === recordId);
      if (record) {
        setRecordOfTheDay(record);
        return;
      }
    }

    // Otherwise, select a new record of the day
    selectNewRecordOfTheDay();
  });

  const selectNewRecordOfTheDay = () => {
    if (releases().length === 0) return;

    const randomIndex = Math.floor(Math.random() * releases().length);
    const record = releases()[randomIndex];

    // Store in localStorage
    localStorage.setItem("recordOfTheDay", record.id.toString());
    localStorage.setItem("recordOfTheDayDate", new Date().toLocaleDateString());
    localStorage.setItem("recordOfTheDayPlayed", "false");

    setRecordOfTheDay(record);
  };

  const generateSuggestions = () => {
    if (releases().length === 0) return;

    let suggested: Release[] = [];

    switch (suggestType()) {
      case "one":
        // Suggest one random record
        const randomIndex = Math.floor(Math.random() * releases().length);
        suggested = [releases()[randomIndex]];
        break;

      case "several":
        // Suggest three random records
        suggested = getRandomReleases(3);
        break;

      case "leastPlayed":
        // Create a weighted list that favors less-played records
        const weightedReleases = [...releases()].map((release) => {
          // Calculate a weight inversely proportional to play count
          // Records with 0 plays get the highest weight
          const playCount = release.playHistory?.length || 0;
          // Add a small random factor to avoid always selecting the same records
          const weight = 100 - Math.min(playCount * 10, 95) + Math.random() * 5;

          return { release, weight };
        });

        // Sort by weight descending (higher weight = more likely to be selected)
        weightedReleases.sort((a, b) => b.weight - a.weight);

        // Take top weighted records, but add their play count for display
        suggested = weightedReleases.slice(0, 3).map((item) => {
          const playCount = item.release.playHistory?.length || 0;
          // Create a new object to avoid modifying original data
          return { ...item.release, playCountDisplay: playCount };
        });
        break;

      case "randomGenre":
        // Get a random genre
        const allGenres = new Set<string>();
        releases().forEach((release) => {
          release.genres.forEach((genre) => {
            allGenres.add(genre.name);
          });
        });

        const genres = Array.from(allGenres);
        if (genres.length === 0) {
          suggested = getRandomReleases(3);
          break;
        }

        const randomGenre = genres[Math.floor(Math.random() * genres.length)];

        // Get all records from that genre
        suggested = releases().filter((release) =>
          release.genres.some((genre) => genre.name === randomGenre),
        );

        // Add genre to the first record for display purposes
        if (suggested.length > 0) {
          (suggested[0] as Release & { selectedGenre: string }).selectedGenre =
            randomGenre;
        }
        break;
    }

    setSuggestedReleases(suggested);
    setShowSuggestions(true);
  };

  const getRandomReleases = (count: number): Release[] => {
    const allReleases = [...releases()];
    const result: Release[] = [];

    // Safety check
    if (allReleases.length === 0) return [];

    // Get random records without duplicates
    while (result.length < count && allReleases.length > 0) {
      const randomIndex = Math.floor(Math.random() * allReleases.length);
      result.push(allReleases[randomIndex]);
      allReleases.splice(randomIndex, 1);
    }

    return result;
  };

  const handleLogPlay = async (release: Release) => {
    try {
      const primaryStylus = styluses().find((stylus) => stylus.primary);

      const result = await createPlayHistory({
        releaseId: release.id,
        playedAt: new Date().toISOString(),
        stylusId: primaryStylus?.id,
      });

      if (result.success) {
        showSuccess(`Logged play for "${release.title}"`);
        setKleioStore(result.data);

        // If this was the record of the day, update the UI
        if (recordOfTheDay()?.id === release.id) {
          localStorage.setItem("recordOfTheDayPlayed", "true");
          // Force a re-render
          setRecordOfTheDay({ ...release, played: true });
        }
      } else {
        throw new Error("Failed to log play");
      }
    } catch (error) {
      console.error("Error logging play:", error);
      showError("Failed to log play. Please try again.");
    }
  };

  const recordOfTheDayPlayed = () => {
    if (!recordOfTheDay()) return false;

    // Check localStorage first for performance
    if (localStorage.getItem("recordOfTheDayPlayed") === "true") return true;

    // Check if it's been played today
    const record = recordOfTheDay();
    if (!record || !record.playHistory || record.playHistory.length === 0)
      return false;

    const today = new Date();
    today.setHours(0, 0, 0, 0);

    return record.playHistory.some((play) => {
      const playDate = new Date(play.playedAt);
      playDate.setHours(0, 0, 0, 0);
      return playDate.getTime() === today.getTime();
    });
  };

  return (
    <div class={styles.subNavbar}>
      <div class={styles.container}>
        {/* Record of the day section */}
        <div class={styles.recordOfTheDay}>
          <Show
            when={recordOfTheDay()}
            fallback={<span class={styles.noRecord}>Loading...</span>}
          >
            <div class={styles.rotdAlbumCover} title="Record of the Day">
              {recordOfTheDay()?.thumb ? (
                <img
                  src={recordOfTheDay()?.thumb}
                  alt={recordOfTheDay()?.title}
                />
              ) : (
                <div class={styles.noImage}>
                  <BsVinylFill size={24} />
                </div>
              )}
            </div>

            <div class={styles.rotdDetails}>
              <span class={styles.rotdTitle}>{recordOfTheDay()?.title}</span>
              <span class={styles.rotdArtist}>
                {recordOfTheDay()
                  ?.artists.filter((artist) => artist.role !== "Producer")
                  .map((artist) => artist.artist?.name)
                  .join(", ")}
              </span>
            </div>

            <Show when={!recordOfTheDayPlayed()}>
              <button
                class={styles.playButton}
                onClick={() => handleLogPlay(recordOfTheDay())}
                title="Log play for record of the day"
              >
                Play
              </button>
            </Show>
            <Show when={recordOfTheDayPlayed()}>
              <div
                class={styles.playedBadge}
                title="You've played this record today"
              >
                Played
              </div>
            </Show>
          </Show>
        </div>

        {/* Random suggestion controls */}
        <div class={styles.suggestControls}>
          <select
            class={styles.suggestSelect}
            value={suggestType()}
            onChange={(e) => setSuggestType(e.target.value)}
          >
            <option value="one">Suggest One</option>
            <option value="several">Suggest Several</option>
            <option value="leastPlayed">Least Played</option>
            <option value="randomGenre">Random Genre</option>
          </select>

          <button class={styles.suggestButton} onClick={generateSuggestions}>
            Suggest
          </button>
        </div>
      </div>

      {/* Suggestions dropdown */}
      <Show when={showSuggestions() && suggestedReleases().length > 0}>
        <div class={styles.suggestionsDropdown}>
          <div class={styles.suggestionsHeader}>
            <h3>
              {suggestType() === "randomGenre" &&
              (suggestedReleases()[0] as Release & { selectedGenre: string })
                ?.selectedGenre
                ? `${(suggestedReleases()[0] as Release & { selectedGenre: string }).selectedGenre} Albums`
                : "Suggested Albums"}
            </h3>
            <button
              class={styles.closeButton}
              onClick={() => setShowSuggestions(false)}
            >
              ×
            </button>
          </div>

          <div class={styles.suggestionsList}>
            <For each={suggestedReleases()}>
              {(release) => (
                <div class={styles.suggestionItem}>
                  <div class={styles.albumImage}>
                    {release.thumb ? (
                      <img src={release.thumb} alt={release.title} />
                    ) : (
                      <div class={styles.noImage}>No Image</div>
                    )}
                  </div>

                  <div class={styles.releaseInfo}>
                    <div class={styles.releaseTitle}>{release.title}</div>
                    <div class={styles.releaseArtist}>
                      {release.artists
                        .filter((artist) => artist.role !== "Producer")
                        .map((artist) => artist.artist?.name)
                        .join(", ")}
                    </div>

                    <Show when={suggestType() === "leastPlayed"}>
                      <div class={styles.playCount}>
                        Play count: {release.playHistory?.length || 0}
                      </div>
                    </Show>
                  </div>

                  <button
                    class={styles.playButton}
                    onClick={() => handleLogPlay(release)}
                    title="Log play for this record"
                  >
                    <BsVinylFill size={16} />
                    Play
                  </button>
                </div>
              )}
            </For>
          </div>
        </div>
      </Show>
    </div>
  );
};

export default SubNavbar;
