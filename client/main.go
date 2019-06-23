package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	protoCinemaHall "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_hall/proto"
	protoMovie "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/movie/proto"
)

func main() {
	service := micro.NewService(micro.Name("cinema.client"))
	service.Init()

	movieService := protoMovie.NewMovieService("cinema.movie.service", service.Client())
	cinemaHallService := protoCinemaHall.NewCinemaHallService("cinema.cinema_hall.service", service.Client())

	createResponseMovie, _ := movieService.Create(context.TODO(), &protoMovie.CreateMovieRequest{Title: "Fight Club"})
	fmt.Printf("Created %v\n", createResponseMovie)
	createResponseMovie, _ = movieService.Create(context.TODO(), &protoMovie.CreateMovieRequest{Title: "Se7en"})
	fmt.Printf("Created %v\n", createResponseMovie)
	findAllResponseMovie, _ := movieService.FindAll(context.TODO(), &protoMovie.FindAllMoviesRequest{})
	fmt.Printf("Found %v\n", findAllResponseMovie)
	deleteResponseMovie, _ := movieService.Delete(context.TODO(), &protoMovie.DeleteMovieRequest{Id: 1})
	fmt.Printf("Deleted %v\n", deleteResponseMovie)
	findAllResponseMovie, _ = movieService.FindAll(context.TODO(), &protoMovie.FindAllMoviesRequest{})
	fmt.Printf("Found %v\n", findAllResponseMovie)

	createResponseHall, _ := cinemaHallService.Create(context.TODO(), &protoCinemaHall.CreateCinemaHallRequest{Name: "Alpha", Rows: 10, Cols: 20})
	fmt.Printf("Created %v\n", createResponseHall)
	createResponseHall, _ = cinemaHallService.Create(context.TODO(), &protoCinemaHall.CreateCinemaHallRequest{Name: "Beta", Rows: 20, Cols: 10})
	fmt.Printf("Created %v\n", createResponseHall)
	findAllResponseHall, _ := cinemaHallService.FindAll(context.TODO(), &protoCinemaHall.FindAllCinemaHallsRequest{})
	fmt.Printf("Found %v\n", findAllResponseHall)
	deleteResponseHall, _ := cinemaHallService.Delete(context.TODO(), &protoCinemaHall.DeleteCinemaHallRequest{Id: 1})
	fmt.Printf("Deleted %v\n", deleteResponseHall)
	findAllResponseHall, _ = cinemaHallService.FindAll(context.TODO(), &protoCinemaHall.FindAllCinemaHallsRequest{})
	fmt.Printf("Found %v\n", findAllResponseHall)

}
