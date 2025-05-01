package domain

import "time"

type Servico struct {
	Nome      string  `json:"nome"`
	Descricao string  `json:"descricao"`
	Preco     float64 `json:"preco"`
}

type Metadados struct {
	DataCriacao       time.Time `json:"dataCriacao"`
	UltimaAtualizacao time.Time `json:"ultimaAtualizacao"`
	Origem            string    `json:"origem"`
}

type Prompt struct {
	ClienteID    string     `json:"clienteID"`
	NomeEmpresa  string     `json:"nomeEmpresa"`
	SystemPrompt string     `json:"systemPrompt"`
	Servicos     []Servico  `json:"servicos"`
	Idiomas      []string   `json:"idiomas"`
	ModeloIA     string     `json:"modeloIA"`
	Metadados    Metadados  `json:"metadados"`
	VersaoPrompt int        `json:"versaoPrompt"`
}