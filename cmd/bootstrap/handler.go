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

const (
	defaultRegion = "eu-west-1"
	ssmPath       = "/personal/daily_schedule/config"
	weeklySheet   = "Weekly"
	scheduleSheet = "Future"
)

func handler(ctx context.Context) error {
	action, _ := os.LookupEnv("LIFEOS_ACTION")
	if !slices.Contains([]string{"daily_schedule", "event_schedule", "event_notification", "invoice_reminder"}, action) {
		return fmt.Errorf("please set LIFEOS_ACTION environment variable")
	}

	cfg, err := config.NewConfig(defaultRegion, ssmPath)
	if err != nil {
		return fmt.Errorf("failed to create application config: %w", err)
	}

	ctx = config.ContextWithConfig(ctx, cfg)

	// Weekday starts at 0, our table starts at 1
	dayOfWeek := time.Now().Weekday() - 1

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
