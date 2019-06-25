package main

import (
	"context"
	"errors"
	"sync"

	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	protoReservation "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/reservation/proto"
	protoUser "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/user/proto"
)

type UserHandler struct {
	mutex              sync.Mutex
	idCounter          int64
	users              map[int64]user
	reservationService protoReservation.ReservationService
}

func NewUserHandler(reservationService protoReservation.ReservationService) *UserHandler {
	handler := &UserHandler{}
	handler.idCounter = 1
	handler.users = make(map[int64]user)
	handler.reservationService = reservationService
	return handler
}

type user struct {
	id   int64
	name string
}

func (handler *UserHandler) Create(
	ctx context.Context,
	req *protoUser.CreateUserRequest,
	res *protoUser.CreateUserResponse,
) error {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	id := handler.idCounter
	handler.idCounter++

	user := user{
		id:   id,
		name: req.Name,
	}
	handler.users[id] = user

	res.User = &protoUser.User{
		Id:   user.id,
		Name: user.name,
	}

	return nil
}

func (handler *UserHandler) Delete(
	ctx context.Context,
	req *protoUser.DeleteUserRequest,
	res *protoUser.DeleteUserResponse,
) error {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()
	// check if user exists
	if _, ok := handler.users[req.Id]; !ok {
		return errors.New("no such id")
	}
	// check if user has reservations
	reservations, err := handler.reservationService.FindAll(
		context.TODO(),
		&protoReservation.FindAllReservationsRequest{},
	)
	if err != nil {
		return err
	}
	for _, reservation := range reservations.Reservations {
		if reservation.User == req.Id {
			return errors.New("this user still has reservations")
		}
	}
	// delete user and track user
	user := handler.users[req.Id]
	delete(handler.users, req.Id)
	res.User = &protoUser.User{
		Id:   user.id,
		Name: user.name,
	}

	return nil
}

func (handler *UserHandler) FindAll(
	ctx context.Context,
	req *protoUser.FindAllUsersRequest,
	res *protoUser.FindAllUsersResponse,
) error {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	users := make([]*protoUser.User, 0)
	for _, user := range handler.users {
		users = append(users, &protoUser.User{
			Id:   user.id,
			Name: user.name,
		})
	}
	res.Users = users

	return nil
}

func main() {
	service := micro.NewService(micro.Name("cinema.user.service"))
	service.Init()

	reservationService := protoReservation.NewReservationService("cinema.reservation.service", service.Client())

	handler := NewUserHandler(reservationService)

	err := protoUser.RegisterUserServiceHandler(service.Server(), handler)
	if err != nil {
		log.Fatal(err)
	}

	err = service.Run()
	if err != nil {
		log.Fatal(err)
	}
}
