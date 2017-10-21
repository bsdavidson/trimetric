import {assert} from "chai";
import {mount} from "enzyme";
import {createStore} from "redux";
// import {Router} from "react-router-dom";
import {MemoryRouter} from "react-router-dom";
import {reducer} from "../src/store";
import {render} from "../src/router";
import {configure} from "enzyme";
import Adapter from "enzyme-adapter-react-16";

configure({adapter: new Adapter()});

describe("App", () => {
  it("Render a container", function() {
    let store = createStore(reducer, {stops: [], vehicles: {arrivals: []}});
    let wrapper = mount(render(store, {Router: MemoryRouter}));

    assert.equal(wrapper.find(".app").length, 1);
  });
});
