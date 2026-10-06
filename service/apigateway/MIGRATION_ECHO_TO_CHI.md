# Migrasi Apigateway: Echo → Chi

> **Lingkup:** `microservice-payment-gateway-grpc/service/apigateway` (+ harness `tests/`)
> **Router:** `github.com/labstack/echo/v4` → `github.com/go-chi/chi/v5`
> **Status:** ✅ Selesai — build ✅ · vet ✅ · unit test ✅
> **Pola:** sama dengan migrasi apigateway `monolith-payment-gateway-grpc` (lihat `modular-monolith-grpc-sqlc-restapi/monolith-payment-gateway-grpc/service/apigateway/MIGRATION_ECHO_TO_CHI.md`) serta `microservice-pointofsale-grpc` dan `microservice-ecommerce-grpc` di workspace ini.

---

## 1. Ringkasan

| Item | Jumlah |
|---|---|
| File yang import echo (sebelum) | 49 (38 handler, 8 middleware + test, `app/client.go`) |
| Paket handler domain yang dimigrasi | 11 (`auth`, `card`, `merchant`, `merchantdocument`, `role`, `saldo`, `topup`, `transaction`, `transfer`, `user`, `withdraw`) — struktur subpackage (`handler.go`/`query.go`/`command.go`/`stats.go`/`dashboard.go`/`transaction.go`), 39 file |
| Rute terdaftar (terkonversi) | 247 (165 GET · 69 POST · 11 DELETE · 1 PUT + hello/health) — terverifikasi identik dengan baseline echo |
| Middleware global | 8 (`Recoverer → RequestID → Logger(JSON) → Pyroscope → CORS → Compress(5) → Secure → JWTAuth`) |
| Middleware per-route | API-key validator Kafka (`r.With(...)`) di create/update transaction; role & require-role di role |
| Test harness `tests/` yang dikonversi | 10 file `handler_api_test.go` |

## 2. Karakteristik Khusus Payment-Gateway

1. **JWT menyimpan dua context key** — `JWTAuth()` chi menyimpan claim `sub` ke key `"userId"` sekaligus `"user_id"` (string), menggantikan `WebSecurityConfig` echo-jwt; role-protected handler membaca `"user_id"`.
2. **API-key validator Kafka dipertahankan** — `middlewares.NewApiKeyValidator(...)` (request/response topic Kafka) dikonversi ke `func(http.Handler) http.Handler` dan dipasang via `routerTransaction.With(transactionMiddleware.Middleware()).Post(...)` untuk create/update transaction; handler membaca `httpx.Get(r, "apiKey")`.
3. **Middleware khas lain** — `RoleValidator` (Kafka) + `RequireRoles` (chi middleware, `r.With(roleMiddlewareChain, requireAdmin)`) di paket role; merchant API-key middleware (`MaskIdentifier` dibuat lokal sebagai `maskIdentifier` karena helper itu hanya ada di shared monolith).
4. **File stats/dashboard** — paket `card` (dashboard.go, stats.go), `merchant` (stats.go, transaction.go), `saldo`/`topup`/`transaction`/`transfer`/`withdraw` (stats.go) dikonversi dengan pola `deps.Router.Route("/api/<domain>/stats", func(g chi.Router) {...})`.
5. **Tidak ada upload file** — `pkg/` tidak menyentuh echo sama sekali di repo ini; `shared/errors/api_error.go` juga sudah net/http murni (hanya import `fmt`).

## 3. Dependensi

**Ditambah:** `go-chi/chi/v5 v5.3.2`, `go-chi/cors v1.2.2`, `swaggo/http-swagger v1.3.4` (via `go mod edit` + `go.sum` dari monolith karena repo `MamangRust/*` privat tidak bisa di-`go get`).
**Dihapus total:** `labstack/echo/v4`, `echo-jwt/v4`, `echo-swagger` — tidak ada lagi di `go.mod`/`go.sum` modul mana pun (dibersihkan dengan `go mod tidy`).

## 4. Struktur Baru

```
service/apigateway/
├── apierror/          # NEW — ApiHandler versi net/http
├── httpx/             # NEW — JSON/Bind/RealIP, request-scoped values, error adapter
├── app/client.go      # bootstrap chi + http.Server (Recover → RequestID → Logger → ... → JWT)
├── handler/           # 11 subpackage + handle.go → chi.Router
├── middlewares/       # auth (JWT), merchant (API-key Kafka), role, required_role, rate_limiter, pyroscope
└── ...
tests/                 # 10 file harness terkonversi
```

## 5. Verifikasi

```bash
cd microservice-payment-gateway-grpc/service/apigateway
go build ./... && go vet ./...   # ✅
cd ../../tests && go build ./... && go vet ./... && go test ./...
```

- Registrasi rute dihitung ulang dengan metrik `Handle(`/`httpx.Handler(`: **247 = 247** (baseline echo dari git HEAD), termasuk registrasi multi-line di paket role dan chain `With()` di transaction.
- Harness `tests/` (10 file): `server *echo.Echo` → `chi.Router`, Deps `E:` → `Router:`, `errors.NewApiHandler` → `apierror.NewApiHandler`, bypass middleware `c.Set(...)` → `context.WithValue`.
- `shared/errors/api_error.go` dan seluruh service gRPC tidak tersentuh.

## 6. Catatan Perilaku

1. **Respons error identik** — `*AppError` dipetakan `apierror.HandleApiError` (status asli + trace_id); fallback `echo.NewHTTPError` → `httpx.WriteHTTPError`; 404/405 default echo ditiru via `r.NotFound`/`r.MethodNotAllowed` (JSON `{"message": ...}`).
2. **Sintaks path param** — `:id` → `{id}`; `GET("")` → `Get("/")`.
3. **Swagger** — `echoSwagger.WrapHandler` → `httpSwagger.WrapHandler` di `/swagger/*`; anotasi `@Router` tidak diubah.
4. **Logger JSON** — field-layout echo (`time/id/remote_ip/host/method/uri/status/error/latency/...`) dipertahankan via `chi_mw.RequestLogger` + formatter custom.
5. **Echo sudah tidak ada** di dependency graph (`go.mod`/`go.sum`) seluruh modul.

## 7. Follow-up (Opsional)

- [ ] Smoke test `hurl/` atau full stack docker.
- [x] Migrasi `shared/errors/api_error.go` penuh ke net/http → echo hilang total dari repo.
- [ ] Seluruh apigateway di kedua workspace (monolith & microservice: POS, ecommerce, payment-gateway) kini konsisten memakai chi.
