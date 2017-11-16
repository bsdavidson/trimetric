import React from "react";
import PropTypes from "prop-types";
import {connect} from "react-redux";
import {withRouter} from "react-router-dom";

import StopListItem from "./stop_list_item";
import {updateHomeLocation} from "../actions";

export class StopList extends React.Component {
  constructor(props) {
    super(props);

    this.handleCurrentLocationClick = this.handleCurrentLocationClick.bind(
      this
    );
  }

  createStops(stops) {
    stops = stops
      .filter(s => s.arrivals.length > 0)
      .map(s => <StopListItem key={s.id} stop={s} />);
    if (!stops.length) {
      return "Sorry, no buses are running near you. Better start walking or call an Uber.";
    }
    return stops;
  }

  handleCurrentLocationClick() {
    if ("geolocation" in navigator) {
      navigator.geolocation.watchPosition(position => {
        this.props.onLocationClick(position);
      });
    } else {
      alert("Your browser doesn't support geolocation");
    }
  }

  render() {
    return (
      <div className="stop-list">
        <div
          className="use-current"
          title="Use Current GPS position"
          onClick={this.handleCurrentLocationClick}>
          <span className="fui-location" />
        </div>
        <h3>Nearby Stops</h3>
        <div className="stop-list-items">
          {this.createStops(this.props.stops)}
        </div>
      </div>
    );
  }
}

StopList.propTypes = {
  stops: PropTypes.arrayOf(
    PropTypes.shape({
      arrivals: PropTypes.array.isRequired,
      id: PropTypes.string.isRequired
    })
  ).isRequired
};

function mapDispatchToProps(dispatch) {
  return {
    onLocationClick: position => {
      dispatch(
        updateHomeLocation(position.coords.latitude, position.coords.longitude)
      );
    }
  };
}

export default withRouter(connect(null, mapDispatchToProps)(StopList));
