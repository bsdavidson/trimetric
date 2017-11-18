import React from "react";
import {connect} from "react-redux";

function Header(props) {
  return (
    <div className="header">
      <h1 className="header-title">Trimetric</h1>

      <div className="header-metrics">
        <div className="vehicle-metric">
          There are currently {props.vehicles.length} vehicles active.
        </div>
      </div>
    </div>
  );
}

function mapStateToProps(state) {
  return {
    vehicles: state.vehicles
  };
}

export default connect(mapStateToProps)(Header);
