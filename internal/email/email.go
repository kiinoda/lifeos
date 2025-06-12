package email

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/kiinoda/lifeos/internal/config"
	"github.com/kiinoda/lifeos/internal/events"
)

var (
	ErrNoReminder        = errors.New("No reminder found")
	ErrNoEvents          = errors.New("No event found")
	ErrNoScheduledEvents = errors.New("No scheduled event found")
)

func CreateEventScheduleMessageBody(dayOfWeek time.Weekday, eventList []events.ScheduledEvent) (string, string, error) {
	pastEvents := ""
	upcomingEvents := ""
	text := ""
	slices.SortFunc(eventList, func(a, b events.ScheduledEvent) int {
		return int(a.Time.Sub(b.Time))
	})

	for _, event := range eventList {
		if event.Alertable {
			if event.Time.Before(time.Now().Add(-24 * time.Hour)) {
				pastEvents = pastEvents + fmt.Sprintf("%8s * %s\n", event.Time.Format("20060102"), event.Desc)
			}
			if event.Time.After(time.Now().Add(-24*time.Hour)) && event.Time.Before(time.Now().Add(7*24*time.Hour)) {
				upcomingEvents = upcomingEvents + fmt.Sprintf("%8s * %s\n", event.Time.Format("20060102"), event.Desc)
			}
		}
	}
	if pastEvents == "" && upcomingEvents == "" {
		return "", "", ErrNoScheduledEvents
	}

	if upcomingEvents != "" && pastEvents != "" {
		text = "Upcoming Events\n" + upcomingEvents + "\nPast Events\n" + pastEvents
	} else {
		if upcomingEvents != "" {
			text = "Upcoming Events\n" + upcomingEvents
		}
		if pastEvents != "" {
			text = "Past Events\n" + pastEvents
		}
	}

	html := fmt.Sprintf("<html><pre>%s</pre></html>", text)
	return text, html, nil
}

func CreateReminderMessageBody(dayOfWeek time.Weekday, events []events.Event) (string, string, error) {
	found := false
	text := ""
	for _, event := range events {
		if event.Days[dayOfWeek] != "" {
			// TODO: extract TZ to AWS SSM Parameter Store
			alerting_location, err := time.LoadLocation("Europe/Bucharest")
			if err != nil {
				return "", "", fmt.Errorf("Cannot load timezone: %w", err)
			}
			_, alerting_offset := time.Now().In(alerting_location).Zone()
			_, server_offset := time.Now().Zone()
			difference := time.Until(event.Time.UTC()).Seconds() - float64(alerting_offset-server_offset)
			if -120 < difference && difference < 120 {
				text = text + fmt.Sprintf("%4s * %s\n", event.GetTimePlaceholder(), event.Desc)
				found = true
				break
			}
		}
	}
	if !found {
		return "", "", ErrNoReminder
	}

	html := fmt.Sprintf("<html><pre>%s</pre></html>", text)
	return text, html, nil
}

func CreateInvoiceReminderMessageBody() (string, string, error) {
	text := "The following companies do not email invoices. Please save them manually.\n\n"

	for _, c := range []string{"OpenAI", "Anthropic", "Amp"} {
		text = text + fmt.Sprintf("* %s\n", c)
	}

	html := fmt.Sprintf("<html><pre>%s</pre></html>", text)
	return text, html, nil
}

func CreateDailyMessageBody(dayOfWeek time.Weekday, events []events.Event) (string, string, error) {
	found := false
	text := ""
	for _, event := range events {
		if event.Days[dayOfWeek] != "" {
			text = text + fmt.Sprintf("%s %4s %s\n", event.Days[dayOfWeek], event.GetTimePlaceholder(), event.Desc)
			found = true
		}
	}
	if !found {
		return "", "", ErrNoEvents
	}

	html := fmt.Sprintf("<html><pre>%s</pre></html>", text)
	return text, html, nil
}

func SendEmail(ctx context.Context, fromLabel string, subject string, textBody string, htmlBody string) error {
	awsConfig, err := awscfg.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("Unable to load AWS default config")
	}
	client := ses.NewFromConfig(awsConfig)

	cfg, err := config.ConfigFromContext(ctx)
	if err != nil {
		return fmt.Errorf("Unable to get config from context: %w", err)
	}

	source := fmt.Sprintf("%s <%s>", fromLabel, cfg.Sender)
	recipient := cfg.Recipient

	emailInput := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{recipient},
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Data: &textBody,
				},
				Html: &types.Content{
					Data: &htmlBody,
				},
			},
			Subject: &types.Content{
				Data: aws.String(subject),
			},
		},
		Source: &source,
	}
	_, err = client.SendEmail(context.Background(), emailInput)
	if err != nil {
		return err
	}
	return nil
}
