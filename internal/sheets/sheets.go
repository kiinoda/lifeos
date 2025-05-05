package sheets

import (
	"context"
	"fmt"

	"github.com/kiinoda/lifeos/internal/config"
	"github.com/kiinoda/lifeos/internal/events"
	"google.golang.org/api/option"
	gs "google.golang.org/api/sheets/v4"
)

func GetEvents(ctx context.Context, sheetName string) ([]events.Event, error) {
	result := []events.Event{}

	cfg, err := config.ConfigFromContext(ctx)
	if err != nil {
		return result, err
	}

	apiKey := cfg.ApiKey
	if apiKey == "" {
		return nil, fmt.Errorf("apiKey not found")
	}

	spreadsheetId := cfg.SpreadsheetId
	if spreadsheetId == "" {
		return nil, fmt.Errorf("spreadsheetId not found")
	}

	// Create the sheets service using API key
	service, err := gs.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("Unable to create sheets service: %w", err)
	}

	// Read the data from the sheet
	readRange := fmt.Sprintf("%s!A1:Z1000", sheetName)
	resp, err := service.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve data from sheet: %w", err)
	}

	// Process the data
	if len(resp.Values) == 0 {
		return result, fmt.Errorf("Got empty response from upstream")
	}

	// Actual data starts at row 3; the rest are headers
	for i := 3; i < len(resp.Values); i++ {
		if len(resp.Values[i]) > 0 {
			event, err := events.NewEvent(resp.Values[i])
			if err != nil {
				return result, err
			}
			result = append(result, event)
		}
	}
	return result, nil
}

func GetEventSchedule(ctx context.Context, sheetName string) ([]events.ScheduledEvent, error) {
	result := []events.ScheduledEvent{}

	cfg, err := config.ConfigFromContext(ctx)
	if err != nil {
		return result, err
	}

	apiKey := cfg.ApiKey
	if apiKey == "" {
		return nil, fmt.Errorf("apiKey not found")
	}

	spreadsheetId := cfg.SpreadsheetId
	if spreadsheetId == "" {
		return nil, fmt.Errorf("spreadsheetId not found")
	}

	// Create the sheets service using API key
	service, err := gs.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("Unable to create sheets service: %w", err)
	}

	// Read the data from the sheet
	readRange := fmt.Sprintf("%s!A1:Z1000", sheetName)
	resp, err := service.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve data from sheet: %w", err)
	}

	// Process the data
	if len(resp.Values) == 0 {
		return result, fmt.Errorf("Got empty response from upstream")
	}

	// Actual future events start at row 2; row 1 is header
	for i := 1; i < len(resp.Values); i++ {
		if len(resp.Values[i]) > 0 {
			event, err := events.NewScheduledEvent(resp.Values[i])
			if err != nil {
				return result, err
			}
			result = append(result, event)
		}
	}
	return result, nil
}
