// src/pages/Collection/Collection.tsx
import { Component, createSignal, createEffect, For, Show } from "solid-js";
import { useAppContext } from "../../provider/Provider";
import styles from "./Collection.module.scss";
import { Release } from "../../types";
import RecordActionModal from "../../components/RecordActionModal/RecordActionModal";
import { AiOutlineSearch } from "solid-icons/ai";
import { BiRegularFilterAlt } from "solid-icons/bi";
import { BsGrid, BsGrid3x3Gap } from "solid-icons/bs";
import { fuzzySearchReleases } from "../../utils/fuzzy";

const Collection: Component = () => {
  const { releases } = useAppContext();

  // State
  const [filteredReleases, setFilteredReleases] = createSignal<Release[]>([]);
  const [searchTerm, setSearchTerm] = createSignal("");
  const [selectedRelease, setSelectedRelease] = createSignal<Release | null>(
    null,
  );
  const [isModalOpen, setIsModalOpen] = createSignal(false);
  const [sortBy, setSortBy] = createSignal("artist");
  const [gridSize, setGridSize] = createSignal<"small" | "medium" | "large">(
    "medium",
  );
  const [showFilters, setShowFilters] = createSignal(false);
  const [genreFilter, setGenreFilter] = createSignal<string[]>([]);

  // Derived state for unique genres
  const availableGenres = () => {
    const genreSet = new Set<string>();
    releases().forEach((release) => {
      release.genres.forEach((genre) => {
        genreSet.add(genre.name);
      });
    });
    return Array.from(genreSet).sort();
  };

  // Filter releases based on search term, genres, and sort them
  createEffect(() => {
    let filtered = releases();

    // Apply genre filter first if any genres are selected
    if (genreFilter().length > 0) {
      filtered = filtered.filter((release) =>
        release.genres.some((genre) => genreFilter().includes(genre.name)),
      );
    }

    // Apply fuzzy search if search term is provided
    if (searchTerm()) {
      filtered = fuzzySearchReleases(filtered, searchTerm());
    }

    // Apply sorting
    filtered = sortReleases(filtered, sortBy());

    setFilteredReleases(filtered);
  });

  // Sort releases based on selected sort option
  const sortReleases = (
    releasesToSort: Release[],
    sortOption: string,
  ): Release[] => {
    switch (sortOption) {
      case "artist":
        // Sort by primary artist
        return [...releasesToSort].sort((a, b) => {
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
        return [...releasesToSort].sort((a, b) =>
          a.title.localeCompare(b.title),
        );

      case "year":
        // Sort by release year (newest first)
        return [...releasesToSort].sort((a, b) => {
          const yearA = a.year || 0;
          const yearB = b.year || 0;
          return yearB - yearA;
        });

      default:
        // Default sort by title
        return [...releasesToSort].sort((a, b) =>
          a.title.localeCompare(b.title),
        );
    }
  };

  const handleReleaseClick = (release: Release) => {
    setSelectedRelease(release);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
  };

  const toggleGenre = (genre: string) => {
    if (genreFilter().includes(genre)) {
      setGenreFilter((prev) => prev.filter((g) => g !== genre));
    } else {
      setGenreFilter((prev) => [...prev, genre]);
    }
  };

  const clearFilters = () => {
    setGenreFilter([]);
    setSearchTerm("");
  };

  return (
    <div class={styles.container}>
      <h1 class={styles.title}>Your Collection</h1>

      <div class={styles.controls}>
        <div class={styles.searchBar}>
          <AiOutlineSearch size={20} class={styles.searchIcon} />
          <input
            type="text"
            placeholder="Search by album or artist..."
            value={searchTerm()}
            onInput={(e) => setSearchTerm(e.target.value)}
            class={styles.searchInput}
          />
        </div>

        <div class={styles.controlButtons}>
          <button
            class={styles.filterButton}
            onClick={() => setShowFilters(!showFilters())}
            classList={{ [styles.active]: showFilters() }}
          >
            <BiRegularFilterAlt size={20} />
            <span>Filter</span>
          </button>

          <div class={styles.sortContainer}>
            <select
              value={sortBy()}
              onInput={(e) => setSortBy(e.target.value)}
              class={styles.sortSelect}
            >
              <option value="album">Album (A-Z)</option>
              <option value="artist">Artist (A-Z)</option>
              <option value="year">Year (newest first)</option>
            </select>
          </div>

          <div class={styles.gridSizeToggle}>
            <button
              onClick={() => setGridSize("small")}
              classList={{ [styles.active]: gridSize() === "small" }}
              title="Small grid"
            >
              <BsGrid3x3Gap size={18} />
            </button>
            <button
              onClick={() => setGridSize("medium")}
              classList={{ [styles.active]: gridSize() === "medium" }}
              title="Medium grid"
            >
              <BsGrid3x3Gap size={22} />
            </button>
            <button
              onClick={() => setGridSize("large")}
              classList={{ [styles.active]: gridSize() === "large" }}
              title="Large grid"
            >
              <BsGrid size={22} />
            </button>
          </div>
        </div>
      </div>

      <Show when={showFilters()}>
        <div class={styles.filterPanel}>
          <div class={styles.filterSection}>
            <h3 class={styles.filterTitle}>Genres</h3>
            <div class={styles.filterOptions}>
              <For each={availableGenres()}>
                {(genre) => (
                  <label class={styles.filterOption}>
                    <input
                      type="checkbox"
                      checked={genreFilter().includes(genre)}
                      onChange={() => toggleGenre(genre)}
                    />
                    <span>{genre}</span>
                  </label>
                )}
              </For>
            </div>
          </div>

          <div class={styles.filterActions}>
            <button class={styles.clearButton} onClick={clearFilters}>
              Clear All Filters
            </button>
          </div>
        </div>
      </Show>

      <Show when={filteredReleases().length === 0}>
        <div class={styles.noResults}>
          <p>No albums match your search or filters.</p>
          <button class={styles.clearButton} onClick={clearFilters}>
            Clear All Filters
          </button>
        </div>
      </Show>

      <div class={`${styles.albumGrid} ${styles[gridSize()]}`}>
        <For each={filteredReleases()}>
          {(release) => (
            <div
              class={styles.albumCard}
              onClick={() => handleReleaseClick(release)}
            >
              <div class={styles.albumArtwork}>
                {release.coverImage ? (
                  <img
                    src={release.coverImage}
                    alt={release.title}
                    class={styles.albumImage}
                  />
                ) : (
                  <div class={styles.noImage}>
                    <span>{release.title}</span>
                  </div>
                )}

                <div class={styles.albumHover}>
                  <div class={styles.trackList}>
                    <h4 class={styles.trackListTitle}>Tracks</h4>
                    <ol class={styles.tracks}>
                      <Show
                        when={release.tracks && release.tracks.length > 0}
                        fallback={
                          <li class={styles.noTracks}>
                            No track data available
                          </li>
                        }
                      >
                        <For each={release.tracks}>
                          {(track) => (
                            <li>
                              <span class={styles.trackPosition}>
                                {track.position}
                              </span>
                              {track.title}
                              <Show when={track.durationText}>
                                <span class={styles.trackDuration}>
                                  {track.durationText}
                                </span>
                              </Show>
                            </li>
                          )}
                        </For>
                      </Show>
                    </ol>
                  </div>
                </div>
              </div>

              <div class={styles.albumInfo}>
                <h3 class={styles.albumTitle}>{release.title}</h3>
                <p class={styles.albumArtist}>
                  {release.artists
                    .filter((artist) => artist.role !== "Producer")
                    .map((artist) => artist.artist?.name)
                    .join(", ")}
                </p>
              </div>
            </div>
          )}
        </For>
      </div>

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

export default Collection;
