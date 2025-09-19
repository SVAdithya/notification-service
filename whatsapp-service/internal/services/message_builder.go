package services

import (
	"whatsapp-service/internal/models"
	"whatsapp-service/internal/utils"
)

// MessageBuilder handles creation of different WhatsApp message types
type MessageBuilder struct{}

// NewMessageBuilder creates a new message builder
func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{}
}

// BuildMessage creates a WhatsApp message based on notification payload
func (mb *MessageBuilder) BuildMessage(payload *models.NotificationPayload, to string) (*models.WhatsAppMessage, error) {
	switch {
	case payload.MediaURL != "":
		return mb.buildMediaMessage(to, payload)
	case payload.TemplateName != "":
		return mb.buildTemplateMessage(to, payload)
	default:
		return mb.buildTextMessage(to, payload)
	}
}

// buildTextMessage creates a text message
func (mb *MessageBuilder) buildTextMessage(to string, payload *models.NotificationPayload) (*models.WhatsAppMessage, error) {
	body := utils.RenderTemplate(payload.TemplateBody, payload.Params)
	if body == "" {
		return nil, models.ErrMissingContent
	}

	return &models.WhatsAppMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "text",
		Text: &models.TextMessage{
			PreviewURL: true,
			Body:       body,
		},
	}, nil
}

// buildMediaMessage creates a media message (image, document, audio, video)
func (mb *MessageBuilder) buildMediaMessage(to string, payload *models.NotificationPayload) (*models.WhatsAppMessage, error) {
	if payload.MediaURL == "" {
		return nil, models.ErrMissingContent
	}

	message := &models.WhatsAppMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             string(payload.MediaType),
	}

	mediaMsg := &models.MediaMessage{
		Link:    payload.MediaURL,
		Caption: utils.RenderTemplate(payload.TemplateBody, payload.Params),
	}

	// Add filename for documents if provided
	if payload.MediaType == models.MediaTypeDocument {
		if filename, exists := payload.Params["filename"]; exists {
			mediaMsg.Filename = filename
		}
	}

	// Set the appropriate media field based on type
	switch payload.MediaType {
	case models.MediaTypeImage:
		message.Image = mediaMsg
	case models.MediaTypeDocument:
		message.Document = mediaMsg
	case models.MediaTypeAudio:
		message.Audio = mediaMsg
	case models.MediaTypeVideo:
		message.Video = mediaMsg
	default:
		// Default to image if invalid type
		message.Type = "image"
		message.Image = mediaMsg
	}

	return message, nil
}

// buildTemplateMessage creates a template message
func (mb *MessageBuilder) buildTemplateMessage(to string, payload *models.NotificationPayload) (*models.WhatsAppMessage, error) {
	if payload.TemplateName == "" {
		return nil, models.ErrInvalidTemplate
	}

	message := &models.WhatsAppMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "template",
		Template: &models.TemplateMessage{
			Name: payload.TemplateName,
			Language: models.TemplateLanguage{
				Code: utils.GetLanguageCode(payload.Locale),
			},
		},
	}

	// Add parameters if available
	if len(payload.Params) > 0 {
		var parameters []models.TemplateParameter
		for _, value := range payload.Params {
			parameters = append(parameters, models.TemplateParameter{
				Type: "text",
				Text: value,
			})
		}

		message.Template.Components = []models.TemplateComponent{
			{
				Type:       "body",
				Parameters: parameters,
			},
		}
	}

	return message, nil
}