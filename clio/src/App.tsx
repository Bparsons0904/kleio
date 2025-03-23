// src/App.tsx
import { Component, createResource, Match, Switch } from "solid-js";
// import { Route, Router } from "@solidjs/router";
// import About from "./pages/About";
// import NotFound from "./pages/NotFound";
import "./App.scss";
import GetToken from "./components/GetToken/GetToken";
import Navbar from "./components/layout/Navbar/Navbar";
import { fetchApi } from "./utils/api";
import Home from "./pages/Home/Home";

const App: Component = () => {
  const [auth] = createResource("auth", fetchApi);

  return (
    <>
      <Navbar />
      <Switch>
        <Match when={auth.loading}>Loading...</Match>
        <Match when={auth.error}>Error: {auth.error.message}</Match>
        <Match when={!!auth().data?.token}>
          <Home />
        </Match>
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
