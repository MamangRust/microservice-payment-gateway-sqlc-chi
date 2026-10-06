# 📦 Package `merchanthandler`

**Source Path:** `apigateway/internal/handler/merchant`

## 🧩 Types

### `DepsMerchant`

```go
type DepsMerchant struct {
	Client *grpc.ClientConn
	Router chi.Router
	Logger logger.LoggerInterface
}
```

### `merchantCommandHandleApi`

```go
type merchantCommandHandleApi struct {
	client pb.MerchantCommandServiceClient
	logger logger.LoggerInterface
	mapper apimapper.MerchantCommandResponseMapper
	trace trace.Tracer
	requestCounter *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}
```

#### Methods

##### `Create`

Create godoc
@Summary Create a new merchant
@Tags Merchant
@Security Bearer
@Description Create a new merchant with the given name and user ID
@Accept json
@Produce json
@Param body body requests.CreateMerchantRequest true "Create merchant request"
@Success 200 {object} response.ApiResponseMerchant "Created merchant"
@Failure 400 {object} response.ErrorResponse "Bad request or validation error"
@Failure 500 {object} response.ErrorResponse "Failed to create merchant"
@Router /api/merchant-command/create [post]

```go
func (h *merchantCommandHandleApi) Create(w http.ResponseWriter, r *http.Request) error
```

##### `Delete`

Delete godoc
@Summary Delete a merchant permanently
@Tags Merchant
@Security Bearer
@Description Delete a merchant by its ID permanently
@Accept json
@Produce json
@Param id path int true "Merchant ID"
@Success 200 {object} response.ApiResponseMerchantDelete "Deleted merchant"
@Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
@Failure 500 {object} response.ErrorResponse "Failed to delete merchant"
@Router /api/merchant-command/{id} [delete]

```go
func (h *merchantCommandHandleApi) Delete(w http.ResponseWriter, r *http.Request) error
```

##### `DeleteAllMerchantPermanent`

DeleteAllMerchantPermanent godoc.
@Summary Permanently delete all merchant records
@Tags Merchant
@Security Bearer
@Description Permanently delete all merchant records from the database.
@Accept json
@Produce json
@Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted all merchant records permanently"
@Failure 500 {object} response.ErrorResponse "Failed to permanently delete all merchant records"
@Router /api/merchant-command/permanent/all [post]

```go
func (h *merchantCommandHandleApi) DeleteAllMerchantPermanent(w http.ResponseWriter, r *http.Request) error
```

##### `RestoreAllMerchant`

RestoreAllMerchant godoc.
@Summary Restore all merchant records
@Tags Merchant
@Security Bearer
@Description Restore all merchant records that were previously deleted.
@Accept json
@Produce json
@Success 200 {object} response.ApiResponseMerchantAll "Successfully restored all merchant records"
@Failure 500 {object} response.ErrorResponse "Failed to restore all merchant records"
@Router /api/merchant-command/restore/all [post]

```go
func (h *merchantCommandHandleApi) RestoreAllMerchant(w http.ResponseWriter, r *http.Request) error
```

##### `RestoreMerchant`

RestoreMerchant godoc
@Summary Restore a merchant
@Tags Merchant
@Security Bearer
@Description Restore a merchant by its ID
@Accept json
@Produce json
@Param id path int true "Merchant ID"
@Success 200 {object} response.ApiResponseMerchant "Restored merchant"
@Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
@Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
@Router /api/merchant-command/restore/{id} [post]

```go
func (h *merchantCommandHandleApi) RestoreMerchant(w http.ResponseWriter, r *http.Request) error
```

##### `TrashedMerchant`

TrashedMerchant godoc
@Summary Trashed a merchant
@Tags Merchant
@Security Bearer
@Description Trashed a merchant by its ID
@Accept json
@Produce json
@Param id path int true "Merchant ID"
@Success 200 {object} response.ApiResponseMerchant "Trashed merchant"
@Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
@Failure 500 {object} response.ErrorResponse "Failed to trashed merchant"
@Router /api/merchant-command/trashed/{id} [post]

```go
func (h *merchantCommandHandleApi) TrashedMerchant(w http.ResponseWriter, r *http.Request) error
```

##### `Update`

