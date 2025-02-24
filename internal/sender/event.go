package sender

import "ozon-omp-api/internal/model"

type EventSender interface {
	Send(model *model.CourseEvent) error
}
