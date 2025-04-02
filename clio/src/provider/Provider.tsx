import {
  createContext,
  createEffect,
  createSignal,
  onCleanup,
  ParentProps,
  useContext,
  Show,
} from "solid-js";
import { Folder, PlayHistory, Release, Stylus } from "../types";
import { fetchApi } from "../utils/api";
import Toast, { ToastType } from "../components/layout/Toast/Toast";

export interface Payload {
  isSyncing: boolean;
  lastSync: string;
  releases: Release[];
  stylus: Stylus[];
  playHistory: PlayHistory[];
  folders: Folder[];
}

interface ToastState {
  message: string;
  type: ToastType;
  duration?: number;
}

type AppStore = {
  folders: () => Folder[];
  releases: () => Release[];
  styluses: () => Stylus[];
  playHistory: () => PlayHistory[];
  lastSynced: () => string;
  isSyncing: () => boolean;
  selectedFolderId: () => number;

  setFolders: (value: Folder[]) => void;
  setReleases: (value: Release[]) => void;
  setStyluses: (value: Stylus[]) => void;
  setPlayHistory: (value: PlayHistory[]) => void;
  setLastSynced: (value: string) => void;
  setIsSyncing: (value: boolean) => void;
  setSelectedFolderId: (value: number) => void;

  setKleioStore: (payload: Payload) => void;

  showToast: (message: string, type: ToastType, duration?: number) => void;
  showSuccess: (message: string, duration?: number) => void;
  showError: (message: string, duration?: number) => void;
  showInfo: (message: string, duration?: number) => void;
  clearToast: () => void;
};

const defaultStore: Partial<AppStore> = {};
export const AppContext = createContext<AppStore>(defaultStore as AppStore);

const loadSelectedFolderId = (): number => {
  const savedId = localStorage.getItem("selectedFolderId");
  return savedId ? parseInt(savedId, 10) : 0; // Default to folder ID 0 (All)
};

export function AppProvider(props: ParentProps) {
  const [folders, setFolders] = createSignal<Folder[]>([]);
  const [releases, setReleases] = createSignal<Release[]>([]);
  const [rawReleases, setRawReleases] = createSignal<Release[]>([]);
  const [styluses, setStyluses] = createSignal<Stylus[]>([]);
  const [playHistory, setPlayHistory] = createSignal<PlayHistory[]>([]);
  const [lastSynced, setLastSynced] = createSignal("");
  const [isSyncing, setIsSyncing] = createSignal(false);
  const [toast, setToast] = createSignal<ToastState | null>(null);
  const [selectedFolderId, setSelectedFolderId] = createSignal(
    loadSelectedFolderId(),
  );

  const filterReleases = (rawReleases: Release[]) => {
    const folderId = selectedFolderId();
    if (folderId === 0) {
      return rawReleases;
    }
    return rawReleases.filter((release) => release.folderId === folderId);
  };

  createEffect(() => {
    localStorage.setItem("selectedFolderId", selectedFolderId().toString());
  });

  createEffect(() => {
    if (!isSyncing()) return;

    const pollInterval = setInterval(async () => {
      try {
        const response = await fetchApi("collection/sync");

        if (response.data?.status === "complete") {
          setIsSyncing(false);
          clearInterval(pollInterval);
          const result = await fetchApi("collection");
          if (result.data) {
            console.log("Fetched data from server:", result.data);
            setKleioStore(result.data);
          }
        }
      } catch (error) {
        console.error("Error polling sync status:", error);
        setIsSyncing(false);
        clearInterval(pollInterval);
      }
    }, 1000);

    onCleanup(() => {
      clearInterval(pollInterval);
    });
  });

  const setKleioStore = (payload: Payload) => {
    if (!payload) return;
    setIsSyncing(payload.isSyncing);
    setLastSynced(payload.lastSync);
    setRawReleases(payload.releases);
    setStyluses(payload.stylus);
    setPlayHistory(payload.playHistory);
    setFolders(payload.folders);
    setReleases(filterReleases(payload.releases));
  };

  const showToast = (message: string, type: ToastType, duration?: number) => {
    setToast({ message, type, duration });
  };

  const showSuccess = (message: string, duration?: number) => {
    showToast(message, "success", duration);
  };

  const showError = (message: string, duration?: number) => {
    showToast(message, "error", duration || 5000);
  };

  const showInfo = (message: string, duration?: number) => {
    showToast(message, "info", duration);
  };

  const clearToast = () => {
    setToast(null);
  };

  const store: AppStore = {
    folders,
    releases,
    styluses,
    playHistory,
    lastSynced,
    isSyncing,
    selectedFolderId,

    setFolders,
    setReleases,
    setStyluses,
    setPlayHistory,
    setLastSynced,
    setIsSyncing,
    setSelectedFolderId,

    setKleioStore,

    // Toast functions
    showToast,
    showSuccess,
    showError,
    showInfo,
    clearToast,
  };

  return (
    <AppContext.Provider value={store}>
      {props.children}
      <Show when={toast()}>
        {(activeToast) => (
          <Toast
            message={activeToast().message}
            type={activeToast().type}
            duration={activeToast().duration}
            onClose={clearToast}
          />
        )}
      </Show>
    </AppContext.Provider>
  );
}

export function useAppContext() {
  return useContext(AppContext);
}
