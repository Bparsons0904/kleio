// pages/NotFound.tsx
import { Component } from "solid-js";
import { A } from "@solidjs/router";

const NotFound: Component = () => {
  return (
    <div class="not-found-page">
      <h1>404 - Page Not Found</h1>
      <p>The page you're looking for doesn't exist.</p>
      <A href="/">Go back home</A>
    </div>
  );
};

export default NotFound;
