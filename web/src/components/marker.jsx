import deepEqual from "deep-equal";
import React from "react";
import PropTypes from "prop-types";

export default class Marker extends React.Component {
  constructor(props) {
    super(props);
    this.prevOpts = null;
    this.state = {
      marker: null
    };
  }

  componentDidMount() {
    this.createOrUpdateMarker();
  }

  componentDidUpdate() {
    this.createOrUpdateMarker();
  }

  componentWillUnmount() {
    if (this.state.marker) {
      this.state.marker.setMap(null);
    }
  }

  createOrUpdateMarker() {
    let opts = Object.assign({}, this.props.opts, {
      map: this.props.map
    });
    if (!opts.position.lat || !opts.position.lng) {
      return;
    }
    if (this.state.marker) {
      if (!deepEqual(this.prevOpts, opts)) {
        this.prevOpts = opts;
        this.state.marker.setOptions(opts);
      }
    } else if (this.props.google) {
      let marker = new this.props.google.maps.Marker(opts);
      this.setState({
        marker
      });
      this.prevOpts = opts;
    }
  }

  render() {
    return null;
  }
}

Marker.propTypes = {
  google: PropTypes.object,
  map: PropTypes.object,
  opts: PropTypes.shape({
    position: PropTypes.object.isRequired
  }).isRequired
};
