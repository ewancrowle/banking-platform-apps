package main

import (
	accountv1 "account-svc/gen/account/v1"
	"account-svc/gen/account/v1/accountv1connect"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"identity-svc/gen/identity/v1/identityv1connect"
	"log"
	"net/http"
	v1 "oauth-svc/gen/oauth/v1"
	"oauth-svc/gen/oauth/v1/oauthv1connect"
	"oauth-svc/pkg/model/device"
	"oauth-svc/pkg/model/token"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kelseyhightower/envconfig"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"google.golang.org/protobuf/types/known/emptypb"
)

type config struct {
	port int `default:"8080"`
	// e.g. "http://localhost:8080"
	identityServiceAddr string `required:"true" split_words:"true"`
	accountServiceAddr  string `required:"true" split_words:"true"`
	jwtSecret           string `required:"true" split_words:"true"`
	// e.g. "postgres://postgres:@localhost:5432/test?sslmode=disable"
	dbAddr string `required:"true" split_words:"true"`
}

type service struct {
	oauthv1connect.OAuthServiceHandler
	db                    *bun.DB
	identityServiceClient identityv1connect.IdentityServiceClient
	accountServiceClient  accountv1connect.AccountServiceClient
	jwtSecret             string
}

var (
	accessTokenLifetime  = 30 * time.Minute
	refreshTokenLifetime = 7 * 24 * time.Hour
)

func (s service) Token(ctx context.Context, request *v1.TokenRequest) (*v1.TokenResponse, error) {
	response, err := s.accountServiceClient.VerifyCredentials(ctx, &accountv1.VerifyCredentialsRequest{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if response.Id == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("incorrect email or password"))
	}

	id, err := s.identityServiceClient.ID(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	d := &device.Device{
		ID:        id.Id,
		AccountID: response.Id,
		IPAddress: request.IpAddress,
		UserAgent: request.UserAgent,
	}

	err = d.Insert(ctx, s.db)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	id, err = s.identityServiceClient.ID(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	accessTokenExpiration := time.Now().Add(accessTokenLifetime)
	accessToken := token.AccessToken{
		Token: token.Token{
			ID:        id.Id,
			AccountID: response.Id,
			DeviceID:  &d.ID,
			ExpiresAt: accessTokenExpiration,
		},
	}

	err = accessToken.Insert(ctx, s.db)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	id, err = s.identityServiceClient.ID(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	refreshTokenString := base64.RawURLEncoding.EncodeToString(b)

	hash := sha256.Sum256([]byte(refreshTokenString))
	hashString := hex.EncodeToString(hash[:])

	refreshTokenExpiration := time.Now().Add(refreshTokenLifetime)
	refreshToken := token.RefreshToken{
		Token: token.Token{
			ID:        id.Id,
			AccountID: response.Id,
			DeviceID:  &d.ID,
			ExpiresAt: refreshTokenExpiration,
		},
		Hash: hashString,
	}

	err = refreshToken.Insert(ctx, s.db)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	c := jwt.RegisteredClaims{
		Issuer:    "oauth-svc",
		Subject:   strconv.FormatInt(*response.Id, 10),
		ExpiresAt: jwt.NewNumericDate(accessTokenExpiration),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        strconv.FormatInt(accessToken.ID, 10),
	}

	jwtAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &c)
	signedJwt, err := jwtAccessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &v1.TokenResponse{
		AccessToken:  signedJwt,
		ExpiresIn:    int32(accessTokenLifetime.Seconds()),
		RefreshToken: refreshTokenString,
	}, nil
}

func (s service) Refresh(ctx context.Context, request *v1.RefreshRequest) (*v1.RefreshResponse, error) {
	hash := sha256.Sum256([]byte(request.RefreshToken))
	hashString := hex.EncodeToString(hash[:])

	refreshToken, err := token.SelectRefreshTokenByHash(ctx, s.db, hashString)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("refresh token not found"))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if refreshToken.IsRevoked {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("refresh token is revoked"))
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("refresh token expired"))
	}

	id, err := s.identityServiceClient.ID(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	accessTokenExpiration := time.Now().Add(accessTokenLifetime)
	accessToken := token.AccessToken{
		Token: token.Token{
			ID:        id.Id,
			AccountID: refreshToken.AccountID,
			DeviceID:  refreshToken.DeviceID,
			ExpiresAt: accessTokenExpiration,
		},
	}

	err = accessToken.Insert(ctx, s.db)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	id, err = s.identityServiceClient.ID(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	refreshTokenString := base64.RawURLEncoding.EncodeToString(b)

	hash = sha256.Sum256([]byte(refreshTokenString))
	hashString = hex.EncodeToString(hash[:])

	refreshTokenExpiration := time.Now().Add(refreshTokenLifetime)
	refreshToken = &token.RefreshToken{
		Token: token.Token{
			ID:        id.Id,
			AccountID: refreshToken.AccountID,
			DeviceID:  refreshToken.DeviceID,
			ExpiresAt: refreshTokenExpiration,
		},
		Hash: hashString,
	}

	err = refreshToken.Insert(ctx, s.db)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	c := jwt.RegisteredClaims{
		Issuer:    "oauth-svc",
		Subject:   strconv.FormatInt(*refreshToken.AccountID, 10),
		ExpiresAt: jwt.NewNumericDate(accessTokenExpiration),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        strconv.FormatInt(accessToken.ID, 10),
	}

	jwtAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &c)
	signedJwt, err := jwtAccessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &v1.RefreshResponse{
		AccessToken:  signedJwt,
		ExpiresIn:    int32(accessTokenLifetime.Seconds()),
		RefreshToken: refreshTokenString,
	}, nil
}

func (s service) Introspect(ctx context.Context, request *v1.IntrospectRequest) (*v1.IntrospectResponse, error) {
	claims := &jwt.RegisteredClaims{}

	t, err := jwt.ParseWithClaims(request.AccessToken, claims, func(_ *jwt.Token) (any, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if !t.Valid {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("access t is invalid"))
	}

	accessTokenID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	accessToken, err := token.SelectAccessToken(ctx, s.db, accessTokenID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if accessToken.IsRevoked {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("access token is revoked"))
	}

	if accessToken.ExpiresAt.Before(time.Now()) {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("access token expired"))
	}

	return &v1.IntrospectResponse{
		Email: accessToken.Account.Email,
	}, nil
}

func main() {
	var c config
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(c.dbAddr)))

	db := bun.NewDB(sqlDB, pgdialect.New())
	db.WithQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv(),
	))

	identityServiceClient := identityv1connect.NewIdentityServiceClient(
		http.DefaultClient,
		c.identityServiceAddr,
	)

	accountServiceClient := accountv1connect.NewAccountServiceClient(
		http.DefaultClient,
		c.accountServiceAddr,
	)

	svc := service{
		db:                    db,
		identityServiceClient: identityServiceClient,
		accountServiceClient:  accountServiceClient,
		jwtSecret:             c.jwtSecret,
	}

	path, handler := oauthv1connect.NewOAuthServiceHandler(svc, connect.WithInterceptors(validate.NewInterceptor()))

	mux := http.NewServeMux()
	mux.Handle(path, handler)

	p := new(http.Protocols)
	p.SetHTTP1(true)
	// Use h2c so we can serve HTTP/2 without TLS.
	p.SetUnencryptedHTTP2(true)
	s := http.Server{
		Addr:      "localhost:8080",
		Handler:   mux,
		Protocols: p,
	}
	s.ListenAndServe()
}
