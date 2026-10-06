package transferhandler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MamangRust/microservice-payment-gateway-grpc/pb/transfer"
	pbTransferStats "github.com/MamangRust/microservice-payment-gateway-grpc/pb/stats/transfer"
	"github.com/MamangRust/microservice-payment-gateway-grpc/pkg/logger"
	"github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/apierror"
	"github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/httpx"
	stats_cache "github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/redis/api/stats"
	"github.com/MamangRust/microservice-payment-gateway-grpc/shared/errors"
	transferapimapper "github.com/MamangRust/microservice-payment-gateway-grpc/shared/mapper/transfer"
	"github.com/go-chi/chi/v5"
)

type transferStatsHandleApi struct {
	logger logger.LoggerInterface
	cache  stats_cache.StatsCache

	apiHandler apierror.ApiHandler

	// Transfer Clients
	transferAmount pbTransferStats.TransferStatsAmountServiceClient
	transferStatus pbTransferStats.TransferStatsStatusServiceClient
}

func setupTransferStatsHandler(deps *DepsTransfer) func() {
	return func() {
		statsCache := stats_cache.NewStatsCache(deps.Cache)
		h := &transferStatsHandleApi{
			logger:     deps.Logger,
			cache:      statsCache,
			apiHandler: deps.ApiHandler,

			transferAmount: pbTransferStats.NewTransferStatsAmountServiceClient(deps.StatsClient),
			transferStatus: pbTransferStats.NewTransferStatsStatusServiceClient(deps.StatsClient),
		}

		deps.Router.Route("/api/transfer/stats", func(g chi.Router) {

			// Amount Handlers
			g.Get("/amount/monthly", h.apiHandler.Handle("find-monthly-transfer-amounts", h.FindMonthlyTransferAmounts))
			g.Get("/amount/yearly", h.apiHandler.Handle("find-yearly-transfer-amounts", h.FindYearlyTransferAmounts))
			g.Get("/amount/monthly/sender/{card_number}", h.apiHandler.Handle("find-monthly-transfer-amounts-by-sender-card-number", h.FindMonthlyTransferAmountsBySenderCardNumber))
			g.Get("/amount/monthly/receiver/{card_number}", h.apiHandler.Handle("find-monthly-transfer-amounts-by-receiver-card-number", h.FindMonthlyTransferAmountsByReceiverCardNumber))
			g.Get("/amount/yearly/sender/{card_number}", h.apiHandler.Handle("find-yearly-transfer-amounts-by-sender-card-number", h.FindYearlyTransferAmountsBySenderCardNumber))
			g.Get("/amount/yearly/receiver/{card_number}", h.apiHandler.Handle("find-yearly-transfer-amounts-by-receiver-card-number", h.FindYearlyTransferAmountsByReceiverCardNumber))

			// Status Handlers
			g.Get("/status/monthly/success", h.apiHandler.Handle("find-monthly-transfer-status-success", h.FindMonthlyTransferStatusSuccess))
			g.Get("/status/yearly/success", h.apiHandler.Handle("find-yearly-transfer-status-success", h.FindYearlyTransferStatusSuccess))
			g.Get("/status/monthly/failed", h.apiHandler.Handle("find-monthly-transfer-status-failed", h.FindMonthlyTransferStatusFailed))
			g.Get("/status/yearly/failed", h.apiHandler.Handle("find-yearly-transfer-status-failed", h.FindYearlyTransferStatusFailed))
			g.Get("/status/monthly/success/{card_number}", h.apiHandler.Handle("find-monthly-transfer-status-success-by-card-number", h.FindMonthlyTransferStatusSuccessByCardNumber))
			g.Get("/status/yearly/success/{card_number}", h.apiHandler.Handle("find-yearly-transfer-status-success-by-card-number", h.FindYearlyTransferStatusSuccessByCardNumber))
			g.Get("/status/monthly/failed/{card_number}", h.apiHandler.Handle("find-monthly-transfer-status-failed-by-card-number", h.FindMonthlyTransferStatusFailedByCardNumber))
			g.Get("/status/yearly/failed/{card_number}", h.apiHandler.Handle("find-yearly-transfer-status-failed-by-card-number", h.FindYearlyTransferStatusFailedByCardNumber))
		})
	}
}

