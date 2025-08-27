// src/components/layout/Navbar/Navbar.tsx
import { Component, createEffect } from "solid-js";
import styles from "./Navbar.module.scss";
import { useAppContext } from "../../../provider/Provider";
import { useFormattedShortDate } from "../../../utils/dates";
import { useNavigate, useLocation } from "@solidjs/router";
import { refreshCollection } from "../../../utils/api";

const Navbar: Component = () => {
  const { isSyncing, lastSynced, setIsSyncing } = useAppContext();
  const navigate = useNavigate();
  const location = useLocation();

  const isActive = (path: string) => {
    return location.pathname === path ? styles.active : "";
  };

  // Add a small effect to make sure body has proper padding when navbar changes
  createEffect(() => {
    const navbar = document.querySelector(`.${styles.navbar}`);
    if (navbar) {
      // const navbarHeight = navbar.clientHeight;
      // document.body.style.paddingTop = `${navbarHeight}px`;
    }
  });

  const handleResync = async () => {
    try {
      console.log("User clicked sync button");
      setIsSyncing(true); // Set syncing state immediately for UI feedback
      
      const response = await refreshCollection();
      console.log("Sync response:", response.data);
      
      if (response.status === 200) {
        console.log("Sync started successfully");
        // Keep syncing state true - it will be cleared when sync completes
      } else {
        console.error("Unexpected response status:", response.status);
        setIsSyncing(false);
      }
    } catch (error) {
      console.error("Error starting sync:", error);
      setIsSyncing(false); // Clear syncing state on error
      // You might want to show a toast notification here
    }
  };

  return (
    <nav class={styles.navbar}>
      <div class={styles.logo} onclick={() => navigate("/")}>
        Kleio
      </div>

      <div class={styles.navLinks}>
        <a
          class={`${styles.navLink} ${isActive("/log")}`}
          onclick={(e) => {
            e.preventDefault();
            navigate("/log");
          }}
        >
          Log
        </a>

        <a
          class={`${styles.navLink} ${isActive("/collection")}`}
          onclick={(e) => {
            e.preventDefault();
            navigate("/collection");
          }}
        >
          Collection
        </a>

        <a
          class={`${styles.navLink} ${isActive("/playHistory")}`}
          onclick={(e) => {
            e.preventDefault();
            navigate("/playHistory");
          }}
        >
          History
        </a>

        <a
          class={`${styles.navLink} ${isActive("/analytics")}`}
          onclick={(e) => {
            e.preventDefault();
            navigate("/analytics");
          }}
        >
          Analytics
        </a>
      </div>

      <div class={styles.syncStatus}>
        {isSyncing() && (
          <div class={styles.syncIndicator}>
            <span class={styles.syncSpinner}></span>
            <span>Syncing...</span>
          </div>
        )}
        {!isSyncing() && (
          <div class={styles.syncControls}>
            <button class={styles.syncButton} onclick={handleResync}>
              Sync Now
            </button>
            {lastSynced() && (
              <div class={styles.lastSyncInfo}>
                Last: {useFormattedShortDate(lastSynced())}
              </div>
            )}
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