Update godoc
@Summary Update a merchant
@Tags Merchant
@Security Bearer
@Description Update a merchant with the given ID
@Accept json
@Produce json
@Param body body requests.UpdateMerchantRequest true "Update merchant request"
@Success 200 {object} response.ApiResponseMerchant "Updated merchant"
@Failure 400 {object} response.ErrorResponse "Bad request or validation error"
@Failure 500 {object} response.ErrorResponse "Failed to update merchant"
@Router /api/merchant-command/update/{id} [post]

```go
func (h *merchantCommandHandleApi) Update(w http.ResponseWriter, r *http.Request) error
```

##### `UpdateStatus`

UpdateStatus godoc
@Summary Update merchant status
@Tags Merchant
@Security Bearer
@Description Update the status of a merchant with the given ID
@Accept json
@Produce json
@Param id path int true "Merchant ID"
@Param body body requests.UpdateMerchantStatusRequest true "Update merchant status request"
@Success 200 {object} response.ApiResponseMerchant "Updated merchant status"
@Failure 400 {object} response.ErrorResponse "Bad request or validation error"
@Failure 500 {object} response.ErrorResponse "Failed to update merchant status"
@Router /api/merchant-command/update-status/{id} [post]

```go
func (h *merchantCommandHandleApi) UpdateStatus(w http.ResponseWriter, r *http.Request) error
```

##### `recordMetrics`

```go
func (s *merchantCommandHandleApi) recordMetrics(method string, status string, start time.Time)
```

##### `startTracingAndLogging`

```go
func (s *merchantCommandHandleApi) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (end func(), logSuccess func(string, ...zap.Field), logError func(string, error, ...zap.Field))
```

### `merchantCommandHandleParams`

```go
type merchantCommandHandleParams struct {
	client pb.MerchantCommandServiceClient
	router chi.Router
	logger logger.LoggerInterface
	mapper apimapper.MerchantCommandResponseMapper
}
```

### `merchantQueryHandleApi`

```go
type merchantQueryHandleApi struct {
	client pb.MerchantQueryServiceClient
	logger logger.LoggerInterface
	mapper apimapper.MerchantQueryResponseMapper
	trace trace.Tracer
	requestCounter *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}
```

#### Methods

##### `FindAll`

FindAll godoc
@Summary Find all merchants
@Tags Merchant Query
@Security Bearer
@Description Retrieve a list of all merchants
@Accept json
@Produce json
@Param page query int false "Page number" default(1)
@Param page_size query int false "Number of items per page" default(10)
@Param search query string false "Search query"
@Success 200 {object} response.ApiResponsePaginationMerchant "List of merchants"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
@Router /api/merchant-query [get]

```go
func (h *merchantQueryHandleApi) FindAll(w http.ResponseWriter, r *http.Request) error
```

##### `FindByActive`

FindByActive godoc
@Summary Find active merchants
@Tags Merchant Query
@Security Bearer
@Description Retrieve a list of active merchants
@Accept json
@Produce json
@Param page query int false "Page number" default(1)
@Param page_size query int false "Number of items per page" default(10)
@Param search query string false "Search query"
@Success 200 {object} response.ApiResponsesMerchant "List of active merchants"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
@Router /api/merchant-query/active [get]

```go
func (h *merchantQueryHandleApi) FindByActive(w http.ResponseWriter, r *http.Request) error
```

##### `FindByApiKey`

FindByApiKey godoc
@Summary Find a merchant by API key
@Tags Merchant Query
@Security Bearer
@Description Retrieve a merchant by its API key
@Accept json
@Produce json
@Param api_key query string true "API key"
@Success 200 {object} response.ApiResponseMerchant "Merchant data"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
@Router /api/merchant-query/api-key [get]

```go
func (h *merchantQueryHandleApi) FindByApiKey(w http.ResponseWriter, r *http.Request) error
```

##### `FindById`

FindById godoc
@Summary Find a merchant by ID
@Tags Merchant Query
@Security Bearer
@Description Retrieve a merchant by its ID.
@Accept json
@Produce json
@Param id path int true "Merchant ID"
@Success 200 {object} response.ApiResponseMerchant "Merchant data"
@Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
@Router /api/merchant-query/{id} [get]

```go
func (h *merchantQueryHandleApi) FindById(w http.ResponseWriter, r *http.Request) error
```

##### `FindByMerchantUserId`

