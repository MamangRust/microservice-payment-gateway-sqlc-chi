package authhandler

import (
	pb "github.com/MamangRust/microservice-payment-gateway-grpc/pb"
	"github.com/MamangRust/microservice-payment-gateway-grpc/pkg/logger"
	"github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/apierror"
	auth_cache "github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/redis/api/auth"
	"github.com/MamangRust/microservice-payment-gateway-grpc/shared/cache"
	authapimapper "github.com/MamangRust/microservice-payment-gateway-grpc/shared/mapper/auth"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type DepsAuth struct {
	Client     *grpc.ClientConn
	Router     chi.Router
	Logger     logger.LoggerInterface
	Cache      *cache.CacheStore
	ApiHandler apierror.ApiHandler
}

func RegisterAuthHandler(deps *DepsAuth) {
	mapper := authapimapper.NewAuthResponseMapper()

	cache := auth_cache.NewMencache(deps.Cache)

	NewHandlerAuth(&authHandleParams{
		client:     pb.NewAuthServiceClient(deps.Client),
		router:     deps.Router,
		logger:     deps.Logger,
		mapper:     mapper,
		cache:      cache,
		apiHandler: deps.ApiHandler,
	})
}
