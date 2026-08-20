package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"backend-go/domain"
	"backend-go/internal/handler"
	"backend-go/internal/repository"
	"backend-go/internal/service"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, membaca dari environment system")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL belum diatur di dalam file .env")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	log.Println("Berhasil terhubung ke database Neon!")

	err = db.AutoMigrate(
		&domain.User{}, &domain.UserAddress{}, &domain.Store{},
		&domain.ProductCategory{}, &domain.Product{}, &domain.ProductImage{}, &domain.ProductStock{}, &domain.Cart{},
		&domain.Order{}, &domain.OrderItem{}, &domain.StockJournal{},
	)
	if err != nil {
		log.Printf("Peringatan: Gagal melakukan AutoMigrate: %v", err)
	} else {
		log.Println("AutoMigrate berhasil dijalankan!")
	}

	checkoutRepo := repository.NewCheckoutRepository(db)
	checkoutSvc := service.NewCheckoutService(checkoutRepo)
	checkoutHandler := handler.NewCheckoutHandler(checkoutSvc)
	cartRepo := repository.NewCartRepository(db)
	cartSvc := service.NewCartService(cartRepo, checkoutRepo)
	cartHandler := handler.NewCartHandler(cartSvc)
	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productSvc)
	addressRepo := repository.NewAddressRepository(db)
	addressSvc := service.NewAddressService(addressRepo)
	addressHandler := handler.NewAddressHandler(addressSvc)
	orderRepo := repository.NewOrderRepository(db)
	orderSvc := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderSvc)
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)
	authSvc := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authSvc)

	e := echo.New()

	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api/v1")

	mockAuthMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// UUID dari tabel neon
			c.Set("user_id", "edbbc9ed-f0f3-4e61-a1af-d4592f02e541")
			return next(c)
		}
	}

	api.POST("/checkout", checkoutHandler.CreateOrder, mockAuthMiddleware)

	api.GET("/checkout", checkoutHandler.GetUserOrders, mockAuthMiddleware)

	api.PUT("/checkout/:id/cancel", checkoutHandler.CancelOrder, mockAuthMiddleware)
	api.GET("/cart", cartHandler.GetCart, mockAuthMiddleware)
	api.POST("/cart", cartHandler.AddToCart, mockAuthMiddleware)
	api.PUT("/cart/:id", cartHandler.UpdateCart, mockAuthMiddleware)
	api.DELETE("/cart/:id", cartHandler.RemoveCartItem, mockAuthMiddleware)
	api.GET("/products", productHandler.GetProducts)
	api.GET("/products/:id", productHandler.GetProductDetail)
	api.GET("/addresses", addressHandler.GetAddresses, mockAuthMiddleware)
	api.POST("/addresses", addressHandler.AddAddress, mockAuthMiddleware)
	api.PUT("/addresses/:id/primary", addressHandler.SetPrimary, mockAuthMiddleware)
	api.DELETE("/addresses/:id", addressHandler.DeleteAddress, mockAuthMiddleware)
	api.GET("/orders", orderHandler.GetUserOrders, mockAuthMiddleware)
	api.GET("/orders/:id", orderHandler.GetOrderDetail, mockAuthMiddleware)
	api.GET("/profile", userHandler.GetProfile, mockAuthMiddleware)
	api.PUT("/profile", userHandler.UpdateProfile, mockAuthMiddleware)
	api.POST("/api/auth/register", authHandler.Register)
	api.POST("/api/auth/login", authHandler.Login)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server Backend E-Commerce (Golang) jalan!")
	})

	log.Println("Server API berjalan di port 8080...")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
