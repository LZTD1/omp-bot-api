package model

type Course struct {
	ID          uint64
	Title       string
	Description string
}

type EventType uint8
type EventStatus uint8

const (
	Created EventType = iota
	Updated
	Removed
)
const (
	Deferrer EventStatus = iota
	Processed
)

type CourseEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Entity *Course
}
