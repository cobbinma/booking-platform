package venueAPI

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cobbinma/booking/lib/table_api/models"
	"io/ioutil"
	"net/http"
)

var errVenueNotFound = fmt.Errorf("venue not found")

type venueAPI struct {
	client  *http.Client
	baseURL string
}

func NewVenueAPI(baseURL string) models.VenueClient {
	client := http.DefaultClient
	return &venueAPI{client: client, baseURL: baseURL}
}

func (t venueAPI) GetVenue(ctx context.Context, id models.VenueID) (*models.Venue, error) {
	resp, err := t.client.Get(fmt.Sprintf("%s/venues/%v", t.baseURL, id))
	if err != nil {
		return nil, fmt.Errorf("%s : %w", "could not perform get request", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errVenueNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("incorrect response from api : %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body")
	}
	defer resp.Body.Close()

	venue := &models.Venue{}
	if err := json.Unmarshal(body, venue); err != nil {
		return nil, fmt.Errorf("%s : %w", "could not unmarshall body", err)
	}

	return venue, nil
}

func ErrVenueNotFound(err error) bool {
	return errors.Is(err, errVenueNotFound)
}
