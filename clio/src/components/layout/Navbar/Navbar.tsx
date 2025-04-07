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
      const response = await refreshCollection();
      if (response.status === 200) {
        setIsSyncing(true);
      }
    } catch (error) {
      console.error("Error resyncing:", error);
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
        {!isSyncing() && lastSynced() && (
          <div class={styles.lastSync} on:click={handleResync}>
            Last sync: {useFormattedShortDate(lastSynced())}
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
