package main

import (
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"shop/api"
	"shop/docs/openapi"
	"shop/internal/repository"
	"shop/internal/web/handler"
	"shop/internal/web/middleware"
)

func newAuthClient(conn *grpc.ClientConn) api.AuthServiceClient {
	return api.NewAuthServiceClient(conn)
}

func startServer(authHandler *handler.AuthHandler, addressHandler *handler.AddressHandler, clientHandler *handler.ClientHandler, productHandler *handler.ProductHandler, supplierHandler *handler.SupplierHandler, imageHandler *handler.ImageHandler, authenticator *middleware.UserAuthenticator) error {
	mux := http.NewServeMux()
	mux.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.InstanceName(openapi.SwaggerInfo.InstanceName()),
	))

	mux.HandleFunc("POST /api/v1/register", authHandler.Signup)
	mux.HandleFunc("POST /api/v1/auth", authHandler.Signin)
	mux.HandleFunc("POST /api/v1/reset", authHandler.ResetPassword)
	mux.HandleFunc("POST /api/v1/refresh", authHandler.Refresh)
	mux.HandleFunc("POST /api/v1/change_password", authHandler.ChangePassword)

	mux.Handle("POST /api/v1/addresses", authenticator.Authenticate(http.HandlerFunc(addressHandler.Create)))
	mux.Handle("GET /api/v1/addresses", authenticator.Authenticate(http.HandlerFunc(addressHandler.List)))
	mux.Handle("PUT /api/v1/addresses/{id}", authenticator.Authenticate(http.HandlerFunc(addressHandler.Update)))
	mux.Handle("DELETE /api/v1/addresses/{id}", authenticator.Authenticate(http.HandlerFunc(addressHandler.Delete)))

	mux.Handle("POST /api/v1/clients", authenticator.Authenticate(http.HandlerFunc(clientHandler.Create)))
	mux.Handle("GET /api/v1/clients", authenticator.Authenticate(http.HandlerFunc(clientHandler.Get)))
	mux.Handle("PUT /api/v1/clients/{id}/address", authenticator.Authenticate(http.HandlerFunc(clientHandler.UpdateAddress)))
	mux.Handle("DELETE /api/v1/clients/{id}", authenticator.Authenticate(http.HandlerFunc(clientHandler.Delete)))

	mux.Handle("POST /api/v1/products", authenticator.Authenticate(http.HandlerFunc(productHandler.Create)))
	mux.Handle("GET /api/v1/products", authenticator.Authenticate(http.HandlerFunc(productHandler.List)))
	mux.Handle("GET /api/v1/products/{id}", authenticator.Authenticate(http.HandlerFunc(productHandler.GetByID)))
	mux.Handle("PATCH /api/v1/products/{id}/stock", authenticator.Authenticate(http.HandlerFunc(productHandler.DecreaseStock)))
	mux.Handle("DELETE /api/v1/products/{id}", authenticator.Authenticate(http.HandlerFunc(productHandler.Delete)))

	mux.Handle("POST /api/v1/suppliers", authenticator.Authenticate(http.HandlerFunc(supplierHandler.Create)))
	mux.Handle("GET /api/v1/suppliers", authenticator.Authenticate(http.HandlerFunc(supplierHandler.List)))
	mux.Handle("GET /api/v1/suppliers/{id}", authenticator.Authenticate(http.HandlerFunc(supplierHandler.GetByID)))
	mux.Handle("PUT /api/v1/suppliers/{id}/address", authenticator.Authenticate(http.HandlerFunc(supplierHandler.UpdateAddress)))
	mux.Handle("DELETE /api/v1/suppliers/{id}", authenticator.Authenticate(http.HandlerFunc(supplierHandler.Delete)))

	mux.Handle("POST /api/v1/products/{id}/image", authenticator.Authenticate(http.HandlerFunc(imageHandler.Create)))
	mux.Handle("GET /api/v1/products/{id}/image", authenticator.Authenticate(http.HandlerFunc(imageHandler.GetByProductID)))
	mux.Handle("GET /api/v1/images/{id}", authenticator.Authenticate(http.HandlerFunc(imageHandler.GetByID)))
	mux.Handle("PUT /api/v1/images/{id}", authenticator.Authenticate(http.HandlerFunc(imageHandler.Update)))
	mux.Handle("DELETE /api/v1/images/{id}", authenticator.Authenticate(http.HandlerFunc(imageHandler.Delete)))

	return http.ListenAndServe(":8080", mux)
}

// @title SHOP
// @version 1.0
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите Bearer <access_token>, используя access_token из ответа POST /auth.
func main() {
	authaddress := os.Getenv("AUTH_ADDRESS")
	conn, err := grpc.NewClient(authaddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fx.New(
		fx.Supply(conn),
		fx.Provide(
			newAuthClient,
			middleware.NewUserAuthenticator,

			handler.NewAuthHandler,
			handler.NewAddressHandler,
			handler.NewClientHandler,
			handler.NewProductHandler,
			handler.NewSupplierHandler,
			handler.NewImageHandler,

			repository.NewPostgresPool,
			repository.NewAddressRepository,
			repository.NewClientRepository,
			repository.NewProductRepository,
			repository.NewSupplierRepository,
			repository.NewImageRepository,
		),
		fx.Invoke(startServer),
	).Run()
}
