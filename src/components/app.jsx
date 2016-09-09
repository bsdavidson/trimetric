import moment from "moment"
import React from "react";
import ReactCSSTransitionGroup from "react-addons-css-transition-group"

import { ArrivalList } from "./arrival_list"
import { ArrivalMarkers } from "./arrival_markers"
import { Map, Marker } from "./map"
import { PanTo } from "./pan_to"
import { StopList } from "./stop_list"


export class App extends React.Component {
  constructor(props) {
    super(props)
    this.handleGoogle = this.handleGoogle.bind(this)
    this.state = {
      google: null
    },
    this.center
  }

  handleGoogle(google) {
    this.setState({
      google
    })
  }

  render() {
    let store = this.props.route.store
    let state = store.getState()
    let {location, locationClicked, stops, queryTime, vehicles} = state
    let stopID
    if (this.props.params) {
      stopID = this.props.params.stopID
    }
    let stop = stops.find(s => s.locid == stopID)
    let page
    let markers = []

    if (stop) {
      markers.push(<ArrivalMarkers key="arrivals" stop={stop} state={state} />)
      page = <ArrivalList key="arrivalPage" stop={stop} state={state} />
    } else {
      page = <StopList key="stopPage" stops={stops} state={state} />
      markers.push(<ArrivalMarkers key="allTrimet" stop={vehicles} state={state} />)
    }

    if (this.state.google) {
      let opts = {
        position: location,
        title: "WeWork",
        animation: this.state.google.maps.Animation.DROP
      }
      markers.push(<Marker key="home" opts={opts} />)
    }
    return (
      <div className="app">
        <nav>
          <div className="app-query-time">
            Updated: {moment(queryTime).format("h:mm:ss a")}
          </div>
        </nav>
        <Map className="app-map" onGoogle={this.handleGoogle}
            apiKey={process.env.GOOGLE_MAPS_API_KEY}
            opts={{ zoom: 16, center: location }}>
          {markers}
          <PanTo location={locationClicked} />
        </Map>
        <ReactCSSTransitionGroup component="div" transitionName="page" transitionEnterTimeout={700} transitionLeaveTimeout={700}>
          {page}
        </ReactCSSTransitionGroup>
      </div>
    )
  }
}

App.propTypes = {
  route: React.PropTypes.object.isRequired,
  params: React.PropTypes.object
}
