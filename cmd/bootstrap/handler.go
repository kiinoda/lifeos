package main

import (
	"context"

	"fmt"
	"log"
	"os"
	"slices"
	"time"

	"github.com/kiinoda/lifeos/internal/config"
)

const defaultRegion = "eu-west-1"

func handler(ctx context.Context) error {
	action, _ := os.LookupEnv("LIFEOS_ACTION")
	if !slices.Contains([]string{"daily_schedule", "event_schedule", "event_notification", "invoice_reminder"}, action) {
		return fmt.Errorf("please set LIFEOS_ACTION environment variable")
	}

	ssmPath, _ := os.LookupEnv("LIFEOS_SSM_PATH")
	if ssmPath == "" {
		return fmt.Errorf("please set LIFEOS_SSM_PATH environment variable")
	}

	// The Lambda runtime always sets AWS_REGION; defaultRegion only kicks in for local runs.
	region, _ := os.LookupEnv("AWS_REGION")
	if region == "" {
		region = defaultRegion
	}

	cfg, err := config.NewConfig(region, ssmPath)
	if err != nil {
		return fmt.Errorf("failed to create application config: %w", err)
	}

	ctx = config.ContextWithConfig(ctx, cfg)

	// Go's time.Weekday() is Sun=0, Mon=1...Sat=6, but our internal enum is Mon=0...Sun=6
	// (matching the Mon-first column order in the spreadsheet). We can't change the enum
	// without reordering the sheet, so we remap via (weekday+6)%7 to wrap Sunday correctly.
	// A naive -1 subtraction breaks on Sundays, producing -1 (out of bounds).
	dayOfWeek := time.Weekday((int(time.Now().Weekday()) + 6) % 7)

	if action == "daily_schedule" {
		log.Println("Daily schedule")
		return dailySchedule(ctx, dayOfWeek)
	}

	if action == "event_schedule" {
		log.Println("Event Schedule")
		return eventSchedule(ctx, dayOfWeek)
	}

	if action == "event_notification" {
		log.Println("Event Notification")
		return eventNotification(ctx, dayOfWeek)
	}

	if action == "invoice_reminder" {
		log.Println("Sending reminder about invoices")
		return invoiceReminder(ctx)
	}

	return nil
}
