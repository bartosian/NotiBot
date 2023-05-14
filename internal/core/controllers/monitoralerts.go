package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"dstwilio/internal/core/domain"
)

func (c *NotifierController) MonitorAlerts() error {
	return nil
}

func newAlertsHandler(isTriggered bool) bool {
	resp, err := http.Get(os.Getenv("ALERT_MANAGER_URL"))
	if err != nil {
		log.Fatal("Error fetching alerts: ", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body: ", err)
	}

	var (
		voiceBody   string
		messageBody string
		alertList   []domain.Alert
	)

	err = json.Unmarshal(data, &alertList)
	if err != nil {
		log.Fatal("error parsing alert list", err)
	}

	if len(alertList) == 0 {
		if isTriggered {
			voiceBody = alertMessageRecoveredTemplate
			messageBody = alertMessageRecoveredTemplate

			isTriggered = false
		} else {
			log.Println("no active alerts found")

			return false
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

	getDelimiter()
	getCurrentTime()
	makePhoneCall(voiceBody)
	sendTextMessage(messageBody)
	getDelimiter()

	return isTriggered
}
