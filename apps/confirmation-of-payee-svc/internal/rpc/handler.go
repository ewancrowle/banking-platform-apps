package rpc

import (
	confirmation_of_payeev1 "confirmation-of-payee-svc/gen/confirmation_of_payee/v1"
	"context"

	"connectrpc.com/connect"
)

type Handler struct {
	svc *Service
}

// ConfirmPayee implements [confirmation_of_payeev1connect.ConfirmationOfPayeeServiceHandler].
func (h *Handler) ConfirmPayee(ctx context.Context, req *confirmation_of_payeev1.ConfirmPayeeRequest) (*confirmation_of_payeev1.ConfirmPayeeResponse, error) {
	if t, err := h.svc.AcquireCOPToken(ctx, AccountInfo{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		AccountNo: req.AccountNum,
	}); err != nil {
		switch err {
		case ErrAccountNotFound:
			return nil, connect.NewError(connect.CodeNotFound, err)
		default:
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	} else {
		return &confirmation_of_payeev1.ConfirmPayeeResponse{ConfirmationOfPayeeToken: t}, nil
	}
}

// IntrospectToken implements [confirmation_of_payeev1connect.ConfirmationOfPayeeServiceHandler].
func (h *Handler) IntrospectToken(ctx context.Context, req *confirmation_of_payeev1.IntrospectTokenRequest) (*confirmation_of_payeev1.IntrospectTokenResponse, error) {
	id, err := h.svc.IntrospectCOPToken(ctx, req.ConfirmationOfPayeeToken)
	if err != nil {
		switch err {
		case ErrTokenNotFound:
			return nil, connect.NewError(connect.CodeNotFound, err)
		case ErrTokenAlreadyUsed:
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		default:
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}
	return &confirmation_of_payeev1.IntrospectTokenResponse{AccountId: *id}, nil
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc}
}
