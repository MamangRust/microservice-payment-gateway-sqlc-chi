package transactionhandler

import (
	pbAISecurity "github.com/MamangRust/microservice-payment-gateway-grpc/pb/ai_security"
	pb "github.com/MamangRust/microservice-payment-gateway-grpc/pb/transaction"
	"github.com/MamangRust/microservice-payment-gateway-grpc/pkg/kafka"
	"github.com/MamangRust/microservice-payment-gateway-grpc/pkg/logger"
	"github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/apierror"
	mencache "github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/redis"
	transaction_cache "github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/redis/api/transaction"
	"github.com/MamangRust/microservice-payment-gateway-grpc/shared/cache"
	apimapper "github.com/MamangRust/microservice-payment-gateway-grpc/shared/mapper/transaction"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type DepsTransaction struct {
	Client *grpc.ClientConn

	Router chi.Router

	Kafka *kafka.Kafka

	Logger logger.LoggerInterface

	CacheApiGateway mencache.CacheApiGateway

	Cache *cache.CacheStore

	ApiHandler apierror.ApiHandler

	AISecurity  *grpc.ClientConn
	StatsClient *grpc.ClientConn
}

func RegisterTransactionHandler(deps *DepsTransaction) {
	mapper := apimapper.NewTransactionResponseMapper()

	cache := transaction_cache.NewTransactionMencache(deps.Cache)

	handlers := []func(){
		setupTransactionQueryHandler(deps, mapper.QueryMapper(), cache),
		setupTransactionCommandHandler(deps, deps.CacheApiGateway, mapper.CommandMapper(), cache),
		setupTransactionStatsHandler(deps),
	}

	for _, h := range handlers {
		h()
	}
}

func setupTransactionQueryHandler(deps *DepsTransaction, mapper apimapper.TransactionQueryResponseMapper, cache transaction_cache.TransactionMencache) func() {
	return func() {
		NewTransactionQueryHandleApi(&transactionQueryHandleDeps{
			client:     pb.NewTransactionQueryServiceClient(deps.Client),
			router:     deps.Router,
			logger:     deps.Logger,
			mapper:     mapper,
			cache:      cache,
			apiHandler: deps.ApiHandler,
		})
	}
}

func setupTransactionCommandHandler(deps *DepsTransaction, cache mencache.MerchantCache, mapper apimapper.TransactionCommandResponseMapper, cache_ transaction_cache.TransactionMencache) func() {
	return func() {
		NewTransactionCommandHandleApi(&transactionCommandHandleDeps{
			kafka:             deps.Kafka,
			client:            pb.NewTransactionCommandServiceClient(deps.Client),
			router:            deps.Router,
			logger:            deps.Logger,
			mapper:            mapper,
			cache:             cache,
			cache_transaction: cache_,
			apiHandler:        deps.ApiHandler,
			aiSecurity:        pbAISecurity.NewAISecurityServiceClient(deps.AISecurity),
		})
	}
}