FindByMerchantUserId godoc.
@Summary Find a merchant by user ID
@Tags Merchant Query
@Security Bearer
@Description Retrieve a merchant by its user ID
@Accept json
@Produce json
@Param id path int true "User ID"
@Success 200 {object} response.ApiResponsesMerchant "Merchant data"
@Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
@Router /api/merchant-query/merchant-user [get]

```go
func (h *merchantQueryHandleApi) FindByMerchantUserId(w http.ResponseWriter, r *http.Request) error
```

##### `FindByTrashed`

FindByTrashed godoc
@Summary Find trashed merchants
@Tags Merchant Query
@Security Bearer
@Description Retrieve a list of trashed merchants
@Accept json
@Produce json
@Param page query int false "Page number" default(1)
@Param page_size query int false "Number of items per page" default(10)
@Param search query string false "Search query"
@Success 200 {object} response.ApiResponsesMerchant "List of trashed merchants"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
@Router /api/merchant-query/trashed [get]

```go
func (h *merchantQueryHandleApi) FindByTrashed(w http.ResponseWriter, r *http.Request) error
```

##### `recordMetrics`

```go
func (s *merchantQueryHandleApi) recordMetrics(method string, status string, start time.Time)
```

##### `startTracingAndLogging`

```go
func (s *merchantQueryHandleApi) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (end func(), logSuccess func(string, ...zap.Field), logError func(string, error, ...zap.Field))
```

### `merchantQueryHandleParams`

```go
type merchantQueryHandleParams struct {
	client pb.MerchantQueryServiceClient
	router chi.Router
	logger logger.LoggerInterface
	mapper apimapper.MerchantQueryResponseMapper
}
```

### `merchantStatsAmountHandleApi`

```go
type merchantStatsAmountHandleApi struct {
	client pb.MerchantStatsAmountServiceClient
	logger logger.LoggerInterface
	mapper apimapper.MerchantStatsAmountResponseMapper
	trace trace.Tracer
	requestCounter *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}
```

#### Methods

##### `FindMonthlyAmountByApikeys`

FindMonthlyAmountByApikeys godoc.
@Summary Find monthly transaction amounts for a specific merchant
@Tags Merchant
@Security Bearer
@Description Retrieve monthly transaction amounts for a specific merchant by year.
@Accept json
@Produce json
@Param merchant_id query int true "Merchant ID"
@Param year query int true "Year"
@Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
@Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
@Router /api/merchant-stats-amount/monthly-amount-by-apikey [get]

```go
func (h *merchantStatsAmountHandleApi) FindMonthlyAmountByApikeys(w http.ResponseWriter, r *http.Request) error
```

##### `FindMonthlyAmountByMerchants`

FindMonthlyAmountByMerchants godoc.
@Summary Find monthly transaction amounts for a specific merchant
@Tags Merchant
@Security Bearer
@Description Retrieve monthly transaction amounts for a specific merchant by year.
@Accept json
@Produce json
@Param merchant_id query int true "Merchant ID"
@Param year query int true "Year"
@Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
@Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
@Router /api/merchant-stats-amount/monthly-amount-by-merchant [get]

```go
func (h *merchantStatsAmountHandleApi) FindMonthlyAmountByMerchants(w http.ResponseWriter, r *http.Request) error
```

##### `FindMonthlyAmountMerchant`

FindMonthlyAmountMerchant godoc
@Summary Find monthly transaction amounts for a merchant
@Tags Merchant
@Security Bearer
@Description Retrieve monthly transaction amounts for a merchant by year.
@Accept json
@Produce json
@Param year query int true "Year"
@Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
@Failure 400 {object} response.ErrorResponse "Invalid year"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
@Router /api/merchant-stats-amount/monthly-amount [get]

```go
func (h *merchantStatsAmountHandleApi) FindMonthlyAmountMerchant(w http.ResponseWriter, r *http.Request) error
```

##### `FindYearlyAmountByApikeys`

FindYearlyAmountByApikeys godoc.
@Summary Find yearly transaction amounts for a specific merchant
@Tags Merchant
@Security Bearer
@Description Retrieve yearly transaction amounts for a specific merchant by year.
@Accept json
@Produce json
@Param merchant_id query int true "Merchant ID"
@Param year query int true "Year"
@Success 200 {object} response.ApiResponseMerchantYearlyAmount "Yearly transaction amounts"
@Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
@Router /api/merchant-stats-amount/yearly-amount-by-apikey [get]

