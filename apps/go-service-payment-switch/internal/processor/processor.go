package processor

import (
	"fmt"
	paymentv1 "payment-svc/gen/payment/v1"
	"payment-svc/pkg/model/payment"
	v1 "payment-switch-svc/gen/payment_switch/v1"
	"strconv"
	"strings"
	"time"

	"github.com/moov-io/iso4217"
)

func BuildRequest(p *v1.Payload) (*paymentv1.AuthorisePaymentRequest, error) {
	if p.MessageTypeId != "0100" {
		// This service only supports authorisation requests
		return nil, ErrUnsupportedMessageType
	}

	// Account type unsupported
	if (p.ProcessingCode[2:4] != "00" && p.ProcessingCode[2:4] != "20") ||
		(p.ProcessingCode[4:6] != "00" && p.ProcessingCode[4:6] != "20") {
		return nil, ErrUnsupportedProcessingCode
	}

	// Parse the account number
	accountNumber, err := strconv.ParseInt(p.AccountNumber, 10, 64)
	if err != nil {
		return nil, ErrInvalidAccountNumber
	}

	// Infer payment type from the transaction type code
	switch p.ProcessingCode[0:2] {
	case "00": // Equivalent to a card payment
		if p.CardAcceptorId == nil {
			return nil, ErrMissingFields
		}

		// Parse the card acceptor id
		cardAcceptorId, err := strconv.ParseInt(*p.CardAcceptorId, 10, 64)
		if err != nil {
			return nil, ErrInvalidID
		}

		if p.TxCurrencyCode == nil {
			return nil, ErrMissingFields
		}

		// Parse the currency code
		cc, ccExists := iso4217.Lookup(strings.ToUpper(*p.TxCurrencyCode))
		if !ccExists {
			return nil, ErrInvalidCurrencyCode
		}

		return &paymentv1.AuthorisePaymentRequest{
			AccountId:    accountNumber,
			MerchantId:   &cardAcceptorId,
			Amount:       p.TxAmount,
			CurrencyCode: cc.Code,
			Description:  "Card Payment",
			Type:         string(payment.TypeCard),
		}, nil

	case "01": // Equivalent to a withdrawal
		if p.TxCurrencyCode == nil {
			return nil, ErrMissingFields
		}

		// Parse the currency code
		cc, ccExists := iso4217.Lookup(strings.ToUpper(*p.TxCurrencyCode))
		if !ccExists {
			return nil, ErrInvalidCurrencyCode
		}

		return &paymentv1.AuthorisePaymentRequest{
			AccountId:    accountNumber,
			Amount:       p.TxAmount,
			CurrencyCode: cc.Code,
			Description:  "Withdrawal",
			Type:         string(payment.TypeWithdrawal),
		}, nil

	case "21": // Equivalent to a deposit
		if p.TxCurrencyCode == nil {
			return nil, ErrMissingFields
		}

		// Parse the currency code
		cc, ccExists := iso4217.Lookup(strings.ToUpper(*p.TxCurrencyCode))
		if !ccExists {
			return nil, ErrInvalidCurrencyCode
		}

		return &paymentv1.AuthorisePaymentRequest{
			AccountId:    accountNumber,
			Amount:       p.TxAmount,
			CurrencyCode: cc.Code,
			Description:  "Deposit",
			Type:         string(payment.TypeDeposit),
		}, nil

	default: // Unsupported processing code
		return nil, ErrUnsupportedProcessingCode
	}
}

func BuildResponse(p *v1.Payload, r *paymentv1.AuthorisePaymentResponse) (*v1.Payload, error) {
	paymentId := fmt.Sprintf("%d", r.PaymentId)
	var settlementDate string

	responseCode := "06" // Error
	if r.Decision == paymentv1.Decision_DECISION_APPROVED {
		responseCode = "00"
		t := time.Now()
		settlementDate = fmt.Sprintf("%02d%02d", t.Month(), t.Day())
	} else if r.Decision != paymentv1.Decision_DECISION_DECLINED {
		responseCode = "01"
	}

	return &v1.Payload{
		MessageTypeId:               "0110", // 0110 is the code for authorisation responses
		AccountNumber:               p.AccountNumber,
		ProcessingCode:              p.ProcessingCode,
		TxAmount:                    p.TxAmount,
		SettlementAmount:            &p.TxAmount,
		BillingAmount:               &p.TxAmount,
		TransmissionDate:            p.TransmissionDate,
		LocalTxTime:                 p.LocalTxTime,
		LocalTxDate:                 p.LocalTxDate,
		SettlementDate:              &settlementDate,
		Mcc:                         p.Mcc,
		PointOfServiceEntryMode:     p.PointOfServiceEntryMode,
		PointOfServiceConditionCode: p.PointOfServiceConditionCode,
		AuthorisationId:             &paymentId,
		ResponseCode:                &responseCode,
		CardAcceptorId:              p.CardAcceptorId,
		TxCurrencyCode:              p.TxCurrencyCode,
		SettlementCurrencyCode:      p.SettlementCurrencyCode,
		FromAccountId:               p.FromAccountId,
		ToAccountId:                 p.ToAccountId,
	}, nil
}
