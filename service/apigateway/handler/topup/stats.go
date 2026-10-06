package topuphandler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MamangRust/microservice-payment-gateway-grpc/pb/topup"
	pbTopupStats "github.com/MamangRust/microservice-payment-gateway-grpc/pb/stats/topup"
	"github.com/MamangRust/microservice-payment-gateway-grpc/pkg/logger"
	stats_cache "github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/redis/api/stats"
	"github.com/MamangRust/microservice-payment-gateway-grpc/shared/errors"
	topupapimapper "github.com/MamangRust/microservice-payment-gateway-grpc/shared/mapper/topup"

	"github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/apierror"
	"github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/httpx"
	"github.com/go-chi/chi/v5"
)

type topupStatsHandleApi struct {
	logger logger.LoggerInterface
	cache  stats_cache.StatsCache

	apiHandler apierror.ApiHandler

	// Topup Clients
	topupAmount pbTopupStats.TopupStatsAmountServiceClient
	topupMethod pbTopupStats.TopupStatsMethodServiceClient
	topupStatus pbTopupStats.TopupStatsStatusServiceClient
}

func setupTopupStatsHandler(deps *DepsTopup) func() {
	return func() {
		statsCache := stats_cache.NewStatsCache(deps.Cache)
		h := &topupStatsHandleApi{
			logger:     deps.Logger,
			cache:      statsCache,
			apiHandler: deps.ApiHandler,

			topupAmount: pbTopupStats.NewTopupStatsAmountServiceClient(deps.StatsClient),
			topupMethod: pbTopupStats.NewTopupStatsMethodServiceClient(deps.StatsClient),
			topupStatus: pbTopupStats.NewTopupStatsStatusServiceClient(deps.StatsClient),
		}

		deps.Router.Route("/api/topup/stats", func(g chi.Router) {

			// Amount Handlers
			g.Get("/amount/monthly", h.apiHandler.Handle("find-monthly-topup-amounts", h.FindMonthlyTopupAmounts))
			g.Get("/amount/yearly", h.apiHandler.Handle("find-yearly-topup-amounts", h.FindYearlyTopupAmounts))
			g.Get("/amount/monthly/{card_number}", h.apiHandler.Handle("find-monthly-topup-amounts-by-card-number", h.FindMonthlyTopupAmountsByCardNumber))
			g.Get("/amount/yearly/{card_number}", h.apiHandler.Handle("find-yearly-topup-amounts-by-card-number", h.FindYearlyTopupAmountsByCardNumber))

			// Method Handlers
			g.Get("/method/monthly", h.apiHandler.Handle("find-monthly-topup-methods", h.FindMonthlyTopupMethods))
			g.Get("/method/yearly", h.apiHandler.Handle("find-yearly-topup-methods", h.FindYearlyTopupMethods))
			g.Get("/method/monthly/{card_number}", h.apiHandler.Handle("find-monthly-topup-methods-by-card-number", h.FindMonthlyTopupMethodsByCardNumber))
			g.Get("/method/yearly/{card_number}", h.apiHandler.Handle("find-yearly-topup-methods-by-card-number", h.FindYearlyTopupMethodsByCardNumber))

			// Status Handlers
			g.Get("/status/monthly/success", h.apiHandler.Handle("find-monthly-topup-status-success", h.FindMonthlyTopupStatusSuccess))
			g.Get("/status/yearly/success", h.apiHandler.Handle("find-yearly-topup-status-success", h.FindYearlyTopupStatusSuccess))
			g.Get("/status/monthly/failed", h.apiHandler.Handle("find-monthly-topup-status-failed", h.FindMonthlyTopupStatusFailed))
			g.Get("/status/yearly/failed", h.apiHandler.Handle("find-yearly-topup-status-failed", h.FindYearlyTopupStatusFailed))
			g.Get("/status/monthly/success/{card_number}", h.apiHandler.Handle("find-monthly-topup-status-success-by-card-number", h.FindMonthlyTopupStatusSuccessByCardNumber))
			g.Get("/status/yearly/success/{card_number}", h.apiHandler.Handle("find-yearly-topup-status-success-by-card-number", h.FindYearlyTopupStatusSuccessByCardNumber))
			g.Get("/status/monthly/failed/{card_number}", h.apiHandler.Handle("find-monthly-topup-status-failed-by-card-number", h.FindMonthlyTopupStatusFailedByCardNumber))
			g.Get("/status/yearly/failed/{card_number}", h.apiHandler.Handle("find-yearly-topup-status-failed-by-card-number", h.FindYearlyTopupStatusFailedByCardNumber))
		})
	}
}

