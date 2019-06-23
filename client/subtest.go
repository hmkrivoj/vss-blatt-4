package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	protoMovie "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/movie/proto"
)

func callBackMovie(ctx context.Context, event *protoMovie.DeleteMovieResponse) error {
	fmt.Printf("Delete event: %v", event)
	return nil
}

func main() {
	service := micro.NewService(micro.Name("test.movie.sub"))
	service.Init()

	_ = micro.RegisterSubscriber("cinema.movie.deleted", service.Server(), callBackMovie, server.SubscriberQueue("test.movie.sub.queue"))

	_ = service.Run()
}
