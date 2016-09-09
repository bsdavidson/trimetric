import React from "react"
import { assert } from "chai"
import { shallow } from "enzyme"
import { ArrivalList, ArrivalListItem } from "../src/components/arrival_list"
import { getMockCombinedData, adjustTime } from "./mock_data"
import { store } from "../src/store"

describe("<ArrivalsList />", () => {
  beforeEach(function() {
    let mockStop = adjustTime(getMockCombinedData()).stops[0]
    let mockState = {}
    this.wrapper = shallow(<ArrivalList stop={mockStop} state={mockState}/>)
  })

  it("should render Stop Title", function() {
    assert.equal(this.wrapper.find(".arrival-list-description").text(), "SW Washington & 3rdGet Directions")
  })

  it("should render a list of <ArrivalListItem />", function() {
    assert.equal(this.wrapper.find(".arrival-list-items").text(), "<ArrivalListItem /><ArrivalListItem />")
  })

  it("should update location when Stop is clicked", function() {
    assert.equal(store.getState().locationClicked, null)
    this.wrapper.find("h3").first().simulate("click");
    assert.equal(store.getState().locationClicked.id, 6158)
  })
})

describe("<ArrivalListItem />", () => {
  beforeEach(function() {
    let mockArrival = adjustTime(getMockCombinedData()).stops[0].arrivals[0]
    let mockArrival2 = adjustTime(getMockCombinedData()).stops[0].arrivals[1]

    let mockState = {}
    let mockGoogle = {}
    let mockArrivalTime = "5 minutes"
    let mockColor = "#000000"

    this.wrapper = shallow(<ArrivalListItem google={mockGoogle} key={1} color={mockColor} arrival={mockArrival} arrivalTime={mockArrivalTime} state={mockState} />)
    this.wrapper2 = shallow(<ArrivalListItem google={mockGoogle} key={1} color={mockColor} arrival={mockArrival2} arrivalTime={mockArrivalTime} state={mockState} />)

  })

  it("should show distance of next bus", function() {
    assert.equal(this.wrapper.find(".arrival-bus-distance").text(), "1.94 miles away")
  })

  it("should show estimated arrival time", function() {
    assert.equal(this.wrapper.find(".arrival-est-time").text(), "5 minutes")

  })

  it("should show full sign info", function() {
    assert.equal(this.wrapper.find(".arrival-name").text(), "15 To SW 5th & Washington")
  })

  it("should show direction traveling", function(){
    assert.equal(this.wrapper.find(".arrival-direction").text(), "Traveling: W")
  })

  it("should Distance in Feet", function() {
    assert.equal(this.wrapper2.find(".arrival-bus-distance").text(), "100 feet away")
  })

})
