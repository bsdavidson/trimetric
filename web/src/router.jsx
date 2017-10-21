import React from "react";
import {Provider} from "react-redux";
import {Route, Switch, BrowserRouter} from "react-router-dom";
import App from "./components/app";

export function render(store, {Router = BrowserRouter, history}) {
  return (
    <Provider store={store}>
      <Router history={history}>
        <Switch>
          <Route exact path="/" component={App} />
          <Route path="/stop/:stopID" component={App} />
        </Switch>
      </Router>
    </Provider>
  );
}
