import { Component, createResource, createEffect } from "solid-js";
import "./App.scss";
import Navbar from "./components/layout/Navbar/Navbar";
import { fetchApi } from "./utils/api";
import { RouteSectionProps, useNavigate } from "@solidjs/router";
import { useAppContext } from "./provider/Provider";

const App: Component<RouteSectionProps<unknown>> = ({ children }) => {
  const navigate = useNavigate();
  const [auth] = createResource("auth", fetchApi);
  const { setAuthPayload } = useAppContext();

  createEffect(() => {
    if (auth.state === "ready") {
      if (auth()?.data?.token) {
        const data = auth().data;

        setAuthPayload({
          isSyncing: data.syncingData,
          lastSynced: data.lastSync,
          releases: data.releases,
          stylus: data.stylus,
        });

        navigate("/");
      } else {
        navigate("/getToken");
      }
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
