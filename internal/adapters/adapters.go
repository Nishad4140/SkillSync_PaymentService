package adapters

import (
	"github.com/Nishad4140/SkillSync_PaymentService/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentAdapter struct {
	DB *gorm.DB
}

func NewPaymentAdapter(db *gorm.DB) *PaymentAdapter {
	return &PaymentAdapter{
		DB: db,
	}
}

func (payment *PaymentAdapter) AddPayment(req entities.Payment) (string, error) {
	id := uuid.New()
	insertQuery := `INSERT INTO payments (id, user_id, payment_ref, time, project_id) VALUES ($1, $2, $3, NOW(), $4)`
	if err := payment.DB.Exec(insertQuery, id, req.UserId, req.PaymentRef, req.ProjectId).Error; err != nil {
		return "", err
	}
	return id.String(), nil
}

func (payment *PaymentAdapter) AddPaymentToFreelancer(req entities.FreelancerAccount) error {
	id := uuid.New()
	insertQuery := `INSERT INTO freelancer_accounts (id, amount, payment_id, freelancer_id) VALUES ($1, $2, $3, $4)`
	if err := payment.DB.Exec(insertQuery, id, req.Amount, req.PaymentId, req.FreelancerId).Error; err != nil {
		return err
	}
	return nil
}

func (payment *PaymentAdapter) AddPaymentToAdmin(req entities.AdminAccount) error {
	id := uuid.New()
	insertQuery := `INSERT INTO admin_accounts (id, amount, payment_id) VALUES ($1, $2, $3)`
	if err := payment.DB.Exec(insertQuery, id, req.Amount, req.PaymentId).Error; err != nil {
		return err
	}
	return nil
}
