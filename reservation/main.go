package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	protoHall "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_hall/proto"
	protoShowing "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_showing/proto"
	protoReservation "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/reservation/proto"
)

type serviceHandler struct {
	mutex             sync.Mutex
	db                *dataBase
	cinemaHallService protoHall.CinemaHallService
	showingService    protoShowing.CinemaShowingService
}

func newReservationHandler(
	showingService protoShowing.CinemaShowingService,
	hallService protoHall.CinemaHallService,
	db *dataBase,
) *serviceHandler {
	handler := &serviceHandler{}
	handler.db = db
	handler.showingService = showingService
	handler.cinemaHallService = hallService
	return handler
}

func (handler *serviceHandler) Create(
	ctx context.Context,
	req *protoReservation.CreateReservationRequest,
	res *protoReservation.CreateReservationResponse,
) error {
	// synchronize this method to guarantee first come first serve
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	// validate seats
	showingRes, err := handler.showingService.Find(
		context.TODO(),
		&protoShowing.FindCinemaShowingRequest{Id: req.Showing},
	)
	if err != nil {
		return err
	}
	hallRes, err := handler.cinemaHallService.Find(
		context.TODO(),
		&protoHall.FindCinemaHallRequest{Id: showingRes.Showing.CinemaHall},
	)
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
				return fmt.Errorf("seat (r%d, s%d) is already reserved", s.Row, s.Col)
			}
		}
		seats = append(seats, seat{col: s.Col, row: s.Row})
	}

	// seats are valid, create reservation
	rsv := &reservation{
		showing:   req.Showing,
		user:      req.User,
		seats:     seats,
		confirmed: false,
	}
	handler.db.create(rsv)
	// map seats to proto seats
	pSeats := make([]*protoReservation.Seat, 0)
	for _, s := range rsv.seats {
		pSeats = append(pSeats, &protoReservation.Seat{Row: s.row, Col: s.col})
	}
	// map reservation to proto reservation
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

func (handler *serviceHandler) Confirm(
	cxt context.Context,
	req *protoReservation.ConfirmReservationRequest,
	res *protoReservation.ConfirmReservationResponse,
) error {
	// compare tokens to confirm reservation
	rsv := handler.db.confirm(req.Id, req.Token)
	// map seats to proto seats
	pSeats := make([]*protoReservation.Seat, 0)
	for _, s := range rsv.seats {
		pSeats = append(pSeats, &protoReservation.Seat{Row: s.row, Col: s.col})
	}
	// map reservation to proto reservation
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

func (handler *serviceHandler) Delete(
	ctx context.Context,
	req *protoReservation.DeleteReservationRequest,
	res *protoReservation.DeleteReservationResponse,
) error {
	rsv, err := handler.db.remove(req.Id)
	if err != nil {
		return err
	}
	// map seats to proto seats
	pSeats := make([]*protoReservation.Seat, 0)
	for _, s := range rsv.seats {
		pSeats = append(pSeats, &protoReservation.Seat{Row: s.row, Col: s.col})
	}

	// map reservation to proto reservation
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

func (handler *serviceHandler) FindAll(
	ctx context.Context,
	req *protoReservation.FindAllReservationsRequest,
	res *protoReservation.FindAllReservationsResponse,
) error {
	// map reservations to proto reservations
	reservations := handler.db.findAll()
	pReservations := make([]*protoReservation.Reservation, 0)
	for _, rsv := range reservations {
		// map seats to proto seats
		pSeats := make([]*protoReservation.Seat, 0)
		for _, s := range rsv.seats {
			pSeats = append(pSeats, &protoReservation.Seat{Row: s.row, Col: s.col})
		}
		// map reservation to proto reservation
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

func main() {
	service := micro.NewService(micro.Name("cinema.reservation.service"))
	service.Init()

	// init dependencies
	cinemaShowingService := protoShowing.NewCinemaShowingService("cinema.cinema_showing.service", service.Client())
	cinemaHallService := protoHall.NewCinemaHallService("cinema.cinema_hall.service", service.Client())
	db := newDataBase()

	// init handlers
	handler := newReservationHandler(cinemaShowingService, cinemaHallService, db)
	deletedShowingHandler := newCinemaShowingDeletedHandler(db)

	// register handlers
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
