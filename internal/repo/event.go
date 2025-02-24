package repo

import "ozon-omp-api/internal/model"

type EventRepo interface {
	Lock(n uint64) ([]model.CourseEvent, error)
	Unlock(eventIDs []uint64) error

	Add(event []model.CourseEvent) error
	Remove(eventIDs []uint64) error
}
