package main

import (
	"log"

	"mini-project-ostore/internal/config"
	"mini-project-ostore/internal/handler"
	"mini-project-ostore/internal/middleware"
	"mini-project-ostore/internal/repository"
	"mini-project-ostore/internal/usecase"
	"mini-project-ostore/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// ------------------------
	// LOAD CONFIG & DB
	// ------------------------
	cfg := config.LoadConfig()

	db, err := database.NewMySQLConnection(cfg.Database)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	database.Migrate(db)

	// ------------------------
	// INITIALIZE REPOSITORIES
	// ------------------------
	userRepo := repository.NewUserRepository(db)
	storeRepo := repository.NewStoreRepository(db)
	addressRepo := repository.NewAddressRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Region API repositories
	provinceAPIRepo := repository.NewProvinceAPIRepository()
	cityAPIRepo := repository.NewCityAPIRepository()
	subdistrictAPIRepo := repository.NewSubdistrictAPIRepository()

	// ------------------------
	// INITIALIZE USECASES
	// ------------------------
	userUC := usecase.NewUserUseCase(userRepo, storeRepo)
	authUC := usecase.NewAuthUseCase(userRepo)
	storeUC := usecase.NewStoreUseCase(storeRepo, userRepo)
	addressUC := usecase.NewAddressUseCase(addressRepo, userRepo)
	categoryUC := usecase.NewCategoryUseCase(categoryRepo)
	productUC := usecase.NewProductUseCase(productRepo, storeRepo, userRepo, categoryRepo)
	transactionUC := usecase.NewTransactionUseCase(transactionRepo, productRepo, userRepo, addressRepo)
	regionUC := usecase.NewRegionUseCase(provinceAPIRepo, cityAPIRepo, subdistrictAPIRepo)

	// ------------------------
	// INITIALIZE HANDLERS
	// ------------------------
	userHandler := handler.NewUserHandler(userUC)
	authHandler := handler.NewAuthHandler(authUC, userUC)
	storeHandler := handler.NewStoreHandler(storeUC)
	addressHandler := handler.NewAddressHandler(addressUC)
	categoryHandler := handler.NewCategoryHandler(categoryUC)
	productHandler := handler.NewProductHandler(productUC)
	transactionHandler := handler.NewTransactionHandler(transactionUC)
	regionHandler := handler.NewRegionHandler(regionUC)

	// ------------------------
	// MIDDLEWARE
	// ------------------------
	authMiddleware := middleware.NewAuthMiddleware(userRepo)

	// ------------------------
	// SETUP ROUTER
	// ------------------------
	r := gin.Default()

	// =========================================================
	// 🔓 PUBLIC ROUTES
	// =========================================================
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	// Category (public)
	r.GET("/category", categoryHandler.GetCategories)
	r.GET("/category/:id", categoryHandler.GetCategoryByID)

	// Product (public)
	r.GET("/product", productHandler.GetProducts)
	r.GET("/product/:id", productHandler.GetProductByID)

	// Region (public)
	r.GET("/regions/provinces", regionHandler.GetProvinces)
	r.GET("/regions/provinces/:provinceID/cities", regionHandler.GetCities)
	r.GET("/regions/cities/:cityID/subdistricts", regionHandler.GetSubdistricts)

	// =========================================================
	// 🔒 PROTECTED ROUTES (need login)
	// =========================================================
	protected := r.Group("")
	protected.Use(authMiddleware.ValidateToken())
	{
		// USER PROFILE
		userGroup := protected.Group("/user")
		{
			userGroup.GET("", userHandler.GetUserProfile)
			userGroup.PUT("", userHandler.UpdateUserProfile)

			// ADDRESS
			userGroup.GET("/alamat", addressHandler.GetUserAddresses)
			userGroup.GET("/alamat/:id", addressHandler.GetAddress)
			userGroup.POST("/alamat", addressHandler.CreateAddress)
			userGroup.PUT("/alamat/:id", addressHandler.UpdateAddress)
			userGroup.DELETE("/alamat/:id", addressHandler.DeleteAddress)
		}

		// STORE (TOKO)
		storeGroup := protected.Group("/toko")
		{
			storeGroup.GET("", storeHandler.GetStores)
			storeGroup.GET("/my", storeHandler.GetMyStore)
			storeGroup.GET("/:id_toko", storeHandler.GetStoreByID)
			storeGroup.PUT("/:id_toko", storeHandler.UpdateStore)
		}

		// TRANSACTION
		transactionGroup := protected.Group("/transaction")
		{
			transactionGroup.POST("", transactionHandler.CreateTransaction)
			transactionGroup.GET("", transactionHandler.GetUserTransactions)
			transactionGroup.GET("/:id", transactionHandler.GetTransaction)
		}

		// =====================================================
		// 🛡️ ADMIN ROUTES (require is_admin = true)
		// =====================================================
		admin := protected.Group("")
		admin.Use(authMiddleware.RequireAdmin())
		{
			admin.POST("/category", categoryHandler.CreateCategory)
			admin.PUT("/category/:id", categoryHandler.UpdateCategory)
			admin.DELETE("/category/:id", categoryHandler.DeleteCategory)
		}
	}

	// =========================================================
	// 🚀 START SERVER
	// =========================================================
	r.Run(":" + cfg.Server.Port)
}
