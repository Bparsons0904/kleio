// src/components/RecordHistoryItem/RecordHistoryItem.tsx
import { Component, Match, Show, Switch, createSignal } from "solid-js";
import styles from "./RecordHistoryItem.module.scss";
import NotesViewPanel from "../NotesViewPanel/NotesViewPanel";
import { BiSolidEdit } from "solid-icons/bi";
import { VsNote } from "solid-icons/vs";
import { FaSolidTrash } from "solid-icons/fa";
import { BsVinylFill } from "solid-icons/bs";
import { TbWashTemperature5 } from "solid-icons/tb";
import {
  deletePlayHistory,
  deleteCleaningHistory,
} from "../../utils/mutations/delete";
import { useAppContext } from "../../provider/Provider";
import ConfirmationModal from "../ConfirmationModal/ConfirmationModal";

export interface HistoryItemProps {
  id: number;
  releaseId: number;
  type: "play" | "cleaning";
  date: string;
  notes?: string;
  stylus?: string;
  stylusId?: number;
  onEdit: (
    id: number,
    type: "play" | "cleaning",
    releaseId: number,
    stylusId?: number,
  ) => void;
}

const RecordHistoryItem: Component<HistoryItemProps> = (props) => {
  const { showSuccess, showError, setKleioStore } = useAppContext();
  const [isNotesPanelOpen, setIsNotesPanelOpen] = createSignal(false);
  const [isDeleteConfirmOpen, setIsDeleteConfirmOpen] = createSignal(false);

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  const openNotesPanel = (e: Event) => {
    e.stopPropagation();
    setIsNotesPanelOpen(true);
  };

  const handleEdit = (e: Event) => {
    e.stopPropagation();
    props.onEdit(props.id, props.type, props.releaseId, props.stylusId);
  };

  const handleDelete = async () => {
    try {
      let result;

      if (props.type === "play") {
        result = await deletePlayHistory(props.id);
      } else {
        result = await deleteCleaningHistory(props.id);
      }

      if (result.success) {
        showSuccess(
          `${props.type === "play" ? "Play" : "Cleaning"} record deleted successfully`,
        );
        setKleioStore(result.data);
      } else {
        throw new Error(`Failed to delete ${props.type} record`);
      }
    } catch (error) {
      console.error("Error deleting record:", error);
      showError(`Failed to delete ${props.type} record`);
    } finally {
      setIsDeleteConfirmOpen(false);
    }
  };

  return (
    <>
      <div
        class={`${styles.historyItem} ${props.type === "play" ? styles.playItem : styles.cleaningItem}`}
      >
        <div class={styles.historyItemContent}>
          <div class={styles.historyItemHeader}>
            <div class={styles.typeAndNotes}>
              <span class={styles.historyItemType}>
                <Switch>
                  <Match when={props.type === "play"}>
                    <span class={styles.historyItems}>
                      <BsVinylFill size={18} /> Played
                    </span>
                  </Match>
                  <Match when={props.type === "cleaning"}>
                    <span class={styles.historyItems}>
                      <TbWashTemperature5 size={20} /> Cleaned
                    </span>
                  </Match>
                </Switch>
              </span>

              <Show when={props.stylus}>
                <div class={styles.historyItemStylus}>
                  Stylus: {props.stylus}
                </div>
              </Show>
              <Show when={props.notes}>
                <button
                  class={styles.noteButton}
                  onClick={openNotesPanel}
                  title="View notes"
                >
                  <VsNote class={styles.noteIcon} size={18} />
                </button>
              </Show>
            </div>

            <div class={styles.dateAndActions}>
              <div class={styles.actionIcons}>
                <button
                  class={styles.editButton}
                  onClick={handleEdit}
                  title="Edit"
                >
                  <BiSolidEdit size={16} />
                </button>
                <button
                  class={styles.deleteButton}
                  onClick={(e) => {
                    e.stopPropagation();
                    setIsDeleteConfirmOpen(true);
                  }}
                  title="Delete"
                >
                  <FaSolidTrash size={16} />
                </button>
              </div>
              <span class={styles.historyItemDate}>
                {formatDate(props.date)}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Notes Panel */}
      <Show when={props.notes}>
        <NotesViewPanel
          isOpen={isNotesPanelOpen()}
          onClose={() => setIsNotesPanelOpen(false)}
          notes={props.notes || ""}
          title={`${props.type === "play" ? "Play" : "Cleaning"} Record Details`}
          date={props.date}
          stylus={props.stylus}
        />
      </Show>

      {/* Delete Confirmation Modal */}
      <ConfirmationModal
        isOpen={isDeleteConfirmOpen()}
        title="Confirm Delete"
        message={`Are you sure you want to delete this ${props.type === "play" ? "play" : "cleaning"} record? This action cannot be undone.`}
        confirmText="Delete"
        onConfirm={handleDelete}
        onCancel={() => setIsDeleteConfirmOpen(false)}
      />
    </>
  );
};

export default RecordHistoryItem;
