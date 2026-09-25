package notification

import (
	"fmt"
	"strings"
	"time"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/domain/notification"
)

func formatMessage(notice *notification.Notification) (string, error) {
	switch notice.Type {
	case "competitor_opened":
		return competitorMessage(notice), nil
	case "event_upcoming":
		return eventMessage(notice)
	default:
		return "", fmt.Errorf(
			"unsupported notification type %q",
			notice.Type,
		)
	}
}

func competitorMessage(notice *notification.Notification) string {
	var message strings.Builder

	fmt.Fprintf(
		&message,
		"Рядом с точкой «%s» открылся новый конкурент\n\n%s",
		notice.TrackedLocation.Name,
		notice.Subject.Title,
	)

	appendLine(&message, "Адрес", notice.Subject.Address)
	appendLine(&message, "Категория", notice.Subject.Category)
	appendLine(
		&message,
		"Дата открытия",
		formatDate(
			notice.Subject.OpenedAt,
			time.DateOnly,
			"02.01.2006",
		),
	)
	appendRoute(&message, notice.Route)
	appendList(&message, "Почему это важно", notice.Assessment.Reasons)
	appendList(
		&message,
		"Рекомендации",
		notice.Assessment.Recommendations,
	)

	return message.String()
}

func eventMessage(notice *notification.Notification) (string, error) {
	var message strings.Builder

	fmt.Fprintf(
		&message,
		"Рядом с точкой «%s» скоро состоится мероприятие\n\n%s",
		notice.TrackedLocation.Name,
		notice.Subject.Title,
	)

	date, err := formatEventDate(notice.Subject.Date)
	if err != nil {
		return "", err
	}

	appendLine(
		&message,
		"Дата",
		date,
	)
	appendLine(&message, "Адрес", notice.Subject.Address)
	appendRoute(&message, notice.Route)
	appendList(&message, "Почему это важно", notice.Assessment.Reasons)
	appendList(
		&message,
		"Рекомендации",
		notice.Assessment.Recommendations,
	)

	return message.String(), nil
}

func appendLine(message *strings.Builder, label, value string) {
	if value != "" {
		fmt.Fprintf(message, "\n%s: %s", label, value)
	}
}

func appendRoute(
	message *strings.Builder,
	route notification.Route,
) {
	fmt.Fprintf(
		message,
		"\nРасстояние: %.0f м, примерно %.0f мин. пешком",
		route.DistanceMeters,
		route.DurationSeconds/60,
	)
}

func appendList(
	message *strings.Builder,
	title string,
	items []string,
) {
	if len(items) == 0 {
		return
	}

	fmt.Fprintf(
		message,
		"\n\n%s:\n• %s",
		title,
		strings.Join(items, "\n• "),
	)
}

func formatDate(
	value string,
	layout string,
	outputLayout string,
) string {
	date, err := time.Parse(layout, value)
	if err != nil {
		return value
	}

	return date.Format(outputLayout)
}

func formatEventDate(value string) (string, error) {
	date, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return "", fmt.Errorf("parse event date: %w", err)
	}

	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return "", fmt.Errorf("load Moscow timezone: %w", err)
	}

	return date.
		In(location).
		Format("15:04 МСК - 02.01.2006"), nil
}
