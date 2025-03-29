import { Component, Show, createSignal } from "solid-js";
import styles from "./RecordActionModal.module.scss";
import { Release, Stylus } from "../../types";
import RecordHistoryItem from "../RecordHistoryItem/RecordHistoryItem";
import ConfirmationModal from "../ConfirmationModal/ConfirmationModal";
import EditHistoryPanel from "../EditHistoryPanel/EditHistoryPanel";
import {
  deletePlayHistory,
  updatePlayHistory,
  deleteCleaningHistory,
  updateCleaningHistory,
} from "../../utils/api";
import { useAppContext } from "../../provider/Provider";

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
  const { showSuccess, showError, setKleioStore } = useAppContext();

  const [isDeleteConfirmOpen, setIsDeleteConfirmOpen] = createSignal(false);
  const [deleteItemId, setDeleteItemId] = createSignal<number | null>(null);
  const [deleteItemType, setDeleteItemType] = createSignal<
    "play" | "cleaning" | null
  >(null);

  const [isEditPanelOpen, setIsEditPanelOpen] = createSignal(false);
  const [editItem, setEditItem] = createSignal<{
    id: number;
    type: "play" | "cleaning";
    date: string;
    notes?: string;
    stylusId?: number;
  } | null>(null);

  // Get only active styluses for the dropdown
  const activeStyluses = () => {
    return (
      props.styluses?.filter((stylus) => stylus.active || stylus.primary) || []
    );
  };

  // Format date helper
  // const formatDisplayDate = (dateString: string) => {
  //   const date = new Date(dateString);
  //   return date.toLocaleDateString("en-US", {
  //     year: "numeric",
  //     month: "short",
  //     day: "numeric",
  //   });
  // };

  const handleEdit = (id: number, type: "play" | "cleaning") => {
    if (type === "play") {
      const playItem = props.release.playHistory.find((item) => item.id === id);
      if (playItem) {
        setEditItem({
          id,
          type,
          date: playItem.playedAt,
          notes: playItem.notes,
          stylusId: playItem.stylusId,
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
          date: cleaningItem.cleanedAt,
          notes: cleaningItem.notes,
        });
        setIsEditPanelOpen(true);
      }
    }
  };

  const handleDelete = (id: number, type: "play" | "cleaning") => {
    setDeleteItemId(id);
    setDeleteItemType(type);
    setIsDeleteConfirmOpen(true);
  };

  const confirmDelete = async () => {
    if (!deleteItemId() || !deleteItemType()) return;

    try {
      if (deleteItemType() === "play") {
        const response = await deletePlayHistory(deleteItemId()!);
        if (response.status === 200) {
          showSuccess("Play record deleted successfully");
          if (response.data) {
            setKleioStore(response.data);
          }
        }
      } else {
        const response = await deleteCleaningHistory(deleteItemId()!);
        if (response.status === 200) {
          showSuccess("Cleaning record deleted successfully");
          if (response.data) {
            setKleioStore(response.data);
          }
        }
      }
    } catch (error) {
      console.error("Error deleting record:", error);
      showError("Failed to delete record");
    } finally {
      setIsDeleteConfirmOpen(false);
      setDeleteItemId(null);
      setDeleteItemType(null);
    }
  };

  const handleSaveEdit = async (data: {
    id: number;
    type: "play" | "cleaning";
    date: string;
    notes: string;
    stylusId?: number;
  }) => {
    try {
      if (data.type === "play") {
        const response = await updatePlayHistory(data.id, {
          releaseId: props.release.id,
          playedAt: data.date,
          stylusId: data.stylusId,
          notes: data.notes,
        });

        if (response.status === 200) {
          showSuccess("Play record updated successfully");
          if (response.data) {
            setKleioStore(response.data);
          }
        }
      } else {
        const response = await updateCleaningHistory(data.id, {
          releaseId: props.release.id,
          cleanedAt: data.date,
          notes: data.notes,
        });

        if (response.status === 200) {
          showSuccess("Cleaning record updated successfully");
          if (response.data) {
            setKleioStore(response.data);
          }
        }
      }
    } catch (error) {
      console.error("Error updating record:", error);
      showError("Failed to update record");
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

            <div class={styles.historyList}>
              {/* Combine play and cleaning history, sort by date */}
              {[
                ...(props.release.playHistory || []).map((play) => ({
                  id: play.id,
                  type: "play" as const,
                  date: new Date(play.playedAt),
                  notes: play.notes || "",
                  stylus: play.stylus?.name,
                  stylusId: play.stylusId,
                })),
                ...(props.release.cleaningHistory || []).map((cleaning) => ({
                  id: cleaning.id,
                  type: "cleaning" as const,
                  date: new Date(cleaning.cleanedAt),
                  notes: cleaning.notes || "",
                })),
              ]
                .sort((a, b) => b.date.getTime() - a.date.getTime())
                .slice(0, 10) // Show only the 10 most recent activities
                .map((item) => (
                  <RecordHistoryItem
                    id={item.id}
                    type={item.type}
                    date={item.date.toISOString()}
                    notes={item.notes}
                    // @ts-expect-error Stylus doesn't exist on cleaning
                    stylus={item.stylus}
                    onEdit={handleEdit}
                    onDelete={handleDelete}
                  />
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

      {/* Confirmation modal for delete */}
      <ConfirmationModal
        isOpen={isDeleteConfirmOpen()}
        title="Confirm Delete"
        message={`Are you sure you want to delete this ${deleteItemType() === "play" ? "play" : "cleaning"} record? This action cannot be undone.`}
        confirmText="Delete"
        onConfirm={confirmDelete}
        onCancel={() => setIsDeleteConfirmOpen(false)}
        isDestructive={true}
      />

      {/* Edit panel */}
      <Show when={editItem()}>
        <EditHistoryPanel
          isOpen={isEditPanelOpen()}
          onClose={() => setIsEditPanelOpen(false)}
          id={editItem()!.id}
          type={editItem()!.type}
          date={editItem()!.date}
          notes={editItem()!.notes}
          stylusId={editItem()!.stylusId}
          styluses={props.styluses}
          onSave={handleSaveEdit}
        />
      </Show>
    </Show>
  );
};

export default RecordActionModal;
