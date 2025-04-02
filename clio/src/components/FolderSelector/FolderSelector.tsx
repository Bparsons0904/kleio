// src/components/FolderSelector/FolderSelector.tsx (updated component)
import {
  Component,
  createEffect,
  createSignal,
  For,
  onCleanup,
  Show,
} from "solid-js";
import { useAppContext } from "../../provider/Provider";
import styles from "./FolderSelector.module.scss";

const FolderSelector: Component = () => {
  const { folders, selectedFolderId, setSelectedFolderId } = useAppContext();
  const [isOpen, setIsOpen] = createSignal(false);

  const toggleDropdown = () => {
    setIsOpen(!isOpen());
  };

  const handleFolderSelect = (folderId: number) => {
    setSelectedFolderId(folderId);
    setIsOpen(false);
  };

  // Get name of currently selected folder
  const selectedFolderName = () => {
    const folder = folders().find((f) => f.id === selectedFolderId());
    return folder ? folder.name : "All";
  };

  // Close dropdown when clicking outside
  const handleClickOutside = (e: MouseEvent) => {
    const selector = document.querySelector(`.${styles.folderSelector}`);
    if (selector && !selector.contains(e.target as Node)) {
      setIsOpen(false);
    }
  };

  // Add event listener when dropdown is open
  createEffect(() => {
    if (isOpen()) {
      document.addEventListener("mousedown", handleClickOutside);
    } else {
      document.removeEventListener("mousedown", handleClickOutside);
    }

    // Cleanup on unmount
    onCleanup(() => {
      document.removeEventListener("mousedown", handleClickOutside);
    });
  });

  return (
    <div class={styles.folderSelector}>
      <button class={styles.selectorButton} onClick={toggleDropdown}>
        <span>Folder: {selectedFolderName()}</span>
        <span class={styles.arrow}></span>
      </button>

      <Show when={isOpen()}>
        <div class={styles.dropdown}>
          {/* Add an "All" option if not already in folders */}
          <Show when={!folders().some((f) => f.id === 0)}>
            <div
              class={`${styles.option} ${selectedFolderId() === 0 ? styles.selected : ""}`}
              onClick={() => handleFolderSelect(0)}
            >
              All Folders
            </div>
          </Show>

          <For each={folders()}>
            {(folder) => (
              <div
                class={`${styles.option} ${folder.id === selectedFolderId() ? styles.selected : ""}`}
                onClick={() => handleFolderSelect(folder.id)}
              >
                {folder.name} {folder.count ? `(${folder.count})` : ""}
              </div>
            )}
          </For>
        </div>
      </Show>
    </div>
  );
};

export default FolderSelector;
