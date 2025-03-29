import { Component, createSignal, For, Show } from "solid-js";
import { useAppContext } from "../../provider/Provider";
import { Stylus } from "../../types";
import { createStylus, deleteStylus, updateStylus } from "../../utils/api";
import styles from "./StylusManager.module.scss";

const StylusManager: Component = () => {
  const { styluses, setStyluses } = useAppContext();
  const [isAddingStylus, setIsAddingStylus] = createSignal(false);
  const [editingStylus, setEditingStylus] = createSignal<Stylus | null>(null);
  const [isLoading, setIsLoading] = createSignal(false);
  const [successMessage, setSuccessMessage] = createSignal("");
  const [errorMessage, setErrorMessage] = createSignal("");
  const [showArchived] = createSignal(false);
  // Form state for new/editing stylus
  const [name, setName] = createSignal("");
  const [manufacturer, setManufacturer] = createSignal("");
  const [expectedLifespan, setExpectedLifespan] = createSignal<
    number | undefined
  >();
  const [purchaseDate, setPurchaseDate] = createSignal("");
  const [isActive, setIsActive] = createSignal(false);
  const [isPrimary, setIsPrimary] = createSignal(false);
  const [isOwned, setIsOwned] = createSignal(true);

  const resetForm = () => {
    setName("");
    setManufacturer("");
    setExpectedLifespan(undefined);
    setPurchaseDate("");
    setIsActive(false);
    setIsPrimary(false);
    setIsOwned(true);
    setEditingStylus(null);
  };

  const startAddStylus = () => {
    resetForm();
    setIsAddingStylus(true);
  };

  const startEditStylus = (stylus: Stylus) => {
    setName(stylus.name);
    setManufacturer(stylus.manufacturer || "");
    setExpectedLifespan(stylus.expectedLifespan);

    if (stylus.purchaseDate) {
      // Convert to YYYY-MM-DD format for date input
      setPurchaseDate(
        new Date(stylus.purchaseDate).toISOString().split("T")[0],
      );
    } else {
      setPurchaseDate("");
    }

    setIsActive(stylus.active);
    setIsPrimary(stylus.primary);
    setIsOwned(stylus.owned !== false); // Default to true if property doesn't exist
    setEditingStylus(stylus);
    setIsAddingStylus(true);
  };

  const cancelForm = () => {
    setIsAddingStylus(false);
    resetForm();
  };

  const copyFromStylus = (e: Event) => {
    const select = e.target as HTMLSelectElement;
    const stylusId = parseInt(select.value);

    // If user selected the default "Select a stylus" option (value=""), reset form
    if (!stylusId) {
      return;
    }

    // Find the stylus by ID
    const stylus = styluses().find((s) => s.id === stylusId);
    if (!stylus) return;

    setName(stylus.name);
    setManufacturer(stylus.manufacturer || "");
    setExpectedLifespan(stylus.expectedLifespan);

    // Don't copy purchase date as it would be new
    setPurchaseDate("");

    // Set default values for new copy
    setIsActive(true);
    setIsPrimary(true);
    setIsOwned(true);
  };

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    if (!name()) {
      setErrorMessage("Name is required");
      return;
    }
    setIsLoading(true);
    setErrorMessage("");
    setSuccessMessage("");
    try {
      const stylusData = {
        name: name(),
        manufacturer: manufacturer() || undefined,
        expectedLifespan: expectedLifespan(),
        purchaseDate: purchaseDate()
          ? new Date(purchaseDate()).toISOString()
          : undefined,
        active: isActive(),
        primary: isPrimary(),
        owned: isOwned(),
      };

      let response;
      if (editingStylus()) {
        // Update existing stylus
        response = await updateStylus(editingStylus()!.id, {
          ...stylusData,
          id: editingStylus()!.id,
        });
        setSuccessMessage(`Updated stylus "${name()}"`);
      } else {
        // Create new stylus
        response = await createStylus(stylusData);
        setSuccessMessage(`Created new stylus "${name()}"`);
      }

      // Update styluses list with the new data from the response
      if (response?.data) {
        setStyluses(response.data);
      }

      // Reset form and exit edit mode
      setIsAddingStylus(false);
      resetForm();
    } catch (error) {
      console.error("Error saving stylus:", error);
      setErrorMessage("Failed to save stylus. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  const handleDelete = async (stylus: Stylus) => {
    if (!confirm(`Are you sure you want to delete "${stylus.name}"?`)) {
      return;
    }
    setIsLoading(true);
    try {
      const response = await deleteStylus(stylus.id);
      setSuccessMessage(`Deleted stylus "${stylus.name}"`);

      // Update styluses list with the new data from the response
      if (response?.data) {
        setStyluses(response.data);
      }
    } catch (error) {
      console.error("Error deleting stylus:", error);
      setErrorMessage("Failed to delete stylus. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  // Filter styluses based on owned/active status
  const activeStyluses = () =>
    styluses().filter((s) => s.owned !== false && (s.active || s.primary));

  const inactiveStyluses = () =>
    styluses().filter((s) => s.owned !== false && !s.active && !s.primary);

  const archivedStyluses = () => styluses().filter((s) => s.owned === false);

  return (
    <div class={styles.container}>
      <div class={styles.header}>
        <h2 class={styles.title}>Stylus Manager</h2>
        <div class={styles.headerButtons}>
          {/* <button */}
          {/*   class={styles.archiveToggle} */}
          {/*   onClick={() => setShowArchived(!showArchived())} */}
          {/* > */}
          {/*   {showArchived() ? "Hide Archived" : "Show Archived"} */}
          {/* </button> */}
          <Show when={!isAddingStylus()}>
            <button
              class={styles.addButton}
              onClick={startAddStylus}
              disabled={isLoading()}
            >
              Add New Stylus
            </button>
          </Show>
        </div>
      </div>

      <Show when={successMessage()}>
        <div class={styles.successMessage}>{successMessage()}</div>
      </Show>

      <Show when={errorMessage()}>
        <div class={styles.errorMessage}>{errorMessage()}</div>
      </Show>

      <Show when={isAddingStylus()}>
        <div class={styles.formContainer}>
          <h3>{editingStylus() ? "Edit Stylus" : "Add New Stylus"}</h3>

          <Show when={!editingStylus() && styluses().length > 0}>
            <div class={styles.copyFromContainer}>
              <label for="copyFromSelect" class={styles.label}>
                Copy from existing stylus:
              </label>
              <select
                id="copyFromSelect"
                class={styles.select}
                onChange={copyFromStylus}
              >
                <option value="">Select a stylus</option>
                <For
                  each={styluses()
                    .filter((s) => s.baseModel === true)
                    .sort((a, b) => {
                      if (a.primary && !b.primary) return -1;
                      if (!a.primary && b.primary) return 1;
                      return a.name.localeCompare(b.name);
                    })}
                >
                  {(stylus) => (
                    <option value={stylus.id}>
                      {stylus.name}{" "}
                      {stylus.manufacturer ? `(${stylus.manufacturer})` : ""}
                    </option>
                  )}
                </For>
              </select>
            </div>
          </Show>

          <form class={styles.form} onSubmit={handleSubmit}>
            <div class={styles.formGroup}>
              <label for="name" class={styles.label}>
                Name *
              </label>
              <input
                type="text"
                id="name"
                class={styles.input}
                value={name()}
                onInput={(e) => setName(e.target.value)}
                required
              />
            </div>

            <div class={styles.formGroup}>
              <label for="manufacturer" class={styles.label}>
                Manufacturer
              </label>
              <input
                type="text"
                id="manufacturer"
                class={styles.input}
                value={manufacturer()}
                onInput={(e) => setManufacturer(e.target.value)}
              />
            </div>

            <div class={styles.formGroup}>
              <label for="expectedLifespan" class={styles.label}>
                Expected Lifespan (hours)
              </label>
              <input
                type="number"
                id="expectedLifespan"
                class={styles.input}
                value={expectedLifespan() || ""}
                onInput={(e) =>
                  setExpectedLifespan(parseInt(e.target.value) || undefined)
                }
                min="0"
              />
            </div>

            <div class={styles.formGroup}>
              <label for="purchaseDate" class={styles.label}>
                Purchase Date
              </label>
              <input
                type="date"
                id="purchaseDate"
                class={styles.input}
                value={purchaseDate()}
                onInput={(e) => setPurchaseDate(e.target.value)}
              />
            </div>

            {/* <div class={styles.checkboxGroup}> */}
            {/* <label class={styles.checkboxLabel}> */}
            {/*   <input */}
            {/*     type="checkbox" */}
            {/*     checked={isOwned()} */}
            {/*     onChange={(e) => setIsOwned(e.target.checked)} */}
            {/*     hidden */}
            {/*   /> */}
            {/*   Currently Owned */}
            {/* </label> */}
            {/* </div> */}

            <div class={styles.checkboxGroup}>
              <label class={styles.checkboxLabel}>
                <input
                  type="checkbox"
                  checked={isActive()}
                  onChange={(e) => setIsActive(e.target.checked)}
                  disabled={!isOwned()}
                />
                Active (Working Stylus / In Rotation)
              </label>
            </div>

            <div class={styles.checkboxGroup}>
              <label class={styles.checkboxLabel}>
                <input
                  type="checkbox"
                  checked={isPrimary()}
                  onChange={(e) => setIsPrimary(e.target.checked)}
                  disabled={!isOwned()}
                />
                Primary (default for new plays)
              </label>
            </div>

            <div class={styles.formActions}>
              <button
                type="button"
                class={styles.cancelButton}
                onClick={cancelForm}
                disabled={isLoading()}
              >
                Cancel
              </button>
              <button
                type="submit"
                class={styles.submitButton}
                disabled={isLoading()}
              >
                {isLoading()
                  ? "Saving..."
                  : editingStylus()
                    ? "Update"
                    : "Create"}
              </button>
            </div>
          </form>
        </div>
      </Show>

      <Show when={!isAddingStylus() && styluses().length === 0}>
        <p class={styles.noStyluses}>
          No styluses found. Click "Add New Stylus" to create one.
        </p>
      </Show>

      <Show when={!isAddingStylus() && activeStyluses().length > 0}>
        <div class={styles.section}>
          <h3 class={styles.sectionTitle}>Active Styluses</h3>
          <div class={styles.stylusList}>
            <For
              each={activeStyluses().sort((a, b) => {
                if (a.primary && !b.primary) return -1;
                if (!a.primary && b.primary) return 1;
                return a.name.localeCompare(b.name);
              })}
            >
              {(stylus) => (
                <div
                  class={`${styles.stylusCard} ${stylus.primary ? styles.primaryCard : ""}`}
                >
                  <div class={styles.stylusInfo}>
                    <h3 class={styles.stylusName}>
                      {stylus.name}
                      <div class={styles.tagsContainer}>
                        {stylus.primary && (
                          <span class={styles.primaryTag}>Primary</span>
                        )}
                      </div>
                    </h3>

                    <Show when={stylus.manufacturer}>
                      <p class={styles.stylusDetail}>
                        <strong>Manufacturer:</strong> {stylus.manufacturer}
                      </p>
                    </Show>

                    <Show when={stylus.expectedLifespan}>
                      <p class={styles.stylusDetail}>
                        <strong>Expected Lifespan:</strong>{" "}
                        {stylus.expectedLifespan} hours
                      </p>
                    </Show>

                    <Show when={stylus.purchaseDate}>
                      <p class={styles.stylusDetail}>
                        <strong>Purchased:</strong>{" "}
                        {new Date(stylus.purchaseDate!).toLocaleDateString()}
                      </p>
                    </Show>
                  </div>

                  <div class={styles.stylusActions}>
                    <button
                      class={styles.editButton}
                      onClick={() => startEditStylus(stylus)}
                      disabled={isLoading()}
                    >
                      Edit
                    </button>
                    <button
                      class={styles.deleteButton}
                      onClick={() => handleDelete(stylus)}
                      disabled={isLoading()}
                    >
                      Delete
                    </button>
                  </div>
                </div>
              )}
            </For>
          </div>
        </div>
      </Show>

      <Show when={!isAddingStylus() && inactiveStyluses().length > 0}>
        <div class={styles.section}>
          <h3 class={styles.sectionTitle}>Inactive Styluses</h3>
          <div class={styles.stylusList}>
            <For each={inactiveStyluses()}>
              {(stylus) => (
                <div class={styles.stylusCard}>
                  <div class={styles.stylusInfo}>
                    <h3 class={styles.stylusName}>{stylus.name}</h3>

                    <Show when={stylus.manufacturer}>
                      <p class={styles.stylusDetail}>
                        <strong>Manufacturer:</strong> {stylus.manufacturer}
                      </p>
                    </Show>

                    <Show when={stylus.expectedLifespan}>
                      <p class={styles.stylusDetail}>
                        <strong>Expected Lifespan:</strong>{" "}
                        {stylus.expectedLifespan} hours
                      </p>
                    </Show>

                    <Show when={stylus.purchaseDate}>
                      <p class={styles.stylusDetail}>
                        <strong>Purchased:</strong>{" "}
                        {new Date(stylus.purchaseDate!).toLocaleDateString()}
                      </p>
                    </Show>
                  </div>

                  <div class={styles.stylusActions}>
                    <button
                      class={styles.editButton}
                      onClick={() => startEditStylus(stylus)}
                      disabled={isLoading()}
                    >
                      Edit
                    </button>
                    <button
                      class={styles.deleteButton}
                      onClick={() => handleDelete(stylus)}
                      disabled={isLoading()}
                    >
                      Delete
                    </button>
                  </div>
                </div>
              )}
            </For>
          </div>
        </div>
      </Show>

      <Show
        when={
          !isAddingStylus() && showArchived() && archivedStyluses().length > 0
        }
      >
        <div class={styles.section}>
          <h3 class={styles.sectionTitle}>Archived Styluses</h3>
          <div class={styles.stylusList}>
            <For each={archivedStyluses()}>
              {(stylus) => (
                <div class={`${styles.stylusCard} ${styles.archivedCard}`}>
                  <div class={styles.stylusInfo}>
                    <h3 class={styles.stylusName}>
                      {stylus.name}
                      <span class={styles.archivedTag}>Archived</span>
                    </h3>

                    <Show when={stylus.manufacturer}>
                      <p class={styles.stylusDetail}>
                        <strong>Manufacturer:</strong> {stylus.manufacturer}
                      </p>
                    </Show>

                    <Show when={stylus.expectedLifespan}>
                      <p class={styles.stylusDetail}>
                        <strong>Expected Lifespan:</strong>{" "}
                        {stylus.expectedLifespan} hours
                      </p>
                    </Show>

                    <Show when={stylus.purchaseDate}>
                      <p class={styles.stylusDetail}>
                        <strong>Purchased:</strong>{" "}
                        {new Date(stylus.purchaseDate!).toLocaleDateString()}
                      </p>
                    </Show>
                  </div>

                  <div class={styles.stylusActions}>
                    <button
                      class={styles.editButton}
                      onClick={() => startEditStylus(stylus)}
                      disabled={isLoading()}
                    >
                      Edit
                    </button>
                    <button
                      class={styles.deleteButton}
                      onClick={() => handleDelete(stylus)}
                      disabled={isLoading()}
                    >
                      Delete
                    </button>
                  </div>
                </div>
              )}
            </For>
          </div>
        </div>
      </Show>
    </div>
  );
};

export default StylusManager;
