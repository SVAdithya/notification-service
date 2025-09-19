package models

// WhatsAppMessage represents a WhatsApp message request
type WhatsAppMessage struct {
	MessagingProduct string           `json:"messaging_product"`
	RecipientType    string           `json:"recipient_type"`
	To               string           `json:"to"`
	Type             string           `json:"type"`
	Text             *TextMessage     `json:"text,omitempty"`
	Image            *MediaMessage    `json:"image,omitempty"`
	Document         *MediaMessage    `json:"document,omitempty"`
	Audio            *MediaMessage    `json:"audio,omitempty"`
	Video            *MediaMessage    `json:"video,omitempty"`
	Template         *TemplateMessage `json:"template,omitempty"`
}

// TextMessage represents a text message
type TextMessage struct {
	PreviewURL bool   `json:"preview_url"`
	Body       string `json:"body"`
}

// MediaMessage represents a media message (image, document, audio, video)
type MediaMessage struct {
	ID       string `json:"id,omitempty"`
	Link     string `json:"link,omitempty"`
	Caption  string `json:"caption,omitempty"`
	Filename string `json:"filename,omitempty"`
}

// TemplateMessage represents a template message
type TemplateMessage struct {
	Name       string              `json:"name"`
	Language   TemplateLanguage    `json:"language"`
	Components []TemplateComponent `json:"components,omitempty"`
}

// TemplateLanguage represents template language settings
type TemplateLanguage struct {
	Code string `json:"code"`
}

// TemplateComponent represents a template component
type TemplateComponent struct {
	Type       string              `json:"type"`
	Parameters []TemplateParameter `json:"parameters,omitempty"`
}

// TemplateParameter represents a template parameter
type TemplateParameter struct {
	Type  string        `json:"type"`
	Text  string        `json:"text,omitempty"`
	Image *MediaMessage `json:"image,omitempty"`
}

// WhatsAppResponse represents the API response
type WhatsAppResponse struct {
	Messages []WhatsAppMessageStatus `json:"messages"`
	Contacts []WhatsAppContact       `json:"contacts"`
}

// WhatsAppMessageStatus represents message status in response
type WhatsAppMessageStatus struct {
	ID string `json:"id"`
}

// WhatsAppContact represents contact information in response
type WhatsAppContact struct {
	Input string `json:"input"`
	WaID  string `json:"wa_id"`
}

// WhatsAppError represents API error response
type WhatsAppError struct {
	Error WhatsAppErrorDetail `json:"error"`
}

// WhatsAppErrorDetail represents error details
type WhatsAppErrorDetail struct {
	Message   string                 `json:"message"`
	Type      string                 `json:"type"`
	Code      int                    `json:"code"`
	ErrorData WhatsAppErrorData      `json:"error_data"`
}

// WhatsAppErrorData represents additional error data
type WhatsAppErrorData struct {
	Details string `json:"details"`
}
