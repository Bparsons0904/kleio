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
import { Payload, useAppContext } from "../../provider/Provider";
import ConfirmationModal from "../ConfirmationModal/ConfirmationModal";
import { EditItem } from "../../types";
import { useFormattedShortDate } from "../../utils/dates";

export interface HistoryItemProps {
  item: EditItem;
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

  const openNotesPanel = (e: Event) => {
    e.stopPropagation();
    setIsNotesPanelOpen(true);
  };

  const handleEdit = (e: Event) => {
    e.stopPropagation();
    props.onEdit(
      props.item.id,
      props.item.type,
      props.item.releaseId,
      props.item.stylusId,
    );
  };

  const handleDelete = async () => {
    try {
      let result:
        | { success: boolean; data: Payload; error?: undefined }
        | { success: boolean; error: Error; data?: undefined };

      if (props.item.type === "play") {
        result = await deletePlayHistory(props.item.id);
      } else {
        result = await deleteCleaningHistory(props.item.id);
      }

      if (result.success) {
        showSuccess(
          `${props.item.type === "play" ? "Play" : "Cleaning"} record deleted successfully`,
        );
        setKleioStore(result.data);
      } else {
        throw new Error(`Failed to delete ${props.item.type} record`);
      }
    } catch (error) {
      console.error("Error deleting record:", error);
      showError(`Failed to delete ${props.item.type} record`);
    } finally {
      setIsDeleteConfirmOpen(false);
    }
  };

  return (
    <>
      <div
        class={`${styles.historyItem} ${props.item.type === "play" ? styles.playItem : styles.cleaningItem}`}
      >
        <div class={styles.historyItemContent}>
          <div class={styles.historyItemHeader}>
            <div class={styles.typeAndNotes}>
              <span class={styles.historyItemType}>
                <Switch>
                  <Match when={props.item.type === "play"}>
                    <span class={styles.historyItems}>
                      <BsVinylFill size={18} /> Played
                    </span>
                  </Match>
                  <Match when={props.item.type === "cleaning"}>
                    <span class={styles.historyItems}>
                      <TbWashTemperature5 size={20} /> Cleaned
                    </span>
                  </Match>
                </Switch>
              </span>

              <Show when={props.item.stylus}>
                <div class={styles.historyItemStylus}>
                  Stylus: {props.item.stylus}
                </div>
              </Show>
              <Show when={props.item.notes}>
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
                {useFormattedShortDate(props.item.date.toISOString())}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Notes Panel */}
      <Show when={props.item.notes}>
        <NotesViewPanel
          isOpen={isNotesPanelOpen()}
          onClose={() => setIsNotesPanelOpen(false)}
          title={`${props.item.type === "play" ? "Play" : "Cleaning"} Record Details`}
          item={props.item}
        />
      </Show>

      {/* Delete Confirmation Modal */}
      <ConfirmationModal
        isOpen={isDeleteConfirmOpen()}
        title="Confirm Delete"
        message={`Are you sure you want to delete this ${props.item.type === "play" ? "play" : "cleaning"} record? This action cannot be undone.`}
        confirmText="Delete"
        onConfirm={handleDelete}
        onCancel={() => setIsDeleteConfirmOpen(false)}
      />
    </>
  );
};

export default RecordHistoryItem;
