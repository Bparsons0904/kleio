import { Component, createResource, createEffect } from "solid-js";
import "./App.scss";
import Navbar from "./components/layout/Navbar/Navbar";
import { fetchApi } from "./utils/api";
import { RouteSectionProps, useNavigate } from "@solidjs/router";
import { AppProvider, useAppContext } from "./provider/Provider";

const App: Component<RouteSectionProps<unknown>> = ({ children }) => {
  const navigate = useNavigate();
  const [auth] = createResource("auth", fetchApi);
  const store = useAppContext();

  createEffect(() => {
    if (auth.state === "ready" && auth()?.data?.token) {
      const data = auth().data;

      store.setIsSyncing(data.syncingData);
      store.setLastSynced(data.lastSync);
      store.setReleases(data.releases);
      store.setSyluses(data.stylus);

      navigate("/");
    }
  });
  return (
    <AppProvider>
      <Navbar />
      {children}
    </AppProvider>
  );
};

export default App;
