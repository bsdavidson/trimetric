import React from "react";
import {assert} from "chai";
import {shallow} from "enzyme";
import {ArrivalMarkers} from "../src/components/arrival_markers";
import {getMockCombinedData} from "./mock_data";

describe("<ArrivalMarkers />", () => {
  it("should render an array of Markers", function() {
    let size = "";
    let mockGoogle = {
      maps: {
        Size: function(w, h) {
          size = {w: w, h: h};
        },
        Animation: {
          DROP: 1
        },
        Marker: function() {
          return "Marker";
        }
      }
    };

    let mockMap = {};
    let mockState = {};
    let stop = Object.assign({}, getMockCombinedData().stops[0]);
    shallow(
      <ArrivalMarkers
        google={mockGoogle}
        map={mockMap}
        stop={stop}
        state={mockState}
      />
    );
    assert.equal(size.w, 25);
    assert.equal(size.h, 25);
  });
});
