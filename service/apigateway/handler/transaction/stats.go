package transactionhandler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MamangRust/microservice-payment-gateway-grpc/pb/transaction"
	pbTransactionStats "github.com/MamangRust/microservice-payment-gateway-grpc/pb/stats/transaction"
	"github.com/MamangRust/microservice-payment-gateway-grpc/pkg/logger"
	"github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/apierror"
	"github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/httpx"
	stats_cache "github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/redis/api/stats"
	"github.com/MamangRust/microservice-payment-gateway-grpc/shared/errors"
	transactionapimapper "github.com/MamangRust/microservice-payment-gateway-grpc/shared/mapper/transaction"
	"github.com/go-chi/chi/v5"
)

type transactionStatsHandleApi struct {
	logger logger.LoggerInterface
	cache  stats_cache.StatsCache

	apiHandler apierror.ApiHandler

	// Transaction Clients
	transAmount pbTransactionStats.TransactionStatsAmountServiceClient
	transMethod pbTransactionStats.TransactionStatsMethodServiceClient
	transStatus pbTransactionStats.TransactionStatsStatusServiceClient
}

func setupTransactionStatsHandler(deps *DepsTransaction) func() {
	return func() {
		statsCache := stats_cache.NewStatsCache(deps.Cache)
		h := &transactionStatsHandleApi{
			logger:     deps.Logger,
			cache:      statsCache,
			apiHandler: deps.ApiHandler,

			transAmount: pbTransactionStats.NewTransactionStatsAmountServiceClient(deps.StatsClient),
			transMethod: pbTransactionStats.NewTransactionStatsMethodServiceClient(deps.StatsClient),
			transStatus: pbTransactionStats.NewTransactionStatsStatusServiceClient(deps.StatsClient),
		}

		deps.Router.Route("/api/transaction/stats", func(g chi.Router) {

			// Amount Handlers
			g.Get("/amount/monthly", h.apiHandler.Handle("find-monthly-transaction-amounts", h.FindMonthlyTransactionAmounts))
			g.Get("/amount/yearly", h.apiHandler.Handle("find-yearly-transaction-amounts", h.FindYearlyTransactionAmounts))
			g.Get("/amount/monthly/{card_number}", h.apiHandler.Handle("find-monthly-transaction-amounts-by-card-number", h.FindMonthlyTransactionAmountsByCardNumber))
			g.Get("/amount/yearly/{card_number}", h.apiHandler.Handle("find-yearly-transaction-amounts-by-card-number", h.FindYearlyTransactionAmountsByCardNumber))

			// Method Handlers
			g.Get("/method/monthly", h.apiHandler.Handle("find-monthly-transaction-methods", h.FindMonthlyTransactionMethods))
			g.Get("/method/yearly", h.apiHandler.Handle("find-yearly-transaction-methods", h.FindYearlyTransactionMethods))
			g.Get("/method/monthly/{card_number}", h.apiHandler.Handle("find-monthly-transaction-methods-by-card-number", h.FindMonthlyTransactionMethodsByCardNumber))
			g.Get("/method/yearly/{card_number}", h.apiHandler.Handle("find-yearly-transaction-methods-by-card-number", h.FindYearlyTransactionMethodsByCardNumber))

			// Status Handlers
			g.Get("/status/monthly/success", h.apiHandler.Handle("find-monthly-transaction-status-success", h.FindMonthlyTransactionStatusSuccess))
			g.Get("/status/yearly/success", h.apiHandler.Handle("find-yearly-transaction-status-success", h.FindYearlyTransactionStatusSuccess))
			g.Get("/status/monthly/failed", h.apiHandler.Handle("find-monthly-transaction-status-failed", h.FindMonthlyTransactionStatusFailed))
			g.Get("/status/yearly/failed", h.apiHandler.Handle("find-yearly-transaction-status-failed", h.FindYearlyTransactionStatusFailed))
			g.Get("/status/monthly/success/{card_number}", h.apiHandler.Handle("find-monthly-transaction-status-success-by-card-number", h.FindMonthlyTransactionStatusSuccessByCardNumber))
			g.Get("/status/yearly/success/{card_number}", h.apiHandler.Handle("find-yearly-transaction-status-success-by-card-number", h.FindYearlyTransactionStatusSuccessByCardNumber))
			g.Get("/status/monthly/failed/{card_number}", h.apiHandler.Handle("find-monthly-transaction-status-failed-by-card-number", h.FindMonthlyTransactionStatusFailedByCardNumber))
			g.Get("/status/yearly/failed/{card_number}", h.apiHandler.Handle("find-yearly-transaction-status-failed-by-card-number", h.FindYearlyTransactionStatusFailedByCardNumber))
		})
	}
}

