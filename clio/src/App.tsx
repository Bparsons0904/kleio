import { Component, createResource, createEffect } from "solid-js";
import "./App.scss";
import Navbar from "./components/layout/Navbar/Navbar";
import { fetchApi } from "./utils/api";
import { RouteSectionProps, useNavigate } from "@solidjs/router";
import { useAppContext } from "./provider/Provider";

const App: Component<RouteSectionProps<unknown>> = ({ children }) => {
  const navigate = useNavigate();
  const [auth] = createResource("auth", fetchApi);
  const { setKleioStore: setAuthPayload } = useAppContext();

  createEffect(() => {
    if (auth.state === "ready") {
      if (!auth()?.data?.token) {
        navigate("/getToken");
        return;
      }

      const data = auth().data;
      setAuthPayload({
        isSyncing: data.syncingData,
        lastSynced: data.lastSync,
        releases: data.releases,
        stylus: data.stylus,
      });
    }
  });
  return (
    <>
      <Navbar />
      {children}
    </>
  );
};

export default App;