```go
func (h *merchantStatsAmountHandleApi) FindYearlyAmountByApikeys(w http.ResponseWriter, r *http.Request) error
```

##### `FindYearlyAmountByMerchants`

FindYearlyAmountByMerchants godoc.
@Summary Find yearly transaction amounts for a specific merchant
@Tags Merchant
@Security Bearer
@Description Retrieve yearly transaction amounts for a specific merchant by year.
@Accept json
@Produce json
@Param merchant_id query int true "Merchant ID"
@Param year query int true "Year"
@Success 200 {object} response.ApiResponseMerchantYearlyAmount "Yearly transaction amounts"
@Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
@Router /api/merchant-stats-amount/yearly-amount-by-merchant [get]

```go
func (h *merchantStatsAmountHandleApi) FindYearlyAmountByMerchants(w http.ResponseWriter, r *http.Request) error
```

##### `FindYearlyAmountMerchant`

FindYearlyAmountMerchant godoc.
@Summary Find yearly transaction amounts for a merchant
@Tags Merchant
@Security Bearer
@Description Retrieve yearly transaction amounts for a merchant by year.
@Accept json
@Produce json
@Param year query int true "Year"
@Success 200 {object} response.ApiResponseYearlyAmount "Yearly transaction amounts"
@Failure 400 {object} response.ErrorResponse "Invalid year"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
@Router /api/merchant-stats-amount/yearly-amount [get]

```go
func (h *merchantStatsAmountHandleApi) FindYearlyAmountMerchant(w http.ResponseWriter, r *http.Request) error
```

##### `recordMetrics`

```go
func (s *merchantStatsAmountHandleApi) recordMetrics(method string, status string, start time.Time)
```

##### `startTracingAndLogging`

```go
func (s *merchantStatsAmountHandleApi) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (end func(), logSuccess func(string, ...zap.Field), logError func(string, error, ...zap.Field))
```

### `merchantStatsAmountHandleParams`

```go
type merchantStatsAmountHandleParams struct {
	client pb.MerchantStatsAmountServiceClient
	router chi.Router
	logger logger.LoggerInterface
	mapper apimapper.MerchantStatsAmountResponseMapper
}
```

### `merchantStatsMethodHandleApi`

```go
type merchantStatsMethodHandleApi struct {
	client pb.MerchantStatsMethodServiceClient
	logger logger.LoggerInterface
	mapper apimapper.MerchantStatsMethodResponseMapper
	trace trace.Tracer
	requestCounter *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}
```

#### Methods

##### `FindMonthlyPaymentMethodByApikeys`

FindMonthlyPaymentMethodByApikeys godoc.
@Summary Find monthly payment methods for a specific merchant
@Tags Merchant
@Security Bearer
@Description Retrieve monthly payment methods for a specific merchant by year.
@Accept json
@Produce json
@Param merchant_id query int true "Merchant ID"
@Param year query int true "Year"
@Success 200 {object} response.ApiResponseMerchantMonthlyPaymentMethod "Monthly payment methods"
@Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly payment methods"
@Router /api/merchant-stats-amount/monthly-payment-methods-by-apikey [get]

```go
func (h *merchantStatsMethodHandleApi) FindMonthlyPaymentMethodByApikeys(w http.ResponseWriter, r *http.Request) error
```

##### `FindMonthlyPaymentMethodByMerchants`

FindMonthlyPaymentMethodByMerchants godoc.
@Summary Find monthly payment methods for a specific merchant
@Tags Merchant
@Security Bearer
@Description Retrieve monthly payment methods for a specific merchant by year.
@Accept json
@Produce json
@Param merchant_id query int true "Merchant ID"
@Param year query int true "Year"
@Success 200 {object} response.ApiResponseMerchantMonthlyPaymentMethod "Monthly payment methods"
@Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly payment methods"
@Router /api/merchant-stats-amount/monthly-payment-methods-by-merchant [get]

```go
func (h *merchantStatsMethodHandleApi) FindMonthlyPaymentMethodByMerchants(w http.ResponseWriter, r *http.Request) error
```

##### `FindMonthlyPaymentMethodsMerchant`

