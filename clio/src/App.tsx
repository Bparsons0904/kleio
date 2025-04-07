import { Component, createResource, createEffect } from "solid-js";
import "./App.scss";
import Navbar from "./components/layout/Navbar/Navbar";
import { fetchApi } from "./utils/api";
import { RouteSectionProps, useNavigate } from "@solidjs/router";
import { useAppContext } from "./provider/Provider";
import SubNavbar from "./components/layout/SubNavbar/SubNavbar";

const App: Component<RouteSectionProps<unknown>> = ({ children }) => {
  const navigate = useNavigate();
  const [auth] = createResource("auth", fetchApi);
  const { setKleioStore } = useAppContext();

  createEffect(() => {
    if (auth.state === "ready") {
      if (!auth()?.data?.token) {
        navigate("/getToken");
        return;
      }

      setKleioStore(auth().data);
    }
  });

  return (
    <>
      <Navbar />
      <SubNavbar />
      {children}
    </>
  );
};

export default App;
