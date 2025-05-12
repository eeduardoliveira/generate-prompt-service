package http

import (
	"context"
	"encoding/json"
	"generate-prompt-service/app/domain"
	"generate-prompt-service/app/domain/agentnumber"
	"generate-prompt-service/app/usecase"
	usecaseagent "generate-prompt-service/app/usecase/agentnumber"
	"log"
	"net/http"
)

type PromptHandler struct {
	UseCase *usecase.GeneratePromptUseCase
}

func NewPromptHandler(uc *usecase.GeneratePromptUseCase) *PromptHandler {
	return &PromptHandler{UseCase: uc}
}

// Handle godoc
// @Summary Geração de Prompt
// @Description Recebe os dados de prompt e realiza o upload em bucket S3 com metadados automáticos
// @Tags prompts
// @Accept json
// @Produce json
// @Param bucket query string true "Nome do bucket"
// @Param prompt body domain.Prompt true "Dados do Prompt"
// @Success 201 {object} map[string]string "Prompt gerado com sucesso"
// @Failure 400 {object} map[string]string "Requisição inválida ou campos obrigatórios ausentes"
// @Failure 405 {object} map[string]string "Método HTTP não permitido"
// @Failure 500 {object} map[string]string "Erro interno ao processar prompt"
// @Router /api/v1/generate-prompt [post]
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

type AgentNumberHandler struct {
	UseCase *usecaseagent.CreateAgentNumberUseCase
}

func NewAgentNumberHandler(uc *usecaseagent.CreateAgentNumberUseCase) *AgentNumberHandler {
	return &AgentNumberHandler{UseCase: uc}
}

// Create godoc
// @Summary Criar AgentNumber
// @Description Insere um novo mapeamento client_id e número de WhatsApp
// @Tags agent-number
// @Accept json
// @Produce json
// @Param data body agentnumber.AgentNumber true "Dados do agente"
// @Success 201 {string} string "Criado com sucesso"
// @Failure 400 {object} map[string]string "Payload inválido"
// @Failure 500 {object} map[string]string "Erro interno"
// @Router /api/v1/agent-number [post]
func (h *AgentNumberHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Método não permitido"}`, http.StatusMethodNotAllowed)
		return
	}


	var an agentnumber.AgentNumber
	if err := json.NewDecoder(r.Body).Decode(&an); err != nil {
		http.Error(w, `{"error":"Payload inválido"}`, http.StatusBadRequest)
		return
	}
	log.Println(an)

	if an.ClientID == "" || an.WhatsappNumber == "" || an.CustomerID == 0 {
		http.Error(w, `{"error":"Campos obrigatórios ausentes"}`, http.StatusBadRequest)
		return
	}

	if err := h.UseCase.Execute(r.Context(), an); err != nil {
		http.Error(w, `{"error":"Erro interno"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Criado com sucesso"}`))
}