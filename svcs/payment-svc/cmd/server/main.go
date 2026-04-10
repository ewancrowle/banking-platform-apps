package main

import (
	accountv1 "account-svc/gen/account/v1"
	"account-svc/gen/account/v1/accountv1connect"
	confirmation_of_payeev1 "confirmation-of-payee-svc/gen/confirmation_of_payee/v1"
	"confirmation-of-payee-svc/gen/confirmation_of_payee/v1/confirmation_of_payeev1connect"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"identity-svc/gen/identity/v1/identityv1connect"
	"log"
	merchantv1 "merchant-svc/gen/merchant/v1"
	"merchant-svc/gen/merchant/v1/merchantv1connect"
	"net/http"
	paymentdecisionv1 "payment-decision-svc/gen/payment_decision/v1"
	"payment-decision-svc/gen/payment_decision/v1/payment_decisionv1connect"
	"payment-decision-svc/pkg/model/paymentdecision"
	v1 "payment-svc/gen/payment/v1"
	"payment-svc/gen/payment/v1/paymentv1connect"
	"payment-svc/pkg/config"
	"payment-svc/pkg/model/payment"
	"strings"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/kelseyhightower/envconfig"
	"github.com/moov-io/iso4217"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"google.golang.org/protobuf/types/known/emptypb"
)

type service struct {
	paymentv1connect.PaymentServiceHandler
	db                               *bun.DB
	kafkaCl                          *kgo.Client
	identityServiceClient            identityv1connect.IdentityServiceClient
	accountServiceClient             accountv1connect.AccountServiceClient
	merchantServiceClient            merchantv1connect.MerchantServiceClient
	paymentDecisionClient            payment_decisionv1connect.PaymentDecisionServiceClient
	confirmationOfPayeeServiceClient confirmation_of_payeev1connect.ConfirmationOfPayeeServiceClient
}

