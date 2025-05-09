package main

import (
	"log"
	"net/http"
	"os"

	"generate-prompt-service/app/usecase"
	"generate-prompt-service/dependencies/signer"
	"generate-prompt-service/dependencies/uploader"
	handler "generate-prompt-service/presentation/http"

	httpSwagger "github.com/swaggo/http-swagger"
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
		// Carrega variáveis de ambiente
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	signerClient := signer.NewHTTPClient()
	uploaderClient := uploader.NewHTTPUploader()
	useCase := &usecase.GeneratePromptUseCase{
		Signer:   signerClient,
		Uploader: uploaderClient,
	}

	h := handler.NewPromptHandler(useCase)

	http.HandleFunc("/generate-prompt", h.Handle)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)


	log.Printf("🚀 Servidor rodando na porta %s...", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}