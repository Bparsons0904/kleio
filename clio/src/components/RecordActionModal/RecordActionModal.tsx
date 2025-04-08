// src/components/RecordActionModal/RecordActionModal.tsx
import { Component, Show, createSignal } from "solid-js";
import styles from "./RecordActionModal.module.scss";
import { EditItem, Release, Stylus } from "../../types";
import RecordHistoryItem from "../RecordHistoryItem/RecordHistoryItem";
import EditHistoryPanel from "../EditHistoryPanel/EditHistoryPanel";
import {
  createPlayHistory,
  createCleaningHistory,
  createPlayAndCleaning,
} from "../../utils/mutations/post";
import { useAppContext } from "../../provider/Provider";
import { formatDateForInput } from "../../utils/dates";

interface RecordActionModalProps {
  isOpen: boolean;
  onClose: () => void;
  release: Release;
}

const RecordActionModal: Component<RecordActionModalProps> = (props) => {
  const { styluses, showSuccess, showError, setKleioStore } = useAppContext();

  const [date, setDate] = createSignal(formatDateForInput(new Date()));
  const [selectedStylus, setSelectedStylus] = createSignal<Stylus | null>(
    styluses()?.find((stylus) => stylus.primary),
  );
  const [notes, setNotes] = createSignal("");

  const [isEditPanelOpen, setIsEditPanelOpen] = createSignal(false);
  const [editItem, setEditItem] = createSignal<EditItem | null>(null);

  // Get only active styluses for the dropdown
  const activeStyluses = () => {
    return (
      styluses()?.filter((stylus) => stylus.active || stylus.primary) || []
    );
  };

  const handleLogPlay = async () => {
    try {
      // Use the date from date picker, but add current time
      const playDateTime = new Date(
        `${date()}T${new Date().toTimeString().slice(0, 8)}`,
      );

      const result = await createPlayHistory({
        releaseId: props.release.id,
        playedAt: playDateTime.toISOString(),
        stylusId: selectedStylus()?.id || null,
        notes: notes(),
      });

      if (result.success) {
        showSuccess("Play logged successfully!");
        setNotes("");
        setKleioStore(result.data);
      } else {
        throw new Error("Failed to log play");
      }
    } catch (error) {
      console.error("Error logging play:", error);
      showError("Failed to log play. Please try again.");
    }
  };

  const handleLogCleaning = async () => {
    try {
      // Use the date from date picker, but add current time
      const cleaningDateTime = new Date(
        `${date()}T${new Date().toTimeString().slice(0, 8)}`,
      );

      const result = await createCleaningHistory({
        releaseId: props.release.id,
        cleanedAt: cleaningDateTime.toISOString(),
        notes: notes(),
      });

      if (result.success) {
        showSuccess("Cleaning logged successfully!");
        setNotes("");
        setKleioStore(result.data);
      } else {
        throw new Error("Failed to log cleaning");
      }
    } catch (error) {
      console.error("Error logging cleaning:", error);
      showError("Failed to log cleaning. Please try again.");
    }
  };

  const handleLogBoth = async () => {
    try {
      // Create base datetime with date from picker and current time
      const baseDateTime = new Date(
        `${date()}T${new Date().toTimeString().slice(0, 8)}`,
      );

      // Create cleaning timestamp (done first)
      const cleaningTimestamp = baseDateTime.toISOString();

      // Create play timestamp 1 second later
      baseDateTime.setSeconds(baseDateTime.getSeconds() + 1);
      const playTimestamp = baseDateTime.toISOString();

      const result = await createPlayAndCleaning(
        {
          releaseId: props.release.id,
          playedAt: playTimestamp,
          stylusId: selectedStylus()?.id || null,
          notes: notes(),
        },
        {
          releaseId: props.release.id,
          cleanedAt: cleaningTimestamp,
          notes: notes(),
        },
      );

      if (result.success) {
        showSuccess("Both cleaning and play logged successfully!");
        setNotes("");
        setKleioStore(result.data);
      } else {
        throw new Error("Failed to log both activities");
      }
    } catch (error) {
      console.error("Error logging both activities:", error);
      showError("Failed to log both activities. Please try again.");
    }
  };

  const handleEdit = (
    id: number,
    type: "play" | "cleaning",
    releaseId: number,
  ) => {
    if (type === "play") {
      const playItem = props.release.playHistory.find((item) => item.id === id);
      if (playItem) {
        setEditItem({
          id,
          type,
          date: new Date(playItem.playedAt),
          notes: playItem.notes,
          stylusId: playItem.stylusId,
          releaseId: releaseId,
        });
        setIsEditPanelOpen(true);
      }
    } else {
      const cleaningItem = props.release.cleaningHistory.find(
        (item) => item.id === id,
      );
      if (cleaningItem) {
        setEditItem({
          id,
          type,
          date: new Date(cleaningItem.cleanedAt),
          notes: cleaningItem.notes,
          releaseId: releaseId,
        });
        setIsEditPanelOpen(true);
      }
    }
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
                  value={date()}
                  onInput={(e) => setDate(e.target.value)}
                />
              </div>

              <div class={styles.formGroup}>
                <label class={styles.label} for="stylusSelect">
                  Stylus Used
                </label>
                <select
                  id="stylusSelect"
                  class={styles.select}
                  value={selectedStylus()?.id || ""}
                  onChange={(e) => {
                    const id = parseInt(e.target.value);
                    const stylus = styluses().find((s) => s.id === id);
                    setSelectedStylus(stylus || null);
                  }}
                >
                  <option value="">None</option>
                  {activeStyluses().map((stylus) => (
                    <option
                      value={stylus.id}
                      selected={stylus.id === selectedStylus()?.id}
                    >
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
                value={notes()}
                onInput={(e) => setNotes(e.target.value)}
                placeholder="Enter any notes about this play or cleaning..."
                rows="3"
              />
            </div>
          </div>

          <div class={styles.actionButtons}>
            <button class={styles.playButton} onClick={handleLogPlay}>
              Log Play
            </button>
            <button class={styles.bothButton} onClick={handleLogBoth}>
              Log Both
            </button>
            <button class={styles.cleaningButton} onClick={handleLogCleaning}>
              Log Cleaning
            </button>
          </div>

          <div class={styles.historySection}>
            <h3 class={styles.historyTitle}>Record History</h3>

            <div class={styles.historyList}>
              {/* Combine play and cleaning history, sort by date */}
              {[
                ...(props.release.playHistory || []).map((play) => ({
                  id: play.id,
                  releaseId: props.release.id,
                  type: "play" as const,
                  date: new Date(play.playedAt),
                  notes: play.notes || "",
                  stylus: play.stylus?.name,
                  stylusId: play.stylusId,
                })),
                ...(props.release.cleaningHistory || []).map((cleaning) => ({
                  id: cleaning.id,
                  releaseId: props.release.id,
                  type: "cleaning" as const,
                  date: new Date(cleaning.cleanedAt),
                  notes: cleaning.notes || "",
                })),
              ]
                .sort((a, b) => b.date.getTime() - a.date.getTime())
                // .slice(0, 10) // Show only the 10 most recent activities
                .map((item) => (
                  <RecordHistoryItem item={item} onEdit={handleEdit} />
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

      {/* Edit panel */}
      <Show when={editItem()}>
        <EditHistoryPanel
          isOpen={isEditPanelOpen()}
          onClose={() => setIsEditPanelOpen(false)}
          editItem={editItem()}
          styluses={styluses()}
        />
      </Show>
    </Show>
  );
};

export default RecordActionModal;
