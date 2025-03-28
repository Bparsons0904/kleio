import { Component, Show } from "solid-js";
import styles from "./RecordActionModal.module.scss";
import { Release, Stylus } from "../../types";

interface RecordActionModalProps {
  isOpen: boolean;
  onClose: () => void;
  release: Release;
  date: string;
  setDate: (date: string) => void;
  selectedStylus: Stylus | null;
  setSelectedStylus: (stylus: Stylus | null) => void;
  styluses: Stylus[];
  notes: string;
  setNotes: (notes: string) => void;
  onLogPlay: () => void;
  onLogCleaning: () => void;
  onLogBoth: () => void;
  isSubmittingPlay: boolean;
  isSubmittingCleaning: boolean;
  isSubmittingBoth: boolean;
}

const RecordActionModal: Component<RecordActionModalProps> = (props) => {
  // Get only active styluses for the dropdown
  const activeStyluses = () => {
    return (
      props.styluses?.filter((stylus) => stylus.active || stylus.primary) || []
    );
  };

  // Format date helper
  const formatDisplayDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  return (
    <Show when={props.isOpen}>
      <div class={styles.modalOverlay} onClick={props.onClose}>
        <div class={styles.modal} onClick={(e) => e.stopPropagation()}>
          <div class={styles.modalHeader}>
            <button class={styles.closeButton} onClick={props.onClose}>
              ×
            </button>
            <h2 class={styles.modalTitle}>Record Actions</h2>
          </div>

          <div class={styles.recordDetails}>
            <div class={styles.recordImage}>
              {props.release.thumb ? (
                <img src={props.release.thumb} alt={props.release.title} />
              ) : (
                <div class={styles.noImage}>No Image</div>
              )}
            </div>
            <div class={styles.recordInfo}>
              <h3 class={styles.recordTitle}>{props.release.title}</h3>
              <p class={styles.recordArtist}>
                {props.release.artists
                  .filter((artist) => artist.role !== "Producer")
                  .map((artist) => artist.artist?.name)
                  .join(", ")}
              </p>
              {props.release.year && (
                <p class={styles.recordYear}>{props.release.year}</p>
              )}
            </div>
          </div>

          <div class={styles.formSection}>
            <div class={styles.formRow}>
              <div class={styles.formGroup}>
                <label class={styles.label} for="actionDate">
                  Date
                </label>
                <input
                  type="date"
                  id="actionDate"
                  class={styles.dateInput}
                  value={props.date}
                  onInput={(e) => props.setDate(e.target.value)}
                />
              </div>

              <div class={styles.formGroup}>
                <label class={styles.label} for="stylusSelect">
                  Stylus Used
                </label>
                <select
                  id="stylusSelect"
                  class={styles.select}
                  value={props.selectedStylus?.id || ""}
                  onChange={(e) => {
                    const id = parseInt(e.target.value);
                    const stylus = props.styluses.find((s) => s.id === id);
                    props.setSelectedStylus(stylus || null);
                  }}
                >
                  <option value="">None</option>
                  {activeStyluses().map((stylus) => (
                    <option value={stylus.id}>
                      {stylus.name} {stylus.primary ? "(Primary)" : ""}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            <div class={styles.formGroup}>
              <label class={styles.label} for="notes">
                Notes
              </label>
              <textarea
                id="notes"
                class={styles.textarea}
                value={props.notes}
                onInput={(e) => props.setNotes(e.target.value)}
                placeholder="Enter any notes about this play or cleaning..."
                rows="3"
              />
            </div>
          </div>

          <div class={styles.actionButtons}>
            <button
              class={styles.playButton}
              onClick={props.onLogPlay}
              disabled={props.isSubmittingPlay}
            >
              {props.isSubmittingPlay ? "Logging..." : "Log Play"}
            </button>
            <button
              class={styles.bothButton}
              onClick={props.onLogBoth}
              disabled={props.isSubmittingBoth}
            >
              {props.isSubmittingBoth ? "Logging..." : "Log Both"}
            </button>
            <button
              class={styles.cleaningButton}
              onClick={props.onLogCleaning}
              disabled={props.isSubmittingCleaning}
            >
              {props.isSubmittingCleaning ? "Logging..." : "Log Cleaning"}
            </button>
          </div>

          <div class={styles.historySection}>
            <h3 class={styles.historyTitle}>Record History</h3>
            {/* <div class={styles.historyTabs}> */}
            {/*   <button class={`${styles.historyTab} ${styles.active}`}> */}
            {/*     All */}
            {/*   </button> */}
            {/*   <button class={styles.historyTab}>Plays</button> */}
            {/*   <button class={styles.historyTab}>Cleanings</button> */}
            {/* </div> */}

            <div class={styles.historyList}>
              {/* Show plays and cleanings in chronological order */}
              {[
                ...(props.release.playHistory || []).map((play) => ({
                  type: "play",
                  date: new Date(play.playedAt),
                  notes: play.notes || "",
                  stylus: play.stylus?.name,
                })),
                ...(props.release.cleaningHistory || []).map((cleaning) => ({
                  type: "cleaning",
                  date: new Date(cleaning.cleanedAt),
                  notes: cleaning.notes || "",
                })),
              ]
                .sort((a, b) => b.date.getTime() - a.date.getTime())
                .slice(0, 10) // Show only the 10 most recent activities
                .map((item) => (
                  <div
                    class={`${styles.historyItem} ${item.type === "play" ? styles.playItem : styles.cleaningItem}`}
                  >
                    <div class={styles.historyItemHeader}>
                      <span class={styles.historyItemType}>
                        {item.type === "play" ? "► Played" : "✓ Cleaned"}
                      </span>
                      <span class={styles.historyItemDate}>
                        {formatDisplayDate(item.date.toISOString())}
                      </span>
                    </div>
                    {/* @ts-expect-error Cleaning history doesn't have a stylus, but is not a problem */}
                    {item.stylus && (
                      <div class={styles.historyItemStylus}>
                        {/* @ts-expect-error Cleaning history doesn't have a stylus, but is not a problem */}
                        Stylus: {item.stylus}
                      </div>
                    )}
                    {item.notes && (
                      <div class={styles.historyItemNotes}>
                        Notes: {item.notes}
                      </div>
                    )}
                  </div>
                ))}

              {!props.release.playHistory?.length &&
                !props.release.cleaningHistory?.length && (
                  <div class={styles.noHistory}>
                    No play or cleaning history for this record yet.
                  </div>
                )}
            </div>
          </div>
        </div>
      </div>
    </Show>
  );
};

export default RecordActionModal;
