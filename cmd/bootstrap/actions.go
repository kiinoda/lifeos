package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kiinoda/lifeos/internal/email"
	"github.com/kiinoda/lifeos/internal/sheets"
)

func invoiceReminder(ctx context.Context) error {
	textBody, htmlBody, err := email.CreateInvoiceReminderMessageBody()
	if err != nil {
		return err
	}
	err = email.SendEmail(ctx, "LifeOS", "Invoice Reminder", textBody, htmlBody)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func eventNotification(ctx context.Context, dayOfWeek time.Weekday) error {
	events, err := sheets.GetEvents(ctx, weeklySheet)
	if err != nil {
		return fmt.Errorf("failed to extract events: %w", err)
	}
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
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func dailySchedule(ctx context.Context, dayOfWeek time.Weekday) error {
	events, err := sheets.GetEvents(ctx, weeklySheet)
	if err != nil {
		return fmt.Errorf("failed to extract events: %w", err)
	}
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
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func eventSchedule(ctx context.Context, dayOfWeek time.Weekday) error {
	events, err := sheets.GetEventSchedule(ctx, scheduleSheet)
	if err != nil {
		return fmt.Errorf("failed to extract future events: %w", err)
	}
	textBody, htmlBody, err := email.CreateEventScheduleMessageBody(dayOfWeek, events)
	if err != nil {
		if errors.Is(err, email.ErrNoScheduledEvents) {
			return nil
		} else {
			return err
		}
	}
	err = email.SendEmail(ctx, "LifeOS Event Schedule Bot", "Event Schedule", textBody, htmlBody)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
