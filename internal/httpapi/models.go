package httpapi

type TextMessage struct {
	Text string `json:"text" validate:"required,notblank"`
}

type DataMessage struct {
	Data string `json:"data" validate:"required,notblank"`
	Port int    `json:"port" validate:"gte=0,lte=65535"`
}

type SendSMSRequest struct {
	DataMessage        *DataMessage `json:"dataMessage,omitempty" validate:"required_without_all=Message TextMessage,omitempty,excluded_with=Message TextMessage"`
	DeviceID           string       `json:"deviceId,omitempty"`
	ID                 string       `json:"id,omitempty"`
	Message            string       `json:"message,omitempty" validate:"required_without_all=TextMessage DataMessage,omitempty,notblank,excluded_with=TextMessage DataMessage"`
	PhoneNumbers       []string     `json:"phoneNumbers" validate:"required,min=1,dive,required,notblank"`
	Priority           *int         `json:"priority,omitempty"`
	ScheduleAt         string       `json:"scheduleAt,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	SimNumber          *int         `json:"simNumber,omitempty"`
	TextMessage        *TextMessage `json:"textMessage,omitempty" validate:"required_without_all=Message DataMessage,omitempty,excluded_with=Message DataMessage"`
	TTL                *int         `json:"ttl,omitempty" validate:"omitempty,excluded_with=ValidUntil"`
	ValidUntil         string       `json:"validUntil,omitempty" validate:"omitempty,excluded_with=TTL,datetime=2006-01-02T15:04:05Z07:00"`
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
