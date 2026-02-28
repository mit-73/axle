package handler

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	bffv1 "github.com/ApeironFoundation/axle/contracts/generated/go/bff/v1"
	"github.com/ApeironFoundation/axle/contracts/generated/go/bff/v1/bffv1connect"
)

// Compile-time interface check.
var _ bffv1connect.UserServiceHandler = (*UsersHandler)(nil)

// UsersHandler implements the bff.v1.UserService ConnectRPC methods.
type UsersHandler struct{}

func (h *UsersHandler) ListUsers(
	_ context.Context,
	_ *bffv1.ListUsersRequest,
) (*bffv1.ListUsersResponse, error) {
	return &bffv1.ListUsersResponse{
		Users: []*bffv1.User{},
		Total: 0,
	}, nil
}

func (h *UsersHandler) GetUser(
	_ context.Context,
	req *bffv1.GetUserRequest,
) (*bffv1.GetUserResponse, error) {
	return nil, connect.NewError(connect.CodeNotFound,
		errors.New("user not found: "+req.GetId()))
}

func (h *UsersHandler) GetMe(
	_ context.Context,
	_ *bffv1.GetMeRequest,
) (*bffv1.GetMeResponse, error) {
	// TODO: extract user from auth context (Ory integration).
	return &bffv1.GetMeResponse{
		User: &bffv1.User{
			Id:    "00000000-0000-0000-0000-000000000000",
			Name:  "Demo User",
			Email: "demo@github.com/ApeironFoundation/axle",
			Role:  bffv1.UserRole_USER_ROLE_ADMIN,
		},
	}, nil
}

func (h *UsersHandler) UpdateUser(
	_ context.Context,
	_ *bffv1.UpdateUserRequest,
) (*bffv1.UpdateUserResponse, error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}
