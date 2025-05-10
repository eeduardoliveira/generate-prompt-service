package main

import (
	"log"
	"net/http"
	"os"

	"generate-prompt-service/app/usecase"
	"generate-prompt-service/dependencies/signer"
	"generate-prompt-service/dependencies/uploader"
	handler "generate-prompt-service/presentation/http"

	_ "generate-prompt-service/docs"               // Swagger docs (gerado pelo swag init)
	httpSwagger "github.com/swaggo/http-swagger"  // Swagger handler
	"github.com/joho/godotenv"
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

	// Handlers
	h := handler.NewPromptHandler(useCase)

	// Rotas
	http.HandleFunc("/generate-prompt", h.Handle)
	
	// Rota do Swagger
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Printf("🚀 Servidor rodando na porta %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}