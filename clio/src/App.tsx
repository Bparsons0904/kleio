// src/App.tsx
import { Component, createResource, Match, Switch } from "solid-js";
// import { Route, Router } from "@solidjs/router";
// import Home from "./pages/Home";
// import About from "./pages/About";
// import NotFound from "./pages/NotFound";
// import Navbar from "./components/Navbar";
import "./App.scss";
import GetToken from "./components/GetToken/GetToken";
import Navbar from "./components/layout/Navbar/Navbar";
import { fetchApi } from "./utils/api";

const App: Component = () => {
  const [auth] = createResource("health", fetchApi);

  return (
    <>
      <Navbar />
      <Switch>
        <Match when={auth.loading}>Loading...</Match>
        <Match when={auth.error}>Error: {auth.error.message}</Match>
        <Match when={!!auth().data?.token}>Your in!</Match>
        <Match when={!auth().data?.token}>
          <GetToken />
        </Match>
      </Switch>
    </>

    // return (
    //   <>
    //     <Navbar />
    //     <main>
    //       <Router>
    //         <Route path="/" component={Home} />
    //         <Route path="/about" component={About} />
    //         <Route path="*" component={NotFound} />
    //       </Router>
    //     </main>
    //   </>
  );
};

export default App;
