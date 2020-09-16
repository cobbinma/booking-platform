package tableAPI

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cobbinma/booking/lib/booking_api/config"
	"github.com/cobbinma/booking/lib/booking_api/models"
	"io/ioutil"
	"net/http"
)

var _ models.TableClient = (*tableAPI)(nil)

type tableAPI struct {
	client *http.Client
}

func NewTableAPI() models.TableClient {
	client := http.DefaultClient
	return &tableAPI{client: client}
}

func (t tableAPI) GetTable(ctx context.Context, id models.TableID) (*models.Table, error) {
	venue, ok := ctx.Value(models.VenueCtxKey).(models.Venue)
	if !ok {
		return nil, fmt.Errorf("venue was not in context")
	}

	resp, err := t.client.Get(fmt.Sprintf("%s/venues/%v/tables/%v", config.TableAPIRoot(), venue.ID, id))
	if err != nil {
		return nil, fmt.Errorf("%s : %w", "could not perform get request", err)
	}

	if resp.StatusCode != http.StatusOK || resp.Body == nil {
		message := fmt.Sprintf("incorrect response from api")
		return nil, fmt.Errorf("%s : %v", message, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body")
	}

	table := &models.Table{}
	if err := json.Unmarshal(body, table); err != nil {
		return nil, fmt.Errorf("%s : %w", "could not unmarshall body", err)
	}

	return table, nil
}

func (t tableAPI) GetTablesWithCapacity(ctx context.Context, capacity int) ([]models.Table, error) {
	venue, ok := ctx.Value(models.VenueCtxKey).(models.Venue)
	if !ok {
		return nil, fmt.Errorf("venue was not in context")
	}

	resp, err := t.client.Get(fmt.Sprintf("%s/venues/%v/tables/capacity/%v", config.TableAPIRoot(), venue.ID, capacity))
	if err != nil {
		return nil, fmt.Errorf("%s : %w", "could not perform get request", err)
	}

	if resp.StatusCode != http.StatusOK || resp.Body == nil {
		message := fmt.Sprintf("incorrect response from api")
		return nil, fmt.Errorf("%s : %v", message, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body")
	}

	tables := []models.Table{}
	if err := json.Unmarshal(body, &tables); err != nil {
		return nil, fmt.Errorf("%s : %w", "could not unmarshall body", err)
	}

	return tables, nil
}
