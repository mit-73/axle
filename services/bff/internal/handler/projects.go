package handler

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgxpool"

	bffv1 "github.com/ApeironFoundation/axle/contracts/go/bff/v1"
	"github.com/ApeironFoundation/axle/contracts/go/bff/v1/gen_bff_v1connect"
)

// Compile-time interface check.
var _ gen_bff_v1connect.ProjectServiceHandler = (*ProjectsHandler)(nil)

// ProjectsHandler implements the bff.v1.ProjectService ConnectRPC methods.
// Stub implementation — real DB queries wired when migrations are applied.
type ProjectsHandler struct {
	Pool *pgxpool.Pool
}

func (h *ProjectsHandler) ListProjects(
	_ context.Context,
	req *bffv1.ListProjectsRequest,
) (*bffv1.ListProjectsResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 20
	}
	// TODO: query pool for real data after goose migrations are applied.
	return &bffv1.ListProjectsResponse{
		Projects: []*bffv1.Project{},
		Total:    0,
	}, nil
}

func (h *ProjectsHandler) GetProject(
	_ context.Context,
	req *bffv1.GetProjectRequest,
) (*bffv1.GetProjectResponse, error) {
	return nil, connect.NewError(connect.CodeNotFound,
		errors.New("project not found: "+req.GetId()))
}

func (h *ProjectsHandler) CreateProject(
	_ context.Context,
	req *bffv1.CreateProjectRequest,
) (*bffv1.CreateProjectResponse, error) {
	// TODO: insert into DB.
	return &bffv1.CreateProjectResponse{
		Project: &bffv1.Project{
			Id:          "00000000-0000-0000-0000-000000000000",
			Name:        req.GetName(),
			Description: req.GetDescription(),
			Status:      bffv1.ProjectStatus_PROJECT_STATUS_ACTIVE,
		},
	}, nil
}

func (h *ProjectsHandler) UpdateProject(
	_ context.Context,
	_ *bffv1.UpdateProjectRequest,
) (*bffv1.UpdateProjectResponse, error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

func (h *ProjectsHandler) DeleteProject(
	_ context.Context,
	_ *bffv1.DeleteProjectRequest,
) (*bffv1.DeleteProjectResponse, error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}
