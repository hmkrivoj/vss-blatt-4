package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	proto "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_hall/proto"
)

func callBack(ctx context.Context, event *proto.DeleteResponse) error {
	fmt.Printf("Delete event: %v", event)
	return nil
}

func main() {
	service := micro.NewService(micro.Name("cinema.sub"))
	service.Init()

	_ = micro.RegisterSubscriber("cinema.cinema_hall.deleted", service.Server(), callBack, server.SubscriberQueue("cinema.sub.queue"))

	_ = service.Run()
}
