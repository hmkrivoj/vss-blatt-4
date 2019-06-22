package main

import (
	"context"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_hall/proto"
	"sync"
)

type CinemaHallHandler struct {
	mutex       sync.Mutex
	idCounter   int64
	cinemaHalls map[int64]cinemaHall
}

func NewCinemaHallHandler() *CinemaHallHandler {
	handler := &CinemaHallHandler{}
	handler.idCounter = 1
	handler.cinemaHalls = make(map[int64]cinemaHall)
	return handler
}

type cinemaHall struct {
	id   int64
	name string
	rows int32
	cols int32
}

func (handler *CinemaHallHandler) Create(ctx context.Context, req *proto.CreateRequest, res *proto.CreateResponse) error {
	handler.mutex.Lock()
	id := handler.idCounter
	handler.idCounter++
	handler.mutex.Unlock()

	hall := cinemaHall{
		id:   id,
		name: req.Name,
		rows: req.Rows,
		cols: req.Cols,
	}
	handler.cinemaHalls[id] = hall

	res.Hall = &proto.CinemaHall{
		Id:   hall.id,
		Name: hall.name,
		Rows: hall.rows,
		Cols: hall.cols,
	}

	return nil
}

func (*CinemaHallHandler) Delete(context.Context, *proto.DeleteRequest, *proto.DeleteResponse) error {
	panic("implement me")
}

func (handler *CinemaHallHandler) FindAll(ctx context.Context, req *proto.FindAllRequest, res *proto.FindAllResponse) error {
	halls := make([]*proto.CinemaHall, 0)
	for _, hall := range handler.cinemaHalls {
		halls = append(halls, &proto.CinemaHall{
			Id:   hall.id,
			Name: hall.name,
			Rows: hall.rows,
			Cols: hall.cols,
		})
	}
	res.Halls = halls

	return nil
}

func main() {
	service := micro.NewService(micro.Name("cinema.cinema_hall"))
	service.Init()

	err := proto.RegisterCinemaHallServiceHandler(service.Server(), NewCinemaHallHandler())
	if err != nil {
		log.Fatal(err)
	}

	err = service.Run()
	if err != nil {
		log.Fatal(err)
	}
}
