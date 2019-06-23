// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: user.proto

package user

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

// Client API for UserService service

type UserService interface {
	Create(ctx context.Context, in *CreateUserRequest, opts ...client.CallOption) (*CreateUserResponse, error)
	Delete(ctx context.Context, in *DeleteUserRequest, opts ...client.CallOption) (*DeleteUserResponse, error)
	FindAll(ctx context.Context, in *FindAllUsersRequest, opts ...client.CallOption) (*FindAllUsersResponse, error)
}

type userService struct {
	c    client.Client
	name string
}

func NewUserService(name string, c client.Client) UserService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "userservice"
	}
	return &userService{
		c:    c,
		name: name,
	}
}

func (c *userService) Create(ctx context.Context, in *CreateUserRequest, opts ...client.CallOption) (*CreateUserResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.Create", in)
	out := new(CreateUserResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) Delete(ctx context.Context, in *DeleteUserRequest, opts ...client.CallOption) (*DeleteUserResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.Delete", in)
	out := new(DeleteUserResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) FindAll(ctx context.Context, in *FindAllUsersRequest, opts ...client.CallOption) (*FindAllUsersResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.FindAll", in)
	out := new(FindAllUsersResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserService service

type UserServiceHandler interface {
	Create(context.Context, *CreateUserRequest, *CreateUserResponse) error
	Delete(context.Context, *DeleteUserRequest, *DeleteUserResponse) error
	FindAll(context.Context, *FindAllUsersRequest, *FindAllUsersResponse) error
}

func RegisterUserServiceHandler(s server.Server, hdlr UserServiceHandler, opts ...server.HandlerOption) error {
	type userService interface {
		Create(ctx context.Context, in *CreateUserRequest, out *CreateUserResponse) error
		Delete(ctx context.Context, in *DeleteUserRequest, out *DeleteUserResponse) error
		FindAll(ctx context.Context, in *FindAllUsersRequest, out *FindAllUsersResponse) error
	}
	type UserService struct {
		userService
	}
	h := &userServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&UserService{h}, opts...))
}

type userServiceHandler struct {
	UserServiceHandler
}

func (h *userServiceHandler) Create(ctx context.Context, in *CreateUserRequest, out *CreateUserResponse) error {
	return h.UserServiceHandler.Create(ctx, in, out)
}

func (h *userServiceHandler) Delete(ctx context.Context, in *DeleteUserRequest, out *DeleteUserResponse) error {
	return h.UserServiceHandler.Delete(ctx, in, out)
}

func (h *userServiceHandler) FindAll(ctx context.Context, in *FindAllUsersRequest, out *FindAllUsersResponse) error {
	return h.UserServiceHandler.FindAll(ctx, in, out)
}