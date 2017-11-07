import React from "react";
import PropTypes from "prop-types";
import {withRouter} from "react-router-dom";
import {connect} from "react-redux";

import {degreeToCompass} from "../helpers/directions";
import {updateLocation, LocationTypes} from "../actions";

export class ArrivalListItem extends React.Component {
  constructor(props) {
    super(props);
    this.handleVehicleClick = this.handleVehicleClick.bind(this);
  }

  componentWillUnmount() {
    if (
      this.props.locationClicked &&
      this.props.arrival.vehicleID === this.props.locationClicked.id
    ) {
      this.props.clearLocation(this.props.location);
    }
  }

  handleVehicleClick() {
    if (window) {
      window.scrollTo(0, 0);
    }
    this.props.onVehicleClick(LocationTypes.VEHICLE, this.props.arrival);
  }

  render() {
    let routeClass = "";
    if (
      this.props.locationClicked &&
      this.props.locationClicked.following &&
      this.props.locationClicked.id === this.props.arrival.vehicleID
    ) {
      routeClass = "active";
    }

    let routeStyle = {
      backgroundColor: this.props.color
    };

    let vehicleDistance;
    if (this.props.arrival.feet > 500) {
      vehicleDistance =
        Math.round(this.props.arrival.feet / 5280 * 100) / 100 + " miles away";
    } else {
      vehicleDistance = this.props.arrival.feet + " feet away";
    }
    return (
      <div
        className={"arrival-list-item " + routeClass}
        onClick={this.handleVehicleClick}>
        <div className="arrival-id" style={routeStyle}>
          {this.props.arrival.route}
        </div>
        <div className="arrival-name">{this.props.arrival.shortSign}</div>
        <div className="arrival-metrics">
          <div className="arrival-metric arrival-bus-distance">
            {vehicleDistance}
          </div>
          <div className="arrival-metric arrival-est-time">
            {this.props.arrivalTime}
          </div>
          <div className="arrival-metric arrival-direction">
            Traveling: {degreeToCompass(this.props.arrival.bearing)}
          </div>
        </div>
      </div>
    );
  }
}

ArrivalListItem.propTypes = {
  arrival: PropTypes.shape({
    bearing: PropTypes.number,
    feet: PropTypes.number.isRequired,
    latitude: PropTypes.number.isRequired,
    longitude: PropTypes.number.isRequired,
    route: PropTypes.number.isRequired,
    shortSign: PropTypes.string.isRequired,
    vehicleID: PropTypes.string
  }).isRequired,
  arrivalTime: PropTypes.string.isRequired,
  location: PropTypes.object,
  locationClicked: PropTypes.object
};

function mapDispatchToProps(dispatch) {
  return {
    onVehicleClick: (type, arrival) => {
      dispatch(
        updateLocation(
          type,
          arrival.vehicleID,
          arrival.latitude,
          arrival.longitude,
          true
        )
      );
    },
    clearLocation: location => {
      dispatch(
        updateLocation(
          LocationTypes.HOME,
          null,
          location.lat,
          location.lng,
          false
        )
      );
    }
  };
}

function mapStateToProps(state) {
  return {
    location: state.location,
    locationClicked: state.locationClicked
  };
}

export default withRouter(
  connect(mapStateToProps, mapDispatchToProps)(ArrivalListItem)
);
