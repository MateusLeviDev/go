package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

type MetaData struct {
	Source  string
	Type    string
	Subject string
}

func NewEvent(meta MetaData, data any) (cloudevents.Event, error) {
		event := cloudevents.NewEvent()

	event.SetSpecVersion(cloudevents.VersionV1)
	event.SetID(uuid.NewString())
	event.SetSource(meta.Source)
	event.SetType(meta.Type)
	event.SetSubject(meta.Subject)
	event.SetTime(time.Now())
	event.SetDataContentType(cloudevents.ApplicationJSON)

	if err := event.SetData(cloudevents.ApplicationJSON, data); err != nil {
		return cloudevents.Event{}, err
	}

	return event, nil
}
