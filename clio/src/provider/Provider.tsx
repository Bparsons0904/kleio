import {
  createContext,
  createEffect,
  createSignal,
  onCleanup,
  ParentProps,
  useContext,
} from "solid-js";
import { Folder, Release, Stylus } from "../types";
import { fetchApi } from "../utils/api";

interface AuthPayload {
  isSyncing: boolean;
  lastSynced: string;
  releases: Release[];
  stylus: Stylus[];
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

  setAuthPayload: (payload: AuthPayload) => void;
};

const defaultStore: Partial<AppStore> = {};
export const AppContext = createContext<AppStore>(defaultStore as AppStore);

export function AppProvider(props: ParentProps) {
  const [folders, setFolders] = createSignal<Folder[]>([]);
  const [releases, setReleases] = createSignal<Release[]>([]);
  const [styluses, setStyluses] = createSignal<Stylus[]>([]);
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

  const setAuthPayload = (payload: AuthPayload) => {
    setIsSyncing(payload.isSyncing);
    setLastSynced(payload.lastSynced);
    setReleases(payload.releases);
    setStyluses(payload.stylus);
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

    setAuthPayload,
  };

  return (
    <AppContext.Provider value={store}>{props.children}</AppContext.Provider>
  );
}

export function useAppContext() {
  return useContext(AppContext);
}
