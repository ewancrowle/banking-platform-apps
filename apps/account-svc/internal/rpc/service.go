package rpc

import (
	"account-svc/internal/luhn"
	"account-svc/pkg/model/account"
	"context"
	"database/sql"
	"errors"
	identityv1c "identity-svc/gen/identity/v1/identityv1connect"
	"strings"

	"connectrpc.com/connect"
	"github.com/alexedwards/argon2id"
	"github.com/moov-io/iso4217"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	repo        Repo
	identitySvc identityv1c.IdentityServiceClient
}

func NewService(repo Repo, identitySvc identityv1c.IdentityServiceClient) *Service {
	return &Service{repo, identitySvc}
}

type NewAccountInfo struct {
	FirstName   string
	MiddleNames string
	LastName    string
	Email       string
	Password    string
	Line1       string
	Line2       string
	Town        string
	Postcode    string
}

func (s *Service) CreateAccount(ctx context.Context, i *NewAccountInfo) (*account.Account, error) {
	if exists, err := s.repo.ExistsAccountByEmail(ctx, i.Email); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrEmailExists
	}

	id, err := s.identitySvc.ID(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	n, err := luhn.Generate(8)
	if err != nil {
		return nil, err
	}

	pwd, err := argon2id.CreateHash(i.Password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	return s.repo.CreateAccount(ctx, &account.Account{
		ID:           id.Id,
		AccountNo:    n,
		FirstName:    i.FirstName,
		MiddleNames:  i.MiddleNames,
		LastName:     i.LastName,
		Email:        strings.ToLower(i.Email),
		PasswordHash: pwd,
		Line1:        i.Line1,
		Line2:        i.Line2,
		Town:         i.Town,
		Postcode:     i.Postcode,
		IsFrozen:     false,
		CurrencyCode: iso4217.GBP.Code,
	})
}

func (s *Service) GetAccount(ctx context.Context, id int64) (*account.Account, error) {
	if a, err := s.repo.ReadAccountByID(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrIDNotFound
		}
		return nil, err
	} else {
		return a, nil
	}
}

func (s *Service) VerifyCredentials(ctx context.Context, email, password string) (*account.Account, error) {
	a, err := s.repo.ReadAccountByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEmailNotFound
		}
		return nil, err
	}

	if match, err := argon2id.ComparePasswordAndHash(password, a.PasswordHash); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	} else if !match {
		return nil, ErrPasswordMismatch
	}

	return a, nil
}