// FindMonthlyTransferAmounts godoc
// @Summary Get monthly transfer amounts
// @Tags Stats Transfer
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferMonthAmount
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transfer/stats/amount/monthly [get]
func (h *transferStatsHandleApi) FindMonthlyTransferAmounts(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transfer:amount:monthly:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transferAmount.FindMonthlyTransferAmounts(ctx, &transfer.FindYearTransferStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transferapimapper.NewTransferStatsAmountResponseMapper()
	apiRes := mapper.ToApiResponseTransferMonthAmount(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTransferAmounts godoc
// @Summary Get yearly transfer amounts
// @Tags Stats Transfer
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferYearAmount
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transfer/stats/amount/yearly [get]
func (h *transferStatsHandleApi) FindYearlyTransferAmounts(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transfer:amount:yearly:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transferAmount.FindYearlyTransferAmounts(ctx, &transfer.FindYearTransferStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transferapimapper.NewTransferStatsAmountResponseMapper()
	apiRes := mapper.ToApiResponseTransferYearAmount(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTransferAmountsBySenderCardNumber godoc
// @Summary Get monthly transfer amounts by sender card number
// @Tags Stats Transfer
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferMonthAmount
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transfer/stats/amount/monthly/sender/{card_number} [get]
func (h *transferStatsHandleApi) FindMonthlyTransferAmountsBySenderCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transfer:amount:monthly:sender:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transferAmount.FindMonthlyTransferAmountsBySenderCardNumber(ctx, &transfer.FindByCardNumberTransferRequest{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transferapimapper.NewTransferStatsAmountResponseMapper()
	apiRes := mapper.ToApiResponseTransferMonthAmount(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTransferAmountsByReceiverCardNumber godoc
// @Summary Get monthly transfer amounts by receiver card number
// @Tags Stats Transfer
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferMonthAmount
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transfer/stats/amount/monthly/receiver/{card_number} [get]
func (h *transferStatsHandleApi) FindMonthlyTransferAmountsByReceiverCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transfer:amount:monthly:receiver:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transferAmount.FindMonthlyTransferAmountsByReceiverCardNumber(ctx, &transfer.FindByCardNumberTransferRequest{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transferapimapper.NewTransferStatsAmountResponseMapper()
	apiRes := mapper.ToApiResponseTransferMonthAmount(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTransferAmountsBySenderCardNumber godoc
// @Summary Get yearly transfer amounts by sender card number
// @Tags Stats Transfer
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferYearAmount
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transfer/stats/amount/yearly/sender/{card_number} [get]
func (h *transferStatsHandleApi) FindYearlyTransferAmountsBySenderCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transfer:amount:yearly:sender:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transferAmount.FindYearlyTransferAmountsBySenderCardNumber(ctx, &transfer.FindByCardNumberTransferRequest{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transferapimapper.NewTransferStatsAmountResponseMapper()
	apiRes := mapper.ToApiResponseTransferYearAmount(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTransferAmountsByReceiverCardNumber godoc
// @Summary Get yearly transfer amounts by receiver card number
// @Tags Stats Transfer
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferYearAmount
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transfer/stats/amount/yearly/receiver/{card_number} [get]
func (h *transferStatsHandleApi) FindYearlyTransferAmountsByReceiverCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transfer:amount:yearly:receiver:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transferAmount.FindYearlyTransferAmountsByReceiverCardNumber(ctx, &transfer.FindByCardNumberTransferRequest{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transferapimapper.NewTransferStatsAmountResponseMapper()
	apiRes := mapper.ToApiResponseTransferYearAmount(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTransferStatusSuccess godoc
// @Summary Get monthly transfer status success stats
// @Tags Stats Transfer
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransferMonthStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transfer/stats/status/monthly/success [get]
func (h *transferStatsHandleApi) FindMonthlyTransferStatusSuccess(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transfer:status:monthly:success:%d:%d", year, month)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transferStatus.FindMonthlyTransferStatusSuccess(ctx, &transfer.FindMonthlyTransferStatus{Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transferapimapper.NewTransferStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTransferMonthStatusSuccess(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTransferStatusSuccess godoc
// @Summary Get yearly transfer status success stats
// @Tags Stats Transfer
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferYearStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transfer/stats/status/yearly/success [get]
func (h *transferStatsHandleApi) FindYearlyTransferStatusSuccess(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transfer:status:yearly:success:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transferStatus.FindYearlyTransferStatusSuccess(ctx, &transfer.FindYearTransferStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transferapimapper.NewTransferStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTransferYearStatusSuccess(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTransferStatusFailed godoc
// @Summary Get monthly transfer status failed stats
// @Tags Stats Transfer
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransferMonthStatusFailed
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transfer/stats/status/monthly/failed [get]
func (h *transferStatsHandleApi) FindMonthlyTransferStatusFailed(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transfer:status:monthly:failed:%d:%d", year, month)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transferStatus.FindMonthlyTransferStatusFailed(ctx, &transfer.FindMonthlyTransferStatus{Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transferapimapper.NewTransferStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTransferMonthStatusFailed(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTransferStatusFailed godoc
// @Summary Get yearly transfer status failed stats
// @Tags Stats Transfer
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferYearStatusFailed
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transfer/stats/status/yearly/failed [get]
func (h *transferStatsHandleApi) FindYearlyTransferStatusFailed(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transfer:status:yearly:failed:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transferStatus.FindYearlyTransferStatusFailed(ctx, &transfer.FindYearTransferStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transferapimapper.NewTransferStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTransferYearStatusFailed(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTransferStatusSuccessByCardNumber godoc
// @Summary Get monthly transfer status success stats by card number
// @Tags Stats Transfer
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransferMonthStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transfer/stats/status/monthly/success/{card_number} [get]
func (h *transferStatsHandleApi) FindMonthlyTransferStatusSuccessByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transfer:status:monthly:success:%s:%d:%d", cardNumber, year, month)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transferStatus.FindMonthlyTransferStatusSuccessByCardNumber(ctx, &transfer.FindMonthlyTransferStatusCardNumber{Year: int32(year), Month: int32(month), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transferapimapper.NewTransferStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTransferMonthStatusSuccess(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTransferStatusSuccessByCardNumber godoc
// @Summary Get yearly transfer status success stats by card number
// @Tags Stats Transfer
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferYearStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transfer/stats/status/yearly/success/{card_number} [get]
func (h *transferStatsHandleApi) FindYearlyTransferStatusSuccessByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transfer:status:yearly:success:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transferStatus.FindYearlyTransferStatusSuccessByCardNumber(ctx, &transfer.FindYearTransferStatusCardNumber{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transferapimapper.NewTransferStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTransferYearStatusSuccess(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTransferStatusFailedByCardNumber godoc
// @Summary Get monthly transfer status failed stats by card number
// @Tags Stats Transfer
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransferMonthStatusFailed
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transfer/stats/status/monthly/failed/{card_number} [get]
func (h *transferStatsHandleApi) FindMonthlyTransferStatusFailedByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transfer:status:monthly:failed:%s:%d:%d", cardNumber, year, month)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transferStatus.FindMonthlyTransferStatusFailedByCardNumber(ctx, &transfer.FindMonthlyTransferStatusCardNumber{Year: int32(year), Month: int32(month), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transferapimapper.NewTransferStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTransferMonthStatusFailed(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTransferStatusFailedByCardNumber godoc
// @Summary Get yearly transfer status failed stats by card number
// @Tags Stats Transfer
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferYearStatusFailed
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/transfer/stats/status/yearly/failed/{card_number} [get]
func (h *transferStatsHandleApi) FindYearlyTransferStatusFailedByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:transfer:status:yearly:failed:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.transferStatus.FindYearlyTransferStatusFailedByCardNumber(ctx, &transfer.FindYearTransferStatusCardNumber{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := transferapimapper.NewTransferStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTransferYearStatusFailed(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}
