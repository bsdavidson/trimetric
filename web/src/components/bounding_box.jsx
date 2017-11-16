import React from "react";
import {connect} from "react-redux";
import PropTypes from "prop-types";

import {updateBoundingBox} from "../actions";

class BoundingBox extends React.Component {
  constructor(props) {
    super(props);

    this.timeout = null;
    this.boundsChangedListener = null;
    this.updateBoundingBox = this.updateBoundingBox.bind(this);
    this.addListener = this.addListener.bind(this);
    this.removeListener = this.removeListener.bind(this);
  }

  removeListener() {
    if (!this.boundsChangedListener || !this.props.google) {
      return;
    }
    this.props.google.maps.event.removeListener(this.boundsChangedListener);
  }

  addListener(map) {
    if (!map) {
      return;
    }
    this.removeListener();
    this.boundsChangedListener = this.props.google.maps.event.addListener(
      map,
      "bounds_changed",
      () => this.updateBoundingBox(map.getBounds().toJSON())
    );
  }

  updateBoundingBox(bounds) {
    if (this.timeout) {
      clearTimeout(this.timeout);
    }
    this.timeout = setTimeout(() => {
      this.props.onUpdate(bounds);
    }, 1000);
  }

  componentWillUpdate(nextProps) {
    this.addListener(nextProps.map);
  }

  componentDidMount() {
    this.addListener(this.props.map);
  }

  render() {
    return null;
  }
}

BoundingBox.propTypes = {
  google: PropTypes.object,
  map: PropTypes.object
};

function mapDispatchToProps(dispatch) {
  return {
    onUpdate: bounds => {
      return dispatch(updateBoundingBox(bounds));
    }
  };
}

export default connect(null, mapDispatchToProps)(BoundingBox);
