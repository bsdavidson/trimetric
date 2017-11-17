import React, {Component} from "react";
import ReactCSSTransitionGroup from "react-addons-css-transition-group";
import moment from "moment";
import {Marker} from "react-map-gl";
import {connect} from "react-redux";
import {withRouter} from "react-router-dom";

import ArrivalList from "./arrival_list";
import Map from "./map";
import StopList from "./stop_list";
import {updateViewport} from "../actions";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      mapWidth: 1,
      mapHeight: 1
    };

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
    this.setState({
      mapWidth: mapbox.clientWidth,
      mapHeight: mapbox.clientHeight
    });
  }

  render() {
    let {location, stops, queryTime} = this.props;
    if (!stops) {
      return <div>No stops</div>;
    }
    let stopID;
    if (this.props.match && this.props.match.params) {
      stopID = this.props.match.params.stopID;
    }
    let stop = stops.find(s => s.id == stopID);
    let page;
    let markers = [];
    if (stop) {
      page = <ArrivalList stop={stop} />;
    } else if (this.props.zoom > 15.5) {
      page = <StopList stops={stops} />;
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
        <nav>
          <div className="app-query-time">
            Updated: {moment(queryTime).format("h:mm:ss a")}
          </div>
        </nav>
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
