package rpc

import (
	accountv1 "account-svc/gen/account/v1"
	"context"

	"connectrpc.com/connect"
)

type Handler struct {
	// v1c.AccountServiceHandler
	svc *Service
}

// CreateAccount implements [accountv1connect.AccountServiceHandler].
func (h *Handler) CreateAccount(ctx context.Context, r *accountv1.CreateAccountRequest) (*accountv1.CreateAccountResponse, error) {
	a, err := h.svc.CreateAccount(ctx, &NewAccountInfo{
		FirstName:   r.FirstName,
		MiddleNames: r.GetMiddleNames(),
		LastName:    r.LastName,
		Email:       r.Email,
		Password:    r.Password,
		Line1:       r.Line_1,
		Line2:       r.GetLine_2(),
		Town:        r.Town,
		Postcode:    r.Postcode,
	})
	if err != nil {
		switch err {
		case ErrEmailExists:
			return nil, connect.NewError(connect.CodeAlreadyExists, err)
		default:
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}
	return &accountv1.CreateAccountResponse{Id: a.ID}, nil
}

// GetAccount implements [accountv1connect.AccountServiceHandler].
func (h *Handler) GetAccount(ctx context.Context, r *accountv1.GetAccountRequest) (*accountv1.GetAccountResponse, error) {
	a, err := h.svc.GetAccount(ctx, r.Id)
	if err != nil {
		switch err {
		case ErrIDNotFound:
			return nil, connect.NewError(connect.CodeNotFound, err)
		default:
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}
	return &accountv1.GetAccountResponse{
		Id:           a.ID,
		AccountNum:   a.AccountNo,
		FirstName:    a.FirstName,
		MiddleNames:  a.MiddleNames,
		LastName:     a.LastName,
		Email:        a.Email,
		Line_1:       a.Line1,
		Line_2:       a.Line2,
		Town:         a.Town,
		Postcode:     a.Postcode,
		IsFrozen:     a.IsFrozen,
		CurrencyCode: a.CurrencyCode,
		CreatedAt:    a.CreatedAt.String(),
	}, nil
}

// VerifyCredentials implements [accountv1connect.AccountServiceHandler].
func (h *Handler) VerifyCredentials(ctx context.Context, r *accountv1.VerifyCredentialsRequest) (*accountv1.VerifyCredentialsResponse, error) {
	a, err := h.svc.VerifyCredentials(ctx, r.Email, r.Password)
	if err != nil {
		switch err {
		case ErrEmailNotFound:
			return nil, connect.NewError(connect.CodeNotFound, err)
		case ErrPasswordMismatch:
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		default:
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}
	return &accountv1.VerifyCredentialsResponse{Id: &a.ID}, nil
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc}
}
