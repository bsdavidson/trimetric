package logic

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bsdavidson/trimetric/trimet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProduceTripUpdates(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	var httpCalls int

	pb, err := ioutil.ReadFile("./testdata/trip_updates.pb")
	require.NoError(t, err)

	tuds := TripUpdateSQLDataset{DB: db}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpCalls++
		if httpCalls == 1 {
			cancel()
		}
		w.Write(pb)
	}))

	mp := &mockProducer{T: t}

	err = ProduceTripUpdates(ctx, ts.URL, "123", mp)
	require.NoError(t, err)

	assert.Equal(t, 1, len(mp.bytes))
	var vp trimet.TripUpdatesMsg
	for _, b := range mp.bytes {
		_, err := vp.UnmarshalMsg(b)
		require.NoError(t, err)
	}

	err = tuds.UpdateTripUpdates(vp.TripUpdates)
	require.NoError(t, err)

	tus, err := tuds.FetchTripUpdates()
	require.NoError(t, err)
	assert.Equal(t, 818, len(tus))

}
