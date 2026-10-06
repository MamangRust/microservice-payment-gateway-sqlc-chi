package merchanthandler

import (
	pb "github.com/MamangRust/microservice-payment-gateway-grpc/pb/merchant"
	"github.com/MamangRust/microservice-payment-gateway-grpc/pkg/logger"
	"github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/apierror"
	merchant_cache "github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/redis/api/merchant"
	apimapper "github.com/MamangRust/microservice-payment-gateway-grpc/shared/mapper/merchant"
	"github.com/go-chi/chi/v5"
)

type merchantTransactionHandleApi struct {
	client pb.MerchantTransactionServiceClient

	logger logger.LoggerInterface
	mapper apimapper.MerchantTransactionResponseMapper

	cache merchant_cache.MerchantMencache

	apiHandler apierror.ApiHandler
}

type merchantTransactionHandleDeps struct {
	client pb.MerchantTransactionServiceClient
	router chi.Router

	logger logger.LoggerInterface
	mapper apimapper.MerchantTransactionResponseMapper

	cache merchant_cache.MerchantMencache

	apiHandler apierror.ApiHandler
}

func NewMerchantTransactionHandleApi(params *merchantTransactionHandleDeps) *merchantTransactionHandleApi {

	merchantHandler := &merchantTransactionHandleApi{
		client:     params.client,
		logger:     params.logger,
		mapper:     params.mapper,
		cache:      params.cache,
		apiHandler: params.apiHandler,
	}

	// routerMerchant := params.router.Route("/api/merchant-transaction", func(routerMerchant chi.Router) {})

	// Transaction analytical endpoints are handled by StatsReader
	// This handler is for domain-specific transaction operations if needed

	return merchantHandler
}
