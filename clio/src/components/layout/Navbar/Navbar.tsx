import { Component } from "solid-js";
import styles from "./Navbar.module.scss";
import { useAppContext } from "../../../provider/Provider";
import { useFormattedMediumDate } from "../../../utils/dates";

const Navbar: Component = () => {
  const { isSyncing, lastSynced } = useAppContext();

  return (
    <nav class={styles.navbar}>
      <div class={styles.logo} onclick={() => (window.location.href = "/")}>
        Kleio
      </div>
      <div class={styles.navLinks}>
        {isSyncing() && (
          <div class={styles.syncIndicator}>
            <span class={styles.syncSpinner}></span>
            <span>Syncing...</span>
          </div>
        )}
        {!isSyncing() && lastSynced() && (
          <div class={styles.lastSync}>
            Last sync: {useFormattedMediumDate(lastSynced())}
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
