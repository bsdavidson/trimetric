import React from "react";
import {connect} from "react-redux";
// import PropTypes from "prop-types";
import {withRouter} from "react-router-dom";

// import {TrimetricPropTypes} from "./prop_types";
// import {Marker} from "react-map-gl";

export class StopMarkers extends React.Component {
  constructor(props) {
    super(props);
    this.markers = [];
  }

  render() {
    return null;
    // let {stops} = this.props;

    // if (!stops || stops.length > 300) {
    //   return null;
    // }

    // let fontSize = 17;
    // 559.6438 -
    // 72.33213 * this.context.viewport.zoom +
    // 2.38382 * Math.pow(this.context.viewport.zoom, 2);

    //   return (
    //     <div>
    //       {stops.map(s => (
    //         <Marker
    //           key={"stop2-" + s.id}
    //           latitude={s.lat}
    //           longitude={s.lng}
    //           offsetLeft={-12}
    //           offsetTop={-12}>
    //           <span style={{fontSize: fontSize}} className="fui-location" />
    //         </Marker>
    //       ))}
    //     </div>
    //   );
  }
}

// StopMarkers.contextTypes = {
//   viewport: PropTypes.object
// };

// StopMarkers.propTypes = {
//   stops: TrimetricPropTypes.stops
// };

function mapStateToProps(state) {
  return {
    locationClicked: state.locationClicked
  };
}

export default withRouter(connect(mapStateToProps)(StopMarkers));
