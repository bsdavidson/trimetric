import { assert } from "chai"

import { combineResponses } from "../src/data"
import { getMockStopsResponse, getMockArrivalsResponse, getMockVehiclesResponse, getMockCombinedData } from "./mock_data"

describe("Data", () => {
  it("Should take two datasets and output formatted json", function() {
    let output = combineResponses(getMockStopsResponse(), getMockArrivalsResponse(), getMockVehiclesResponse())
    assert.deepEqual(output, getMockCombinedData())
  })
})
