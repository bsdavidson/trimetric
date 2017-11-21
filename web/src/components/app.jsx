import React, {Component} from "react";
import ReactCSSTransitionGroup from "react-addons-css-transition-group";
import {Marker} from "react-map-gl";
import {connect} from "react-redux";
import {withRouter} from "react-router-dom";

import Info from "./info";
import ArrivalList from "./arrival_list";
import Map from "./map";
import StopList from "./stop_list";
import {updateViewport} from "../actions";

export class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      mapWidth: 1,
      mapHeight: 1
    };

    this.selectedStop = null;

    this.handleResize = this.handleResize.bind(this);
  }

  componentDidMount() {
    window.addEventListener("resize", this.handleResize, false);
    this.handleResize();
  }

  componentWillUnmount() {
    window.removeEventListener("resize", this.handleResize, false);
  }

  handleResize() {
    let mapbox = document.getElementById("mapbox");
    if (!mapbox) {
      return;
    }
    this.setState({
      mapWidth: mapbox.clientWidth,
      mapHeight: mapbox.clientHeight
    });
  }

  componentWillReceiveProps(nextProps) {
    let newStopID =
      (nextProps.match &&
        nextProps.match.params &&
        nextProps.match.params.stopID) ||
      null;
    let selectedStop = this.selectedStop ? this.selectedStop.id : null;

    if (!newStopID) {
      this.selectedStop = null;
    }

    if (newStopID !== selectedStop) {
      this.selectedStop = this.props.stops.find(s => s.id == newStopID);
      this.props.onStopChange(this.selectedStop);
    }
  }

  render() {
    let {location, stops} = this.props;

    if (!stops) {
      return <div>No stops</div>;
    }

    let page;
    let markers = [];
    if (this.selectedStop && this.props.zoom > 13.5) {
      page = <ArrivalList key="transition-stops" stop={this.selectedStop} />;
    } else if (this.props.zoom > 15.5) {
      page = <StopList key="transition-stops" stops={stops} />;
    } else {
      page = <Info key="transition-info" />;
    }

    markers.push(
      <Marker
        key="home"
        latitude={location.lat}
        longitude={location.lng}
        offsetLeft={-12}
        offsetTop={-12}>
        <span className="fui-user" />
      </Marker>
    );

    return (
      <div className="app">
        <Map
          onViewportChange={this.props.onViewportChange}
          width={this.state.mapWidth}
          height={this.state.mapHeight}>
          {markers}
        </Map>
        <ReactCSSTransitionGroup
          component="div"
          transitionName="page"
          transitionEnterTimeout={700}
          transitionLeaveTimeout={700}>
          {page}
        </ReactCSSTransitionGroup>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    location: state.location,
    locationClicked: state.locationClicked,
    stops: state.stops,
    queryTime: state.queryTime,
    vehicles: state.vehicles,
    zoom: state.zoom
  };
}

function mapDispatchToProps(dispatch) {
  return {
    onViewportChange: (bounds, zoom) => {
      dispatch(updateViewport(bounds, zoom));
    }
  };
}

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(App));
