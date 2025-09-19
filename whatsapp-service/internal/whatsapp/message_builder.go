package whatsapp

import (
	"strings"

	"whatsapp-service/internal/models"
)

// MessageBuilder builds different types of WhatsApp messages
type MessageBuilder struct{}

// NewMessageBuilder creates a new message builder
func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{}
}

// BuildTextMessage creates a text message
func (mb *MessageBuilder) BuildTextMessage(to, body string) *models.WhatsAppMessage {
	return &models.WhatsAppMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               CleanPhoneNumber(to),
		Type:             "text",
		Text: &models.TextMessage{
			PreviewURL: true,
			Body:       body,
		},
	}
}

// BuildMediaMessage creates a media message
func (mb *MessageBuilder) BuildMediaMessage(to string, payload *models.NotificationPayload) *models.WhatsAppMessage {
	message := &models.WhatsAppMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               CleanPhoneNumber(to),
		Type:             string(payload.MediaType),
	}

	caption := mb.renderTemplate(payload.TemplateBody, payload.Params)
	mediaMsg := &models.MediaMessage{
		Link:    payload.MediaURL,
		Caption: caption,
	}

	switch payload.MediaType {
	case models.MediaTypeImage:
		message.Image = mediaMsg
	case models.MediaTypeDocument:
		message.Document = mediaMsg
		if filename, ok := payload.Params["filename"]; ok {
			message.Document.Filename = filename
		}
	case models.MediaTypeAudio:
		message.Audio = mediaMsg
	case models.MediaTypeVideo:
		message.Video = mediaMsg
	default:
		message.Type = "image"
		message.Image = mediaMsg
	}

	return message
}

// BuildTemplateMessage creates a template message
func (mb *MessageBuilder) BuildTemplateMessage(to string, payload *models.NotificationPayload) *models.WhatsAppMessage {
	components := []models.TemplateComponent{}

	// Add body parameters if available
	if len(payload.Params) > 0 {
		parameters := make([]models.TemplateParameter, 0, len(payload.Params))
		for _, value := range payload.Params {
			parameters = append(parameters, models.TemplateParameter{
				Type: "text",
				Text: value,
			})
		}

		components = append(components, models.TemplateComponent{
			Type:       "body",
			Parameters: parameters,
		})
	}

	return &models.WhatsAppMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               CleanPhoneNumber(to),
		Type:             "template",
		Template: &models.TemplateMessage{
			Name: payload.TemplateName,
			Language: models.TemplateLanguage{
				Code: GetLanguageCode(payload.Locale),
			},
			Components: components,
		},
	}
}

// BuildFromPayload builds appropriate message based on payload type
func (mb *MessageBuilder) BuildFromPayload(payload *models.NotificationPayload) *models.WhatsAppMessage {
	switch {
	case payload.MediaURL != "":
		return mb.BuildMediaMessage(payload.To, payload)
	case payload.TemplateName != "":
		return mb.BuildTemplateMessage(payload.To, payload)
	default:
		body := mb.renderTemplate(payload.TemplateBody, payload.Params)
		return mb.BuildTextMessage(payload.To, body)
	}
}

// renderTemplate replaces template variables with actual values
func (mb *MessageBuilder) renderTemplate(template string, params map[string]string) string {
	result := template
	for key, value := range params {
		placeholder := "{" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}