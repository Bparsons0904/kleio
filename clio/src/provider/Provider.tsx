// src/provider/Provider.tsx (modified)
import {
  createContext,
  createEffect,
  createSignal,
  onCleanup,
  ParentProps,
  useContext,
  Show,
} from "solid-js";
import { Folder, Release, Stylus } from "../types";
import { fetchApi } from "../utils/api";
import Toast, { ToastType } from "../components/layout/Toast/Toast";

export interface Payload {
  isSyncing: boolean;
  lastSynced: string;
  releases: Release[];
  stylus: Stylus[];
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
  lastSynced: () => string;
  isSyncing: () => boolean;

  setFolders: (value: Folder[]) => void;
  setReleases: (value: Release[]) => void;
  setStyluses: (value: Stylus[]) => void;
  setLastSynced: (value: string) => void;
  setIsSyncing: (value: boolean) => void;

  setKleioStore: (payload: Payload) => void;

  // Toast functions
  showToast: (message: string, type: ToastType, duration?: number) => void;
  showSuccess: (message: string, duration?: number) => void;
  showError: (message: string, duration?: number) => void;
  showInfo: (message: string, duration?: number) => void;
  clearToast: () => void;
};

const defaultStore: Partial<AppStore> = {};
export const AppContext = createContext<AppStore>(defaultStore as AppStore);

export function AppProvider(props: ParentProps) {
  const [folders, setFolders] = createSignal<Folder[]>([]);
  const [releases, setReleases] = createSignal<Release[]>([]);
  const [styluses, setStyluses] = createSignal<Stylus[]>([]);
  const [lastSynced, setLastSynced] = createSignal("");
  const [isSyncing, setIsSyncing] = createSignal(false);
  const [toast, setToast] = createSignal<ToastState | null>(null);

  createEffect(() => {
    if (!isSyncing()) return;

    const pollInterval = setInterval(async () => {
      try {
        const response = await fetchApi<Payload>("collection/sync");

        if (response.data?.status === "complete") {
          setIsSyncing(false);
          clearInterval(pollInterval);
          const result = await fetchApi("collection");
          if (result.data) {
            setAuthPayload(result.data);
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

  const setAuthPayload = (payload: Payload) => {
    setIsSyncing(payload.isSyncing);
    setLastSynced(payload.lastSynced);
    setReleases(payload.releases);
    setStyluses(payload.stylus);
  };

  // Toast functions
  const showToast = (message: string, type: ToastType, duration?: number) => {
    setToast({ message, type, duration });
  };

  const showSuccess = (message: string, duration?: number) => {
    showToast(message, "success", duration);
  };

  const showError = (message: string, duration?: number) => {
    showToast(message, "error", duration || 5000); // Longer duration for errors
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
    lastSynced,
    isSyncing,

    setFolders,
    setReleases,
    setStyluses,
    setLastSynced,
    setIsSyncing,

    setKleioStore: setAuthPayload,

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
