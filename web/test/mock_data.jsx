// Locations 6158, 6160, 8381, 7586, 7642

export function getMockStopsResponse() {
  return {
    resultSet: {
      queryTime: "2016-07-04T11:00:53.988-0700",
      location: [
        {
          lng: -122.67459934134,
          dir: "Westbound",
          lat: 45.519590161127,
          locid: 6158,
          desc: "SW Washington & 3rd"
        }
      ]
    }
  };
}

export function getMockArrivalsResponse() {
  return {
    resultSet: {
      detour: [
        {
          route: [
            {
              route: 15,
              detour: true,
              type: "B",
              desc: "15-Belmont/NW 23rd"
            }
          ],
          info_link_url: "",
          end: 2136445200000,
          id: "39017",
          begin: 1428832800000,
          desc:
            "Buses continue to use the Hawthorne Bridge due to weight restrictions on the Morrison Bridge."
        }
      ],
      arrival: [
        {
          feet: 10225,
          inCongestion: false,
          departed: true,
          scheduled: 1467656007000,
          loadPercentage: 0,
          shortSign: "15 To SW 5th & Washington",
          estimated: 1467655868000,
          detoured: true,
          tripID: "6484467",
          dir: 0,
          blockID: 1504,
          detour: ["39017"],
          route: 15,
          piece: "1",
          fullSign: "15  To SW 5th & Washington",
          id: "6484467_40407_4",
          vehicleID: "2629",
          locid: 6158,
          newTrip: false,
          status: "estimated"
        },
        {
          feet: 100,
          inCongestion: false,
          departed: true,
          scheduled: 1467656907000,
          loadPercentage: 0,
          shortSign: "15 To NW Yeon-44th",
          estimated: 1467656942000,
          detoured: true,
          tripID: "6484468",
          dir: 0,
          blockID: 1502,
          detour: ["39017"],
          route: 15,
          piece: "1",
          fullSign: "15  NW 23rd Ave to NW Yeon & 44th via Montgomery Park",
          id: "6484468_41307_4",
          vehicleID: null,
          locid: 6158,
          newTrip: false,
          status: "estimated"
        },
        {
          feet: null,
          departed: false,
          scheduled: 1467725809000,
          shortSign: "51 To Dosch Rd",
          detoured: false,
          tripID: "6487892",
          dir: 0,
          blockID: 5103,
          route: 51,
          piece: "1",
          fullSign: "51  Vista to Dosch Rd",
          id: "6487892_23809_5",
          vehicleID: null,
          locid: 6158,
          newTrip: false,
          status: "scheduled"
        },
        {
          feet: null,
          departed: false,
          scheduled: 1467726362000,
          shortSign: "51 Council Crest",
          detoured: false,
          tripID: "6487906",
          dir: 1,
          blockID: 5101,
          route: 51,
          piece: "1",
          fullSign: "51  Vista to Council Crest",
          id: "6487906_24362_5",
          vehicleID: null,
          locid: 6158,
          newTrip: false,
          status: "scheduled"
        }
      ],
      queryTime: 1467655255669,
      location: [
        {
          lng: -122.67459934134,
          passengerCode: "E",
          id: 6158,
          dir: "Westbound",
          lat: 45.519590161127,
          desc: "SW Washington & 3rd"
        }
      ]
    }
  };
}

export function getMockVehiclesResponse() {
  return {
    resultSet: {
      queryTime: 1467954664375,
      vehicle: [
        {
          expires: 1467954885000,
          signMessage: "15 Montgomery Pk",
          serviceDate: 1467874800000,
          loadPercentage: 0,
          latitude: 45.5163309,
          nextStopSeq: 24,
          type: "bus",
          blockID: 1508,
          signMessageLong: "15  Belmont/NW 23rd to Montgomery Park",
          lastLocID: 7899,
          nextLocID: 6446,
          locationInScheduleDay: 79713,
          newTrip: false,
          longitude: -122.5871591,
          direction: 0,
          inCongestion: false,
          routeNumber: 15,
          bearing: 270,
          garage: "CENTER",
          tripID: "6484445",
          delay: -132,
          extraBlockID: null,
          messageCode: 138,
          lastStopSeq: 23,
          vehicleID: 2629,
          time: 1467954646616,
          offRoute: false
        }
      ]
    }
  };
}

export function adjustTime(combinedData) {
  let newCombinedData = Object.assign({}, combinedData);
  let {queryTime} = newCombinedData;
  let now = new Date().getTime();

  let newStops = newCombinedData.stops.map(function(stop) {
    let newStop = Object.assign({}, stop);
    let newArrivals = newStop.arrivals.map(function(arrival) {
      let newArrival = Object.assign({}, arrival);
      newArrival.estimated = newArrival.estimated - queryTime + now;
      newArrival.scheduled = newArrival.scheduled - queryTime + now;
      return newArrival;
    });
    newStop.arrivals = newArrivals;
    return newStop;
  });
  newCombinedData.stops = newStops;
  return newCombinedData;
}

export function getMockState() {
  return {
    queryTime: 1467655255669,
    stops: {},
    vehicles: {},
    stopID: null,
    location: {lat: 45.519316, lng: -122.6755836},
    locationClicked: null
  };
}

export function getMockCombinedData() {
  return {
    queryTime: 1467655255669,
    stops: [
      {
        lng: -122.67459934134,
        lat: 45.519590161127,
        locid: 6158,
        desc: "SW Washington & 3rd",
        arrivals: [
          {
            feet: 10225,
            scheduled: 1467656007000,
            shortSign: "15 To SW 5th & Washington",
            estimated: 1467655868000,
            route: 15,
            id: "6484467_40407_4",
            status: "estimated",
            latitude: 45.5163309,
            longitude: -122.5871591,
            signMessage: "15 Montgomery Pk",
            type: "bus",
            vehicleID: "2629",
            bearing: 270
          },
          {
            feet: 100,
            scheduled: 1467656907000,
            shortSign: "15 To NW Yeon-44th",
            estimated: 1467656942000,
            route: 15,
            id: "6484468_41307_4",
            status: "estimated",
            latitude: 0,
            longitude: 0,
            signMessage: undefined,
            type: undefined,
            vehicleID: null,
            bearing: undefined
          }
        ]
      }
    ],
    vehicles: {
      arrivals: [
        {
          bearing: 270,
          blockID: 1508,
          delay: -132,
          direction: 0,
          expires: 1467954885000,
          extraBlockID: null,
          garage: "CENTER",
          inCongestion: false,
          lastLocID: 7899,
          lastStopSeq: 23,
          latitude: 45.5163309,
          loadPercentage: 0,
          locationInScheduleDay: 79713,
          longitude: -122.5871591,
          messageCode: 138,
          newTrip: false,
          nextLocID: 6446,
          nextStopSeq: 24,
          offRoute: false,
          routeNumber: 15,
          serviceDate: 1467874800000,
          signMessage: "15 Montgomery Pk",
          signMessageLong: "15  Belmont/NW 23rd to Montgomery Park",
          time: 1467954646616,
          tripID: "6484445",
          type: "bus",
          vehicleID: 2629
        }
      ]
    }
  };
}
