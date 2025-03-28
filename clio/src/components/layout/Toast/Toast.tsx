import { Component, createEffect, onCleanup } from "solid-js";
import styles from "./Toast.module.scss";

export type ToastType = "success" | "error" | "info";

export interface ToastProps {
  message: string;
  type: ToastType;
  duration?: number;
  onClose: () => void;
}

const Toast: Component<ToastProps> = (props) => {
  const duration = props.duration || 3000; // Default 3 seconds

  createEffect(() => {
    if (props.message) {
      const timer = setTimeout(() => {
        props.onClose();
      }, duration);

      onCleanup(() => clearTimeout(timer));
    }
  });

  return (
    <div class={`${styles.toast} ${styles[props.type]}`}>
      <div class={styles.content}>{props.message}</div>
    </div>
  );
};

export default Toast;
