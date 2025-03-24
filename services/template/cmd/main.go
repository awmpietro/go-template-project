package main

import (
	"log"
	"net/http"
	"os"

	handlers "github.com/nuhorizon/go-project-template/services/template/internal/delivery/handlers"
	routes "github.com/nuhorizon/go-project-template/services/template/internal/delivery/routes"
	firebase "github.com/nuhorizon/go-project-template/services/template/internal/infra/firebase"
	pg "github.com/nuhorizon/go-project-template/services/template/internal/infra/postgres"
	"github.com/nuhorizon/go-project-template/services/template/internal/ports/infrastructure"
	servicesPorts "github.com/nuhorizon/go-project-template/services/template/internal/ports/services"
	pgRepositories "github.com/nuhorizon/go-project-template/services/template/internal/repository/postgres"
	"github.com/nuhorizon/go-project-template/services/template/internal/services"
	"github.com/nuhorizon/go-project-template/services/template/internal/usecases"
	middlewares "github.com/nuhorizon/go-project-template/services/template/pkg/middlewares"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	// _ "github.com/nuhorizon/go-project-template/services/template/docs" //uncomment when swagger is ready
)

// @title template API
// @version 1.0.0
// @contact.name Arthur Mastropietro <amcod3>
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath /

func main() {
	err := Run()
	if err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
		return err
	}

	// Initialize Postgres
	db := pg.NewPGSql()
	if err = db.InitDB(); err != nil {
		log.Println("failed to initialize database:", err)
		return err
	}
	defer db.CloseDB()

	// Initialize Firebase
	firebaseClient, err := firebase.NewFirebaseClient()
	if err != nil {
		log.Println("Failed to initialize Firebase:", err)
		return err
	}

	firebaseService := services.NewFirebaseAuthService(firebaseClient)

	// Initialize JWT Service
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("JWT_SECRET is missing from environment variables")
		return err
	}
	jwtService := services.NewJWTService(jwtSecret)

	// Router
	mux := chi.NewRouter()
	initializeMux(mux)

	// Dependency Injection for Handlers
	setupHandlers(mux, db, firebaseService, jwtService)

	log.Printf("🚀 CatWise API running on port %s", os.Getenv("PORT"))

	if err = http.ListenAndServe(":"+os.Getenv("PORT"), mux); err != nil {
		log.Println("failed to start server:", err)
		return err
	}

	return nil
}

func initializeMux(mux *chi.Mux) {
	mux.Use(chiMiddleware.RequestID)
	mux.Use(chiMiddleware.RealIP)
	mux.Use(chiMiddleware.Logger)
	mux.Use(chiMiddleware.Recoverer)
	mux.Use(chiMiddleware.Heartbeat("/ping"))
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Swagger
	mux.Route("/swagger", func(r chi.Router) {
		r.Use(middlewares.BasicAuthMiddleware)
		r.Get("/*", httpSwagger.Handler(
			httpSwagger.URL(os.Getenv("LOCALHOST")+"/swagger/doc.json"),
		))
	})
}

func setupHandlers(
	mux *chi.Mux,
	db infrastructure.SQLConnector,
	firebaseService servicesPorts.FirebaseAuthService,
	jwtService servicesPorts.JWTService,
) {
	// Repositórios
	userRepo := pgRepositories.NewUserPostgres(db.GetDB())

	// Use Cases
	authUseCase := usecases.NewAuthUseCase(userRepo, firebaseService, jwtService)

	// Handlers
	authHandler := handlers.NewAuthHandler(authUseCase)

	// Registro de rotas
	routes.RegisterAuthRoutes(mux, authHandler)

}
