package withdrawhandler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MamangRust/microservice-payment-gateway-grpc/pb/withdraw"
	pbWithdrawStats "github.com/MamangRust/microservice-payment-gateway-grpc/pb/stats/withdraw"
	"github.com/MamangRust/microservice-payment-gateway-grpc/pkg/logger"
	"github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/apierror"
	"github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/httpx"
	stats_cache "github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/redis/api/stats"
	"github.com/MamangRust/microservice-payment-gateway-grpc/shared/errors"
	withdrawapimapper "github.com/MamangRust/microservice-payment-gateway-grpc/shared/mapper/withdraw"
	"github.com/go-chi/chi/v5"
)

type withdrawStatsHandleApi struct {
	logger logger.LoggerInterface
	cache  stats_cache.StatsCache

	apiHandler apierror.ApiHandler

	// Withdraw Clients
	withdrawAmount pbWithdrawStats.WithdrawStatsAmountServiceClient
	withdrawStatus pbWithdrawStats.WithdrawStatsStatusServiceClient
}

func setupWithdrawStatsHandler(deps *DepsWithdraw) func() {
	return func() {
		statsCache := stats_cache.NewStatsCache(deps.Cache)
		h := &withdrawStatsHandleApi{
			logger:     deps.Logger,
			cache:      statsCache,
			apiHandler: deps.ApiHandler,

			withdrawAmount: pbWithdrawStats.NewWithdrawStatsAmountServiceClient(deps.StatsClient),
			withdrawStatus: pbWithdrawStats.NewWithdrawStatsStatusServiceClient(deps.StatsClient),
		}

		deps.Router.Route("/api/withdraw/stats", func(g chi.Router) {

			// Amount Handlers
			g.Get("/amount/monthly", h.apiHandler.Handle("find-monthly-withdraws", h.FindMonthlyWithdraws))
			g.Get("/amount/yearly", h.apiHandler.Handle("find-yearly-withdraws", h.FindYearlyWithdraws))
			g.Get("/amount/monthly/{card_number}", h.apiHandler.Handle("find-monthly-withdraws-by-card-number", h.FindMonthlyWithdrawsByCardNumber))
			g.Get("/amount/yearly/{card_number}", h.apiHandler.Handle("find-yearly-withdraws-by-card-number", h.FindYearlyWithdrawsByCardNumber))

			// Status Handlers
			g.Get("/status/monthly/success", h.apiHandler.Handle("find-monthly-withdraw-status-success", h.FindMonthlyWithdrawStatusSuccess))
			g.Get("/status/yearly/success", h.apiHandler.Handle("find-yearly-withdraw-status-success", h.FindYearlyWithdrawStatusSuccess))
			g.Get("/status/monthly/failed", h.apiHandler.Handle("find-monthly-withdraw-status-failed", h.FindMonthlyWithdrawStatusFailed))
			g.Get("/status/yearly/failed", h.apiHandler.Handle("find-yearly-withdraw-status-failed", h.FindYearlyWithdrawStatusFailed))
			g.Get("/status/monthly/success/{card_number}", h.apiHandler.Handle("find-monthly-withdraw-status-success-by-card-number", h.FindMonthlyWithdrawStatusSuccessByCardNumber))
			g.Get("/status/yearly/success/{card_number}", h.apiHandler.Handle("find-yearly-withdraw-status-success-by-card-number", h.FindYearlyWithdrawStatusSuccessByCardNumber))
			g.Get("/status/monthly/failed/{card_number}", h.apiHandler.Handle("find-monthly-withdraw-status-failed-by-card-number", h.FindMonthlyWithdrawStatusFailedByCardNumber))
			g.Get("/status/yearly/failed/{card_number}", h.apiHandler.Handle("find-yearly-withdraw-status-failed-by-card-number", h.FindYearlyWithdrawStatusFailedByCardNumber))
		})
	}
}

