import React from "react"
import { assert } from "chai"
import { mount, shallow } from "enzyme"
import { createStore } from "redux"
import { Router, createMemoryHistory } from "react-router"

import { App } from "../src/components/app"
import { StopList, StopListItem } from "../src/components/stop_list"
import { reducer, store } from "../src/store"
import { updateHomeLocation } from "../src/actions"
import { getMockCombinedData, getMockState, adjustTime } from "./mock_data"

describe("<StopList />", () => {
  it("should render a list of <StopListItem />", function() {
    let mockStops = adjustTime(getMockCombinedData())
    let mockState = getMockState()
    this.wrapper = mount(<StopList stops={mockStops.stops} state={mockState} />)
    assert.equal(this.wrapper.find(".stop-link").length, 1 )
  })

  it("should handle no stops", function() {
    let mockStops = {stops:[]}
    let mockState = getMockState()
    this.wrapper = shallow(<StopList stops={mockStops.stops} state={mockState} />)
    assert.equal(
      this.wrapper.find(".stop-list-items").first().text(),
      "Sorry, no buses are running near you. Better start walking or call an Uber." )
  })

  it("should update current Location", function() {
    assert.equal(store.getState().location.lat, 45.5247402)
    store.dispatch(updateHomeLocation(123, 456, true))
    assert.equal(store.getState().location.lat, 123)
  })

})

describe("<StopListItem />", () => {
  beforeEach(function() {
    let mockStops = adjustTime(getMockCombinedData())
    let stop = mockStops.stops[0]
    let state = getMockState()
    this.wrapper = shallow(<StopListItem key={123} stop={stop} state={state} />)
  })

  it("should render a single stop item", function(){
    assert.equal(this.wrapper.find(".stop-list-item").length, 1)
  })

  it("should render next arrival time", function() {
    assert.equal(this.wrapper.find(".stop-estimate").first().text(), "in 10 minutes ")
  })

  it("should render distance to stop", function() {
    assert.equal(this.wrapper.find(".stop-distance").first().text(), "271")
  })

  it("should render a stop ID", function() {
    assert.equal(this.wrapper.find(".stop-id").first().text(), "6158")
  })

})

describe("<StopList /> routing", () => {
  it("should render <Arrivals /> when URL is /stop/6158 and can goBack", function() {

    let mockData = getMockCombinedData()
    let mockStore = createStore(reducer, {
      stops: mockData.stops,
      vehicles: {arrivals: []}
    })

    let mockHistory = createMemoryHistory("/")
    const ROUTES = [
      {
        path: "/",
        component: App,
        store: mockStore
      },
      {
        path: "/stop/:stopID",
        component: App,
        store: mockStore
      }
    ]

    this.wrapper = mount(<Router history={mockHistory} routes={ROUTES} />)
    let firstLink = this.wrapper.find(".stop-link").first()
    assert.equal(firstLink.html().includes('<a class="stop-link" href="/stop/6158">'), true);
    mockHistory.push({
      pathname: "/stop/6158"
    })

    assert.equal(this.wrapper.find(".arrival-name").first().text(), "15 To SW 5th & Washington");
    this.wrapper.find(".back-button").first().simulate("click")
    assert.equal(firstLink.html().includes('<a class="stop-link" href="/stop/6158">'), true);
  })
})
