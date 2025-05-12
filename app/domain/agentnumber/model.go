package agentnumber


type AgentNumber struct {
    ID             int64  `json:"id,omitempty"`
    ClientID       string `json:"client_id"`
    WhatsappNumber string `json:"whatsapp_number"`
    Description    string `json:"description,omitempty"`
    CustomerID     int64  `json:"customer_id"`
}