FindMonthlyPaymentMethodsMerchant godoc
@Summary Find monthly payment methods for a merchant
@Tags Merchant
@Security Bearer
@Description Retrieve monthly payment methods for a merchant by year.
@Accept json
@Produce json
@Param year query int true "Year"
@Success 200 {object} response.ApiResponseMerchantMonthlyPaymentMethod "Monthly payment methods"
@Failure 400 {object} response.ErrorResponse "Invalid year"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly payment methods"
@Router /api/merchant-stats-amount/monthly-payment-methods [get]

```go
func (h *merchantStatsMethodHandleApi) FindMonthlyPaymentMethodsMerchant(w http.ResponseWriter, r *http.Request) error
```

##### `FindYearlyPaymentMethodByApikeys`

FindYearlyPaymentMethodByApikeys godoc.
@Summary Find yearly payment methods for a specific merchant
@Tags Merchant
@Security Bearer
@Description Retrieve yearly payment methods for a specific merchant by year.
@Accept json
@Produce json
@Param merchant_id query int true "Merchant ID"
@Param year query int true "Year"
@Success 200 {object} response.ApiResponseMerchantYearlyPaymentMethod "Yearly payment methods"
@Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly payment methods"
@Router /api/merchant-stats-amount/yearly-payment-methods-by-apikey [get]

```go
func (h *merchantStatsMethodHandleApi) FindYearlyPaymentMethodByApikeys(w http.ResponseWriter, r *http.Request) error
```

##### `FindYearlyPaymentMethodByMerchants`

FindYearlyPaymentMethodByMerchants godoc.
@Summary Find yearly payment methods for a specific merchant
@Tags Merchant
@Security Bearer
@Description Retrieve yearly payment methods for a specific merchant by year.
@Accept json
@Produce json
@Param merchant_id query int true "Merchant ID"
@Param year query int true "Year"
@Success 200 {object} response.ApiResponseMerchantYearlyPaymentMethod "Yearly payment methods"
@Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly payment methods"
@Router /api/merchant-stats-amount/yearly-payment-methods-by-merchant [get]

```go
func (h *merchantStatsMethodHandleApi) FindYearlyPaymentMethodByMerchants(w http.ResponseWriter, r *http.Request) error
```

##### `FindYearlyPaymentMethodMerchant`

FindYearlyPaymentMethodMerchant godoc.
@Summary Find yearly payment methods for a merchant
@Tags Merchant
@Security Bearer
@Description Retrieve yearly payment methods for a merchant by year.
@Accept json
@Produce json
@Param year query int true "Year"
@Success 200 {object} response.ApiResponseMerchantYearlyPaymentMethod "Yearly payment methods"
@Failure 400 {object} response.ErrorResponse "Invalid year"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly payment methods"
@Router /api/merchant-stats-amount/yearly-payment-methods [get]

```go
func (h *merchantStatsMethodHandleApi) FindYearlyPaymentMethodMerchant(w http.ResponseWriter, r *http.Request) error
```

##### `recordMetrics`

```go
func (s *merchantStatsMethodHandleApi) recordMetrics(method string, status string, start time.Time)
```

##### `startTracingAndLogging`

```go
func (s *merchantStatsMethodHandleApi) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (end func(), logSuccess func(string, ...zap.Field), logError func(string, error, ...zap.Field))
```

### `merchantStatsMethodHandleParams`

```go
type merchantStatsMethodHandleParams struct {
	client pb.MerchantStatsMethodServiceClient
	router chi.Router
	logger logger.LoggerInterface
	mapper apimapper.MerchantStatsMethodResponseMapper
}
```

### `merchantStatsTotalAmountHandleApi`

```go
type merchantStatsTotalAmountHandleApi struct {
	client pb.MerchantStatsTotalAmountServiceClient
	logger logger.LoggerInterface
	mapper apimapper.MerchantStatsTotalAmountResponseMapper
	trace trace.Tracer
	requestCounter *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}
```

#### Methods

##### `FindMonthlyTotalAmountByApikeys`

FindMonthlyAmountByApikeys godoc.
@Summary Find monthly transaction amounts for a specific merchant
@Tags Merchant
@Security Bearer
@Description Retrieve monthly transaction amounts for a specific merchant by year.
@Accept json
@Produce json
@Param merchant_id query int true "Merchant ID"
@Param year query int true "Year"
@Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
@Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
@Router /api/merchant-stats-totalamount/monthly-totalamount-by-apikey [get]

