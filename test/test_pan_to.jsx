import React from "react";
import {assert} from "chai";
import {shallow} from "enzyme";
import {PanTo} from "../src/components/pan_to";

describe("<PanTo />", function() {
  it("should add and remove event listeners when map changes", function() {
    let removedListener;
    let addListenerEvent;
    let addListenerCallback;
    let mockGoogle = {
      maps: {
        event: {
          removeListener: function(listener) {
            removedListener = listener;
          }
        }
      }
    };
    let mockMap = {
      addListener: function(event, callback) {
        if (event) {
          addListenerEvent = "event";
        }
        if (callback) {
          addListenerCallback = "called";
        }
        return {mock: "mockMap1"};
      }
    };
    let mockMap2 = {
      addListener: function(event, callback) {
        if (event) {
          addListenerEvent = "event";
        }
        if (callback) {
          addListenerCallback = "called";
        }
        return {mock: "mockMap2"};
      }
    };
    let wrapper = shallow(<PanTo location={null} />);
    wrapper.setProps({google: mockGoogle, map: mockMap});
    assert.equal(addListenerEvent, "event");
    assert.equal(addListenerCallback, "called");
    wrapper.setProps({google: mockGoogle, map: mockMap2});
    assert.equal(removedListener.mock, "mockMap1");
  });
});
