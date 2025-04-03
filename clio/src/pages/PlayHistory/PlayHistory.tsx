import { Component, createMemo, createSignal, For, Show } from "solid-js";
import { useAppContext } from "../../provider/Provider";
import styles from "./PlayHistory.module.scss";
import { useFormattedShortDate } from "../../utils/dates";
import { VsNote } from "solid-icons/vs";
import { TbWashTemperature5 } from "solid-icons/tb";
import { PlayHistory, Release } from "../../types";
import RecordActionModal from "../../components/RecordActionModal/RecordActionModal";
import { fuzzySearchPlayHistory } from "../../utils/fuzzy";

const PlayHistoryPage: Component = () => {
  const { playHistory, releases } = useAppContext();
  const [timeFilter, setTimeFilter] = createSignal("month");
  const [searchTerm, setSearchTerm] = createSignal("");
  const [groupBy, setGroupBy] = createSignal("date");

  // Modal state
  const [selectedRelease, setSelectedRelease] = createSignal<Release | null>(
    null,
  );
  const [isModalOpen, setIsModalOpen] = createSignal(false);

  // Check if a cleaning was done on the same day as the play
  const hasCleaning = (release: Release, playDate: string) => {
    if (!release.cleaningHistory || release.cleaningHistory.length === 0) {
      return false;
    }

    const playDateStr = new Date(playDate).toISOString().split("T")[0];

    return release.cleaningHistory.some((cleaning) => {
      const cleaningDateStr = new Date(cleaning.cleanedAt)
        .toISOString()
        .split("T")[0];
      return cleaningDateStr === playDateStr;
    });
  };

  const getFilteredDate = () => {
    const now = new Date();
    switch (timeFilter()) {
      case "week":
        const lastWeek = new Date();
        lastWeek.setDate(now.getDate() - 7);
        return lastWeek;
      case "month":
        const lastMonth = new Date();
        lastMonth.setMonth(now.getMonth() - 1);
        return lastMonth;
      case "year":
        const lastYear = new Date();
        lastYear.setFullYear(now.getFullYear() - 1);
        return lastYear;
      default:
        return new Date(0); // Jan 1, 1970
    }
  };

  // Filter and sort play history
  const filteredHistory = createMemo(() => {
    let filtered = [...playHistory()];

    // Filter by date
    const filterDate = getFilteredDate();
    filtered = filtered.filter((play) => new Date(play.playedAt) >= filterDate);

    if (searchTerm().trim()) {
      filtered = fuzzySearchPlayHistory(filtered, searchTerm());
    }

    return filtered;
  });

  // Group history if needed
  const groupedHistory = createMemo(() => {
    const history = filteredHistory();
    if (groupBy() === "none") return { "": history };

    const grouped = history.reduce((acc, play) => {
      let key: string;
      if (groupBy() === "date") {
        // Group by day
        key = new Date(play.playedAt).toISOString().split("T")[0];
      } else if (groupBy() === "artist") {
        // Group by first artist
        key = play.release.artists[0]?.artist?.name || "Unknown Artist";
      } else if (groupBy() === "album") {
        // Group by album
        key = play.release.title;
      }

      if (!acc[key]) acc[key] = [];
      acc[key].push(play);
      return acc;
    }, {});

    // Sort keys by the most recent play within each group
    return Object.fromEntries(
      Object.entries(grouped).sort(([, playsA], [, playsB]) => {
        const latestA = new Date(playsA[0].playedAt).getTime();
        const latestB = new Date(playsB[0].playedAt).getTime();
        return latestB - latestA; // newest first
      }),
    );
  });

  // Handle item click to open modal
  const handleItemClick = (play: PlayHistory) => {
    // Find the full release object from the releases array
    const fullRelease = releases().find((r) => r.id === play.release.id);
    if (fullRelease) {
      setSelectedRelease(fullRelease);
      setIsModalOpen(true);
    }
  };

  // Handle modal close
  const handleCloseModal = () => {
    setIsModalOpen(false);
  };

  return (
    <div class={styles.container}>
      <h1 class={styles.title}>Play History</h1>

      <div class={styles.filters}>
        <div class={styles.filterGroup}>
          <label class={styles.label}>Time Period:</label>
          <select
            class={styles.select}
            value={timeFilter()}
            onInput={(e) => setTimeFilter(e.target.value)}
          >
            <option value="all">All Time</option>
            <option value="week">Last Week</option>
            <option value="month">Last Month</option>
            <option value="year">Last Year</option>
          </select>
        </div>

        <div class={styles.filterGroup}>
          <label class={styles.label}>Group By:</label>
          <select
            class={styles.select}
            value={groupBy()}
            onInput={(e) => setGroupBy(e.target.value)}
          >
            <option value="none">None</option>
            <option value="date">Date</option>
            <option value="artist">Artist</option>
            <option value="album">Album</option>
          </select>
        </div>

        <div class={styles.searchBox}>
          <input
            type="text"
            class={styles.searchInput}
            placeholder="Search by artist, album, stylus or notes..."
            value={searchTerm()}
            onInput={(e) => setSearchTerm(e.target.value)}
          />
        </div>
      </div>

      <Show when={filteredHistory().length === 0}>
        <div class={styles.noResults}>
          <p>No play history found for the selected filters.</p>
        </div>
      </Show>

      <div class={styles.historyList}>
        <For each={Object.entries(groupedHistory())}>
          {([groupName, plays]: [string, PlayHistory[]]) => (
            <>
              <Show when={groupName && groupBy() !== "none"}>
                <div class={styles.groupHeader}>
                  {groupBy() === "date"
                    ? useFormattedShortDate(groupName)
                    : groupName}
                </div>
              </Show>

              <div style="display: flex; flex-wrap: wrap; gap: 1rem; width: 100%;">
                <For each={plays}>
                  {(play) => (
                    <div
                      class={styles.playItem}
                      onClick={() => handleItemClick(play)}
                    >
                      <div class={styles.albumArt}>
                        {play.release.thumb ? (
                          <img
                            src={play.release.thumb}
                            alt={play.release.title}
                            class={styles.albumImage}
                          />
                        ) : (
                          <div class={styles.noImage}>No Image</div>
                        )}
                      </div>

                      <div class={styles.playDetails}>
                        <h3 class={styles.albumTitle}>{play.release.title}</h3>
                        <p class={styles.artistName}>
                          {play.release.artists
                            .filter((artist) => artist.role !== "Producer")
                            .map((artist) => artist.artist?.name)
                            .join(", ")}
                        </p>

                        <div class={styles.playInfoRow}>
                          <p class={styles.playDate}>
                            Played: {useFormattedShortDate(play.playedAt)}
                          </p>

                          <Show when={play.stylus?.name}>
                            <p class={styles.stylus}>
                              Stylus: {play.stylus.name}
                            </p>
                          </Show>
                        </div>

                        <div class={styles.indicators}>
                          <Show when={play.notes && play.notes.trim() !== ""}>
                            <span
                              class={`${styles.indicator} ${styles.hasNotes}`}
                            >
                              <VsNote size={14} />
                              <span>Notes</span>
                            </span>
                          </Show>

                          <Show when={hasCleaning(play.release, play.playedAt)}>
                            <span
                              class={`${styles.indicator} ${styles.hasCleaning}`}
                            >
                              <TbWashTemperature5 size={14} />
                              <span>Cleaned</span>
                            </span>
                          </Show>
                        </div>
                      </div>
                    </div>
                  )}
                </For>
              </div>
            </>
          )}
        </For>
      </div>

      {/* Record Action Modal */}
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

export default PlayHistoryPage;
