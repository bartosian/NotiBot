package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"dstwilio/internal/core/domain"
)

const (
	alertVoiceTriggeredTemplate   = "NEW ALERT TRIGGERED!\n\n*ALERT NAME:* %s\n*SUMMARY:* %s\n*START TIME:* %s\n*END TIME:* %s\n*STATUS:* %s\n\n"
	alertMessageTriggeredTemplate = "🚨 NEW ALERT TRIGGERED!\n\nALERT NAME: %s\n\nSUMMARY: %s\n\nSTART TIME: %s\n\nEND TIME: %s\n\nSTATUS: %s\n\n"
	alertMessageRecoveredTemplate = "✅ ALL ALERTS HAVE BEEN RESOLVED AND THE SYSTEM IS NOW OPERATING NORMALLY."

	alertIntervalCheck          = 1 * time.Minute
	alertTriggeredIntervalCheck = 5 * time.Minute
)

// MonitorAlerts continuously checks for new alerts at an interval. If an alert is triggered, the interval is increased.
func (c *NotifierController) MonitorAlerts() error {
	var (
		alertTriggered bool
		err            error
	)

	intervalCheck := alertIntervalCheck

	for {
		alertTriggered, err = c.newAlertsHandler(alertTriggered)
		if err != nil {
			return err
		}

		if alertTriggered {
			intervalCheck = alertTriggeredIntervalCheck
		} else {
			intervalCheck = alertIntervalCheck
		}

		time.Sleep(intervalCheck)
	}
}

// newAlertsHandler fetches alerts from the ALERT_MANAGER_URL. If a new alert is found, it is handled accordingly and
// notifications are sent via voice and message. If no new alerts are found and a previous alert was triggered, a
// recovery message is sent. Returns a boolean indicating whether an alert was triggered.
func (c *NotifierController) newAlertsHandler(isTriggered bool) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, os.Getenv("ALERT_MANAGER_URL"), nil)
	if err != nil {
		return false, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Error("error fetching alerts", err, nil)

		return false, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Error("error closing response body", err, nil)

			os.Exit(1)
		}
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("error reading response body", err, nil)

		return false, err
	}

	var (
		voiceBody   string
		messageBody string
		alertList   []domain.Alert
	)

	err = json.Unmarshal(data, &alertList)
	if err != nil {
		c.logger.Error("error parsing alert list", err, nil)

		return false, err
	}

	if len(alertList) == 0 {
		if isTriggered {
			voiceBody = alertMessageRecoveredTemplate
			messageBody = alertMessageRecoveredTemplate

			isTriggered = false
		} else {
			c.logger.Info("no active alerts found", nil)

			return false, nil
		}
	} else {
		isTriggered = true

		for _, alert := range alertList {
			voiceBody += fmt.Sprintf(alertVoiceTriggeredTemplate,
				alert.Labels.AlertName,
				alert.Annotations.Summary,
				alert.StartsAt,
				alert.EndsAt,
				alert.Status.State,
			)

			messageBody += fmt.Sprintf(alertMessageTriggeredTemplate,
				alert.Labels.AlertName,
				alert.Annotations.Summary,
				alert.StartsAt,
				alert.EndsAt,
				alert.Status.State,
			)
		}
	}

	err = c.notifierGateway.CreateCall(voiceBody)
	if err != nil {
		c.logger.Error("error creating call", err, nil)

		return false, err
	}

	err = c.notifierGateway.SendMessage(messageBody)
	if err != nil {
		c.logger.Error("error sending message", err, nil)

		return false, err
	}

	return isTriggered, nil
}
