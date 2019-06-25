package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"sync"
)

type seat struct {
	col int32
	row int32
}

type reservation struct {
	id        int64
	token     string
	showing   int64
	user      int64
	confirmed bool

	seats []seat
}

type dataBase struct {
	mutex        sync.Mutex
	idCounter    int64
	reservations map[int64]*reservation
}

func newDataBase() *dataBase {
	db := &dataBase{}
	db.idCounter = 1
	db.reservations = make(map[int64]*reservation)

	return db
}

func (db *dataBase) create(rsv *reservation) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	id := db.idCounter
	db.idCounter++

	token := make([]byte, 4)
	_, _ = rand.Read(token)
	tokenString := fmt.Sprintf("%x", token)

	rsv.id = id
	rsv.token = tokenString
	rsv.confirmed = false
	db.reservations[id] = rsv
}

func (db *dataBase) remove(id int64) (*reservation, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	if _, ok := db.reservations[id]; !ok {
		return nil, errors.New("no such id")
	}
	rsv := db.reservations[id]
	delete(db.reservations, id)
	return rsv, nil
}

func (db *dataBase) removeAllWhereShowingID(showing int64) ([]*reservation, error) {
	reservations := db.findAll()
	toBeRemoved := make([]int64, 0)
	for _, reservation := range reservations {
		if reservation.showing == showing {
			toBeRemoved = append(toBeRemoved, reservation.id)
		}
	}
	removed := make([]*reservation, 0)
	for _, id := range toBeRemoved {
		showing, err := db.remove(id)
		if err != nil {
			return removed, err
		}
		removed = append(removed, showing)
	}
	return removed, nil
}

func (db *dataBase) findAll() []*reservation {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	reservations := make([]*reservation, 0)
	for _, rsv := range db.reservations {
		reservations = append(reservations, rsv)
	}
	return reservations
}

func (db *dataBase) findAllReservedSeats() []seat {
	reservations := db.findAll()

	seats := make([]seat, 0)
	for _, rsv := range reservations {
		seats = append(seats, rsv.seats...)
	}
	return seats
}

func (db *dataBase) confirm(id int64, token string) *reservation {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if rsv, ok := db.reservations[id]; ok && rsv.token == token {
		rsv.confirmed = true
		return rsv
	}
	return nil
}
