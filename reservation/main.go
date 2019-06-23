package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	protoCinemaShowing "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_showing/proto"
	protoReservation "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/reservation/proto"
	"sync"
)

type seat struct {
	col int32
	row int32
}

type reservation struct {
	id      int64
	token   string
	showing int64
	user    int64

	seats []seat
}

type dataBase struct {
	mutex        sync.Mutex
	idCounter    int64
	reservations map[int64]reservation
}

func newDataBase() *dataBase {
	db := &dataBase{}
	db.idCounter = 1
	db.reservations = make(map[int64]reservation)

	return db
}

func (db *dataBase) create(rsv reservation) reservation {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	id := db.idCounter
	db.idCounter++

	token := make([]byte, 4)
	_, _ = rand.Read(token)
	tokenString := fmt.Sprintf("%x", token)

	created := reservation{
		id:      id,
		token:   tokenString,
		showing: rsv.showing,
		user:    rsv.user,
		seats:   rsv.seats,
	}
	db.reservations[id] = created

	return created
}

func (db *dataBase) remove(id int64) (reservation, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	if _, ok := db.reservations[id]; !ok {
		return reservation{}, errors.New("no such id")
	}
	showing := db.reservations[id]
	delete(db.reservations, id)
	return showing, nil
}

func (db *dataBase) removeAllWhereShowingId(showing int64) []reservation {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	reservations := db.findAll()
	toBeRemoved := make([]int64, 0)
	for _, reservation := range reservations {
		if reservation.showing == showing {
			toBeRemoved = append(toBeRemoved, reservation.id)
		}
	}
	removed := make([]reservation, 0)
	for _, id := range toBeRemoved {
		showing, _ := db.remove(id)
		removed = append(removed, showing)
	}
	return removed
}

func (db *dataBase) findAll() []reservation {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	reservations := make([]reservation, 0)
	for _, rsv := range db.reservations {
		reservations = append(reservations, rsv)
	}
	return reservations
}

type cinemaShowingDeletedHandler struct {
	db *dataBase
}

func NewCinemaHallDeletedHandler(db *dataBase) *cinemaShowingDeletedHandler {
	return &cinemaShowingDeletedHandler{db: db}
}

type serviceHandler struct {
	db *dataBase
}

func NewReservationHandler(db *dataBase) *serviceHandler {
	handler := &serviceHandler{}
	handler.db = db
	return handler
}

func (handler *serviceHandler) Confirm(cxt context.Context, req *protoReservation.ConfirmReservationRequest, res *protoReservation.ConfirmReservationResponse) error {
	panic("implement me")
}

func (handler *serviceHandler) Create(ctx context.Context, req *protoReservation.CreateReservationRequest, res *protoReservation.CreateReservationResponse) error {
	seats := make([]seat, 0)
	for _, s := range req.Seats {
		seats = append(seats, seat{col: s.Col, row: s.Row})
	}
	rsv := reservation{
		showing: req.Showing,
		user:    req.User,
		seats:   seats,
	}
	rsv = handler.db.create(rsv)
	pSeats := make([]*protoReservation.Seat, 0)
	for _, s := range rsv.seats {
		pSeats = append(pSeats, &protoReservation.Seat{Row: s.row, Col: s.col})
	}

	res.Reservation = &protoReservation.Reservation{
		Id:      rsv.id,
		Token:   rsv.token,
		Showing: rsv.showing,
		User:    rsv.user,
		Seats:   pSeats,
	}

	return nil
}

func (handler *serviceHandler) Delete(ctx context.Context, req *protoReservation.DeleteReservationRequest, res *protoReservation.DeleteReservationResponse) error {
	rsv, err := handler.db.remove(req.Id)
	if err != nil {
		return err
	}
	pSeats := make([]*protoReservation.Seat, 0)
	for _, s := range rsv.seats {
		pSeats = append(pSeats, &protoReservation.Seat{Row: s.row, Col: s.col})
	}

	res.Reservation = &protoReservation.Reservation{
		Id:      rsv.id,
		Token:   rsv.token,
		Showing: rsv.showing,
		User:    rsv.user,
		Seats:   pSeats,
	}

	return nil
}

func (handler *serviceHandler) FindAll(ctx context.Context, req *protoReservation.FindAllReservationsRequest, res *protoReservation.FindAllReservationsResponse) error {
	reservations := handler.db.findAll()
	pReservations := make([]*protoReservation.Reservation, 0)
	for _, rsv := range reservations {
		pSeats := make([]*protoReservation.Seat, 0)
		for _, s := range rsv.seats {
			pSeats = append(pSeats, &protoReservation.Seat{Row: s.row, Col: s.col})
		}
		pReservations = append(pReservations, &protoReservation.Reservation{
			Id:      rsv.id,
			Token:   rsv.token,
			Showing: rsv.showing,
			User:    rsv.user,
			Seats:   pSeats,
		})
	}
	res.Reservations = pReservations
	return nil
}

func (handler *cinemaShowingDeletedHandler) MovieDeleted(cxt context.Context, event *protoCinemaShowing.DeleteCinemaShowingResponse) error {
	handler.db.removeAllWhereShowingId(event.Showing.Id)
	return nil
}

func main() {
	service := micro.NewService(micro.Name("cinema.cinema_showing.service"))
	service.Init()

	db := newDataBase()
	handler := NewReservationHandler(db)
	deletedHallHandler := NewCinemaHallDeletedHandler(db)

	err := micro.RegisterSubscriber("cinema.cinema_hall.deleted", service.Server(), deletedHallHandler)
	if err != nil {
		log.Fatal(err)
	}

	err = protoReservation.RegisterReservationServiceHandler(service.Server(), handler)
	if err != nil {
		log.Fatal(err)
	}

	err = service.Run()
	if err != nil {
		log.Fatal(err)
	}
}
