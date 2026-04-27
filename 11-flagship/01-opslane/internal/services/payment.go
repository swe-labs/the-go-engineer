// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/db"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
	paymentflow "github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/payment"
)

var (
	ErrInvalidPayment  = errors.New("invalid payment")
	ErrPaymentNotFound = errors.New("payment not found")
)

const defaultPaymentGatewayTimeout = 2 * time.Second
const defaultPaymentAttempts = 2

// PaymentRepository is the persistence surface the payment workflow needs.
type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *models.Payment) error
	GetPaymentByProviderReference(ctx context.Context, tenantID int64, providerReference string) (models.Payment, error)
	UpdatePaymentStatus(ctx context.Context, tenantID int64, providerReference string, status models.PaymentStatus, failureReason string) (models.Payment, error)
}

// PaymentOrderReader lets the payment service validate tenant-scoped order facts.
type PaymentOrderReader interface {
	GetOrderByID(ctx context.Context, tenantID, orderID int64) (models.Order, error)
}

// PaymentOrderWorkflow is the order state machine seam used by payments.
type PaymentOrderWorkflow interface {
	TransitionOrder(ctx context.Context, req TransitionOrderRequest) (models.Order, error)
}

type PaymentServiceOptions struct {
	GatewayTimeout time.Duration
	MaxAttempts    int
}

type ProcessPaymentResult struct {
	Payment models.Payment
	Created bool
}

type ReconcilePaymentInput struct {
	TenantID          int64
	ProviderReference string
	Status            models.PaymentStatus
	FailureReason     string
}

// PaymentService owns duplicate protection, gateway calls, and payment-driven order transitions.
type PaymentService struct {
	payments      PaymentRepository
	orders        PaymentOrderReader
	orderWorkflow PaymentOrderWorkflow
	gateway       paymentflow.Gateway
	options       PaymentServiceOptions
}

func NewPaymentService(
	payments PaymentRepository,
	orders PaymentOrderReader,
	orderWorkflow PaymentOrderWorkflow,
	gateway paymentflow.Gateway,
	options PaymentServiceOptions,
) *PaymentService {
	if options.GatewayTimeout <= 0 {
		options.GatewayTimeout = defaultPaymentGatewayTimeout
	}
	if options.MaxAttempts <= 0 {
		options.MaxAttempts = defaultPaymentAttempts
	}

	return &PaymentService{
		payments:      payments,
		orders:        orders,
		orderWorkflow: orderWorkflow,
		gateway:       gateway,
		options:       options,
	}
}

