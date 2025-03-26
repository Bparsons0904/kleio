import { createContext, createSignal, ParentProps, useContext } from "solid-js";
import { Folder, Release } from "../types";

// Define the store type with proper TypeScript types
type AppStore = {
  // Accessor functions
  folders: () => Folder[];
  releases: () => Release[];
  lastSynced: () => string;
  isSyncing: () => boolean;

  // Setter functions
  setFolders: (value: Folder[]) => void;
  setReleases: (value: Release[]) => void;
  setLastSynced: (value: string) => void;
  setIsSyncing: (value: boolean) => void;
};

// Create the context with a partial implementation that satisfies TypeScript
const defaultStore: Partial<AppStore> = {};
export const AppContext = createContext<AppStore>(defaultStore as AppStore);

// Provider component
export function AppProvider(props: ParentProps) {
  // Create all your signals
  const [folders, setFolders] = createSignal<Folder[]>([]);
  const [releases, setReleases] = createSignal<Release[]>([]);
  const [lastSynced, setLastSynced] = createSignal("");
  const [isSyncing, setIsSyncing] = createSignal(false);

  const store: AppStore = {
    // Signals
    folders,
    releases,
    lastSynced,
    isSyncing,

    // Setters
    setFolders,
    setReleases,
    setLastSynced,
    setIsSyncing,
  };

  return (
    <AppContext.Provider value={store}>{props.children}</AppContext.Provider>
  );
}

// Hook to use the context
export function useAppContext() {
  return useContext(AppContext);
}
