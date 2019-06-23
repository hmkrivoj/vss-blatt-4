package main

import (
	"context"
	"errors"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/user/proto"
	"sync"
)

type UserHandler struct {
	mutex     sync.Mutex
	idCounter int64
	users     map[int64]user
}

func NewUserHandler() *UserHandler {
	handler := &UserHandler{}
	handler.idCounter = 1
	handler.users = make(map[int64]user)
	return handler
}

type user struct {
	id   int64
	name string
}

func (handler *UserHandler) Create(ctx context.Context, req *proto.CreateUserRequest, res *proto.CreateUserResponse) error {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	id := handler.idCounter
	handler.idCounter++

	user := user{
		id:   id,
		name: req.Name,
	}
	handler.users[id] = user

	res.User = &proto.User{
		Id:   user.id,
		Name: user.name,
	}

	return nil
}

func (handler *UserHandler) Delete(ctx context.Context, req *proto.DeleteUserRequest, res *proto.DeleteUserResponse) error {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	// TODO make sure only user with no reservations can be deleted
	if _, ok := handler.users[req.Id]; !ok {
		return errors.New("no such id")
	}

	user := handler.users[req.Id]
	delete(handler.users, req.Id)

	res.User = &proto.User{
		Id:   user.id,
		Name: user.name,
	}

	return nil
}

func (handler *UserHandler) FindAll(ctx context.Context, req *proto.FindAllUsersRequest, res *proto.FindAllUsersResponse) error {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	users := make([]*proto.User, 0)
	for _, user := range handler.users {
		users = append(users, &proto.User{
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

	handler := NewUserHandler()

	err := proto.RegisterUserServiceHandler(service.Server(), handler)
	if err != nil {
		log.Fatal(err)
	}

	err = service.Run()
	if err != nil {
		log.Fatal(err)
	}
}
