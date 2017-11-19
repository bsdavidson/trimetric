import React from "react";
import {assert} from "chai";
import {shallow, mount} from "enzyme";
import {createStore} from "redux";
import {Router} from "react-router";
import {createMemoryHistory} from "history";

import {StopList} from "../src/components/stop_list";
import {StopListItem} from "../src/components/stop_list_item";
import {reducer, store, DEFAULT_LOCATION} from "../src/store";
import {render} from "../src/router";
import {updateHomeLocation, LocationTypes} from "../src/actions";
import {
  getMockState,
  getMockStopsResponse,
  getMockArrivalsResponse
} from "./mock_data";

describe("<StopList />", () => {
  it("should render a list of <StopListItem />", function() {
    this.wrapper = shallow(<StopList stops={getMockStopsResponse()} />);
    assert.equal(
      this.wrapper.find(".stop-list-items").text(),
      "<withRouter(Connect(StopListItem)) />"
    );
  });

  it("should handle no stops", function() {
    this.wrapper = shallow(<StopList stops={[]} />);
    assert.equal(
      this.wrapper
        .find("h3")
        .first()
        .text(),
      "Visible Stops (0)"
    );
    assert.equal(this.wrapper.find(".stop-list-item").length, 0);
  });

  it("should update current Location", function() {
    assert.equal(store.getState().location, DEFAULT_LOCATION);
    store.dispatch(updateHomeLocation(123, 456, true));
    assert.deepEqual(store.getState().location, {
      locationType: LocationTypes.HOME,
      lat: 123,
      lng: 456,
      gps: true
    });
  });
});

describe("<StopListItem />", () => {
  beforeEach(function() {
    this.wrapper = shallow(
      <StopListItem
        key={123}
        stop={getMockStopsResponse()[0]}
        location={getMockState().location}
      />
    );
  });

  it("should render a single stop item", function() {
    assert.equal(this.wrapper.find(".stop-list-item").length, 1);
  });

  it("should render distance to stop", function() {
    assert.equal(
      this.wrapper
        .find(".stop-distance")
        .first()
        .text(),
      "271"
    );
  });

  it("should render a stop ID", function() {
    assert.equal(
      this.wrapper
        .find(".stop-id")
        .first()
        .text(),
      "6158"
    );
  });
});

describe("<StopList /> routing", () => {
  it("should render <Arrivals /> when URL is /stop/6158 and can goBack", function() {
    let stops = getMockStopsResponse();
    let arrivals = getMockArrivalsResponse();

    let mockStore = createStore(reducer, {
      arrivals,
      stops,
      boundingBox: {ne: stops[0], sw: stops[0]},
      zoom: 16
    });

    let mockHistory = createMemoryHistory("/");

    this.wrapper = mount(
      render(mockStore, () => {}, {Router: Router, history: mockHistory})
    );

    let firstLink = this.wrapper.find(".stop-link").first();
    assert.equal(
      firstLink.html().includes('<a class="stop-link" href="/stop/6158">'),
      true
    );
    // mockHistory.push({
    //   pathname: "/stop/6158"
    // });
    firstLink.simulate("click", {button: 0});

    this.wrapper.update();
    assert.equal(
      this.wrapper
        .find(".arrival-name")
        .first()
        .text(),
      arrivals[0].vehicle_label
    );
    this.wrapper
      .find(".back-button")
      .first()
      .simulate("click", {button: 0});

    this.wrapper.update();
    firstLink = this.wrapper.find(".stop-link").first();

    assert.equal(
      firstLink.html().includes('<a class="stop-link" href="/stop/6158">'),
      true
    );
  });
});
