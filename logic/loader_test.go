package logic

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoader(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, cancel := context.WithCancel(context.Background())
	var httpCalls int

	zd, err := ioutil.ReadFile("./testdata/gtfs.zip")
	require.NoError(t, err)

	lds := LoaderSQLDataset{DB: db}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpCalls++
		if httpCalls == 1 {
			cancel()
		}
		w.Write(zd)
	}))

	err = lds.LoadGTFSData(ts.URL)
	require.NoError(t, err)

}
