// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: reservation.proto

package reservation

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for ReservationService service

type ReservationService interface {
	Create(ctx context.Context, in *CreateReservationRequest, opts ...client.CallOption) (*CreateReservationResponse, error)
	Confirm(ctx context.Context, in *ConfirmReservationRequest, opts ...client.CallOption) (*ConfirmReservationResponse, error)
	Delete(ctx context.Context, in *DeleteReservationRequest, opts ...client.CallOption) (*DeleteReservationResponse, error)
	FindAll(ctx context.Context, in *FindAllReservationsRequest, opts ...client.CallOption) (*FindAllReservationsResponse, error)
}

type reservationService struct {
	c    client.Client
	name string
}

func NewReservationService(name string, c client.Client) ReservationService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "reservationservice"
	}
	return &reservationService{
		c:    c,
		name: name,
	}
}

func (c *reservationService) Create(ctx context.Context, in *CreateReservationRequest, opts ...client.CallOption) (*CreateReservationResponse, error) {
	req := c.c.NewRequest(c.name, "ReservationService.Create", in)
	out := new(CreateReservationResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reservationService) Confirm(ctx context.Context, in *ConfirmReservationRequest, opts ...client.CallOption) (*ConfirmReservationResponse, error) {
	req := c.c.NewRequest(c.name, "ReservationService.Confirm", in)
	out := new(ConfirmReservationResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reservationService) Delete(ctx context.Context, in *DeleteReservationRequest, opts ...client.CallOption) (*DeleteReservationResponse, error) {
	req := c.c.NewRequest(c.name, "ReservationService.Delete", in)
	out := new(DeleteReservationResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reservationService) FindAll(ctx context.Context, in *FindAllReservationsRequest, opts ...client.CallOption) (*FindAllReservationsResponse, error) {
	req := c.c.NewRequest(c.name, "ReservationService.FindAll", in)
	out := new(FindAllReservationsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ReservationService service

type ReservationServiceHandler interface {
	Create(context.Context, *CreateReservationRequest, *CreateReservationResponse) error
	Confirm(context.Context, *ConfirmReservationRequest, *ConfirmReservationResponse) error
	Delete(context.Context, *DeleteReservationRequest, *DeleteReservationResponse) error
	FindAll(context.Context, *FindAllReservationsRequest, *FindAllReservationsResponse) error
}

func RegisterReservationServiceHandler(s server.Server, hdlr ReservationServiceHandler, opts ...server.HandlerOption) error {
	type reservationService interface {
		Create(ctx context.Context, in *CreateReservationRequest, out *CreateReservationResponse) error
		Confirm(ctx context.Context, in *ConfirmReservationRequest, out *ConfirmReservationResponse) error
		Delete(ctx context.Context, in *DeleteReservationRequest, out *DeleteReservationResponse) error
		FindAll(ctx context.Context, in *FindAllReservationsRequest, out *FindAllReservationsResponse) error
	}
	type ReservationService struct {
		reservationService
	}
	h := &reservationServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&ReservationService{h}, opts...))
}

type reservationServiceHandler struct {
	ReservationServiceHandler
}

func (h *reservationServiceHandler) Create(ctx context.Context, in *CreateReservationRequest, out *CreateReservationResponse) error {
	return h.ReservationServiceHandler.Create(ctx, in, out)
}

func (h *reservationServiceHandler) Confirm(ctx context.Context, in *ConfirmReservationRequest, out *ConfirmReservationResponse) error {
	return h.ReservationServiceHandler.Confirm(ctx, in, out)
}

func (h *reservationServiceHandler) Delete(ctx context.Context, in *DeleteReservationRequest, out *DeleteReservationResponse) error {
	return h.ReservationServiceHandler.Delete(ctx, in, out)
}

func (h *reservationServiceHandler) FindAll(ctx context.Context, in *FindAllReservationsRequest, out *FindAllReservationsResponse) error {
	return h.ReservationServiceHandler.FindAll(ctx, in, out)
}