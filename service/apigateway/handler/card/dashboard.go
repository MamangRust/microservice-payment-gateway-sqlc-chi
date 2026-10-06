package cardhandler

import (
	"net/http"

	"github.com/MamangRust/microservice-payment-gateway-grpc/pkg/logger"
	"github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/apierror"
	"github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/httpx"
	card_cache "github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/redis/api/card"
	apimapper "github.com/MamangRust/microservice-payment-gateway-grpc/shared/mapper/card"
	"github.com/go-chi/chi/v5"
)

type cardDashboardHandleApi struct {
	logger logger.LoggerInterface
	mapper apimapper.CardDashboardResponseMapper
	cache  card_cache.CardMencache

	apiHandler apierror.ApiHandler
}

type cardDashboardHandleApiDeps struct {
	router     chi.Router
	logger     logger.LoggerInterface
	mapper     apimapper.CardDashboardResponseMapper
	cache      card_cache.CardMencache
	apiHandler apierror.ApiHandler
}

func NewCardDashboardHandleApi(params *cardDashboardHandleApiDeps) *cardDashboardHandleApi {
	h := &cardDashboardHandleApi{
		logger:     params.logger,
		mapper:     params.mapper,
		cache:      params.cache,
		apiHandler: params.apiHandler,
	}

	// Dashboard endpoints are now handled by StatsHandler's gRPC client to stats-reader.
	// This domain handler exists to maintain apigateway structure and registry compatibility.

	return h
}

// FindDashboard godoc
// @Summary Retrieve card dashboard (analytical)
// @Description This endpoint is deprecated in favor of /api/stats/card/dashboard
// @Tags Card Query
// @Security Bearer
// @Produce json
// @Success 200 {object} response.ApiResponseDashboardCard
// @Router /api/card-query/dashboard [get]
func (h *cardDashboardHandleApi) FindDashboard(w http.ResponseWriter, r *http.Request) error {
	h.apiHandler.Handle("find-dashboard", func(w http.ResponseWriter, r *http.Request) error {
		return httpx.JSON(w, http.StatusMovedPermanently, map[string]string{
			"message": "Please use /api/stats/card/dashboard",
		})
	})(w, r)
	return nil
}
