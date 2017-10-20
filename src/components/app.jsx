import moment from "moment";
import React, {Component} from "react";
import {withRouter} from "react-router-dom";
// import PropTypes from "prop-types";
import {connect} from "react-redux";
import ReactCSSTransitionGroup from "react-addons-css-transition-group";

import ArrivalList from "./arrival_list";
import {ArrivalMarkers} from "./arrival_markers";
import {Map} from "./map";
import Marker from "./marker";
import {PanTo} from "./pan_to";
import StopList from "./stop_list";

class App extends Component {
  constructor(props) {
    super(props);
    this.handleGoogle = this.handleGoogle.bind(this);
    this.state = {
      google: null
    };
  }

  handleGoogle(google) {
    this.setState({
      google
    });
  }

  render() {
    let {location, locationClicked, stops, queryTime, vehicles} = this.props;
    if (!stops) {
      return <div>No stops</div>;
    }
    let stopID;
    if (this.props.match && this.props.match.params) {
      stopID = this.props.match.params.stopID;
    }
    let stop = stops.find(s => s.locid == stopID);
    let page;
    let markers = [];

    if (stop) {
      if (stop.arrivals) {
        markers.push(<ArrivalMarkers key="arrivals" stop={stop} />);
      }
      page = <ArrivalList key="arrivalPage" stop={stop} />;
    } else {
      page = <StopList key="stopPage" stops={stops} />;
      if (vehicles.arrivals) {
        markers.push(<ArrivalMarkers key="allTrimet" stop={vehicles} />);
      }
    }

    if (this.state.google) {
      let opts = {
        position: location,
        title: "WeWork",
        animation: this.state.google.maps.Animation.DROP
      };
      markers.push(<Marker key="home" opts={opts} />);
    }
    return (
      <div className="app">
        <nav>
          <div className="app-query-time">
            Updated: {moment(queryTime).format("h:mm:ss a")}
          </div>
        </nav>
        <Map
          className="app-map"
          onGoogle={this.handleGoogle}
          apiKey={process.env.GOOGLE_MAPS_API_KEY}
          opts={{zoom: 16, center: location}}>
          {markers}
          <PanTo location={locationClicked} />
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

// AppX.propTypes = {
//   route: PropTypes.object.isRequired,
//   params: PropTypes.object
// };

function mapStateToProps(state) {
  return {
    location: state.location,
    locationClicked: state.locationClicked,
    stops: state.stops,
    queryTime: state.queryTime,
    vehicles: state.vehicles
  };
}

export default withRouter(connect(mapStateToProps)(App));
