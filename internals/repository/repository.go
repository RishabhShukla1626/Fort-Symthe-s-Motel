package repository

import (
	"time"

	"github.com/TheDevCarnage/FortSmythesMotel/internals/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservations) (int, error)
	InsertRoomRestriction(r models.RoomRestrictions) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time)([]models.Rooms, error)
	GetRoomByID(id int) (models.Rooms, error)
}