```go
func (h *merchantStatsTotalAmountHandleApi) FindMonthlyTotalAmountByApikeys(w http.ResponseWriter, r *http.Request) error
```

##### `FindMonthlyTotalAmountByMerchants`

FindMonthlyAmountByMerchants godoc.
@Summary Find monthly transaction amounts for a specific merchant
@Tags Merchant
@Security Bearer
@Description Retrieve monthly transaction amounts for a specific merchant by year.
@Accept json
@Produce json
@Param merchant_id query int true "Merchant ID"
@Param year query int true "Year"
@Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
@Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
@Router /api/merchant-stats-totalamount/monthly-totalamount-by-merchant [get]

```go
func (h *merchantStatsTotalAmountHandleApi) FindMonthlyTotalAmountByMerchants(w http.ResponseWriter, r *http.Request) error
```

##### `FindMonthlyTotalAmountMerchant`

FindMonthlyAmountMerchant godoc
@Summary Find monthly transaction amounts for a merchant
@Tags Merchant
@Security Bearer
@Description Retrieve monthly transaction amounts for a merchant by year.
@Accept json
@Produce json
@Param year query int true "Year"
@Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
@Failure 400 {object} response.ErrorResponse "Invalid year"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
@Router /api/merchant-stats-totalamount/monthly-total-amount [get]

```go
func (h *merchantStatsTotalAmountHandleApi) FindMonthlyTotalAmountMerchant(w http.ResponseWriter, r *http.Request) error
```

##### `FindYearlyTotalAmountByApikeys`

FindYearlyAmountByApikeys godoc.
@Summary Find yearly transaction amounts for a specific merchant
@Tags Merchant
@Security Bearer
@Description Retrieve yearly transaction amounts for a specific merchant by year.
@Accept json
@Produce json
@Param merchant_id query int true "Merchant ID"
@Param year query int true "Year"
@Success 200 {object} response.ApiResponseMerchantYearlyAmount "Yearly transaction amounts"
@Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
@Router /api/merchant-stats-totalamount/yearly-totalamount-by-apikey [get]

```go
func (h *merchantStatsTotalAmountHandleApi) FindYearlyTotalAmountByApikeys(w http.ResponseWriter, r *http.Request) error
```

##### `FindYearlyTotalAmountByMerchants`

FindYearlyAmountByMerchants godoc.
@Summary Find yearly transaction amounts for a specific merchant
@Tags Merchant
@Security Bearer
@Description Retrieve yearly transaction amounts for a specific merchant by year.
@Accept json
@Produce json
@Param merchant_id query int true "Merchant ID"
@Param year query int true "Year"
@Success 200 {object} response.ApiResponseMerchantYearlyAmount "Yearly transaction amounts"
@Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
@Router /api/merchant-stats-totalamount/yearly-totalamount-by-merchant [get]

```go
func (h *merchantStatsTotalAmountHandleApi) FindYearlyTotalAmountByMerchants(w http.ResponseWriter, r *http.Request) error
```

##### `FindYearlyTotalAmountMerchant`

FindYearlyAmountMerchant godoc.
@Summary Find yearly transaction amounts for a merchant
@Tags Merchant
@Security Bearer
@Description Retrieve yearly transaction amounts for a merchant by year.
@Accept json
@Produce json
@Param year query int true "Year"
@Success 200 {object} response.ApiResponseYearlyAmount "Yearly transaction amounts"
@Failure 400 {object} response.ErrorResponse "Invalid year"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
@Router /api/merchant-stats-totalamount/yearly-total-amount [get]

```go
func (h *merchantStatsTotalAmountHandleApi) FindYearlyTotalAmountMerchant(w http.ResponseWriter, r *http.Request) error
```

##### `recordMetrics`

```go
func (s *merchantStatsTotalAmountHandleApi) recordMetrics(method string, status string, start time.Time)
```

##### `startTracingAndLogging`

```go
func (s *merchantStatsTotalAmountHandleApi) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (end func(), logSuccess func(string, ...zap.Field), logError func(string, error, ...zap.Field))
```

### `merchantStatsTotalAmountHandleParams`