func (s *PaymentService) ProcessPayment(ctx context.Context, job paymentflow.Job) (ProcessPaymentResult, error) {
	if s == nil || s.payments == nil || s.orders == nil || s.orderWorkflow == nil || s.gateway == nil {
		return ProcessPaymentResult{}, fmt.Errorf("payment service is not configured")
	}

	normalized, err := normalizePaymentJob(job)
	if err != nil {
		return ProcessPaymentResult{}, err
	}

	order, err := s.orders.GetOrderByID(ctx, normalized.TenantID, normalized.OrderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ProcessPaymentResult{}, ErrOrderNotFound
		}
		return ProcessPaymentResult{}, fmt.Errorf("get payment order: %w", err)
	}
	if order.TotalCents != normalized.AmountCents {
		return ProcessPaymentResult{}, ErrInvalidPayment
	}

	existing, err := s.payments.GetPaymentByProviderReference(ctx, normalized.TenantID, normalized.ProviderReference)
	switch {
	case err == nil:
		if !samePaymentTarget(existing, normalized) {
			return ProcessPaymentResult{}, ErrInvalidPayment
		}
		if err := s.syncPaymentToOrder(ctx, existing, order); err != nil {
			return ProcessPaymentResult{Payment: existing, Created: false}, err
		}
		return ProcessPaymentResult{Payment: existing, Created: false}, nil
	case !errors.Is(err, sql.ErrNoRows):
		return ProcessPaymentResult{}, fmt.Errorf("lookup payment by provider reference: %w", err)
	}

	if order.Status != models.OrderStatusPending && order.Status != models.OrderStatusFailed {
		return ProcessPaymentResult{}, fmt.Errorf("order is not in a state that accepts new payments: current status %s", order.Status)
	}

	payment := models.Payment{
		TenantID:          normalized.TenantID,
		OrderID:           normalized.OrderID,
		Status:            models.PaymentStatusPending,
		ProviderReference: normalized.ProviderReference,
		AmountCents:       normalized.AmountCents,
	}

	if err := s.payments.CreatePayment(ctx, &payment); err != nil {
		if errors.Is(err, db.ErrInvalidReference) {
			return ProcessPaymentResult{}, ErrOrderNotFound
		}
		if errors.Is(err, db.ErrDuplicateValue) {
			existing, lookupErr := s.payments.GetPaymentByProviderReference(ctx, normalized.TenantID, normalized.ProviderReference)
			if lookupErr != nil {
				return ProcessPaymentResult{}, fmt.Errorf("load duplicate payment by provider reference: %w", lookupErr)
			}
			if !samePaymentTarget(existing, normalized) {
				return ProcessPaymentResult{}, ErrInvalidPayment
			}
			return ProcessPaymentResult{Payment: existing, Created: false}, nil
		}
		return ProcessPaymentResult{}, fmt.Errorf("create payment: %w", err)
	}

	return s.chargeGateway(ctx, payment, true)
}

func (s *PaymentService) ReconcilePayment(ctx context.Context, input ReconcilePaymentInput) (models.Payment, error) {
	if s == nil || s.payments == nil || s.orderWorkflow == nil {
		return models.Payment{}, fmt.Errorf("payment service is not configured")
	}

	input.ProviderReference = strings.TrimSpace(input.ProviderReference)
	input.FailureReason = strings.TrimSpace(input.FailureReason)
	if input.TenantID <= 0 || input.ProviderReference == "" || !isKnownPaymentStatus(input.Status) {
		return models.Payment{}, ErrInvalidPayment
	}

	_, err := s.payments.GetPaymentByProviderReference(ctx, input.TenantID, input.ProviderReference)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Payment{}, ErrPaymentNotFound
		}
		return models.Payment{}, fmt.Errorf("lookup payment by provider reference: %w", err)
	}

	updated, err := s.payments.UpdatePaymentStatus(ctx, input.TenantID, input.ProviderReference, input.Status, input.FailureReason)
	if err != nil {
		return models.Payment{}, fmt.Errorf("update payment status: %w", err)
	}

	if err := s.applyPaymentStatusToOrder(ctx, updated); err != nil {
		return updated, err
	}

	return updated, nil
}

func (s *PaymentService) chargeGateway(ctx context.Context, current models.Payment, created bool) (ProcessPaymentResult, error) {
	if err := s.moveOrderToProcessing(ctx, current); err != nil {
		return ProcessPaymentResult{}, err
	}

	var lastErr error
	for attempt := 0; attempt < s.options.MaxAttempts; attempt++ {
		callCtx, cancel := context.WithTimeout(ctx, s.options.GatewayTimeout)
		result, err := s.gateway.Charge(callCtx, paymentflow.GatewayRequest{
			TenantID:          current.TenantID,
			OrderID:           current.OrderID,
			ProviderReference: current.ProviderReference,
			AmountCents:       current.AmountCents,
		})
		cancel()

		if err == nil {
			if result.Status == "" {
				result.Status = models.PaymentStatusSettled
			}
			if !isKnownPaymentStatus(result.Status) {
				return ProcessPaymentResult{}, ErrInvalidPayment
			}

			updated, err := s.payments.UpdatePaymentStatus(ctx, current.TenantID, current.ProviderReference, result.Status, strings.TrimSpace(result.FailureReason))
			if err != nil {
				return ProcessPaymentResult{}, fmt.Errorf("update payment after gateway result: %w", err)
			}
			if err := s.applyPaymentStatusToOrder(ctx, updated); err != nil {
				return ProcessPaymentResult{Payment: updated, Created: created}, err
			}
			return ProcessPaymentResult{Payment: updated, Created: created}, nil
		}

		lastErr = normalizeGatewayError(err)
	}

	return ProcessPaymentResult{Payment: current, Created: created}, fmt.Errorf("charge payment: %w", lastErr)
}

