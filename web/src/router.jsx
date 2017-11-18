import React from "react";
import {Provider} from "react-redux";
import {Route, Switch, BrowserRouter} from "react-router-dom";
import App from "./components/app";

export function render(store, onStopChange, {Router = BrowserRouter, history}) {
  function renderApp(props) {
    return <App onStopChange={onStopChange} {...props} />;
  }

  return (
    <Provider store={store}>
      <Router history={history}>
        <Switch>
          <Route exact path="/" render={renderApp} />
          <Route path="/stop/:stopID" render={renderApp} />
        </Switch>
      </Router>
    </Provider>
  );
}
