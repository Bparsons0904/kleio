import { Component, Show } from "solid-js";
import styles from "./NotesViewPanel.module.scss";
import { AiOutlineClose } from "solid-icons/ai";
import { AiTwotoneCalendar } from "solid-icons/ai";
import { ImHeadphones } from "solid-icons/im";
import { EditItem } from "../../types";

export interface NotesViewPanelProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  item: EditItem;
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
            <AiOutlineClose size={20} />
          </button>
        </div>

        <div class={styles.panelBody}>
          <div class={styles.metadata}>
            <div class={styles.metadataItem}>
              <AiTwotoneCalendar size={18} />
              <span>{formatDate(props.item.date)}</span>
            </div>

            <Show when={props.item.stylus}>
              <div class={styles.metadataItem}>
                <ImHeadphones size={18} />
                <span>{props.item.stylus}</span>
              </div>
            </Show>
          </div>

          <div class={styles.notesSection}>
            <h3 class={styles.notesSectionTitle}>Notes</h3>
            <div class={styles.notesContent}>{props.item.notes}</div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default NotesViewPanel;
