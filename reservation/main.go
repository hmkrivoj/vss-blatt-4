package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	protoCinemaHall "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_hall/proto"
	protoCinemaShowing "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_showing/proto"
	protoReservation "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/reservation/proto"
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

func (db *dataBase) removeAllWhereShowingId(showing int64) []*reservation {
	reservations := db.findAll()
	toBeRemoved := make([]int64, 0)
	for _, reservation := range reservations {
		if reservation.showing == showing {
			toBeRemoved = append(toBeRemoved, reservation.id)
		}
	}
	removed := make([]*reservation, 0)
	for _, id := range toBeRemoved {
		showing, _ := db.remove(id)
		removed = append(removed, showing)
	}
	return removed
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

type cinemaShowingDeletedHandler struct {
	db *dataBase
}

func NewCinemaShowingDeletedHandler(db *dataBase) *cinemaShowingDeletedHandler {
	return &cinemaShowingDeletedHandler{db: db}
}

type serviceHandler struct {
	mutex                sync.Mutex
	db                   *dataBase
	cinemaHallService    protoCinemaHall.CinemaHallService
	cinemaShowingService protoCinemaShowing.CinemaShowingService
}

func NewReservationHandler(cinemaShowingService protoCinemaShowing.CinemaShowingService, cinemaHallService protoCinemaHall.CinemaHallService, db *dataBase) *serviceHandler {
	handler := &serviceHandler{}
	handler.db = db
	handler.cinemaShowingService = cinemaShowingService
	handler.cinemaHallService = cinemaHallService
	return handler
}

func (handler *serviceHandler) Confirm(cxt context.Context, req *protoReservation.ConfirmReservationRequest, res *protoReservation.ConfirmReservationResponse) error {
	rsv := handler.db.confirm(req.Id, req.Token)
	pSeats := make([]*protoReservation.Seat, 0)
	for _, s := range rsv.seats {
		pSeats = append(pSeats, &protoReservation.Seat{Row: s.row, Col: s.col})
	}

	res.Reservation = &protoReservation.Reservation{
		Id:        rsv.id,
		Token:     rsv.token,
		Showing:   rsv.showing,
		User:      rsv.user,
		Seats:     pSeats,
		Confirmed: rsv.confirmed,
	}

	return nil
}

func (handler *serviceHandler) Create(ctx context.Context, req *protoReservation.CreateReservationRequest, res *protoReservation.CreateReservationResponse) error {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()
	showingRes, err := handler.cinemaShowingService.Find(context.TODO(), &protoCinemaShowing.FindCinemaShowingRequest{Id: req.Showing})
	if err != nil {
		return err
	}
	hallRes, err := handler.cinemaHallService.Find(context.TODO(), &protoCinemaHall.FindCinemaHallRequest{Id: showingRes.Showing.CinemaHall})
	if err != nil {
		return err
	}
	reservedSeats := handler.db.findAllReservedSeats()
	cols := hallRes.Hall.Cols
	rows := hallRes.Hall.Rows
	seats := make([]seat, 0)
	for _, s := range req.Seats {
		if s.Col >= cols || s.Row >= rows {
			return errors.New("invalid seat")
		}
		for _, reservedSeat := range reservedSeats {
			if s.Col == reservedSeat.col || s.Row == reservedSeat.row {
				return errors.New(fmt.Sprintf("seat (r%d, s%d) is already reserved", s.Row, s.Col))
			}
		}
		seats = append(seats, seat{col: s.Col, row: s.Row})
	}
	rsv := &reservation{
		showing:   req.Showing,
		user:      req.User,
		seats:     seats,
		confirmed: false,
	}
	handler.db.create(rsv)
	pSeats := make([]*protoReservation.Seat, 0)
	for _, s := range rsv.seats {
		pSeats = append(pSeats, &protoReservation.Seat{Row: s.row, Col: s.col})
	}

	res.Reservation = &protoReservation.Reservation{
		Id:        rsv.id,
		Token:     rsv.token,
		Showing:   rsv.showing,
		User:      rsv.user,
		Seats:     pSeats,
		Confirmed: rsv.confirmed,
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
		Id:        rsv.id,
		Token:     rsv.token,
		Showing:   rsv.showing,
		User:      rsv.user,
		Seats:     pSeats,
		Confirmed: rsv.confirmed,
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
			Id:        rsv.id,
			Token:     rsv.token,
			Showing:   rsv.showing,
			User:      rsv.user,
			Seats:     pSeats,
			Confirmed: rsv.confirmed,
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
	service := micro.NewService(micro.Name("cinema.reservation.service"))
	service.Init()

	cinemaShowingService := protoCinemaShowing.NewCinemaShowingService("cinema.cinema_showing.service", service.Client())
	cinemaHallService := protoCinemaHall.NewCinemaHallService("cinema.cinema_hall.service", service.Client())
	db := newDataBase()

	handler := NewReservationHandler(cinemaShowingService, cinemaHallService, db)
	deletedShowingHandler := NewCinemaShowingDeletedHandler(db)

	err := micro.RegisterSubscriber("cinema.cinema_showing.deleted", service.Server(), deletedShowingHandler)
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
