// components/Navbar.tsx
import { Component } from "solid-js";
// import { A } from "@solidjs/router";
import "./Navbar.scss";

const Navbar: Component = () => {
  return (
    <nav class="navbar">
      <div class="logo">Kleio</div>
      <div class="nav-links">
        {/* <A href="/" end activeClass="active"> */}
        {/*   Home */}
        {/* </A> */}
        {/* <A href="/about" activeClass="active"> */}
        {/*   About */}
        {/* </A> */}
      </div>
    </nav>
  );
};

export default Navbar;
