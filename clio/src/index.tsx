import { render } from "solid-js/web";
import { Router, Route } from "@solidjs/router";

import Home from "./pages/Home/Home";
import GetToken from "./components/GetToken/GetToken";
import LogPlay from "./pages/LogPlay/LogPlay";
import App from "./App";
import StylusManager from "./pages/Stylus/StylusManager";
import { AppProvider } from "./provider/Provider";

render(
  () => (
    <AppProvider>
      <Router root={App}>
        <Route path="/" component={Home} />
        <Route path="/getToken" component={GetToken} />
        <Route path="/log" component={LogPlay} />
        <Route path="/equipment" component={StylusManager} />
      </Router>
    </AppProvider>
  ),
  document.getElementById("app"),
);
