import React, {Component} from "react";
import {connect} from "react-redux";

import {updateHomeLocation} from "../actions";

class Header extends Component {
  constructor(props) {
    super(props);

    this.state = {
      gpsLoading: false
    };

    this.handleClick = this.handleClick.bind(this);
  }

  handleClick() {
    if (!("geolocation" in navigator)) {
      this.setState({gpsLoading: false});
      alert("Your browser doesn't support geolocation");
    }

    this.setState({gpsLoading: true});
    let options = {
      enableHighAccuracy: true,
      timeout: 5000,
      maximumAge: 0
    };

    navigator.geolocation.getCurrentPosition(
      position => {
        this.setState({gpsLoading: false});
        this.props.onPositionChange(position);
      },
      err => {
        this.setState({gpsLoading: false});
        throw err;
      },
      options
    );
  }

  render() {
    return (
      <div className="header">
        <h1 className="header-title">Trimetric</h1>
        <div
          className="use-current"
          title="Use Current GPS position"
          onClick={this.handleClick}>
          {this.state.gpsLoading ? (
            <div className="spinner" />
          ) : (
            <span className="fui-location" />
          )}
        </div>

        <div className="header-metrics">
          <div className="vehicle-metric">
            There are currently {this.props.vehicles.length} vehicles active.
          </div>
        </div>
      </div>
    );
  }
}

function mapDispatchToProps(dispatch) {
  return {
    onPositionChange: position => {
      dispatch(
        updateHomeLocation(
          position.coords.latitude,
          position.coords.longitude,
          true
        )
      );
    }
  };
}

function mapStateToProps(state) {
  return {
    vehicles: state.vehicles
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Header);
