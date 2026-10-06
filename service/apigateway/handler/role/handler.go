package rolehandler

import (
	"time"

	pb "github.com/MamangRust/microservice-payment-gateway-grpc/pb/role"
	pbuserrole "github.com/MamangRust/microservice-payment-gateway-grpc/pb/user_role"
	"github.com/MamangRust/microservice-payment-gateway-grpc/pkg/kafka"
	"github.com/MamangRust/microservice-payment-gateway-grpc/pkg/logger"
	"github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/middlewares"
	mencache "github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/redis"
	role_cache "github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/redis/api/role"
	"github.com/MamangRust/microservice-payment-gateway-grpc/shared/cache"
	apimapper "github.com/MamangRust/microservice-payment-gateway-grpc/shared/mapper/role"

	"github.com/MamangRust/microservice-payment-gateway-grpc/service/apigateway/apierror"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type DepsRole struct {
	Client *grpc.ClientConn

	Kafka *kafka.Kafka

	Router chi.Router

	Logger logger.LoggerInterface

	Cache mencache.RoleCache

	CacheStore *cache.CacheStore

	ApiHandler apierror.ApiHandler
}

func RegisterRoleHandler(deps *DepsRole) {
	mapper := apimapper.NewRoleResponseMapper()
	cache := role_cache.NewRoleMencache(deps.CacheStore)

	// Single shared RoleValidator for both command and query handlers.
	// Creating one validator per handler spawned two Kafka consumers in the
	// same consumer group on the same single-partition response topic, so the
	// partition was assigned to only one member and correlation responses
	// landed on the wrong member's channel map ("No waiting channel" -> timeout).
	roleValidator := middlewares.NewRoleValidator(deps.Kafka, "request-role", "response-role", 5*time.Second, deps.Logger, deps.Cache)

	handlers := []func(){
		setupRoleQueryHandler(deps, mapper.QueryMapper(), cache, roleValidator),
		setupRoleCommandHandler(deps, mapper.CommandMapper(), cache, roleValidator),
	}

	for _, h := range handlers {
		h()
	}
}

func setupRoleQueryHandler(deps *DepsRole, mapper apimapper.RoleQueryResponseMapper, cache role_cache.RoleMencache, roleValidator *middlewares.RoleValidator) func() {
	return func() {
		NewRoleQueryHandleApi(&roleQueryHandleDeps{
			client:         pb.NewRoleQueryServiceClient(deps.Client),
			userRoleClient: pbuserrole.NewUserRoleServiceClient(deps.Client),
			router:         deps.Router,
			logger:         deps.Logger,
			mapper:         mapper,
			roleValidator:  roleValidator,
			cache:          cache,
			apiHandler:     deps.ApiHandler,
		})
	}
}

func setupRoleCommandHandler(deps *DepsRole, mapper apimapper.RoleCommandResponseMapper, cache role_cache.RoleMencache, roleValidator *middlewares.RoleValidator) func() {
	return func() {
		NewRoleCommandHandleApi(&roleCommandHandleDeps{
			client:        pb.NewRoleCommandServiceClient(deps.Client),
			router:        deps.Router,
			logger:        deps.Logger,
			mapper:        mapper,
			roleValidator: roleValidator,
			cache:         cache,
			apiHandler:    deps.ApiHandler,
		})
	}
}
