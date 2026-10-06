package userrolegrpcerrors

import (
	"github.com/MamangRust/microservice-payment-gateway-grpc/shared/errors"
)

// ErrAssignRoleToUser is returned when assigning a role to a user fails.
var ErrAssignRoleToUser = errors.ErrInternal.WithMessage("failed to assign role to user")

// ErrRemoveRoleFromUser is returned when removing a role from a user fails.
var ErrRemoveRoleFromUser = errors.ErrInternal.WithMessage("failed to remove role from user")

// ErrFindUserRole is returned when resolving the roles of a user fails.
var ErrFindUserRole = errors.ErrInternal.WithMessage("failed to find roles for user")
