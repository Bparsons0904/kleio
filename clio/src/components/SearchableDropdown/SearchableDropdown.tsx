// src/components/SearchableDropdown/SearchableDropdown.tsx
import { Component, createSignal, For, Show, createEffect } from "solid-js";
import styles from "./SearchableDropdown.module.scss";

export interface DropdownOption {
  value: string;
  label: string;
  disabled?: boolean;
  group?: string;
}

interface SearchableDropdownProps {
  options: DropdownOption[];
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  label?: string;
  isSearchable?: boolean;
}

const SearchableDropdown: Component<SearchableDropdownProps> = (props) => {
  const [isOpen, setIsOpen] = createSignal(false);
  const [searchQuery, setSearchQuery] = createSignal("");
  const [highlightedIndex, setHighlightedIndex] = createSignal(-1);

  // Create a ref for the dropdown element to handle outside clicks
  let dropdownRef: HTMLDivElement | undefined;

  // Selected option label
  const selectedLabel = () => {
    const selectedOption = props.options.find(
      (opt) => opt.value === props.value,
    );
    return selectedOption?.label || props.placeholder || "Select an option";
  };

  // Filter options based on search query
  const filteredOptions = () => {
    const query = searchQuery().toLowerCase();
    if (!query) return props.options;

    return props.options.filter(
      (option) =>
        !option.disabled && option.label.toLowerCase().includes(query),
    );
  };

  // Handle option selection
  const handleSelectOption = (value: string) => {
    props.onChange(value);
    setIsOpen(false);
    setSearchQuery("");
  };

  // Handle click outside to close dropdown
  const handleOutsideClick = (e: MouseEvent) => {
    if (dropdownRef && !dropdownRef.contains(e.target as Node)) {
      setIsOpen(false);
    }
  };

  // Handle keyboard navigation
  const handleKeyDown = (e: KeyboardEvent) => {
    const options = filteredOptions().filter((opt) => !opt.disabled);

    switch (e.key) {
      case "ArrowDown":
        e.preventDefault();
        setHighlightedIndex((prev) =>
          prev < options.length - 1 ? prev + 1 : 0,
        );
        break;
      case "ArrowUp":
        e.preventDefault();
        setHighlightedIndex((prev) =>
          prev > 0 ? prev - 1 : options.length - 1,
        );
        break;
      case "Enter":
        e.preventDefault();
        if (highlightedIndex() >= 0 && options[highlightedIndex()]) {
          handleSelectOption(options[highlightedIndex()].value);
        }
        break;
      case "Escape":
        e.preventDefault();
        setIsOpen(false);
        break;
    }
  };

  // Set up and clean up event listeners
  createEffect(() => {
    if (isOpen()) {
      document.addEventListener("mousedown", handleOutsideClick);
      document.addEventListener("keydown", handleKeyDown);
    } else {
      document.removeEventListener("mousedown", handleOutsideClick);
      document.removeEventListener("keydown", handleKeyDown);
    }

    return () => {
      document.removeEventListener("mousedown", handleOutsideClick);
      document.removeEventListener("keydown", handleKeyDown);
    };
  });

  // Reset search and highlighted index when dropdown opens/closes
  createEffect(() => {
    if (!isOpen()) {
      setSearchQuery("");
      setHighlightedIndex(-1);
    }
  });

  return (
    <div class={styles.dropdownContainer} ref={dropdownRef}>
      {props.label && <label class={styles.label}>{props.label}</label>}

      <div
        class={`${styles.dropdownTrigger} ${isOpen() ? styles.open : ""}`}
        onClick={() => setIsOpen(!isOpen())}
      >
        <span>{selectedLabel()}</span>
        <span class={styles.arrow}></span>
      </div>

      <Show when={isOpen()}>
        <div class={styles.dropdownMenu}>
          <Show when={props.isSearchable}>
            <div class={styles.searchContainer}>
              <input
                type="text"
                class={styles.searchInput}
                placeholder="Search..."
                value={searchQuery()}
                onInput={(e) => setSearchQuery(e.currentTarget.value)}
                // Stop propagation to prevent dropdown from closing
                onClick={(e) => e.stopPropagation()}
                // Initial focus when dropdown opens
                ref={(el) => {
                  if (isOpen()) {
                    el.focus();
                  }
                }}
              />
            </div>
          </Show>

          <div class={styles.optionsList}>
            <Show
              when={filteredOptions().length > 0}
              fallback={<div class={styles.noResults}>No results found</div>}
            >
              <For each={filteredOptions()}>
                {(option, index) => {
                  // Skip rendering for options with the HEADER prefix if search is active
                  if (searchQuery() && option.value.startsWith("HEADER:")) {
                    return null;
                  }

                  const isHeader = option.value.startsWith("HEADER:");
                  const isHighlighted =
                    !isHeader && highlightedIndex() === index();
                  const nonDisabledIndex = filteredOptions()
                    .filter((opt) => !opt.disabled)
                    .findIndex((opt) => opt.value === option.value);

                  return (
                    <div
                      class={`
                        ${styles.option} 
                        ${isHeader ? styles.header : ""}
                        ${isHighlighted ? styles.highlighted : ""}
                        ${option.disabled ? styles.disabled : ""}
                        ${props.value === option.value ? styles.selected : ""}
                      `}
                      onClick={(e) => {
                        e.stopPropagation();
                        if (!isHeader && !option.disabled) {
                          handleSelectOption(option.value);
                        }
                      }}
                      onMouseEnter={() => {
                        if (!isHeader && !option.disabled) {
                          setHighlightedIndex(nonDisabledIndex);
                        }
                      }}
                    >
                      {option.label}
                    </div>
                  );
                }}
              </For>
            </Show>
          </div>
        </div>
      </Show>
    </div>
  );
};

export default SearchableDropdown;
