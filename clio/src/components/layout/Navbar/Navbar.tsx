import { Component } from "solid-js";
import "./Navbar.scss";
import { useAppContext } from "../../../provider/Provider";
import { useFormattedMediumDate } from "../../../utils/dates";

const Navbar: Component = () => {
  const { isSyncing, lastSynced } = useAppContext();

  return (
    <nav class="navbar">
      <div class="logo" onclick={() => (window.location.href = "/")}>
        Kleio
      </div>
      <div class="nav-links">
        {isSyncing() && (
          <div class="sync-indicator">
            <span class="sync-spinner"></span>
            <span>Syncing...</span>
          </div>
        )}
        {!isSyncing() && lastSynced() && (
          <div class="last-sync">
            Last sync: {useFormattedMediumDate(lastSynced())}
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
