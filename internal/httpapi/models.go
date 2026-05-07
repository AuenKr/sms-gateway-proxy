package httpapi

type TextMessage struct {
	Text string `json:"text"`
}

type DataMessage struct {
	Data string `json:"data"`
	Port int    `json:"port"`
}

type SendSMSRequest struct {
	DataMessage        *DataMessage `json:"dataMessage,omitempty"`
	DeviceID           string       `json:"deviceId,omitempty"`
	ID                 string       `json:"id,omitempty"`
	Message            string       `json:"message,omitempty"`
	PhoneNumbers       []string     `json:"phoneNumbers"`
	Priority           *int         `json:"priority,omitempty"`
	ScheduleAt         string       `json:"scheduleAt,omitempty"`
	SimNumber          *int         `json:"simNumber,omitempty"`
	TextMessage        *TextMessage `json:"textMessage,omitempty"`
	TTL                *int         `json:"ttl,omitempty"`
	ValidUntil         string       `json:"validUntil,omitempty"`
	WithDeliveryReport *bool        `json:"withDeliveryReport,omitempty"`
}

type GatewaySendRequest struct {
	DataMessage        *DataMessage `json:"dataMessage,omitempty"`
	DeviceID           string       `json:"deviceId,omitempty"`
	ID                 string       `json:"id,omitempty"`
	IsEncrypted        bool         `json:"isEncrypted"`
	PhoneNumbers       []string     `json:"phoneNumbers"`
	Priority           *int         `json:"priority,omitempty"`
	ScheduleAt         string       `json:"scheduleAt,omitempty"`
	SimNumber          *int         `json:"simNumber,omitempty"`
	TextMessage        *TextMessage `json:"textMessage,omitempty"`
	TTL                *int         `json:"ttl,omitempty"`
	ValidUntil         string       `json:"validUntil,omitempty"`
	WithDeliveryReport *bool        `json:"withDeliveryReport,omitempty"`
}
