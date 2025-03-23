import { defineConfig } from "eslint/config";
import globals from "globals";
import js from "@eslint/js";
import tseslint from "typescript-eslint";

export default defineConfig([
  // Base settings for all files
  {
    files: ["**/*.{js,mjs,cjs,ts}"],
    ignores: ["**/node_modules/**", "main", "*.db", "tmp/**"],
  },

  // Browser environment for client code
  {
    files: ["clio/**/*.{js,mjs,cjs,ts,tsx}"],
    languageOptions: { globals: globals.browser },
  },

  // JS configuration
  {
    files: ["**/*.{js,mjs,cjs}"],
    plugins: { js },
    extends: ["js/recommended"],
  },

  // TypeScript configuration
  {
    ...tseslint.configs.recommended,
    files: ["**/*.{ts,tsx}"],
  },
]);
