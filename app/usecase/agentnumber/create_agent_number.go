package agentnumber

import (
	"context"
	"generate-prompt-service/app/domain/agentnumber"
)

type Repository interface {
	InsertAgentNumber(ctx context.Context, an agentnumber.AgentNumber) error
}

type CreateAgentNumberUseCase struct {
	repo Repository
}

func NewCreateAgentNumberUseCase(repo Repository) *CreateAgentNumberUseCase {
	return &CreateAgentNumberUseCase{repo: repo}
}

func (uc *CreateAgentNumberUseCase) Execute(ctx context.Context, an agentnumber.AgentNumber) error {
	return uc.repo.InsertAgentNumber(ctx, an)
}