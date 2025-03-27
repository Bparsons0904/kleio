import {
  createContext,
  createEffect,
  createSignal,
  onCleanup,
  ParentProps,
  useContext,
} from "solid-js";
import { Folder, Release } from "../types";
import { fetchApi } from "../utils/api";

type AppStore = {
  folders: () => Folder[];
  releases: () => Release[];
  lastSynced: () => string;
  isSyncing: () => boolean;

  setFolders: (value: Folder[]) => void;
  setReleases: (value: Release[]) => void;
  setLastSynced: (value: string) => void;
  setIsSyncing: (value: boolean) => void;
};

const defaultStore: Partial<AppStore> = {};
export const AppContext = createContext<AppStore>(defaultStore as AppStore);

export function AppProvider(props: ParentProps) {
  const [folders, setFolders] = createSignal<Folder[]>([]);
  const [releases, setReleases] = createSignal<Release[]>([]);
  const [lastSynced, setLastSynced] = createSignal("");
  const [isSyncing, setIsSyncing] = createSignal(false);

  createEffect(() => {
    if (!isSyncing()) return;

    const pollInterval = setInterval(async () => {
      try {
        const response = await fetchApi("collection/sync");

        if (response.data?.status === "complete") {
          setIsSyncing(false);
          clearInterval(pollInterval);
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

  const store: AppStore = {
    folders,
    releases,
    lastSynced,
    isSyncing,

    setFolders,
    setReleases,
    setLastSynced,
    setIsSyncing,
  };

  return (
    <AppContext.Provider value={store}>{props.children}</AppContext.Provider>
  );
}

export function useAppContext() {
  return useContext(AppContext);
}
