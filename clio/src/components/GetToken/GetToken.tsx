import styles from "./GetToken.module.scss";
import { postApi } from "../../utils/api";
import { Show, createSignal } from "solid-js";

const GetToken = () => {
  const [token, setToken] = createSignal("");
  const [isSaved, setIsSaved] = createSignal(false);
  const [isError, setIsError] = createSignal(false);
  const [errorMessage, setErrorMessage] = createSignal("");

  // Function to save token using the API client
  const saveToken = async (tokenValue) => {
    try {
      // Use your existing API client to save the token
      const response = await postApi("discogs/token", { token: tokenValue });

      if (response.status !== 200) {
        throw new Error("Failed to save token");
      }

      setIsSaved(true);
      setIsError(false);
      return true;
    } catch (error) {
      console.error("Error saving token:", error);
      setIsError(true);
      setErrorMessage(
        error.response?.data?.message ||
          error.message ||
          "Failed to save token",
      );
      return false;
    }
  };

  // Handle form submission
  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!token().trim()) {
      setIsError(true);
      setErrorMessage("Please enter a valid token");
      return;
    }

    await saveToken(token());
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
