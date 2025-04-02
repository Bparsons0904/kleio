// src/components/layout/Navbar/Navbar.tsx
import { Component } from "solid-js";
import styles from "./Navbar.module.scss";
import { useAppContext } from "../../../provider/Provider";
import { useFormattedShortDate } from "../../../utils/dates";
import { useNavigate, useLocation } from "@solidjs/router";

const Navbar: Component = () => {
  const { isSyncing, lastSynced } = useAppContext();
  const navigate = useNavigate();
  const location = useLocation();

  const isActive = (path: string) => {
    return location.pathname === path ? styles.active : "";
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
          <div class={styles.lastSync}>
            Last sync: {useFormattedShortDate(lastSynced())}
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
