// src/components/EditHistoryPanel/EditHistoryPanel.tsx
import { Component, createEffect, createSignal, Show } from "solid-js";
import styles from "./EditHistoryPanel.module.scss";
import { EditItem, Stylus } from "../../types";
import { AiOutlineClose } from "solid-icons/ai";
import {
  updatePlayHistory,
  updateCleaningHistory,
} from "../../utils/mutations/put";
import { formatDateForInput } from "../../utils/dates";
import { useAppContext } from "../../provider/Provider";

export interface EditHistoryPanelProps {
  isOpen: boolean;
  onClose: () => void;
  editItem: EditItem | null;
  styluses: Stylus[];
}

const EditHistoryPanel: Component<EditHistoryPanelProps> = (props) => {
  const { showSuccess, showError, setKleioStore } = useAppContext();

  const [date, setDate] = createSignal(
    formatDateForInput(props.editItem?.date) || "",
  );
  const [notes, setNotes] = createSignal(props.editItem?.notes || "");
  const [stylusId, setStylusId] = createSignal(props.editItem?.stylusId);
  const [isSubmitting, setIsSubmitting] = createSignal(false);

  createEffect(() => {
    if (props.editItem) {
      setDate(formatDateForInput(props.editItem.date));
      setNotes(props.editItem.notes || "");
      setStylusId(props.editItem.stylusId);
    }
  });

  const handleSave = async () => {
    if (!props.editItem) return;

    setIsSubmitting(true);

    try {
      if (props.editItem.type === "play") {
        const result = await updatePlayHistory(props.editItem.id, {
          releaseId: props.editItem.releaseId,
          playedAt: new Date(date()).toISOString(),
          stylusId: stylusId(),
          notes: notes(),
        });

        if (result.success) {
          showSuccess("Play record updated successfully");
          setKleioStore(result.data);
          props.onClose();
        } else {
          throw new Error("Failed to update play record");
        }
      } else {
        const result = await updateCleaningHistory(props.editItem.id, {
          releaseId: props.editItem.releaseId,
          cleanedAt: new Date(date()).toISOString(),
          notes: notes(),
        });

        if (result.success) {
          showSuccess("Cleaning record updated successfully");
          setKleioStore(result.data);
          props.onClose();
        } else {
          throw new Error("Failed to update cleaning record");
        }
      }
    } catch (error) {
      console.error("Error updating record:", error);
      showError("Failed to update record");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div class={`${styles.panelWrapper} ${props.isOpen ? styles.open : ""}`}>
      <div class={styles.overlay} onClick={props.onClose}></div>

      <div class={styles.panel}>
        <div class={styles.panelHeader}>
          <h2 class={styles.panelTitle}>
            Edit {props.editItem?.type === "play" ? "Play" : "Cleaning"} Record
          </h2>
          <button class={styles.closeButton} onClick={props.onClose}>
            <AiOutlineClose size={20} />
          </button>
        </div>

        <div class={styles.panelBody}>
          <div class={styles.formGroup}>
            <label class={styles.label} for="editDate">
              Date
            </label>
            <input
              type="date"
              id="editDate"
              class={styles.input}
              value={date()}
              onInput={(e) => setDate(e.target.value)}
            />
          </div>

          <Show when={props.editItem?.type === "play"}>
            <div class={styles.formGroup}>
              <label class={styles.label} for="editStylus">
                Stylus
              </label>
              <select
                id="editStylus"
                class={styles.select}
                value={stylusId()}
                onChange={(e) =>
                  setStylusId(parseInt(e.target.value) || undefined)
                }
              >
                <option value="">None</option>
                {props.styluses
                  .filter((s) => s.active || s.primary)
                  .map((stylus) => (
                    <option value={stylus.id}>
                      {stylus.name} {stylus.primary ? "(Primary)" : ""}
                    </option>
                  ))}
              </select>
            </div>
          </Show>

          <div class={styles.formGroup}>
            <label class={styles.label} for="editNotes">
              Notes
            </label>
            <textarea
              id="editNotes"
              class={styles.textarea}
              value={notes()}
              onInput={(e) => setNotes(e.target.value)}
              rows={4}
              placeholder="Add notes about this record..."
            />
          </div>
        </div>

        <div class={styles.panelFooter}>
          <button
            class={styles.cancelButton}
            onClick={props.onClose}
            disabled={isSubmitting()}
          >
            Cancel
          </button>
          <button
            class={styles.saveButton}
            onClick={handleSave}
            disabled={isSubmitting()}
          >
            {isSubmitting() ? "Saving..." : "Save Changes"}
          </button>
        </div>
      </div>
    </div>
  );
};

export default EditHistoryPanel;