func (s service) AuthorisePayment(ctx context.Context, request *v1.AuthorisePaymentRequest) (*v1.AuthorisePaymentResponse, error) {
	if request.Amount == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("amount cannot be zero"))
	}

	if cc, exists := iso4217.Lookup(strings.ToUpper(request.CurrencyCode)); !exists {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("invalid currency code"))
	} else if cc != iso4217.GBP {
		//TODO FX
		return nil, connect.NewError(connect.CodeUnimplemented, errors.New("unsupported currency"))
	}

	paymentType, err := payment.GetType(request.Type)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if _, err := s.accountServiceClient.GetAccount(ctx, &accountv1.GetAccountRequest{
		Id: request.AccountId,
	}); err != nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("account not found"))
	}

	if paymentType == payment.TypeCard {
		if request.MerchantId == nil {
			return nil, connect.NewError(connect.CodeUnimplemented, errors.New("merchant id is required"))
		}

		if _, err := s.merchantServiceClient.GetMerchant(ctx, &merchantv1.GetMerchantRequest{Id: request.GetMerchantId()}); err != nil {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("merchant not found"))
		}
	}

	correctedAmount := paymentType.GetCorrectDirection(request.Amount)
	var otherAccountId *int64

	if paymentType == payment.TypeOutboundTransfer {
		if request.ConfirmationOfPayeeToken == nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("confirmation of payee token is required"))
		}

		i, err := s.confirmationOfPayeeServiceClient.IntrospectToken(ctx, &confirmation_of_payeev1.IntrospectTokenRequest{
			ConfirmationOfPayeeToken: request.GetConfirmationOfPayeeToken(),
		})
		if err != nil {
			if connect.CodeOf(err) == connect.CodeInternal {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("confirmation of payee token is invalid"))
		}

		otherAccountId = &i.AccountId

		if request.AccountId == *otherAccountId {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("both account ids cannot be the same"))
		}

		if _, err := s.accountServiceClient.GetAccount(ctx, &accountv1.GetAccountRequest{
			Id: request.AccountId,
		}); err != nil {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("other account not found"))
		}
	}

	id, err := s.identityServiceClient.ID(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	p := payment.Payment{
		ID:             id.Id,
		AccountID:      request.AccountId,
		MerchantID:     request.MerchantId,
		OtherAccountID: otherAccountId,
		Amount:         correctedAmount,
		CurrencyCode:   request.CurrencyCode,
		Description:    request.Description,
		Type:           paymentType,
		Status:         payment.StatusReceived,
		CreatedAt:      time.Now(),
	}

	if err = p.Insert(ctx, s.db); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	b, err := json.Marshal(p)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := s.kafkaCl.ProduceSync(ctx, &kgo.Record{
		Value: b,
		Topic: "payment_progress",
	}).FirstErr(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	d, err := s.paymentDecisionClient.DecidePayment(ctx, &paymentdecisionv1.DecidePaymentRequest{
		PaymentId:      id.Id,
		AccountId:      request.AccountId,
		MerchantId:     request.MerchantId,
		OtherAccountId: otherAccountId,
		Amount:         correctedAmount,
		CurrencyCode:   request.CurrencyCode,
		Description:    request.Description,
		Type:           request.Type,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	switch d.Decision {
	case paymentdecisionv1.Decision_DECISION_UNSPECIFIED:
		return nil, connect.NewError(connect.CodeInternal, errors.New("unspecified decision"))

	case paymentdecisionv1.Decision_DECISION_DECLINED:
		p.Status = payment.StatusDeclined

		dr := paymentdecision.DeclineReason(d.DeclineReason.Number())
		p.DeclineReason = &dr
		if err = p.SetDeclineReason(ctx, s.db, dr); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		p.UpdatedAt = time.Now()
		if err = p.SetUpdatedAt(ctx, s.db, p.UpdatedAt); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	default:
		p.Status = payment.StatusAuthorised
	}

	if err = p.SetStatus(ctx, s.db, p.Status); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	p.UpdatedAt = time.Now()
	if err = p.SetUpdatedAt(ctx, s.db, p.UpdatedAt); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	b, err = json.Marshal(p)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := s.kafkaCl.ProduceSync(ctx, &kgo.Record{
		Value: b,
		Topic: "payment_progress",
	}).FirstErr(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if p.Type == payment.TypeOutboundTransfer && p.Status == payment.StatusAuthorised {
		creditId, err := s.identityServiceClient.ID(ctx, &emptypb.Empty{})
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		creditPayment := payment.Payment{
			ID:             creditId.Id,
			AccountID:      *p.OtherAccountID,
			OtherAccountID: &p.AccountID,
			Amount:         payment.TypeInboundTransfer.GetCorrectDirection(request.Amount),
			CurrencyCode:   request.CurrencyCode,
			Description:    request.Description,
			Type:           payment.TypeInboundTransfer,
			Status:         payment.StatusCaptured,
			CreatedAt:      time.Now(),
		}

		if err = creditPayment.Insert(ctx, s.db); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		b, err = json.Marshal(creditPayment)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		if err := s.kafkaCl.ProduceSync(ctx, &kgo.Record{
			Value: b,
			Topic: "payment_progress",
		}).FirstErr(); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		p.Status = payment.StatusCaptured
		if err := p.SetStatus(ctx, s.db, p.Status); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		p.UpdatedAt = time.Now()
		if err = p.SetUpdatedAt(ctx, s.db, p.UpdatedAt); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		b, err = json.Marshal(p)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		if err := s.kafkaCl.ProduceSync(ctx, &kgo.Record{
			Value: b,
			Topic: "payment_progress",
		}).FirstErr(); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	return &v1.AuthorisePaymentResponse{
		PaymentId:  id.Id,
		Decision:   v1.Decision(d.Decision),
		DecisionId: d.DecisionId,
	}, nil
}

func (s service) IncrementPayment(ctx context.Context, request *v1.IncrementPaymentRequest) (*emptypb.Empty, error) {
	p, err := payment.Select(ctx, s.db, request.PaymentId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	p.Amount = request.ReplacementAmount
	if err := p.SetAmount(ctx, s.db, p.Amount); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	p.UpdatedAt = time.Now()
	if err := p.SetUpdatedAt(ctx, s.db, p.UpdatedAt); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	p.Status = payment.StatusIncremented
	// Incremented status is not stored in the SQL database, so we don't need to update it here.

	b, err := json.Marshal(p)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := s.kafkaCl.ProduceSync(ctx, &kgo.Record{
		Value: b,
		Topic: "payment_progress",
	}).FirstErr(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &emptypb.Empty{}, nil
}

func (s service) CapturePayment(ctx context.Context, request *v1.CapturePaymentRequest) (*emptypb.Empty, error) {
	p, err := payment.Select(ctx, s.db, request.PaymentId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if p.Status != payment.StatusAuthorised {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("payment is not authorised"))
	} else if p.Status == payment.StatusCaptured {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("payment is already captured"))
	}

	if request.ReplacementAmount != nil {
		p.Amount = p.Type.GetCorrectDirection(*request.ReplacementAmount)
		if err := p.SetAmount(ctx, s.db, p.Amount); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		p.UpdatedAt = time.Now()
		if err := p.SetUpdatedAt(ctx, s.db, p.UpdatedAt); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	p.Status = payment.StatusCaptured
	if err := p.SetStatus(ctx, s.db, p.Status); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	p.UpdatedAt = time.Now()
	if err := p.SetUpdatedAt(ctx, s.db, p.UpdatedAt); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	b, err := json.Marshal(p)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := s.kafkaCl.ProduceSync(ctx, &kgo.Record{
		Value: b,
		Topic: "payment_progress",
	}).FirstErr(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &emptypb.Empty{}, nil
}

func (s service) ExpirePayment(ctx context.Context, request *v1.ExpirePaymentRequest) (*emptypb.Empty, error) {
	p, err := payment.Select(ctx, s.db, request.PaymentId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	p.Status = payment.StatusExpired
	if err := p.SetStatus(ctx, s.db, p.Status); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	p.UpdatedAt = time.Now()
	if err := p.SetUpdatedAt(ctx, s.db, p.UpdatedAt); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	b, err := json.Marshal(p)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := s.kafkaCl.ProduceSync(ctx, &kgo.Record{
		Value: b,
		Topic: "payment_progress",
	}).FirstErr(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &emptypb.Empty{}, nil
}

func (s service) VoidPayment(ctx context.Context, request *v1.VoidPaymentRequest) (*emptypb.Empty, error) {
	p, err := payment.Select(ctx, s.db, request.PaymentId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	p.Status = payment.StatusVoided
	if err := p.SetStatus(ctx, s.db, p.Status); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	p.UpdatedAt = time.Now()
	if err := p.SetUpdatedAt(ctx, s.db, p.UpdatedAt); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	b, err := json.Marshal(p)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := s.kafkaCl.ProduceSync(ctx, &kgo.Record{
		Value: b,
		Topic: "payment_progress",
	}).FirstErr(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &emptypb.Empty{}, nil
}

func (s service) GetPayments(ctx context.Context, req *v1.GetPaymentsRequest) (*v1.GetPaymentsResponse, error) {
	payments, err := payment.SelectDisplayableByAccountID(ctx, s.db, req.AccountId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	paymentResponses := make([]*v1.Payment, 0, len(payments))
	for _, p := range payments {
		var otherAccountName *string
		if p.OtherAccountID != nil {
			name := fmt.Sprintf("%s %s", p.OtherAccount.FirstName, p.OtherAccount.LastName)
			otherAccountName = &name
		}

		var merchant *v1.Merchant
		if p.Merchant != nil {
			merchant = &v1.Merchant{
				Id:              p.Merchant.ID,
				Descriptor_:     p.Merchant.Descriptor,
				ShortDescriptor: p.Merchant.ShortDescriptor,
				Mcc:             p.Merchant.MCC,
				Line_1:          p.Merchant.Line1,
				Line_2:          p.Merchant.Line2,
				Town:            p.Merchant.Town,
				Postcode:        p.Merchant.Postcode,
				CountryCode:     p.Merchant.CountryCode,
				CreatedAt:       p.Merchant.CreatedAt.String(),
				UpdatedAt:       p.Merchant.UpdatedAt.String(),
			}
		}

		var declineReason v1.DeclineReason
		if p.DeclineReason != nil {
			declineReason = v1.DeclineReason(*p.DeclineReason)
		}

		var otherAccountData v1.OtherAccountData
		if p.OtherAccountID != nil {
			otherAccountData = v1.OtherAccountData{
				FirstName:   p.OtherAccount.FirstName,
				MiddleNames: &p.OtherAccount.MiddleNames,
				LastName:    p.OtherAccount.LastName,
			}
		}

		paymentResponses = append(paymentResponses, &v1.Payment{
			Id:               p.ID,
			AccountId:        p.AccountID,
			MerchantId:       p.MerchantID,
			OtherAccountId:   p.OtherAccountID,
			OtherAccountData: &otherAccountData,
			Amount:           p.Amount,
			CurrencyCode:     p.CurrencyCode,
			Type:             string(p.Type),
			Status:           string(p.Status),
			Description:      p.Description,
			CreatedAt:        p.CreatedAt.String(),
			UpdatedAt:        p.UpdatedAt.String(),
			OtherAccountName: otherAccountName,
			Merchant:         merchant,
			DeclineReason:    &declineReason,
		})
	}

	return &v1.GetPaymentsResponse{Payments: paymentResponses}, nil
}

func main() {
	var c config.Config
	if err := envconfig.Process("", &c); err != nil {
		log.Fatal(err.Error())
	}

	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithAddr(c.DBHost),
		pgdriver.WithDatabase(c.DBName),
		pgdriver.WithUser(c.DBUsername),
		pgdriver.WithPassword(c.DBPassword),
		pgdriver.WithInsecure(true),
	))

	db := bun.NewDB(sqlDB, pgdialect.New()).WithQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(true),
		bundebug.FromEnv(),
	))

	kafkaCl, err := kgo.NewClient(kgo.SeedBrokers(c.KafkaBrokers...))
	if err != nil {
		log.Fatal(err.Error())
	}

	identityServiceClient := identityv1connect.NewIdentityServiceClient(
		http.DefaultClient,
		c.IdentityServiceAddr,
	)

	accountServiceClient := accountv1connect.NewAccountServiceClient(
		http.DefaultClient,
		c.AccountServiceAddr,
	)

	merchantServiceClient := merchantv1connect.NewMerchantServiceClient(
		http.DefaultClient,
		c.MerchantServiceAddr,
	)

	paymentDecisionServiceClient := payment_decisionv1connect.NewPaymentDecisionServiceClient(
		http.DefaultClient,
		c.PaymentDecisionServiceAddr,
	)

	confirmationOfPayeeServiceClient := confirmation_of_payeev1connect.NewConfirmationOfPayeeServiceClient(
		http.DefaultClient,
		c.COPServiceAddr,
	)

	svc := service{
		db:                               db,
		kafkaCl:                          kafkaCl,
		identityServiceClient:            identityServiceClient,
		accountServiceClient:             accountServiceClient,
		merchantServiceClient:            merchantServiceClient,
		paymentDecisionClient:            paymentDecisionServiceClient,
		confirmationOfPayeeServiceClient: confirmationOfPayeeServiceClient,
	}

	path, handler := paymentv1connect.NewPaymentServiceHandler(svc, connect.WithInterceptors(validate.NewInterceptor()))

	mux := http.NewServeMux()
	mux.Handle(path, handler)

	p := new(http.Protocols)
	p.SetHTTP1(true)
	// Use h2c so we can serve HTTP/2 without TLS.
	p.SetUnencryptedHTTP2(true)
	s := http.Server{
		Addr:      fmt.Sprintf(":%d", c.Port),
		Handler:   mux,
		Protocols: p,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
