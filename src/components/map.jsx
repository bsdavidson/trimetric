import deepEqual from "deep-equal"
import GoogleMapsLoader from "google-maps"
import React from "react"
import ReactDOM from "react-dom"


export class Map extends React.Component {
  constructor(props) {
    super(props)
    this.prevOpts = null
    this.state = {
      google: null,
      map: null
    }
  }

  componentDidMount() {
    GoogleMapsLoader.KEY = this.props.apiKey
    GoogleMapsLoader.load(google => {
      if (this.props.onGoogle) {
        this.props.onGoogle(google)
      }
      this.setState({
        google
      })
      this.createOrUpdateMap()
    })
  }

  componentDidUpdate() {
    this.createOrUpdateMap()
  }

  createOrUpdateMap() {
    let opts = Object.assign({}, Map.defaultProps.opts, this.props.opts)
    if (this.state.map) {
      if (!deepEqual(this.prevOpts, opts)) {
        this.prevOpts = opts
        this.state.map.setOptions(opts)
      }
    } else if (this.state.google) {
      let node = ReactDOM.findDOMNode(this.refs.map)
      let map = new this.state.google.maps.Map(node, opts)
      this.setState({
        map
      })
      this.prevOpts = opts
    }
  }

  render() {
    return (
      <div style={this.props.style} className={this.props.className}>
        <div className="map" ref="map">Loading map...</div>
        {this.renderChildren()}
      </div>
    )
  }

  renderChildren() {
    if (!this.state) {
      return this.props.children
    }
    return React.Children.map(this.props.children, (child) => {
      return React.cloneElement(child, {
        map: this.state.map,
        google: this.state.google
      })
    })
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
}

Map.propTypes = {
  apiKey: React.PropTypes.string.isRequired,
  children: React.PropTypes.node,
  className: React.PropTypes.string,
  onGoogle: React.PropTypes.func,
  opts: React.PropTypes.object,
  style: React.PropTypes.object
}

export class Marker extends React.Component {
  constructor(props) {
    super(props)
    this.prevOpts = null
    this.state = {
      marker: null
    }
  }

  componentDidMount() {
    this.createOrUpdateMarker()
  }

  componentDidUpdate() {
    this.createOrUpdateMarker()
  }

  componentWillUnmount() {
    if (this.state.marker) {
      this.state.marker.setMap(null)
    }
  }

  createOrUpdateMarker() {
    let opts = Object.assign({}, this.props.opts, {
      map: this.props.map
    })
    if (!opts.position.lat || !opts.position.lng) {
      return
    }
    if (this.state.marker) {
      if (!deepEqual(this.prevOpts, opts)) {
        this.prevOpts = opts
        this.state.marker.setOptions(opts)
      }
    } else if (this.props.google) {
      let marker = new this.props.google.maps.Marker(opts)
      this.setState({
        marker
      })
      this.prevOpts = opts
    }
  }

  render() {
    return null
  }
}

Marker.propTypes = {
  google: React.PropTypes.object,
  map: React.PropTypes.object,
  opts: React.PropTypes.shape({
    position: React.PropTypes.object.isRequired
  }).isRequired
}
