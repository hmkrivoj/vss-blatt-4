package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	protoCinemaHall "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_hall/proto"
	protoCinemaShowing "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_showing/proto"
	protoMovie "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/movie/proto"
	protoReservation "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/reservation/proto"
	protoUser "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/user/proto"
)

func main() {
	service := micro.NewService(micro.Name("cinema.client"))
	service.Init()

	movieService := protoMovie.NewMovieService("cinema.movie.service", service.Client())
	cinemaHallService := protoCinemaHall.NewCinemaHallService("cinema.cinema_hall.service", service.Client())
	cinemaShowingService := protoCinemaShowing.NewCinemaShowingService("cinema.cinema_showing.service", service.Client())
	userService := protoUser.NewUserService("cinema.user.service", service.Client())
	reservationService := protoReservation.NewReservationService("cinema.reservation.service", service.Client())

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

	createResponseShowing, _ := cinemaShowingService.Create(context.TODO(), &protoCinemaShowing.CreateCinemaShowingRequest{Movie: 2, CinemaHall: 2})
	fmt.Printf("Created %v\n", createResponseShowing)
	findAllResponseShowing, _ := cinemaShowingService.FindAll(context.TODO(), &protoCinemaShowing.FindAllCinemaShowingsRequest{})
	fmt.Printf("Found %v\n", findAllResponseShowing)

	createResponseUser, _ := userService.Create(context.TODO(), &protoUser.CreateUserRequest{Name: "Claire Grube"})
	fmt.Printf("Created %v\n", createResponseUser)
	createResponseUser, _ = userService.Create(context.TODO(), &protoUser.CreateUserRequest{Name: "Axel Schweiss"})
	fmt.Printf("Created %v\n", createResponseUser)
	createResponseUser, _ = userService.Create(context.TODO(), &protoUser.CreateUserRequest{Name: "Anna Bolika"})
	fmt.Printf("Created %v\n", createResponseUser)
	createResponseUser, _ = userService.Create(context.TODO(), &protoUser.CreateUserRequest{Name: "Andi Wand"})
	fmt.Printf("Created %v\n", createResponseUser)
	findAllResponseUsers, _ := userService.FindAll(context.TODO(), &protoUser.FindAllUsersRequest{})
	fmt.Printf("Found %v\n", findAllResponseUsers)

	createResponseReservation, _ := reservationService.Create(context.TODO(), &protoReservation.CreateReservationRequest{Showing: 1, User: 4, Seats: []*protoReservation.Seat{{Col: 3, Row: 2}, {Col: 4, Row: 2}}})
	fmt.Printf("Created %v\n", createResponseReservation)
	findAllResponseReservations, _ := reservationService.FindAll(context.TODO(), &protoReservation.FindAllReservationsRequest{})
	fmt.Printf("Found %v\n", findAllResponseReservations)
	confirmResponse, _ := reservationService.Confirm(context.TODO(), &protoReservation.ConfirmReservationRequest{Id: 1, Token: createResponseReservation.Reservation.Token})
	fmt.Printf("Created %v\n", confirmResponse)
	findAllResponseReservations, _ = reservationService.FindAll(context.TODO(), &protoReservation.FindAllReservationsRequest{})
	fmt.Printf("Found %v\n", findAllResponseReservations)
}
