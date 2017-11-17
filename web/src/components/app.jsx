import React, {Component} from "react";
import ReactCSSTransitionGroup from "react-addons-css-transition-group";
import moment from "moment";
import {Marker} from "react-map-gl";
import {connect} from "react-redux";
import {withRouter} from "react-router-dom";

import ArrivalList from "./arrival_list";
import ArrivalMarkers from "./arrival_markers";
import Map from "./map";
import StopList from "./stop_list";
import VehicleMarkers from "./vehicle_markers";
import {updateBoundingBox} from "../actions";

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
    let {location, stops, queryTime, vehicles} = this.props;
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
      if (stop.arrivals) {
        markers.push(<ArrivalMarkers key="arrivals" stop={stop} />);
      }
      page = <ArrivalList key="arrivalPage" stop={stop} />;
    } else {
      page = <StopList key="stopPage" stops={stops} />;

      stops.forEach(s => {
        markers.push(
          <Marker
            key={"stop-" + s.id}
            latitude={s.lat}
            longitude={s.lng}
            offsetLeft={-12}
            offsetTop={-12}>
            <span className="fui-location" />
          </Marker>
        );
      });

      markers.push(<VehicleMarkers key="allTrimet" vehicles={vehicles} />);
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
          onBoundsChanged={this.props.onBoundsChanged}
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
    vehicles: state.vehicles
  };
}

function mapDispatchToProps(dispatch) {
  return {
    onBoundsChanged: bounds => {
      return dispatch(updateBoundingBox(bounds));
    }
  };
}

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(App));
