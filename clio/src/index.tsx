import { render } from "solid-js/web";
import { Router, Route } from "@solidjs/router";
import "solid-devtools";
import Home from "./pages/Home/Home";
import GetToken from "./components/GetToken/GetToken";
import LogPlay from "./pages/LogPlay/LogPlay";
import App from "./App";
import StylusManager from "./pages/Stylus/StylusManager";
import { AppProvider } from "./provider/Provider";
import "./styles/reset.scss";
import PlayHistoryPage from "./pages/PlayHistory/PlayHistory";
import Analytics from "./pages/Analytics/Analytics";
import Collection from "./pages/Collection/Collection";

render(
  () => (
    <AppProvider>
      <Router root={App}>
        <Route path="/" component={Home} />
        <Route path="/getToken" component={GetToken} />
        <Route path="/log" component={LogPlay} />
        <Route path="/equipment" component={StylusManager} />
        <Route path="/playHistory" component={PlayHistoryPage} />
        <Route path="/analytics" component={Analytics} />
        <Route path="/collection" component={Collection} />
      </Router>
    </AppProvider>
  ),
  document.getElementById("app"),
);
