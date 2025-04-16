package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"time"

	"github.com/kiinoda/lifeos/internal/config"
	"github.com/kiinoda/lifeos/internal/email"
	"github.com/kiinoda/lifeos/internal/sheets"
)

const (
	defaultRegion = "eu-west-1"
	ssmPath       = "/personal/daily_schedule/config"
	sheetName     = "Weekly"
)

func handler(ctx context.Context) error {
	action, _ := os.LookupEnv("LIFEOS_ACTION")
	if !slices.Contains([]string{"daily_schedule", "upcoming_event", "invoice_reminder"}, action) {
		return fmt.Errorf("Please set LIFEOS_ACTION environment variable")
	}

	cfg, err := config.NewConfig(defaultRegion, ssmPath)
	if err != nil {
		return fmt.Errorf("Failed to create application config: %w", err)
	}

	ctx = config.ContextWithConfig(ctx, cfg)

	events, err := sheets.GetEvents(ctx, sheetName)
	if err != nil {
		return fmt.Errorf("Failed to extract events: %w", err)
	}

	// Weekday starts at 0, our table starts at 1
	dayOfWeek := time.Now().Weekday() - 1

	if action == "daily_schedule" {
		log.Println("Daily schedule")
		textBody, htmlBody, err := email.CreateDailyMessageBody(dayOfWeek, events)
		if err != nil {
			if errors.Is(err, email.ErrNoEvents) {
				log.Println(err)
				return nil
			} else {
				return err
			}
		}
		err = email.SendEmail(ctx, "LifeOS", "Daily Schedule", textBody, htmlBody)
		if err != nil {
			return fmt.Errorf("Failed to send email: %w", err)
		}
	}

	if action == "upcoming_event" {
		log.Println("Upcoming event")
		textBody, htmlBody, err := email.CreateReminderMessageBody(dayOfWeek, events)
		if err != nil {
			if errors.Is(err, email.ErrNoReminder) {
				log.Println(err)
				return nil
			} else {
				return err
			}
		}
		err = email.SendEmail(ctx, "LifeOS Reminder Bot", "Upcoming Event", textBody, htmlBody)
		if err != nil {
			return fmt.Errorf("Failed to send email: %w", err)
		}
	}

	if action == "invoice_reminder" {
		log.Println("Sending reminder about invoices")
		textBody, htmlBody, err := email.CreateInvoiceReminderMessageBody()
		if err != nil {
			return err
		}
		err = email.SendEmail(ctx, "LifeOS", "Invoice Reminder", textBody, htmlBody)
		if err != nil {
			return fmt.Errorf("Failed to send email: %w", err)
		}
	}

	return nil
}
