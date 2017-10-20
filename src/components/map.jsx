import deepEqual from "deep-equal";
import GoogleMapsLoader from "google-maps";
import React from "react";
import PropTypes from "prop-types";

export class Map extends React.Component {
  constructor(props) {
    super(props);
    this.prevOpts = null;
    this.state = {
      google: null,
      map: null
    };

    this.setMapNode = this.setMapNode.bind(this);
  }

  componentDidMount() {
    GoogleMapsLoader.KEY = this.props.apiKey;
    GoogleMapsLoader.load(google => {
      if (this.props.onGoogle) {
        this.props.onGoogle(google);
      }
      this.setState({
        google
      });
      this.createOrUpdateMap();
    });
  }

  componentDidUpdate() {
    this.createOrUpdateMap();
  }

  createOrUpdateMap() {
    let opts = Object.assign({}, Map.defaultProps.opts, this.props.opts);
    if (this.state.map) {
      if (!deepEqual(this.prevOpts, opts)) {
        this.prevOpts = opts;
        this.state.map.setOptions(opts);
      }
    } else if (this.state.google) {
      let map = new this.state.google.maps.Map(this.mapNode, opts);
      this.setState({
        map
      });
      this.prevOpts = opts;
    }
  }

  render() {
    return (
      <div style={this.props.style} className={this.props.className}>
        <div className="map" ref={this.setMapNode}>
          Loading map...
        </div>
        {this.renderChildren()}
      </div>
    );
  }

  renderChildren() {
    if (!this.state) {
      return this.props.children;
    }
    return React.Children.map(this.props.children, child => {
      return React.cloneElement(child, {
        map: this.state.map,
        google: this.state.google
      });
    });
  }

  setMapNode(node) {
    this.mapNode = node;
  }
}

Map.defaultProps = {
  opts: {
    center: {
      lat: 53.2238484,
      lng: -4.195443
    },
    zoom: 14,
    fullscreenControl: true
  }
};

Map.propTypes = {
  apiKey: PropTypes.string.isRequired,
  children: PropTypes.node,
  className: PropTypes.string,
  onGoogle: PropTypes.func,
  opts: PropTypes.object,
  style: PropTypes.object
};
