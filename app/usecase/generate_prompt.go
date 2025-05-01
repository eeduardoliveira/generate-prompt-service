package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"generate-prompt-service/app/domain"
	"generate-prompt-service/dependencies/signer"
	"generate-prompt-service/dependencies/uploader"
	"time"
)

type GeneratePromptUseCase struct {
	Signer   signer.Client
	Uploader uploader.Client
}

func (uc *GeneratePromptUseCase) Execute(ctx context.Context, bucket string, input domain.Prompt) error {
	input.Metadados = domain.Metadados{
		DataCriacao:       time.Now().UTC(),
		UltimaAtualizacao: time.Now().UTC(),
		Origem:            "generate-prompt-service",
	}
	input.VersaoPrompt = 1

	jsonData, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar prompt para JSON: %w", err)
	}

	signedURL, err := uc.Signer.GetPresignedUploadURL(ctx, bucket, input.ClienteID)
	if err != nil {
		return fmt.Errorf("erro ao obter URL assinada: %w", err)
	}

	err = uc.Uploader.UploadJSON(ctx, signedURL, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("erro ao fazer upload do JSON: %w", err)
	}

	return nil
}