// FindMonthlyWithdraws godoc
// @Summary Get monthly withdraw amounts
// @Tags Stats Withdraw
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawMonthAmount
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/withdraw/stats/amount/monthly [get]
func (h *withdrawStatsHandleApi) FindMonthlyWithdraws(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:withdraw:amount:monthly:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.withdrawAmount.FindMonthlyWithdraws(ctx, &withdraw.FindYearWithdrawStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := withdrawapimapper.NewWithdrawStatsAmountResponseMapper()
	apiRes := mapper.ToApiResponseWithdrawMonthAmount(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyWithdraws godoc
// @Summary Get yearly withdraw amounts
// @Tags Stats Withdraw
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawYearAmount
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/withdraw/stats/amount/yearly [get]
func (h *withdrawStatsHandleApi) FindYearlyWithdraws(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:withdraw:amount:yearly:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.withdrawAmount.FindYearlyWithdraws(ctx, &withdraw.FindYearWithdrawStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := withdrawapimapper.NewWithdrawStatsAmountResponseMapper()
	apiRes := mapper.ToApiResponseWithdrawYearAmount(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyWithdrawsByCardNumber godoc
// @Summary Get monthly withdraw amounts by card number
// @Tags Stats Withdraw
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawMonthAmount
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/withdraw/stats/amount/monthly/{card_number} [get]
func (h *withdrawStatsHandleApi) FindMonthlyWithdrawsByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:withdraw:amount:monthly:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.withdrawAmount.FindMonthlyWithdrawsByCardNumber(ctx, &withdraw.FindYearWithdrawCardNumber{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := withdrawapimapper.NewWithdrawStatsAmountResponseMapper()
	apiRes := mapper.ToApiResponseWithdrawMonthAmount(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyWithdrawsByCardNumber godoc
// @Summary Get yearly withdraw amounts by card number
// @Tags Stats Withdraw
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawYearAmount
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/withdraw/stats/amount/yearly/{card_number} [get]
func (h *withdrawStatsHandleApi) FindYearlyWithdrawsByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:withdraw:amount:yearly:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.withdrawAmount.FindYearlyWithdrawsByCardNumber(ctx, &withdraw.FindYearWithdrawCardNumber{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := withdrawapimapper.NewWithdrawStatsAmountResponseMapper()
	apiRes := mapper.ToApiResponseWithdrawYearAmount(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyWithdrawStatusSuccess godoc
// @Summary Get monthly withdraw status success stats
// @Tags Stats Withdraw
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseWithdrawMonthStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/withdraw/stats/status/monthly/success [get]
func (h *withdrawStatsHandleApi) FindMonthlyWithdrawStatusSuccess(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:withdraw:status:monthly:success:%d:%d", year, month)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.withdrawStatus.FindMonthlyWithdrawStatusSuccess(ctx, &withdraw.FindMonthlyWithdrawStatus{Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := withdrawapimapper.NewWithdrawStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseWithdrawMonthStatusSuccess(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyWithdrawStatusSuccess godoc
// @Summary Get yearly withdraw status success stats
// @Tags Stats Withdraw
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawYearStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/withdraw/stats/status/yearly/success [get]
func (h *withdrawStatsHandleApi) FindYearlyWithdrawStatusSuccess(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:withdraw:status:yearly:success:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.withdrawStatus.FindYearlyWithdrawStatusSuccess(ctx, &withdraw.FindYearWithdrawStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := withdrawapimapper.NewWithdrawStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseWithdrawYearStatusSuccess(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyWithdrawStatusFailed godoc
// @Summary Get monthly withdraw status failed stats
// @Tags Stats Withdraw
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseWithdrawMonthStatusFailed
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/withdraw/stats/status/monthly/failed [get]
func (h *withdrawStatsHandleApi) FindMonthlyWithdrawStatusFailed(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:withdraw:status:monthly:failed:%d:%d", year, month)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.withdrawStatus.FindMonthlyWithdrawStatusFailed(ctx, &withdraw.FindMonthlyWithdrawStatus{Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := withdrawapimapper.NewWithdrawStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseWithdrawMonthStatusFailed(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyWithdrawStatusFailed godoc
// @Summary Get yearly withdraw status failed stats
// @Tags Stats Withdraw
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawYearStatusFailed
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/withdraw/stats/status/yearly/failed [get]
func (h *withdrawStatsHandleApi) FindYearlyWithdrawStatusFailed(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:withdraw:status:yearly:failed:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.withdrawStatus.FindYearlyWithdrawStatusFailed(ctx, &withdraw.FindYearWithdrawStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := withdrawapimapper.NewWithdrawStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseWithdrawYearStatusFailed(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyWithdrawStatusSuccessByCardNumber godoc
// @Summary Get monthly withdraw status success stats by card number
// @Tags Stats Withdraw
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseWithdrawMonthStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/withdraw/stats/status/monthly/success/{card_number} [get]
func (h *withdrawStatsHandleApi) FindMonthlyWithdrawStatusSuccessByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:withdraw:status:monthly:success:%s:%d:%d", cardNumber, year, month)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.withdrawStatus.FindMonthlyWithdrawStatusSuccessCardNumber(ctx, &withdraw.FindMonthlyWithdrawStatusCardNumber{Year: int32(year), Month: int32(month), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := withdrawapimapper.NewWithdrawStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseWithdrawMonthStatusSuccess(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyWithdrawStatusSuccessByCardNumber godoc
// @Summary Get yearly withdraw status success stats by card number
// @Tags Stats Withdraw
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawYearStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/withdraw/stats/status/yearly/success/{card_number} [get]
func (h *withdrawStatsHandleApi) FindYearlyWithdrawStatusSuccessByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:withdraw:status:yearly:success:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.withdrawStatus.FindYearlyWithdrawStatusSuccessCardNumber(ctx, &withdraw.FindYearWithdrawStatusCardNumber{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := withdrawapimapper.NewWithdrawStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseWithdrawYearStatusSuccess(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyWithdrawStatusFailedByCardNumber godoc
// @Summary Get monthly withdraw status failed stats by card number
// @Tags Stats Withdraw
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseWithdrawMonthStatusFailed
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/withdraw/stats/status/monthly/failed/{card_number} [get]
func (h *withdrawStatsHandleApi) FindMonthlyWithdrawStatusFailedByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:withdraw:status:monthly:failed:%s:%d:%d", cardNumber, year, month)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.withdrawStatus.FindMonthlyWithdrawStatusFailedCardNumber(ctx, &withdraw.FindMonthlyWithdrawStatusCardNumber{Year: int32(year), Month: int32(month), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := withdrawapimapper.NewWithdrawStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseWithdrawMonthStatusFailed(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyWithdrawStatusFailedByCardNumber godoc
// @Summary Get yearly withdraw status failed stats by card number
// @Tags Stats Withdraw
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawYearStatusFailed
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/withdraw/stats/status/yearly/failed/{card_number} [get]
func (h *withdrawStatsHandleApi) FindYearlyWithdrawStatusFailedByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:withdraw:status:yearly:failed:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.withdrawStatus.FindYearlyWithdrawStatusFailedCardNumber(ctx, &withdraw.FindYearWithdrawStatusCardNumber{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := withdrawapimapper.NewWithdrawStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseWithdrawYearStatusFailed(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}
