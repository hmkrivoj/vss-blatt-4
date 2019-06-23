package main

import (
	"context"
	"errors"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/movie/proto"
	"sync"
)

type MovieHandler struct {
	mutex     sync.Mutex
	idCounter int64
	movies    map[int64]movie
	pub       micro.Publisher
}

func NewCinemaHallHandler(publisher micro.Publisher) *MovieHandler {
	handler := &MovieHandler{}
	handler.idCounter = 1
	handler.movies = make(map[int64]movie)
	handler.pub = publisher
	return handler
}

type movie struct {
	id    int64
	title string
}

func (handler *MovieHandler) Create(ctx context.Context, req *proto.CreateMovieRequest, res *proto.CreateMovieResponse) error {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	id := handler.idCounter
	handler.idCounter++

	mov := movie{
		id:    id,
		title: req.Title,
	}
	handler.movies[id] = mov

	res.Movie = &proto.Movie{
		Id:    mov.id,
		Title: mov.title,
	}

	return nil
}

func (handler *MovieHandler) Delete(ctx context.Context, req *proto.DeleteMovieRequest, res *proto.DeleteMovieResponse) error {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	if _, ok := handler.movies[req.Id]; !ok {
		return errors.New("no such id")
	}

	mov := handler.movies[req.Id]
	delete(handler.movies, req.Id)

	res.Movie = &proto.Movie{
		Id:    mov.id,
		Title: mov.title,
	}
	err := handler.pub.Publish(context.Background(), res)

	return err
}

func (handler *MovieHandler) FindAll(ctx context.Context, req *proto.FindAllMoviesRequest, res *proto.FindAllMoviesResponse) error {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	movies := make([]*proto.Movie, 0)
	for _, mov := range handler.movies {
		movies = append(movies, &proto.Movie{
			Id:    mov.id,
			Title: mov.title,
		})
	}
	res.Movies = movies

	return nil
}

func main() {
	service := micro.NewService(micro.Name("cinema.movie.service"))
	service.Init()

	publisher := micro.NewPublisher("cinema.movie.deleted", service.Client())
	handler := NewCinemaHallHandler(publisher)

	err := proto.RegisterMovieServiceHandler(service.Server(), handler)
	if err != nil {
		log.Fatal(err)
	}

	err = service.Run()
	if err != nil {
		log.Fatal(err)
	}
}
