import { Component, Show } from "solid-js";
import styles from "./NotesViewPanel.module.scss";
import { X, Calendar, Headphones } from "lucide-solid";

export interface NotesViewPanelProps {
  isOpen: boolean;
  onClose: () => void;
  notes: string;
  title: string;
  date: string;
  stylus?: string;
}

const NotesViewPanel: Component<NotesViewPanelProps> = (props) => {
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  return (
    <div class={`${styles.panelWrapper} ${props.isOpen ? styles.open : ""}`}>
      <div class={styles.overlay} onClick={props.onClose}></div>

      <div class={styles.panel}>
        <div class={styles.panelHeader}>
          <h2 class={styles.panelTitle}>{props.title}</h2>
          <button class={styles.closeButton} onClick={props.onClose}>
            <X size={20} />
          </button>
        </div>

        <div class={styles.panelBody}>
          <div class={styles.metadata}>
            <div class={styles.metadataItem}>
              <Calendar size={18} />
              <span>{formatDate(props.date)}</span>
            </div>

            <Show when={props.stylus}>
              <div class={styles.metadataItem}>
                <Headphones size={18} />
                <span>{props.stylus}</span>
              </div>
            </Show>
          </div>

          <div class={styles.notesSection}>
            <h3 class={styles.notesSectionTitle}>Notes</h3>
            <div class={styles.notesContent}>{props.notes}</div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default NotesViewPanel;
