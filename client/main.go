package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	proto "github.com/ob-vss-ss19/blatt-4-forever_alone_2_electric_boogaloo/cinema_hall/proto"
)

func main() {
	service := micro.NewService(micro.Name("cinema.client"))
	service.Init()

	client := proto.NewCinemaHallService("cinema.cinema_hall.service", service.Client())

	createResponse, _ := client.Create(context.TODO(), &proto.CreateRequest{Name: "Kino 1", Cols: 12, Rows: 8})
	fmt.Printf("Created %v\n", createResponse)
	createResponse, _ = client.Create(context.TODO(), &proto.CreateRequest{Name: "Kino 2", Cols: 20, Rows: 10})
	fmt.Printf("Created %v\n", createResponse)
	findAllResponse, _ := client.FindAll(context.TODO(), &proto.FindAllRequest{})
	fmt.Printf("Found %v\n", findAllResponse)
	deleteResponse, _ := client.Delete(context.TODO(), &proto.DeleteRequest{Id: 1})
	fmt.Printf("Deleted %v\n", deleteResponse)
	findAllResponse, _ = client.FindAll(context.TODO(), &proto.FindAllRequest{})
	fmt.Printf("Found %v\n", findAllResponse)

}
