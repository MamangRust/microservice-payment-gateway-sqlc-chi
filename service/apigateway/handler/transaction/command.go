package transactionhandler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	pbAISecurity "github.com/MamangRust/microservice-payment-gateway-grpc/pb/ai_security"
	pb "github.com/MamangRust/microservice-payment-gateway-grpc/pb/transaction"
	"github.com/MamangRust/microservice-payment-gateway-grpc/pkg/kafka"
	"github.com/MamangRust/microservice-payment-gateway-grpc/pkg/logger"
	"github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/apierror"
	"github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/httpx"
	"github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/middlewares"
	mencache "github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/redis"
	transaction_cache "github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/redis/api/transaction"
	"github.com/MamangRust/microservice-payment-gateway-grpc/shared/domain/requests"
	"github.com/MamangRust/microservice-payment-gateway-grpc/shared/errors"
	"github.com/MamangRust/microservice-payment-gateway-grpc/shared/idempotency"
	apimapper "github.com/MamangRust/microservice-payment-gateway-grpc/shared/mapper/transaction"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type transactionCommandHandleApi struct {
	kafka *kafka.Kafka

	client pb.TransactionCommandServiceClient

	logger logger.LoggerInterface

	mapper apimapper.TransactionCommandResponseMapper

	cache transaction_cache.TransactionMencache

	apiHandler apierror.ApiHandler

	aiSecurity pbAISecurity.AISecurityServiceClient
}

type transactionCommandHandleDeps struct {
	kafka *kafka.Kafka

	client pb.TransactionCommandServiceClient

	router chi.Router

	logger logger.LoggerInterface

	mapper apimapper.TransactionCommandResponseMapper

	cache mencache.MerchantCache

	cache_transaction transaction_cache.TransactionMencache

	apiHandler apierror.ApiHandler

	aiSecurity pbAISecurity.AISecurityServiceClient
}

func NewTransactionCommandHandleApi(params *transactionCommandHandleDeps) *transactionCommandHandleApi {

	transactionCommandHandleApi := &transactionCommandHandleApi{
		kafka:      params.kafka,
		client:     params.client,
		logger:     params.logger,
		mapper:     params.mapper,
		cache:      params.cache_transaction,
		apiHandler: params.apiHandler,
		aiSecurity: params.aiSecurity,
	}

	transactionMiddleware := middlewares.NewApiKeyValidator(params.kafka, "request-transaction", "response-transaction", 5*time.Second, params.logger, params.cache)

	params.router.Route("/api/transaction-command", func(routerTransaction chi.Router) {

		routerTransaction.With(transactionMiddleware.Middleware()).Post("/create", params.apiHandler.Handle("create-transaction", transactionCommandHandleApi.Create))
		routerTransaction.With(transactionMiddleware.Middleware()).Post("/update/{id}", params.apiHandler.Handle("update-transaction", transactionCommandHandleApi.Update))

		routerTransaction.Post("/restore/{id}", params.apiHandler.Handle("restore-transaction", transactionCommandHandleApi.RestoreTransaction))
		routerTransaction.Post("/trashed/{id}", params.apiHandler.Handle("trash-transaction", transactionCommandHandleApi.TrashedTransaction))
		routerTransaction.Delete("/permanent/{id}", params.apiHandler.Handle("delete-transaction-permanent", transactionCommandHandleApi.DeletePermanent))

		routerTransaction.Post("/restore/all", params.apiHandler.Handle("restore-all-transaction", transactionCommandHandleApi.RestoreAllTransaction))
		routerTransaction.Post("/permanent/all", params.apiHandler.Handle("delete-all-transaction-permanent", transactionCommandHandleApi.DeleteAllTransactionPermanent))

	})

	return transactionCommandHandleApi
}

