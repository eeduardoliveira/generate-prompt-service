package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"generate-prompt-service/app/usecase"
	"generate-prompt-service/dependencies/signer"
	"generate-prompt-service/dependencies/uploader"
	handler "generate-prompt-service/presentation/http"

	_ "generate-prompt-service/docs" // Swagger docs (gerado pelo swag init)

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger" // Swagger handler

	usecaseagent "generate-prompt-service/app/usecase/agentnumber"
	"generate-prompt-service/dependencies/database"
)

// @title Generate Prompt Service API
// @version 1.0
// @description API responsável por gerar prompts e fazer upload no bucket.
// @termsOfService http://swagger.io/terms/

// @contact.name Suporte Syphertech
// @contact.email suporte@syphertech.com.br

// @host generate-prompt.syphertech.com.br
// @BasePath /
func main() {
	// Carrega variáveis de ambiente locais, se existirem
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env não encontrado, seguindo com variáveis de ambiente do sistema.")
	}

	// Define a porta do servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	// Instancia dependências
	signerClient := signer.NewHTTPClient()
	uploaderClient := uploader.NewHTTPUploader()
	useCase := &usecase.GeneratePromptUseCase{
		Signer:   signerClient,
		Uploader: uploaderClient,
	}

	dbConn, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Erro ao conectar no banco de dados: %v", err)
	}
	defer dbConn.Close(context.Background())

	agentRepo := database.NewRepository(dbConn)
	agentUseCase := usecaseagent.NewCreateAgentNumberUseCase(agentRepo)

	// Handlers
	promptHandler := handler.NewPromptHandler(useCase)
	agentNumberHandler := handler.NewAgentNumberHandler(agentUseCase)

	// Rotas
	http.HandleFunc("/api/v1/generate-prompt", promptHandler.Handle)
	http.HandleFunc("/api/v1/agent-number", agentNumberHandler.Create)
	// Rota do Swagger
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)


		// Middleware de CORS
		corsHandler := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Authorization", "Content-Type"},
			AllowCredentials: true,
		}).Handler(http.DefaultServeMux)

	log.Printf("🚀 Servidor rodando na porta %s...", port)
	if err := http.ListenAndServe(":"+port, corsHandler); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}