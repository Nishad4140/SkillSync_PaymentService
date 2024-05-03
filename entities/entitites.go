package entities

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	Id         uuid.UUID `gorm:"primayKey" json:"id,omitempty"`
	UserId     uuid.UUID `json:"user_id,omitempty"`
	PaymentRef string    `json:"payment_ref,omitempty"`
	Time       time.Time `json:"time,omitempty"`
	ProjectId  uuid.UUID `json:"project_id,omitempty"`
}

type AdminAccount struct {
	Id        uuid.UUID `gorm:"primayKey" json:"id,omitempty"`
	Amount    float64
	PaymentId uuid.UUID
	Payment   Payment `gorm:"foreignKey:PaymentId"`
}

type FreelancerAccount struct {
	Id           uuid.UUID `gorm:"primayKey" json:"id,omitempty"`
	Amount       float64
	PaymentId    uuid.UUID
	Payment      Payment `gorm:"foreignKey:PaymentId"`
	FreelancerId uuid.UUID
}
