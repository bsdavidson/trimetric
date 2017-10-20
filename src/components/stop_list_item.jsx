import React from "react";
import PropTypes from "prop-types";
import {connect} from "react-redux";
import {Link, withRouter} from "react-router-dom";

import {formatDistance} from "../helpers/directions";
import {humanTimeUntilEpoch} from "../helpers/times";

export class StopListItem extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    let location = {
      latitude: this.props.stop.lat,
      longitude: this.props.stop.lng
    };
    let arrivalInfo = {
      route: this.props.stop.arrivals[0].route,
      estimated: humanTimeUntilEpoch(this.props.stop.arrivals[0].estimated)
    };
    return (
      <Link className="stop-link" to={`/stop/${this.props.stop.locid}`}>
        <div className="stop-list-item">
          <div className="stop-arrow">
            <span className="fui-arrow-right" />
          </div>
          <h4 className="stop-description">
            {this.props.stop.desc} - #<span className="stop-id">
              {this.props.stop.locid}
            </span>
          </h4>
          <p className="stop-metrics">
            <span className="stop-distance">
              {formatDistance(this.props.location, location)}
            </span>{" "}
            feet away.
            <span className="stop-route"> {arrivalInfo.route} arrives </span>
            <span className="stop-estimate">{arrivalInfo.estimated} </span>
          </p>
        </div>
      </Link>
    );
  }
}

StopListItem.propTypes = {
  location: PropTypes.shape({
    lat: PropTypes.number.isRequired,
    lng: PropTypes.number.isRequired
  }).isRequired,
  stop: PropTypes.shape({
    arrivals: PropTypes.arrayOf(
      PropTypes.shape({
        estimated: PropTypes.number.isRequired,
        route: PropTypes.number.isRequired
      })
    ).isRequired,
    desc: PropTypes.string.isRequired,
    lat: PropTypes.number.isRequired,
    lng: PropTypes.number.isRequired,
    locid: PropTypes.number.isRequired
  }).isRequired
};

function mapStateToProps(state) {
  return {
    location: state.location
  };
}

export default withRouter(connect(mapStateToProps)(StopListItem));
