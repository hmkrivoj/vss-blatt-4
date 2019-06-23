package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	protoMovie "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/movie/proto"
)

func main() {
	service := micro.NewService(micro.Name("cinema.client"))
	service.Init()

	movieService := protoMovie.NewMovieService("cinema.movie.service", service.Client())

	createResponseMovie, _ := movieService.Create(context.TODO(), &protoMovie.CreateRequest{Title: "Fight Club"})
	fmt.Printf("Created %v\n", createResponseMovie)
	createResponseMovie, _ = movieService.Create(context.TODO(), &protoMovie.CreateRequest{Title: "Se7en"})
	fmt.Printf("Created %v\n", createResponseMovie)
	findAllResponseMovie, _ := movieService.FindAll(context.TODO(), &protoMovie.FindAllRequest{})
	fmt.Printf("Found %v\n", findAllResponseMovie)
	deleteResponseMovie, _ := movieService.Delete(context.TODO(), &protoMovie.DeleteRequest{Id: 1})
	fmt.Printf("Deleted %v\n", deleteResponseMovie)
	findAllResponseMovie, _ = movieService.FindAll(context.TODO(), &protoMovie.FindAllRequest{})
	fmt.Printf("Found %v\n", findAllResponseMovie)

}
