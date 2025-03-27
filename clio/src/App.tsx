import {
  Component,
  createResource,
  Match,
  Switch,
  createEffect,
  Suspense,
} from "solid-js";
import "./App.scss";
import GetToken from "./components/GetToken/GetToken";
import Navbar from "./components/layout/Navbar/Navbar";
import { fetchApi } from "./utils/api";
import Home from "./pages/Home/Home";
import { Route, Router } from "@solidjs/router";
import { AppProvider, useAppContext } from "./provider/Provider";

const App: Component = () => {
  return (
    <AppProvider>
      <Navbar />
      <Suspense>
        <Main />
      </Suspense>
    </AppProvider>
  );
};

export default App;

const Main: Component = () => {
  const [auth] = createResource("auth", fetchApi);
  const store = useAppContext();

  createEffect(() => {
    if (auth.state === "ready" && auth()?.data) {
      const data = auth().data;

      store.setIsSyncing(data.syncingData);
      store.setLastSynced(data.lastSync);
    }
  });

  return (
    <Router>
      <Switch>
        <Match when={auth.loading}>Loading...</Match>
        <Match when={auth.error}>Error: {auth.error.message}</Match>
        <Match when={auth.state === "ready" && auth()?.data}>
          <Match when={!auth().data}>
            <GetToken />
          </Match>
          <Route path="/" component={Home} />
        </Match>
      </Switch>
    </Router>
  );
};