// FindMonthlyTransactionAmounts godoc
// @Summary Get monthly transaction amounts
// @Tags Stats Transaction
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthAmount
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transaction/stats/amount/monthly [get]
func (h *transactionStatsHandleApi) FindMonthlyTransactionAmounts(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transaction:amount:monthly:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transAmount.FindMonthlyAmounts(ctx, &transaction.FindYearTransactionStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transactionapimapper.NewTransactionStatsAmountResponseMapper()
	apiRes := mapper.ToApiResponseTransactionMonthAmount(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTransactionAmounts godoc
// @Summary Get yearly transaction amounts
// @Tags Stats Transaction
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearAmount
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transaction/stats/amount/yearly [get]
func (h *transactionStatsHandleApi) FindYearlyTransactionAmounts(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transaction:amount:yearly:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transAmount.FindYearlyAmounts(ctx, &transaction.FindYearTransactionStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transactionapimapper.NewTransactionStatsAmountResponseMapper()
	apiRes := mapper.ToApiResponseTransactionYearAmount(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTransactionAmountsByCardNumber godoc
// @Summary Get monthly transaction amounts by card number
// @Tags Stats Transaction
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthAmount
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transaction/stats/amount/monthly/{card_number} [get]
func (h *transactionStatsHandleApi) FindMonthlyTransactionAmountsByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transaction:amount:monthly:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transAmount.FindMonthlyAmountsByCardNumber(ctx, &transaction.FindByYearCardNumberTransactionRequest{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transactionapimapper.NewTransactionStatsAmountResponseMapper()
	apiRes := mapper.ToApiResponseTransactionMonthAmount(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTransactionAmountsByCardNumber godoc
// @Summary Get yearly transaction amounts by card number
// @Tags Stats Transaction
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearAmount
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transaction/stats/amount/yearly/{card_number} [get]
func (h *transactionStatsHandleApi) FindYearlyTransactionAmountsByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transaction:amount:yearly:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transAmount.FindYearlyAmountsByCardNumber(ctx, &transaction.FindByYearCardNumberTransactionRequest{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transactionapimapper.NewTransactionStatsAmountResponseMapper()
	apiRes := mapper.ToApiResponseTransactionYearAmount(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTransactionMethods godoc
// @Summary Get monthly transaction methods
// @Tags Stats Transaction
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthMethod
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transaction/stats/method/monthly [get]
func (h *transactionStatsHandleApi) FindMonthlyTransactionMethods(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transaction:method:monthly:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transMethod.FindMonthlyPaymentMethods(ctx, &transaction.FindYearTransactionStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transactionapimapper.NewTransactionStatsMethodResponseMapper()
	apiRes := mapper.ToApiResponseTransactionMonthMethod(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTransactionMethods godoc
// @Summary Get yearly transaction methods
// @Tags Stats Transaction
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearMethod
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transaction/stats/method/yearly [get]
func (h *transactionStatsHandleApi) FindYearlyTransactionMethods(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transaction:method:yearly:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transMethod.FindYearlyPaymentMethods(ctx, &transaction.FindYearTransactionStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transactionapimapper.NewTransactionStatsMethodResponseMapper()
	apiRes := mapper.ToApiResponseTransactionYearMethod(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTransactionMethodsByCardNumber godoc
// @Summary Get monthly transaction methods by card number
// @Tags Stats Transaction
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthMethod
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transaction/stats/method/monthly/{card_number} [get]
func (h *transactionStatsHandleApi) FindMonthlyTransactionMethodsByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transaction:method:monthly:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transMethod.FindMonthlyPaymentMethodsByCardNumber(ctx, &transaction.FindByYearCardNumberTransactionRequest{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transactionapimapper.NewTransactionStatsMethodResponseMapper()
	apiRes := mapper.ToApiResponseTransactionMonthMethod(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTransactionMethodsByCardNumber godoc
// @Summary Get yearly transaction methods by card number
// @Tags Stats Transaction
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearMethod
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transaction/stats/method/yearly/{card_number} [get]
func (h *transactionStatsHandleApi) FindYearlyTransactionMethodsByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transaction:method:yearly:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transMethod.FindYearlyPaymentMethodsByCardNumber(ctx, &transaction.FindByYearCardNumberTransactionRequest{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transactionapimapper.NewTransactionStatsMethodResponseMapper()
	apiRes := mapper.ToApiResponseTransactionYearMethod(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTransactionStatusSuccess godoc
// @Summary Get monthly transaction status success stats
// @Tags Stats Transaction
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transaction/stats/status/monthly/success [get]
func (h *transactionStatsHandleApi) FindMonthlyTransactionStatusSuccess(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transaction:status:monthly:success:%d:%d", year, month)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transStatus.FindMonthlyTransactionStatusSuccess(ctx, &transaction.FindMonthlyTransactionStatus{Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transactionapimapper.NewTransactionStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTransactionMonthStatusSuccess(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTransactionStatusSuccess godoc
// @Summary Get yearly transaction status success stats
// @Tags Stats Transaction
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transaction/stats/status/yearly/success [get]
func (h *transactionStatsHandleApi) FindYearlyTransactionStatusSuccess(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transaction:status:yearly:success:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transStatus.FindYearlyTransactionStatusSuccess(ctx, &transaction.FindYearTransactionStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transactionapimapper.NewTransactionStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTransactionYearStatusSuccess(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTransactionStatusFailed godoc
// @Summary Get monthly transaction status failed stats
// @Tags Stats Transaction
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transaction/stats/status/monthly/failed [get]
func (h *transactionStatsHandleApi) FindMonthlyTransactionStatusFailed(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transaction:status:monthly:failed:%d:%d", year, month)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transStatus.FindMonthlyTransactionStatusFailed(ctx, &transaction.FindMonthlyTransactionStatus{Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transactionapimapper.NewTransactionStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTransactionMonthStatusFailed(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTransactionStatusFailed godoc
// @Summary Get yearly transaction status failed stats
// @Tags Stats Transaction
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transaction/stats/status/yearly/failed [get]
func (h *transactionStatsHandleApi) FindYearlyTransactionStatusFailed(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transaction:status:yearly:failed:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transStatus.FindYearlyTransactionStatusFailed(ctx, &transaction.FindYearTransactionStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transactionapimapper.NewTransactionStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTransactionYearStatusFailed(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTransactionStatusSuccessByCardNumber godoc
// @Summary Get monthly transaction status success stats by card number
// @Tags Stats Transaction
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transaction/stats/status/monthly/success/{card_number} [get]
func (h *transactionStatsHandleApi) FindMonthlyTransactionStatusSuccessByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transaction:status:monthly:success:%s:%d:%d", cardNumber, year, month)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transStatus.FindMonthlyTransactionStatusSuccessByCardNumber(ctx, &transaction.FindMonthlyTransactionStatusCardNumber{Year: int32(year), Month: int32(month), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transactionapimapper.NewTransactionStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTransactionMonthStatusSuccess(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTransactionStatusSuccessByCardNumber godoc
// @Summary Get yearly transaction status success stats by card number
// @Tags Stats Transaction
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transaction/stats/status/yearly/success/{card_number} [get]
func (h *transactionStatsHandleApi) FindYearlyTransactionStatusSuccessByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transaction:status:yearly:success:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transStatus.FindYearlyTransactionStatusSuccessByCardNumber(ctx, &transaction.FindYearTransactionStatusCardNumber{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transactionapimapper.NewTransactionStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTransactionYearStatusSuccess(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTransactionStatusFailedByCardNumber godoc
// @Summary Get monthly transaction status failed stats by card number
// @Tags Stats Transaction
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusFailed
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transaction/stats/status/monthly/failed/{card_number} [get]
func (h *transactionStatsHandleApi) FindMonthlyTransactionStatusFailedByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transaction:status:monthly:failed:%s:%d:%d", cardNumber, year, month)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transStatus.FindMonthlyTransactionStatusFailedByCardNumber(ctx, &transaction.FindMonthlyTransactionStatusCardNumber{Year: int32(year), Month: int32(month), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transactionapimapper.NewTransactionStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTransactionMonthStatusFailed(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTransactionStatusFailedByCardNumber godoc
// @Summary Get yearly transaction status failed stats by card number
// @Tags Stats Transaction
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearStatusFailed
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transaction/stats/status/yearly/failed/{card_number} [get]
func (h *transactionStatsHandleApi) FindYearlyTransactionStatusFailedByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transaction:status:yearly:failed:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transStatus.FindYearlyTransactionStatusFailedByCardNumber(ctx, &transaction.FindYearTransactionStatusCardNumber{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transactionapimapper.NewTransactionStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTransactionYearStatusFailed(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}