func (s *PaymentService) moveOrderToProcessing(ctx context.Context, payment models.Payment) error {
	_, err := s.orderWorkflow.TransitionOrder(ctx, TransitionOrderRequest{
		TenantID: payment.TenantID,
		OrderID:  payment.OrderID,
		Status:   models.OrderStatusProcessing,
	})
	if err != nil {
		return fmt.Errorf("move order to processing: %w", err)
	}

	return nil
}

func (s *PaymentService) applyPaymentStatusToOrder(ctx context.Context, payment models.Payment) error {
	var target models.OrderStatus
	switch payment.Status {
	case models.PaymentStatusSettled:
		target = models.OrderStatusPaid
	case models.PaymentStatusFailed:
		target = models.OrderStatusFailed
	default:
		return nil
	}

	_, err := s.orderWorkflow.TransitionOrder(ctx, TransitionOrderRequest{
		TenantID: payment.TenantID,
		OrderID:  payment.OrderID,
		Status:   target,
	})
	if err != nil {
		return fmt.Errorf("apply payment status to order: %w", err)
	}

	return nil
}

func (s *PaymentService) syncPaymentToOrder(ctx context.Context, payment models.Payment, order models.Order) error {
	if payment.Status == models.PaymentStatusPending {
		return nil
	}

	desiredOrderStatus := func() models.OrderStatus {
		switch payment.Status {
		case models.PaymentStatusSettled:
			return models.OrderStatusPaid
		case models.PaymentStatusFailed:
			return models.OrderStatusFailed
		default:
			return ""
		}
	}()

	if desiredOrderStatus == "" {
		return nil
	}

	if order.Status == desiredOrderStatus {
		return nil
	}

	if order.Status != models.OrderStatusProcessing {
		return nil
	}

	_, err := s.orderWorkflow.TransitionOrder(ctx, TransitionOrderRequest{
		TenantID: payment.TenantID,
		OrderID:  payment.OrderID,
		Status:   desiredOrderStatus,
	})
	if err != nil {
		return fmt.Errorf("sync payment to order: %w", err)
	}

	return nil
}

func normalizePaymentJob(job paymentflow.Job) (paymentflow.Job, error) {
	job.ProviderReference = strings.TrimSpace(job.ProviderReference)

	if job.TenantID <= 0 || job.OrderID <= 0 || job.AmountCents <= 0 || job.ProviderReference == "" {
		return paymentflow.Job{}, ErrInvalidPayment
	}

	return job, nil
}

func samePaymentTarget(payment models.Payment, job paymentflow.Job) bool {
	return payment.TenantID == job.TenantID &&
		payment.OrderID == job.OrderID &&
		payment.ProviderReference == job.ProviderReference &&
		payment.AmountCents == job.AmountCents
}

func isKnownPaymentStatus(status models.PaymentStatus) bool {
	switch status {
	case models.PaymentStatusPending,
		models.PaymentStatusAuthorized,
		models.PaymentStatusSettled,
		models.PaymentStatusFailed,
		models.PaymentStatusRefunded:
		return true
	default:
		return false
	}
}

func normalizeGatewayError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, paymentflow.ErrGatewayTimeout) {
		return paymentflow.ErrGatewayTimeout
	}
	if errors.Is(err, paymentflow.ErrGatewayUnavailable) {
		return paymentflow.ErrGatewayUnavailable
	}
	return err
}
