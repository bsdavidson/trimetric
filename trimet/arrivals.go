package trimet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ArrivalResponse ...
type ArrivalResponse struct {
	ResultSet ArrivalResultSet `json:"resultSet"`
}

// ArrivalResultSet ...
type ArrivalResultSet struct {
	QueryTime int64         `json:"queryTime"`
	Arrivals  []ArrivalData `json:"arrival"`
}

// ArrivalData ...
type ArrivalData struct {
	LocationID int    `json:"locid"`
	Data       []byte `json:"-"`
}

type ArrivalDataAlias ArrivalData

// UnmarshalJSON ...
func (t *ArrivalData) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, (*ArrivalDataAlias)(t)); err != nil {
		return err
	}
	t.Data = make([]byte, len(b))
	copy(t.Data, b)
	return nil
}

// RequestArrivals ...
func RequestArrivals(apiKey string, ids []int) ([]byte, error) {
	query := url.Values{}

	var sids []string
	for _, id := range ids {
		sids = append(sids, strconv.Itoa(id))
	}

	query.Set("appID", apiKey)
	query.Set("locIDs", strings.Join(sids, ","))
	query.Set("json", "true")

	resp, err := http.Get(fmt.Sprintf("%s?%s", Arrivals, query.Encode()))
	if err != nil {
		return nil, fmt.Errorf("http.Get: %s", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %s", err)
	}

	return b, nil
}
