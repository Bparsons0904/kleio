import styles from "./GetToken.module.scss";
import { Component, Show, createSignal } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { useAppContext } from "../../provider/Provider";
import { postApi } from "../../utils/mutations/post";

const GetToken: Component = () => {
  const [token, setToken] = createSignal("");
  const [isSaved, setIsSaved] = createSignal(false);
  const [isError, setIsError] = createSignal(false);
  const [errorMessage, setErrorMessage] = createSignal("");
  const navigate = useNavigate();
  const { setKleioStore } = useAppContext();

  const saveToken = async (tokenValue: string) => {
    try {
      const response = await postApi("auth/token", { token: tokenValue });

      if (response.status !== 200) {
        throw new Error("Failed to save token");
      }

      setIsSaved(true);
      setIsError(false);
      return response.data;
    } catch (error) {
      console.error("Error saving token:", error);
      setIsError(true);
      setErrorMessage(
        error.response?.data?.message ||
          error.message ||
          "Failed to save token",
      );
      return null;
    }
  };

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    if (!token().trim()) {
      setIsError(true);
      setErrorMessage("Please enter a valid token");
      return;
    }

    const result = await saveToken(token());
    if (result) {
      setKleioStore(result.data);
      navigate("/");
    }
  };

  return (
    <div class={styles.container}>
      <h1 class={styles.title}>Discogs API Configuration</h1>

      <div class={styles.infoSection}>
        <h2 class={styles.sectionTitle}>What is a Discogs Token?</h2>
        <p class={styles.text}>
          A Discogs API token allows this application to securely access your
          Discogs collection, search the Discogs database, and perform other
          operations without requiring your full Discogs credentials.
        </p>

        <h2 class={styles.sectionTitle}>How to Get Your Token</h2>
        <ol class={styles.instructionList}>
          <li>
            Go to your{" "}
            <a
              href="https://www.discogs.com/settings/developers"
              target="_blank"
              rel="noopener noreferrer"
              class={styles.link}
            >
              Discogs Developer Settings
            </a>
          </li>
          <li>
            Sign in to your Discogs account if you're not already logged in
          </li>
          <li>
            Under "Developer Resources", find your personal access token or
            generate a new one
          </li>
          <li>Copy the token and paste it in the field below</li>
        </ol>
      </div>

      <form onSubmit={handleSubmit} class={styles.form}>
        <div class={styles.formGroup}>
          <label for="token" class={styles.label}>
            Your Discogs API Token
          </label>
          <input
            type="text"
            id="token"
            value={token()}
            onInput={(e) => setToken(e.target.value)}
            placeholder="Paste your token here"
            class={styles.input}
            required
          />
        </div>

        <Show when={isError()}>
          <div class={styles.errorMessage}>{errorMessage()}</div>
        </Show>

        <Show when={isSaved()}>
          <div class={styles.successMessage}>Token saved successfully!</div>
        </Show>

        <button type="submit" class={styles.button}>
          Save Token
        </button>
      </form>

      <div class={styles.footer}>
        <p>
          Your token is stored securely and only used to access the Discogs API
          on your behalf. We never share your token or use it for any other
          purpose.
        </p>
      </div>
    </div>
  );
};

export default GetToken;
