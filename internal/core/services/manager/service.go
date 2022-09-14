package manager

import (
	"context"
	ports "f5ipmanager/internal/core/ports/manager"
	spec "f5ipmanager/internal/core/spec/manager"
)

type service struct {
	repo ports.ManagerRepository
}

type Params struct {
	repo ports.ManagerRepository
}

func NewService(params Params) ports.ManagerService {
	return &service{
		repo: params.repo,
	}
}

func (s *service) AllocateIPAddress(ctx context.Context, req spec.AllocateRequest) (spec.AllocateResponse, error) {
	ipAdd, err := s.repo.GetIPAddress(ctx, req.Label, req.Key)
	var resp = spec.AllocateResponse{}
	if err != nil {
		return resp, err
	}
	if ipAdd == "" {
		ipAdd2, err := s.repo.AllocateIPAddress(ctx, req.Label, req.Key)
		if err != nil {
			return resp, err
		}
		resp.IPAddress = ipAdd2
		resp.Success = true
		return resp, nil
	}
	resp.IPAddress = ipAdd
	resp.Success = true
	return resp, nil
}

func (s *service) DeallocateIPAddress(ctx context.Context, req spec.DeallocateRequest) (spec.DeallocateResponse, error) {
	ipAdd, err := s.repo.GetIPAddress(ctx, req.Label, req.Key)
	var resp = spec.DeallocateResponse{}
	if ipAdd == "" {
		resp.Success = true
		return resp, spec.ErrorNotFound
	}
	if err != nil {
		resp.Success = true
		return resp, err
	}
	err = s.repo.FreeIPAddress(ctx, req.Label, req.Key)
	if err != nil {
		resp.Success = true
		return resp, err
	}
	resp.Success = false
	return resp, nil
}
