package models

import (
	"time"
)

// Reservation: holds reservation data for Reservations table
type Reservations struct {
	FirstName 	string
	LastName  	string
	Email     	string
	Phone    	string
	StartDate 	time.Time
	EndDate		time.Time
	RoomID		int
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
	Room 		Rooms
}

// Users: holds Users data for Users table
type Users struct {
	ID 			int
	FirstName 	string
	LastName 	string
	Email 		string
	Password 	string
	AccessLevel int
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}


// Rooms: holds Rooms data for Rooms table
type Rooms struct{
	ID 			int
	RoomName 	string
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}


// Restrictions: holds Restrictions data for Restrictions table
type Restrictions struct{
	ID 					int
	RestrictionName 	string
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
}


// RoomRestrictions: holds RoomRestrictions data for RoomRestrictions table
type RoomRestrictions struct{
	ID 				int
	StartDate 		time.Time
	EndDate			time.Time
	RoomID			int
	ReservationID 	int
	RestrictionID 	int
	CreatedAt 		time.Time
	UpdatedAt 		time.Time
	Room 			Rooms
	Reservation 	Reservations
	Restriction 	Restrictions
}


//MailData: Holds email messages
type MailData struct {
	To 		 string
	From 	 string
	Subject  string
	Content  string
	Template string
}