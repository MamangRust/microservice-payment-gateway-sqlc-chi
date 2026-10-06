# 📦 Package `auth_errors`

**Source Path:** `shared/errors/auth_errors`

## 🏷️ Variables

**Var:**

ErrApiForgotPassword returns an API error response when forgot password request fails.

```go
var ErrApiForgotPassword = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "failed to forgot password", http.StatusBadRequest)
}
```

**Var:**

ErrApiGetMe returns an API error response when unauthenticated user tries to access their info.

```go
var ErrApiGetMe = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "get user info failed: unauthenticated", http.StatusUnauthorized)
}
```

**Var:**

ErrApiLogin returns an API error response when login fails due to invalid arguments.

```go
var ErrApiLogin = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "login failed: invalid argument provided", http.StatusBadRequest)
}
```

**Var:**

ErrApiRefreshToken returns an API error response when access token is invalid during refresh token flow.

```go
var ErrApiRefreshToken = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "refresh-token failed: invalid access token", http.StatusBadRequest)
}
```

**Var:**

ErrApiRegister returns an API error response when register request fails due to invalid arguments.

```go
var ErrApiRegister = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "register failed: invalid argument", http.StatusBadRequest)
}
```

**Var:**

ErrApiResetPassword returns an API error response when reset password request fails.

```go
var ErrApiResetPassword = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "failed to reset password", http.StatusBadRequest)
}
```

**Var:**

ErrApiVerifyCode returns an API error response when code verification fails.

```go
var ErrApiVerifyCode = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "failed to verify code", http.StatusBadRequest)
}
```

**Var:**

ErrBindForgotPassword returns an API error response for binding failure on forgot password request payload.

```go
var ErrBindForgotPassword = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "binding failed: invalid forgot password request payload", http.StatusBadRequest)
}
```

**Var:**

ErrBindLogin returns an API error response for binding failure on login request payload.

```go
var ErrBindLogin = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "binding failed: invalid login request payload", http.StatusBadRequest)
}
```

**Var:**

ErrBindRefreshToken returns an API error response for binding failure on refresh token request payload.

```go
var ErrBindRefreshToken = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "binding failed: invalid refresh token request payload", http.StatusBadRequest)
}
```

**Var:**

ErrBindRegister returns an API error response for binding failure on register request payload.

```go
var ErrBindRegister = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "binding failed: invalid register request payload", http.StatusBadRequest)
}
```

**Var:**

ErrBindResetPassword returns an API error response for binding failure on reset password request payload.

```go
var ErrBindResetPassword = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "binding failed: invalid reset password request payload", http.StatusBadRequest)
}
```

```go
var ErrGrpcGetMe = response.NewGrpcError("error", "get user info failed: unauthenticated", int(codes.Unauthenticated))
```

```go
var ErrGrpcLogin = response.NewGrpcError("error", "login failed: invalid argument provided", int(codes.InvalidArgument))
```

```go
var ErrGrpcRegisterToken = response.NewGrpcError("error", "register failed: invalid argument", int(codes.InvalidArgument))
```

**Var:**

ErrInvalidAccessToken returns an API error response for an invalid access token.

```go
var ErrInvalidAccessToken = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "invalid access token", http.StatusInternalServerError)
}
```

**Var:**

ErrInvalidLogin returns an API error response for incorrect email or password.

```go
var ErrInvalidLogin = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "invalid email or password", http.StatusUnauthorized)
}
```

**Var:**

ErrValidateForgotPassword returns an API error response for invalid forgot password request validation.

```go
var ErrValidateForgotPassword = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "validation failed: invalid forgot-password request", http.StatusBadRequest)
}
```

**Var:**

ErrValidateLogin returns an API error response for invalid login request validation.

```go
var ErrValidateLogin = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "validation failed: invalid login request", http.StatusBadRequest)
}
```

**Var:**

ErrValidateRefreshToken returns an API error response for invalid refresh token request validation.

```go
var ErrValidateRefreshToken = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "validation failed: invalid refresh-token request", http.StatusBadRequest)
}
```

**Var:**

ErrValidateRegister returns an API error response for invalid register request validation.

```go
var ErrValidateRegister = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "validation failed: invalid register request", http.StatusBadRequest)
}
```

**Var:**

ErrValidateResetPassword returns an API error response for invalid reset password request validation.

```go
var ErrValidateResetPassword = func(w http.ResponseWriter, r *http.Request) error {
	return response.NewApiErrorResponse(w, "error", "validation failed: invalid reset-password request", http.StatusBadRequest)
}
```

