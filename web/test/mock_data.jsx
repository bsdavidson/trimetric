// Locations 6158, 6160, 8381, 7586, 7642

export function getMockStopsResponse() {
  return [
    {
      id: "6158",
      code: "6158",
      name: "SW Washington \u0026 3rd",
      desc: "Westbound stop in Portland (Stop ID 6158)",
      lat: 45.51959,
      lng: -122.674599,
      zone_id: "B",
      stop_url: "http://trimet.org/#tracker/stop/6158",
      location_type: 0,
      parent_station: "",
      direction: "West",
      position: "Nearside",
      wheelchair_boarding: 0,
      distance: 0.032123862
    }
  ];
}

export function getMockArrivalsResponse() {
  return [
    {
      route_id: "15",
      route_short_name: "15",
      route_long_name: "Belmont/NW 23rd",
      route_type: 3,
      route_color: "",
      route_text_color: "",
      trip_id: "7743411",
      stop_id: "6158",
      headsign: "NW Yeon \u0026 44th via Montgomery Park",
      arrival_time: "21:09:27",
      departure_time: "21:09:27",
      vehicle_id: "2307",
      vehicle_label: "15 To NW Yeon",
      vehicle_position: {
        lat: 45.516563,
        lng: -122.61697,
        bearing: 270,
        odometer: 0,
        speed: 0
      },
      date: "2017-11-18T00:00:00Z"
    },
    {
      route_id: "15",
      route_short_name: "15",
      route_long_name: "Belmont/NW 23rd",
      route_type: 3,
      route_color: "",
      route_text_color: "",
      trip_id: "7743412",
      stop_id: "6158",
      headsign: "NW Thurman St",
      arrival_time: "21:49:27",
      departure_time: "21:49:27",
      vehicle_id: "3047",
      vehicle_label: "72 Clackamas TC",
      vehicle_position: {
        lat: 45.530796,
        lng: -122.56394,
        bearing: 350,
        odometer: 0,
        speed: 0
      },
      date: "2017-11-18T00:00:00Z"
    }
  ];
}

export function getMockVehiclesResponse() {
  return [
    {
      trip: {
        trip_id: "7753585",
        route_id: "100"
      },
      vehicle: {
        id: "2629",
        label: "15 Montgomery Pk"
      },
      position: {
        lat: 45.5163309,
        lng: -122.5871591,
        bearing: 270,
        odometer: 0,
        speed: 0
      },
      current_stop_sequence: 9,
      stop_id: "8355",
      current_status: 2,
      timestamp: 1510378074,
      congestion_level: 0,
      occupancy_status: 0,
      route_type: 0
    }
  ];
}

export function getMockState() {
  return {
    queryTime: 1467655255669,
    stops: {},
    vehicles: {},
    stopID: null,
    location: {lat: 45.519316, lng: -122.6755836, locationType: "HOME"},
    locationClicked: null
  };
}
