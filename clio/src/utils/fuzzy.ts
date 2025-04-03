import Fuse from "fuse.js";
import { Release } from "../types";

// Default options for Fuse.js
const defaultOptions = {
  // Keys to search within - these represent the properties in the Release objects
  keys: [
    { name: "title", weight: 2 }, // Title has higher weight
    { name: "artists.artist.name", weight: 1.5 }, // Artist name has medium-high weight
    { name: "genres.name", weight: 1 }, // Genre name has normal weight
    { name: "labels.label.name", weight: 0.8 }, // Label has lower weight
  ],
  // Fuzzy matching settings
  isCaseSensitive: false,
  includeScore: true,
  shouldSort: true,
  threshold: 0.4, // Lower means more strict matching (0 = exact match only)
  distance: 100, // Maximum distance for fuzzy matching
  minMatchCharLength: 2,
};

/**
 * Create a new Fuse instance for fuzzy searching releases
 * @param releases The array of releases to search within
 * @param options Custom Fuse.js options to override defaults
 * @returns A configured Fuse instance
 */
export const createFuseInstance = (releases: Release[], options = {}) => {
  return new Fuse(releases, { ...defaultOptions, ...options });
};

/**
 * Perform a fuzzy search on releases
 * @param releases The array of releases to search within
 * @param searchTerm The search term to look for
 * @param options Custom Fuse.js options to override defaults
 * @returns Filtered array of releases that match the search term
 */
export const fuzzySearchReleases = (
  releases: Release[],
  searchTerm: string,
  options = {},
): Release[] => {
  if (!searchTerm.trim()) {
    return releases;
  }

  const fuse = createFuseInstance(releases, options);
  const results = fuse.search(searchTerm);

  // Return just the items (not the score)
  return results.map((result) => result.item);
};

/**
 * Custom search function that combines filtering and fuzzy search
 * @param releases The array of releases to search within
 * @param searchTerm The search term to look for
 * @param filterFn Optional filter function to apply before fuzzy search
 * @returns Filtered array of releases that match both filter and search
 */
export const customSearchReleases = (
  releases: Release[],
  searchTerm: string,
  filterFn?: (release: Release) => boolean,
): Release[] => {
  // Apply filter function if provided
  const filteredReleases = filterFn ? releases.filter(filterFn) : releases;

  // If search term is empty, return the filtered results
  if (!searchTerm.trim()) {
    return filteredReleases;
  }

  // Apply fuzzy search on the filtered releases
  return fuzzySearchReleases(filteredReleases, searchTerm);
};
