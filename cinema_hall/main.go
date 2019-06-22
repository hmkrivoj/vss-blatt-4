package main

import (
	"context"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_hall/proto"
)

type CinemaHallHandler struct {
}

func (*CinemaHallHandler) Create(context.Context, *proto.CreateRequest, *proto.CreateResponse) error {
	panic("implement me")
}

func (*CinemaHallHandler) Delete(context.Context, *proto.DeleteRequest, *proto.DeleteResponse) error {
	panic("implement me")
}

func (*CinemaHallHandler) FindAll(context.Context, *proto.FindAllRequest, *proto.FindAllResponse) error {
	panic("implement me")
}

func main() {
	service := micro.NewService(micro.Name("cinema_hall"))
	service.Init()

	err := proto.RegisterCinemaHallServiceHandler(service.Server(), new(CinemaHallHandler))
	if err != nil {
		log.Fatal(err)
	}

	err = service.Run()
	if err != nil {
		log.Fatal(err)
	}
}
