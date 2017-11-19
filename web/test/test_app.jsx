import {assert} from "chai";
import {mount} from "enzyme";
import {createStore} from "redux";
import {MemoryRouter} from "react-router-dom";
import {reducer} from "../src/store";
import {render} from "../src/router";
import {configure} from "enzyme";
import Adapter from "enzyme-adapter-react-16";
import {App} from "../src/components/app";

configure({adapter: new Adapter()});

describe("App", () => {
  it("Render a container", function() {
    let handleStopChange = () => {};

    let store = createStore(reducer, {stops: [], vehicles: {arrivals: []}});
    let wrapper = mount(
      render(store, handleStopChange, {Router: MemoryRouter})
    );

    assert.equal(wrapper.find(App).length, 1);

    let app = wrapper.find(App);
    assert.equal(app.instance().props.onStopChange, handleStopChange);
  });
});
