import "whatwg-fetch";
import fetchPonyfill from "fetch-ponyfill";
import ReactDOM from "react-dom";
import {BrowserRouter} from "react-router-dom";

import {render} from "./router";
import {Trimet} from "./data";
import {store} from "./store";

const trimet = new Trimet(store, fetchPonyfill().fetch);
trimet.start();

window.addEventListener("load", () => {
  ReactDOM.render(
    render(store, trimet.handleStopChange, {Router: BrowserRouter}),
    document.getElementById("app")
  );
});
