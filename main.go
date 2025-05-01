package main

import (
	"log"
	"net/http"
	"os"

	"generate-prompt-service/app/usecase"
	"generate-prompt-service/dependencies/signer"
	"generate-prompt-service/dependencies/uploader"
	handler "generate-prompt-service/presentation/http"

	"github.com/joho/godotenv"
)

func main() {
		// Carrega variáveis de ambiente
		if err := godotenv.Load(); err != nil {
			log.Println("⚠️  Aviso: .env não encontrado, usando variáveis do ambiente")
		}
		
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

	log.Printf("🚀 Servidor rodando na porta %s...", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}