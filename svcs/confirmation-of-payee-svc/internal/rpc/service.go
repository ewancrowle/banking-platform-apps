package rpc

import (
	"confirmation-of-payee-svc/pkg/model/token"
	"confirmation-of-payee-svc/pkg/opaquestr"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	identityv1c "identity-svc/gen/identity/v1/identityv1connect"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	repo        Repo
	identitySvc identityv1c.IdentityServiceClient
}

func NewService(repo Repo, identitySvc identityv1c.IdentityServiceClient) *Service {
	return &Service{repo, identitySvc}
}

func (s *Service) AcquireCOPToken(ctx context.Context, info AccountInfo) (string, error) {
	id, err := s.repo.ReadAccountIDByInfo(ctx, info)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrAccountNotFound
		}
		return "", err
	}

	t, err := opaquestr.Generate(32)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(t))
	hashStr := hex.EncodeToString(hash[:])

	tokId, err := s.identitySvc.ID(ctx, &emptypb.Empty{})
	if err != nil {
		return "", err
	}

	if _, err = s.repo.CreateCOPToken(ctx, &token.COPToken{
		ID:        tokId.Id,
		AccountID: id,
		Hash:      hashStr,
	}); err != nil {
		return "", connect.NewError(connect.CodeInternal, err)
	}

	return t, nil
}

func (s *Service) IntrospectCOPToken(ctx context.Context, token string) (*int64, error) {
	hash := sha256.Sum256([]byte(token))
	hashStr := hex.EncodeToString(hash[:])

	t, err := s.repo.ReadCOPTokenByHash(ctx, hashStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	if !t.UsedAt.IsZero() {
		return nil, ErrTokenAlreadyUsed
	}

	if err := s.repo.UpdateCOPTokenUsedAt(ctx, t.ID, t.UsedAt); err != nil {
		return nil, err
	}

	return &t.AccountID, nil
}
