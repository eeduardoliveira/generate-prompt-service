package http

import (
	"context"
	"encoding/json"
	"generate-prompt-service/app/domain"
	"generate-prompt-service/app/usecase"
	"net/http"
)

type PromptHandler struct {
	UseCase *usecase.GeneratePromptUseCase
}

func NewPromptHandler(uc *usecase.GeneratePromptUseCase) *PromptHandler {
	return &PromptHandler{UseCase: uc}
}

func (h *PromptHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Método não permitido"}`, http.StatusMethodNotAllowed)
		return
	}

	var input domain.Prompt
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"JSON inválido"}`, http.StatusBadRequest)
		return
	}

	if input.ClienteID == "" || input.NomeEmpresa == "" || input.SystemPrompt == "" {
		http.Error(w, `{"error":"Campos obrigatórios ausentes"}`, http.StatusBadRequest)
		return
	}

	bucket := r.URL.Query().Get("bucket")
	if bucket == "" {
		http.Error(w, `{"error":"bucket é obrigatório"}`, http.StatusBadRequest)
		return
	}

	err := h.UseCase.Execute(context.Background(), bucket, input)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Prompt gerado com sucesso"}`))
}