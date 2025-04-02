import { Component, Show } from "solid-js";
import {
  getCleanlinessScore,
  getCleanlinessColor,
  getPlayRecencyScore,
  getPlayRecencyColor,
  getPlayRecencyText,
  getCleanlinessText,
  countPlaysSinceCleaning,
  getLastCleaningDate,
  getLastPlayDate,
} from "../../utils/playStatus";
import styles from "./StatusIndicators.module.scss";
import { TbWashTemperature5 } from "solid-icons/tb";
import { ImHeadphones } from "solid-icons/im";

export interface StatusIndicatorProps {
  playHistory?: { playedAt: string }[];
  cleaningHistory?: { cleanedAt: string }[];
  showDetails?: boolean;
  onPlayClick?: () => void;
  onCleanClick?: () => void;
}

export const RecordStatusIndicator: Component<StatusIndicatorProps> = (
  props,
) => {
  const lastPlayDate = () => getLastPlayDate(props.playHistory);
  const lastCleanDate = () => getLastCleaningDate(props.cleaningHistory);

  const playsSinceCleaning = () =>
    countPlaysSinceCleaning(props.playHistory || [], lastCleanDate());

  const cleanlinessScore = () =>
    getCleanlinessScore(lastCleanDate(), playsSinceCleaning());

  const playRecencyScore = () => getPlayRecencyScore(lastPlayDate());

  const formatDate = (date: Date | null) => {
    if (!date) return "Never";
    return date.toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  return (
    <div class={styles.container}>
      {/* Show indicators only when not showing details */}
      <Show when={!props.showDetails}>
        <div class={styles.indicatorRow}>
          <PlayStatusIndicator
            score={playRecencyScore()}
            lastPlayed={lastPlayDate()}
            onClick={props.onPlayClick}
          />
          <CleaningStatusIndicator
            score={cleanlinessScore()}
            lastCleaned={lastCleanDate()}
            playsSinceCleaning={playsSinceCleaning()}
            onClick={props.onCleanClick}
          />
        </div>
      </Show>

      {/* Show details section when requested */}
      <Show when={props.showDetails}>
        <div class={styles.detailsSection}>
          <div class={styles.detailRow}>
            <span class={styles.detailLabel}>Last played:</span>
            <span class={styles.detailValue}>{formatDate(lastPlayDate())}</span>
          </div>
          <div class={styles.detailRow}>
            <span class={styles.detailLabel}>Last cleaned:</span>
            <span class={styles.detailValue}>
              {formatDate(lastCleanDate())}
            </span>
          </div>
          <div class={styles.detailRow}>
            <span class={styles.detailLabel}>Plays since cleaning:</span>
            <span class={styles.detailValue}>{playsSinceCleaning()}</span>
          </div>
        </div>
      </Show>
    </div>
  );
};

interface PlayStatusProps {
  score: number;
  lastPlayed: Date | null;
  showDetails?: boolean;
  onClick?: () => void;
}

const PlayStatusIndicator: Component<PlayStatusProps> = (props) => {
  // More muted colors by adding opacity
  const getColorWithOpacity = (colorHex: string): string => {
    return colorHex + "CC"; // Add 80% opacity (CC in hex)
  };

  const color = () => getColorWithOpacity(getPlayRecencyColor(props.score));
  const text = () => getPlayRecencyText(props.lastPlayed);

  const handleClick = (e: MouseEvent) => {
    if (props.onClick) {
      e.stopPropagation(); // Prevent event from bubbling to parent card
      props.onClick();
    }
  };

  return (
    <div class={styles.indicator}>
      <div
        class={`${styles.iconContainer} ${props.onClick ? styles.clickable : ""}`}
        style={{ "background-color": color() }}
        onClick={handleClick}
        title="Click to log a play"
      >
        <ImHeadphones size={15} color="white" />
      </div>
      <span class={styles.tooltipText}>{text()}</span>
    </div>
  );
};

interface CleaningStatusProps {
  score: number;
  lastCleaned: Date | null;
  playsSinceCleaning: number;
  showDetails?: boolean;
  onClick?: () => void;
}

const CleaningStatusIndicator: Component<CleaningStatusProps> = (props) => {
  const getColorWithOpacity = (colorHex: string): string => {
    return colorHex + "CC";
  };

  const color = () => getColorWithOpacity(getCleanlinessColor(props.score));
  const text = () => getCleanlinessText(props.score);

  const handleClick = (e: MouseEvent) => {
    if (props.onClick) {
      e.stopPropagation(); // Prevent event from bubbling to parent card
      props.onClick();
    }
  };

  return (
    <div class={styles.indicator}>
      <div
        class={`${styles.iconContainer} ${props.onClick ? styles.clickable : ""}`}
        style={{ "background-color": color() }}
        onClick={handleClick}
        title="Click to log a cleaning"
      >
        <TbWashTemperature5 size={20} color="white" />
      </div>
      <span class={styles.tooltipText}>{text()}</span>
    </div>
  );
};