```go
type merchantStatsTotalAmountHandleParams struct {
	client pb.MerchantStatsTotalAmountServiceClient
	router chi.Router
	logger logger.LoggerInterface
	mapper apimapper.MerchantStatsTotalAmountResponseMapper
}
```

### `merchantTransactionHandleApi`

```go
type merchantTransactionHandleApi struct {
	client pb.MerchantTransactionServiceClient
	logger logger.LoggerInterface
	mapper apimapper.MerchantTransactionResponseMapper
	trace trace.Tracer
	requestCounter *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}
```

#### Methods

##### `FindAllTransactionByApikey`

FindAllTransactionByApikey godoc
@Summary Find all transactions by api_key
@Tags Merchant
@Security Bearer
@Description Retrieve a list of transactions for a specific merchant
@Accept json
@Produce json
@Param api_key path string true "Api key"
@Param page query int false "Page number" default(1)
@Param page_size query int false "Number of items per page" default(10)
@Param search query string false "Search query"
@Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
@Router /api/merchant-transactions/api-key/:api_key [get]

```go
func (h *merchantTransactionHandleApi) FindAllTransactionByApikey(w http.ResponseWriter, r *http.Request) error
```

##### `FindAllTransactionByMerchant`

FindAllTransactionByMerchant godoc
@Summary Find all transactions by merchant ID
@Tags Merchant
@Security Bearer
@Description Retrieve a list of transactions for a specific merchant
@Accept json
@Produce json
@Param merchant_id path int true "Merchant ID"
@Param page query int false "Page number" default(1)
@Param page_size query int false "Number of items per page" default(10)
@Param search query string false "Search query"
@Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
@Router /api/merchant-transactions/:merchant_id [get]

```go
func (h *merchantTransactionHandleApi) FindAllTransactionByMerchant(w http.ResponseWriter, r *http.Request) error
```

##### `FindAllTransactions`

FindAllTransactions godoc
@Summary Find all transactions
@Tags Merchant
@Security Bearer
@Description Retrieve a list of all transactions
@Accept json
@Produce json
@Param page query int false "Page number" default(1)
@Param page_size query int false "Number of items per page" default(10)
@Param search query string false "Search query"
@Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
@Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
@Router /api/merchants/transaction [get]

```go
func (h *merchantTransactionHandleApi) FindAllTransactions(w http.ResponseWriter, r *http.Request) error
```

##### `recordMetrics`

```go
func (s *merchantTransactionHandleApi) recordMetrics(method string, status string, start time.Time)
```

##### `startTracingAndLogging`

```go
func (s *merchantTransactionHandleApi) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (end func(), logSuccess func(string, ...zap.Field), logError func(string, error, ...zap.Field))
```

### `merchantTransactionHandleParams`

```go
type merchantTransactionHandleParams struct {
	client pb.MerchantTransactionServiceClient
	router chi.Router
	logger logger.LoggerInterface
	mapper apimapper.MerchantTransactionResponseMapper
}
```

## 🚀 Functions

### `RegisterMerchantHandler`

```go
func RegisterMerchantHandler(deps *DepsMerchant)
```

### `setupMerchantCommandHandler`

```go
func setupMerchantCommandHandler(deps *DepsMerchant, mapper apimapper.MerchantCommandResponseMapper) func()
```

### `setupMerchantQueryHandler`

```go
func setupMerchantQueryHandler(deps *DepsMerchant, mapper apimapper.MerchantQueryResponseMapper) func()
```

### `setupMerchantStatsAmountHandler`

```go
func setupMerchantStatsAmountHandler(deps *DepsMerchant, mapper apimapper.MerchantStatsAmountResponseMapper) func()
```

### `setupMerchantStatsMethodHandler`

```go
func setupMerchantStatsMethodHandler(deps *DepsMerchant, mapper apimapper.MerchantStatsMethodResponseMapper) func()
```

### `setupMerchantStatsTotalAmountHandler`

```go
func setupMerchantStatsTotalAmountHandler(deps *DepsMerchant, mapper apimapper.MerchantStatsTotalAmountResponseMapper) func()
```

### `setupMerchantTransactionHandler`

```go
func setupMerchantTransactionHandler(deps *DepsMerchant, mapper apimapper.MerchantTransactionResponseMapper) func()
```

