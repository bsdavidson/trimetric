import "whatwg-fetch";
import fetchPonyfill from "fetch-ponyfill"
import React from "react"
import ReactDOM from "react-dom"
import { Router, hashHistory } from "react-router"

import { App } from "./components/app"
import { Trimet } from "./data"
import { store } from "./store"


const appElement = document.getElementById("app")

const trimet = new Trimet(store, fetchPonyfill())
trimet.start()

const ROUTES = [
  {
    path: "/",
    component: App,
    store: store
  },
  {
    path: "/stop/:stopID",
    component: App,
    store: store
  }
];

function render() {
  ReactDOM.render(
    <Router history={hashHistory} routes={ROUTES} />,
    appElement)
}

store.subscribe(render)