// @Summary Create a new transaction
// @Tags Transaction Command
// @Security Bearer
// @Description Create a new transaction record with the provided details.
// @Accept json
// @Produce json
// @Param CreateTransactionRequest body requests.CreateTransactionRequest true "Create Transaction Request"
// @Success 200 {object} response.ApiResponseTransaction "Successfully created transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create transaction"
// @Router /api/transaction-command/create [post]
func (h *transactionCommandHandleApi) Create(w http.ResponseWriter, r *http.Request) error {
	var body requests.CreateTransactionRequest

	apiKey, _ := httpx.Get(r, "apiKey").(string)

	if apiKey == "" {
		return errors.NewBadRequestError("apiKey is required")
	}

	if err := httpx.Bind(r, &body); err != nil {
		return errors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	key := r.Header.Get("Idempotency-Key")
	if err := idempotency.ValidateKey(key); err != nil {
		return err
	}

	ctx := r.Context()

	// AI Fraud Detection Check
	fraudRes, err := h.aiSecurity.DetectFraud(ctx, &pbAISecurity.FraudRequest{
		TransactionId: strconv.FormatInt(time.Now().UnixNano(), 10), // Temporary ID if not generated yet
		MerchantId:    int32(*body.MerchantID),
		Amount:        float64(body.Amount),
		PaymentMethod: body.PaymentMethod,
		Timestamp:     body.TransactionTime.Unix(),
		IpAddress:     httpx.RealIP(r),
		DeviceId:      r.Header.Get("X-Device-ID"),
	})

	if err == nil && fraudRes.IsFraudulent {
		h.logger.Warn("Transaction rejected by AI Security",
			zap.Float64("risk_score", fraudRes.RiskScore),
			zap.String("reason", fraudRes.Reason),
		)
		return errors.NewForbiddenError(fmt.Sprintf("Transaction rejected for security reasons: %s", fraudRes.Reason))
	}

	res, err := h.client.CreateTransaction(ctx, &pb.CreateTransactionRequest{
		ApiKey:          apiKey,
		CardNumber:      body.CardNumber,
		Amount:          int64(body.Amount),
		PaymentMethod:   body.PaymentMethod,
		MerchantId:      int32(*body.MerchantID),
		TransactionTime: timestamppb.New(body.TransactionTime),
		IdempotencyKey:  key,
	})

	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapper.ToApiResponseTransaction(res)

	h.cache.SetCachedTransactionCache(ctx, so)

	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Update a transaction
// @Tags Transaction Command
// @Security Bearer
// @Description Update an existing transaction record using its ID
// @Accept json
// @Produce json
// @Param transaction body requests.UpdateTransactionRequest true "Transaction data"
// @Success 200 {object} response.ApiResponseTransaction "Updated transaction data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update transaction"
// @Router /api/transaction-command/update [post]
func (h *transactionCommandHandleApi) Update(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	var body requests.UpdateTransactionRequest

	apiKey, ok := httpx.Get(r, "apiKey").(string)
	if !ok {
		return errors.NewBadRequestError("api-key is required")
	}

	if err := httpx.Bind(r, &body); err != nil {
		return errors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := r.Context()

	res, err := h.client.UpdateTransaction(ctx, &pb.UpdateTransactionRequest{
		TransactionId:   int32(id),
		CardNumber:      body.CardNumber,
		ApiKey:          apiKey,
		Amount:          int64(body.Amount),
		PaymentMethod:   body.PaymentMethod,
		MerchantId:      int32(*body.MerchantID),
		TransactionTime: timestamppb.New(body.TransactionTime),
	})

	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapper.ToApiResponseTransaction(res)

	h.cache.DeleteTransactionCache(ctx, id)
	h.cache.SetCachedTransactionCache(ctx, so)

	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Trash a transaction
// @Tags Transaction Command
// @Security Bearer
// @Description Trash a transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransaction "Successfully trashed transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trashed transaction"
// @Router /api/transaction-command/trashed/{id} [post]
func (h *transactionCommandHandleApi) TrashedTransaction(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := r.Context()

	res, err := h.client.TrashedTransaction(ctx, &pb.FindByIdTransactionRequest{
		TransactionId: int32(idInt),
	})

	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapper.ToApiResponseTransactionDeleteAt(res)

	h.cache.DeleteTransactionCache(ctx, idInt)

	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Restore a trashed transaction
// @Tags Transaction Command
// @Security Bearer
// @Description Restore a trashed transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransaction "Successfully restored transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transaction:"
// @Router /api/transaction-command/restore/{id} [post]
func (h *transactionCommandHandleApi) RestoreTransaction(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := r.Context()

	res, err := h.client.RestoreTransaction(ctx, &pb.FindByIdTransactionRequest{
		TransactionId: int32(idInt),
	})

	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapper.ToApiResponseTransactionDeleteAt(res)

	h.cache.DeleteTransactionCache(ctx, idInt)

	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Permanently delete a transaction
// @Tags Transaction Command
// @Security Bearer
// @Description Permanently delete a transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransactionDelete "Successfully deleted transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transaction:"
// @Router /api/transaction-command/permanent/{id} [delete]
func (h *transactionCommandHandleApi) DeletePermanent(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := r.Context()

	res, err := h.client.DeleteTransactionPermanent(ctx, &pb.FindByIdTransactionRequest{
		TransactionId: int32(idInt),
	})

	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapper.ToApiResponseTransactionDelete(res)

	h.cache.DeleteTransactionCache(ctx, idInt)

	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Restore a trashed transaction
// @Tags Transaction Command
// @Security Bearer
// @Description Restore a trashed transaction all.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransactionAll "Successfully restored transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transaction:"
// @Router /api/transaction-command/restore/all [post]
func (h *transactionCommandHandleApi) RestoreAllTransaction(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	res, err := h.client.RestoreAllTransaction(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Failed to restore all transaction", zap.Error(err))
		return errors.ParseGrpcError(err)
	}

	h.logger.Debug("Successfully restored all transaction")

	so := h.mapper.ToApiResponseTransactionAll(res)

	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Permanently delete a transaction
// @Tags Transaction Command
// @Security Bearer
// @Description Permanently delete a transaction all.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransactionAll "Successfully deleted transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transaction:"
// @Router /api/transaction-command/permanent/all [post]
func (h *transactionCommandHandleApi) DeleteAllTransactionPermanent(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	res, err := h.client.DeleteAllTransactionPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Failed to permanently delete all transaction", zap.Error(err))

		return errors.ParseGrpcError(err)
	}

	h.logger.Debug("Successfully deleted all transaction permanently")

	so := h.mapper.ToApiResponseTransactionAll(res)

	return httpx.JSON(w, http.StatusOK, so)
}

func (h *transactionCommandHandleApi) parseValidationErrors(err error) []errors.ValidationError {
	var validationErrs []errors.ValidationError

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			validationErrs = append(validationErrs, errors.ValidationError{
				Field:   fe.Field(),
				Message: h.getValidationMessage(fe),
			})
		}
		return validationErrs
	}

	return []errors.ValidationError{
		{
			Field:   "general",
			Message: err.Error(),
		},
	}
}

func (h *transactionCommandHandleApi) getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Must be at least %s", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s", fe.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", fe.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", fe.Param())
	default:
		return fmt.Sprintf("Validation failed on '%s' tag", fe.Tag())
	}
}
