import { Component, createSignal, Show } from "solid-js";
import styles from "./EditHistoryPanel.module.scss";
import { X } from "lucide-solid";
import { Stylus } from "../../types";

export interface EditHistoryPanelProps {
  isOpen: boolean;
  onClose: () => void;
  type: "play" | "cleaning";
  id: number;
  date: string;
  notes?: string;
  stylusId?: number;
  styluses: Stylus[];
  onSave: (data: {
    id: number;
    type: "play" | "cleaning";
    date: string;
    notes: string;
    stylusId?: number;
  }) => void;
}

const EditHistoryPanel: Component<EditHistoryPanelProps> = (props) => {
  const [date, setDate] = createSignal(props.date.split("T")[0]);
  const [notes, setNotes] = createSignal(props.notes || "");
  const [stylusId, setStylusId] = createSignal(props.stylusId);
  const [isSubmitting, setIsSubmitting] = createSignal(false);

  const handleSave = async () => {
    setIsSubmitting(true);

    try {
      props.onSave({
        id: props.id,
        type: props.type,
        date: new Date(date()).toISOString(),
        notes: notes(),
        ...(props.type === "play" && { stylusId: stylusId() }),
      });
    } finally {
      setIsSubmitting(false);
      props.onClose();
    }
  };

  return (
    <div class={`${styles.panelWrapper} ${props.isOpen ? styles.open : ""}`}>
      <div class={styles.overlay} onClick={props.onClose}></div>

      <div class={styles.panel}>
        <div class={styles.panelHeader}>
          <h2 class={styles.panelTitle}>
            Edit {props.type === "play" ? "Play" : "Cleaning"} Record
          </h2>
          <button class={styles.closeButton} onClick={props.onClose}>
            <X size={20} />
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

          <Show when={props.type === "play"}>
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
