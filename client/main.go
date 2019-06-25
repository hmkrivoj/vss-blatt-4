package main

import (
	"context"
	"fmt"
	"time"

	"github.com/micro/go-micro"
	protoHall "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_hall/proto"
	protoShowing "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_showing/proto"
	protoMovie "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/movie/proto"
	protoReservation "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/reservation/proto"
	protoUser "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/user/proto"
)

// Demo client for demonstration.
func main() {
	service := micro.NewService(micro.Name("cinema.client"))
	service.Init()

	// Init service bindings
	movieService := protoMovie.NewMovieService("cinema.movie.service", service.Client())
	hallService := protoHall.NewCinemaHallService("cinema.cinema_hall.service", service.Client())
	showingService := protoShowing.NewCinemaShowingService("cinema.cinema_showing.service", service.Client())
	userService := protoUser.NewUserService("cinema.user.service", service.Client())
	reservationService := protoReservation.NewReservationService("cinema.reservation.service", service.Client())

	// Create 4 movies and query them
	createResponseMovie, _ := movieService.Create(context.TODO(), &protoMovie.CreateMovieRequest{Title: "Prestige"})
	fmt.Printf("Created %v\n", createResponseMovie)
	createResponseMovie, _ = movieService.Create(context.TODO(), &protoMovie.CreateMovieRequest{Title: "Dark Knight"})
	fmt.Printf("Created %v\n", createResponseMovie)
	createResponseMovie, _ = movieService.Create(context.TODO(), &protoMovie.CreateMovieRequest{Title: "Inception"})
	fmt.Printf("Created %v\n", createResponseMovie)
	createResponseMovie, _ = movieService.Create(context.TODO(), &protoMovie.CreateMovieRequest{Title: "Interstellar"})
	fmt.Printf("Created %v\n", createResponseMovie)
	findAllResponseMovie, _ := movieService.FindAll(context.TODO(), &protoMovie.FindAllMoviesRequest{})
	fmt.Printf("Found %v\n", findAllResponseMovie)

	// Create 2 cinema halls and query them
	createResponseHall, _ := hallService.Create(
		context.TODO(),
		&protoHall.CreateCinemaHallRequest{
			Name: "Alpha",
			Rows: 10, Cols: 20,
		},
	)
	fmt.Printf("Created %v\n", createResponseHall)
	createResponseHall, _ = hallService.Create(
		context.TODO(),
		&protoHall.CreateCinemaHallRequest{
			Name: "Beta",
			Rows: 20,
			Cols: 10,
		},
	)
	fmt.Printf("Created %v\n", createResponseHall)
	findAllResponseHall, _ := hallService.FindAll(context.TODO(), &protoHall.FindAllCinemaHallsRequest{})
	fmt.Printf("Found %v\n", findAllResponseHall)

	// Create 4 users and query them
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

	// Create 4 showings and query them
	createResponseShowing, _ := showingService.Create(
		context.TODO(),
		&protoShowing.CreateCinemaShowingRequest{
			Movie:      4,
			CinemaHall: 2,
		},
	)
	fmt.Printf("Created %v\n", createResponseShowing)
	createResponseShowing, _ = showingService.Create(
		context.TODO(),
		&protoShowing.CreateCinemaShowingRequest{
			Movie:      3,
			CinemaHall: 2,
		},
	)
	fmt.Printf("Created %v\n", createResponseShowing)
	createResponseShowing, _ = showingService.Create(
		context.TODO(),
		&protoShowing.CreateCinemaShowingRequest{
			Movie:      1,
			CinemaHall: 1,
		},
	)
	fmt.Printf("Created %v\n", createResponseShowing)
	createResponseShowing, _ = showingService.Create(
		context.TODO(),
		&protoShowing.CreateCinemaShowingRequest{
			Movie:      2,
			CinemaHall: 2,
		},
	)
	fmt.Printf("Created %v\n", createResponseShowing)
	findAllResponseShowing, _ := showingService.FindAll(context.TODO(), &protoShowing.FindAllCinemaShowingsRequest{})
	fmt.Printf("Found %v\n", findAllResponseShowing)

	// Try creating 2 conflicting reservations, one invalid and one without problems
	createResponseReservation, _ := reservationService.Create(
		context.TODO(),
		&protoReservation.CreateReservationRequest{
			Showing: 1,
			User:    4,
			Seats: []*protoReservation.Seat{
				{Col: 3, Row: 2},
				{Col: 4, Row: 2},
			},
		},
	)
	fmt.Printf("Created %v\n", createResponseReservation)
	createResponseReservation, _ = reservationService.Create(
		context.TODO(),
		&protoReservation.CreateReservationRequest{
			Showing: 1,
			User:    3,
			Seats: []*protoReservation.Seat{
				{Col: 4, Row: 2},
				{Col: 5, Row: 2},
			},
		},
	)
	fmt.Printf("Created %v\n", createResponseReservation)
	createResponseReservation, _ = reservationService.Create(
		context.TODO(),
		&protoReservation.CreateReservationRequest{
			Showing: 1,
			User:    2,
			Seats: []*protoReservation.Seat{
				{Col: 100, Row: 2000},
			},
		},
	)
	fmt.Printf("Created %v\n", createResponseReservation)
	createResponseReservation, _ = reservationService.Create(
		context.TODO(),
		&protoReservation.CreateReservationRequest{
			Showing: 2,
			User:    1,
			Seats: []*protoReservation.Seat{
				{Col: 5, Row: 5},
			},
		},
	)
	fmt.Printf("Created %v\n", createResponseReservation)
	// Query reservations, only two should be listed
	findAllResponseReservations, _ := reservationService.FindAll(
		context.TODO(),
		&protoReservation.FindAllReservationsRequest{},
	)
	fmt.Printf("Found %v\n", findAllResponseReservations)

	// Do one invalid and one correct confirmation request
	confirmResponse, _ := reservationService.Confirm(
		context.TODO(),
		&protoReservation.ConfirmReservationRequest{
			Id:    2,
			Token: "xxxxxx",
		},
	)
	fmt.Printf("Confirmation %v\n", confirmResponse)
	confirmResponse, _ = reservationService.Confirm(
		context.TODO(),
		&protoReservation.ConfirmReservationRequest{
			Id:    2,
			Token: createResponseReservation.Reservation.Token,
		},
	)
	fmt.Printf("Confirmation %v\n", confirmResponse)

	// Deleting user 1 should not be possible because he has a reservation for showing 2
	deleteResponseUser, _ := userService.Delete(context.TODO(), &protoUser.DeleteUserRequest{Id: 1})
	fmt.Printf("Deleted %v\n", deleteResponseUser)
	findAllResponseUser, _ := userService.FindAll(context.TODO(), &protoUser.FindAllUsersRequest{})
	fmt.Printf("Found %v\n", findAllResponseUser)

	// Deleting movie 3 should eventually delete showing 2 and thus reservation 4, enabling the deletion of user 1
	deleteResponseMovie, _ := movieService.Delete(context.TODO(), &protoMovie.DeleteMovieRequest{Id: 3})
	fmt.Printf("Deleted %v\n", deleteResponseMovie)
	time.Sleep(1 * time.Second)
	deleteResponseUser, _ = userService.Delete(context.TODO(), &protoUser.DeleteUserRequest{Id: 1})
	fmt.Printf("Deleted %v\n", deleteResponseUser)
	findAllResponseUser, _ = userService.FindAll(context.TODO(), &protoUser.FindAllUsersRequest{})
	fmt.Printf("Found %v\n", findAllResponseUser)

	// Deleting hall 1 should eventually delete showing 3
	deleteResponseHall, _ := hallService.Delete(context.TODO(), &protoHall.DeleteCinemaHallRequest{Id: 1})
	fmt.Printf("Deleted %v\n", deleteResponseHall)
	time.Sleep(1 * time.Second)
	findAllResponseShowing, _ = showingService.FindAll(context.TODO(), &protoShowing.FindAllCinemaShowingsRequest{})
	fmt.Printf("Found %v\n", findAllResponseShowing)
}
