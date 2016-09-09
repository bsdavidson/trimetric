import React from "react"
import { assert } from "chai"
import { mount } from "enzyme"
import { createStore } from "redux"
import { Router, hashHistory } from "react-router"
import { reducer } from "../src/store"
import { App } from "../src/components/app"


describe("App", () => {
  it("Render a container", function() {
    let store = createStore(reducer, {stops: [], vehicles: {arrivals: []}})

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

    let wrapper = mount(<Router history={hashHistory} routes={ROUTES} />)
    assert.equal(wrapper.find(".app").length, 1)
  })
})