// FindMonthlyTopupAmounts godoc
// @Summary Get monthly topup amounts
// @Tags Stats Topup
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthAmount
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/topup/stats/amount/monthly [get]
func (h *topupStatsHandleApi) FindMonthlyTopupAmounts(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:topup:amount:monthly:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.topupAmount.FindMonthlyTopupAmounts(ctx, &topup.FindYearTopupStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := topupapimapper.NewTopupStatsAmountResponseMapper()
	apiRes := mapper.ToApiResponseTopupMonthAmount(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTopupAmounts godoc
// @Summary Get yearly topup amounts
// @Tags Stats Topup
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearAmount
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/topup/stats/amount/yearly [get]
func (h *topupStatsHandleApi) FindYearlyTopupAmounts(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:topup:amount:yearly:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.topupAmount.FindYearlyTopupAmounts(ctx, &topup.FindYearTopupStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := topupapimapper.NewTopupStatsAmountResponseMapper()
	apiRes := mapper.ToApiResponseTopupYearAmount(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTopupAmountsByCardNumber godoc
// @Summary Get monthly topup amounts by card number
// @Tags Stats Topup
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthAmount
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/topup/stats/amount/monthly/{card_number} [get]
func (h *topupStatsHandleApi) FindMonthlyTopupAmountsByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:topup:amount:monthly:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.topupAmount.FindMonthlyTopupAmountsByCardNumber(ctx, &topup.FindYearTopupCardNumber{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := topupapimapper.NewTopupStatsAmountResponseMapper()
	apiRes := mapper.ToApiResponseTopupMonthAmount(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTopupAmountsByCardNumber godoc
// @Summary Get yearly topup amounts by card number
// @Tags Stats Topup
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearAmount
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/topup/stats/amount/yearly/{card_number} [get]
func (h *topupStatsHandleApi) FindYearlyTopupAmountsByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:topup:amount:yearly:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.topupAmount.FindYearlyTopupAmountsByCardNumber(ctx, &topup.FindYearTopupCardNumber{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := topupapimapper.NewTopupStatsAmountResponseMapper()
	apiRes := mapper.ToApiResponseTopupYearAmount(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTopupMethods godoc
// @Summary Get monthly topup methods
// @Tags Stats Topup
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthMethod
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/topup/stats/method/monthly [get]
func (h *topupStatsHandleApi) FindMonthlyTopupMethods(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:topup:method:monthly:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.topupMethod.FindMonthlyTopupMethods(ctx, &topup.FindYearTopupStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := topupapimapper.NewTopupStatsMethodResponseMapper()
	apiRes := mapper.ToApiResponseTopupMonthMethod(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTopupMethods godoc
// @Summary Get yearly topup methods
// @Tags Stats Topup
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearMethod
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/topup/stats/method/yearly [get]
func (h *topupStatsHandleApi) FindYearlyTopupMethods(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:topup:method:yearly:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.topupMethod.FindYearlyTopupMethods(ctx, &topup.FindYearTopupStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := topupapimapper.NewTopupStatsMethodResponseMapper()
	apiRes := mapper.ToApiResponseTopupYearMethod(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTopupMethodsByCardNumber godoc
// @Summary Get monthly topup methods by card number
// @Tags Stats Topup
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthMethod
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/topup/stats/method/monthly/{card_number} [get]
func (h *topupStatsHandleApi) FindMonthlyTopupMethodsByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:topup:method:monthly:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.topupMethod.FindMonthlyTopupMethodsByCardNumber(ctx, &topup.FindYearTopupCardNumber{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := topupapimapper.NewTopupStatsMethodResponseMapper()
	apiRes := mapper.ToApiResponseTopupMonthMethod(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTopupMethodsByCardNumber godoc
// @Summary Get yearly topup methods by card number
// @Tags Stats Topup
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthMethod
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/topup/stats/method/yearly/{card_number} [get]
func (h *topupStatsHandleApi) FindYearlyTopupMethodsByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:topup:method:yearly:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.topupMethod.FindYearlyTopupMethodsByCardNumber(ctx, &topup.FindYearTopupCardNumber{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := topupapimapper.NewTopupStatsMethodResponseMapper()
	apiRes := mapper.ToApiResponseTopupYearMethod(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTopupStatusSuccess godoc
// @Summary Get monthly topup status success stats
// @Tags Stats Topup
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTopupMonthStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/topup/stats/status/monthly/success [get]
func (h *topupStatsHandleApi) FindMonthlyTopupStatusSuccess(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:topup:status:monthly:success:%d:%d", year, month)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.topupStatus.FindMonthlyTopupStatusSuccess(ctx, &topup.FindMonthlyTopupStatus{Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := topupapimapper.NewTopupStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTopupMonthStatusSuccess(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTopupStatusSuccess godoc
// @Summary Get yearly topup status success stats
// @Tags Stats Topup
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/topup/stats/status/yearly/success [get]
func (h *topupStatsHandleApi) FindYearlyTopupStatusSuccess(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:topup:status:yearly:success:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.topupStatus.FindYearlyTopupStatusSuccess(ctx, &topup.FindYearTopupStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := topupapimapper.NewTopupStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTopupYearStatusSuccess(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTopupStatusFailed godoc
// @Summary Get monthly topup status failed stats
// @Tags Stats Topup
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTopupMonthStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/topup/stats/status/monthly/failed [get]
func (h *topupStatsHandleApi) FindMonthlyTopupStatusFailed(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:topup:status:monthly:failed:%d:%d", year, month)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.topupStatus.FindMonthlyTopupStatusFailed(ctx, &topup.FindMonthlyTopupStatus{Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := topupapimapper.NewTopupStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTopupMonthStatusFailed(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTopupStatusFailed godoc
// @Summary Get yearly topup status failed stats
// @Tags Stats Topup
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/topup/stats/status/yearly/failed [get]
func (h *topupStatsHandleApi) FindYearlyTopupStatusFailed(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:topup:status:yearly:failed:%d", year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.topupStatus.FindYearlyTopupStatusFailed(ctx, &topup.FindYearTopupStatus{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := topupapimapper.NewTopupStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTopupYearStatusFailed(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTopupStatusSuccessByCardNumber godoc
// @Summary Get monthly topup status success stats by card number
// @Tags Stats Topup
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTopupMonthStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/topup/stats/status/monthly/success/{card_number} [get]
func (h *topupStatsHandleApi) FindMonthlyTopupStatusSuccessByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:topup:status:monthly:success:%s:%d:%d", cardNumber, year, month)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.topupStatus.FindMonthlyTopupStatusSuccessByCardNumber(ctx, &topup.FindMonthlyTopupStatusCardNumber{Year: int32(year), Month: int32(month), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := topupapimapper.NewTopupStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTopupMonthStatusSuccess(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTopupStatusSuccessByCardNumber godoc
// @Summary Get yearly topup status success stats by card number
// @Tags Stats Topup
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthStatusSuccess
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/topup/stats/status/yearly/success/{card_number} [get]
func (h *topupStatsHandleApi) FindYearlyTopupStatusSuccessByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:topup:status:yearly:success:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.topupStatus.FindYearlyTopupStatusSuccessByCardNumber(ctx, &topup.FindYearTopupStatusCardNumber{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := topupapimapper.NewTopupStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTopupYearStatusSuccess(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindMonthlyTopupStatusFailedByCardNumber godoc
// @Summary Get monthly topup status failed stats by card number
// @Tags Stats Topup
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTopupMonthStatusFailed
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/topup/stats/status/monthly/failed/{card_number} [get]
func (h *topupStatsHandleApi) FindMonthlyTopupStatusFailedByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:topup:status:monthly:failed:%s:%d:%d", cardNumber, year, month)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.topupStatus.FindMonthlyTopupStatusFailedByCardNumber(ctx, &topup.FindMonthlyTopupStatusCardNumber{Year: int32(year), Month: int32(month), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := topupapimapper.NewTopupStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTopupMonthStatusFailed(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}

// FindYearlyTopupStatusFailedByCardNumber godoc
// @Summary Get yearly topup status failed stats by card number
// @Tags Stats Topup
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearStatusFailed
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/topup/stats/status/yearly/failed/{card_number} [get]
func (h *topupStatsHandleApi) FindYearlyTopupStatusFailedByCardNumber(w http.ResponseWriter, r *http.Request) error {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	cardNumber := chi.URLParam(r, "card_number")
	ctx := r.Context()
	cacheKey := fmt.Sprintf("stats:topup:status:yearly:failed:%s:%d", cardNumber, year)

	if cached, found := h.cache.GetCache(ctx, cacheKey); found {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	resp, err := h.topupStatus.FindYearlyTopupStatusFailedByCardNumber(ctx, &topup.FindYearTopupStatusCardNumber{Year: int32(year), CardNumber: cardNumber})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	mapper := topupapimapper.NewTopupStatsStatusResponseMapper()
	apiRes := mapper.ToApiResponseTopupYearStatusFailed(resp)

	h.cache.SetCache(ctx, cacheKey, apiRes)
	return httpx.JSON(w, http.StatusOK, apiRes)
}
