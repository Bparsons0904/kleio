import { render } from "solid-js/web";
import { Router, Route } from "@solidjs/router";

import Home from "./pages/Home/Home";
import GetToken from "./components/GetToken/GetToken";
import LogPlay from "./pages/LogPlay/LogPlay";
import App from "./App";

render(
  () => (
    <Router root={App}>
      <Route path="/" component={Home} />
      <Route path="/getToken" component={GetToken} />
      <Route path="/log" component={LogPlay} />
    </Router>
  ),
  document.getElementById("app"),
);
