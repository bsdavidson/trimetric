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
      this.props.arrival.vehicle_id === this.props.locationClicked.id
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
      this.props.locationClicked.id === this.props.arrival.vehicle_id
    ) {
      routeClass = "active";
    }

    let routeStyle = {
      backgroundColor: this.props.color
    };

    return (
      <div
        className={"arrival-list-item " + routeClass}
        onClick={this.handleVehicleClick}>
        <div className="arrival-id" style={routeStyle}>
          {this.props.arrival.route_id}
        </div>
        <div className="arrival-name">{this.props.arrival.headsign}</div>
        <div className="arrival-metrics">
          <div className="arrival-metric arrival-bus-distance">
            {this.props.arrival.vehicle_position.lat},{" "}
            {this.props.arrival.vehicle_position.lng}
          </div>
          <div className="arrival-metric arrival-est-time">
            {this.props.arrival.vehicle_id} - {this.props.arrivalTime}
          </div>
          <div className="arrival-metric arrival-direction">
            Traveling:{" "}
            {degreeToCompass(this.props.arrival.vehicle_position.bearing)}
          </div>
        </div>
      </div>
    );
  }
}

ArrivalListItem.propTypes = {
  arrival: PropTypes.shape({
    vehicle_position: PropTypes.shape({
      lat: PropTypes.number.isRequired,
      lng: PropTypes.number.isRequired,
      bearing: PropTypes.number
    }).isRequired,
    feet: PropTypes.number.isRequired,
    route_id: PropTypes.string.isRequired,
    headsign: PropTypes.string.isRequired,
    vehicle_id: PropTypes.string
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
          arrival.vehicle_id,
          arrival.vehicle_position.lat,
          arrival.vehicle_position.lng,
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
    locationClicked: state.locationClicked,
    vehicles: state.vehicles
  };
}

export default withRouter(
  connect(mapStateToProps, mapDispatchToProps)(ArrivalListItem)
);
