package main

import ( // gin-swagger middleware

	"github.com/ardi7923/go-hello-api/config"
	"github.com/ardi7923/go-hello-api/controllers"
	"github.com/ardi7923/go-hello-api/docs"
	_ "github.com/ardi7923/go-hello-api/docs"
	"github.com/ardi7923/go-hello-api/middlewares"
	"github.com/ardi7923/go-hello-api/repositories"
	"github.com/ardi7923/go-hello-api/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
) // swagger embed files)

var (
	db             *gorm.DB                    = config.SetUpDatabaseConnection()
	userRepository repositories.UserRepository = repositories.NewUserRepository(db)
	jwtService     services.JWTService         = services.NewJWTService()
	authService    services.AuthService        = services.NewAuthService(userRepository)
	userService    services.UserService        = services.NewUserService(userRepository)
	authController controllers.AuthController  = controllers.NewAuthController(authService, jwtService)
	userController controllers.UserController  = controllers.NewUserController(userService, jwtService)
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {

	// Swagger Meta Information
	docs.SwaggerInfo.Title = "Hello API Documentation"
	docs.SwaggerInfo.Description = "Contoh Api dengan golang bosku"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8000"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Failed To Load ENV")
	}

	defer config.CloseDatabaseConnection(db)
	router := gin.New()
	baseUrl := "api/v1"

	// authentication Route
	authRoutes := router.Group(baseUrl + "/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	// User Profile Route
	userRoutes := router.Group(baseUrl+"/user", middlewares.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("profile", userController.Profile)
		userRoutes.PUT("profile", userController.Update)
	}

	url := ginSwagger.URL("http://localhost:8000/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	router.Run(":8000")
}
