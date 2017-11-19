import React from "react";
import {assert} from "chai";
import {shallow} from "enzyme";
import {ArrivalList} from "../src/components/arrival_list";
import {ArrivalListItem} from "../src/components/arrival_list_item";
import {getMockStopsResponse, getMockArrivalsResponse} from "./mock_data";

describe("<ArrivalsList />", () => {
  beforeEach(function() {
    this.stop = getMockStopsResponse()[0];
    this.clickedStop = null;
    this.wrapper = shallow(
      <ArrivalList
        stop={this.stop}
        arrivals={getMockArrivalsResponse()}
        onRouteNameClick={stop => {
          this.clickedStop = stop;
        }}
      />
    );
  });

  it("should render Stop Title", function() {
    assert.equal(
      this.wrapper.find(".arrival-list-description").text(),
      this.stop.name
    );
  });

  it("should render a list of <ArrivalListItem />", function() {
    assert.equal(
      this.wrapper.find(".arrival-list-items").text(),
      "<withRouter(Connect(ArrivalListItem)) /><withRouter(Connect(ArrivalListItem)) />"
    );
  });

  it("should update location when Stop is clicked", function() {
    let stopTitle = this.wrapper.find("h3").first();
    stopTitle.simulate("click", {button: 0});
    assert.equal(this.clickedStop.id, "6158");
  });
});

describe("<ArrivalListItem />", () => {
  beforeEach(function() {
    this.arrivals = getMockArrivalsResponse();

    let mockArrivalTime = "5 minutes";
    let mockColor = "#000000";

    this.wrapper = shallow(
      <ArrivalListItem
        key={1}
        color={mockColor}
        arrival={this.arrivals[0]}
        arrivalTime={mockArrivalTime}
      />
    );
    this.wrapper2 = shallow(
      <ArrivalListItem
        key={1}
        color={mockColor}
        arrival={this.arrivals[1]}
        arrivalTime={mockArrivalTime}
      />
    );
  });

  it("should show estimated arrival time", function() {
    assert.equal(this.wrapper.find(".arrival-est-time").text(), "5 minutes");
  });

  it("should show full sign info", function() {
    assert.equal(
      this.wrapper.find(".arrival-name").text(),
      this.arrivals[0].vehicle_label
    );
  });

  it("should show direction traveling", function() {
    assert.equal(
      this.wrapper.find(".arrival-direction").text(),
      "Traveling: W"
    );
  });
});
