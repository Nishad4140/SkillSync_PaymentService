package usecase

import (
	"github.com/Nishad4140/SkillSync_PaymentService/entities"
	"github.com/Nishad4140/SkillSync_PaymentService/internal/adapters"
)

type PaymentUsecase struct {
	adapters adapters.PaymentAdapterInterface
}

func NewPaymentUsecase(adapters adapters.PaymentAdapterInterface) *PaymentUsecase {
	return &PaymentUsecase{
		adapters: adapters,
	}
}

func (payment *PaymentUsecase) AddPayment(req entities.Payment) (string, error) {
	return payment.adapters.AddPayment(req)
}

func (payment *PaymentUsecase) AddPaymentToFreelancer(req entities.FreelancerAccount) error {
	return payment.adapters.AddPaymentToFreelancer(req)
}

func (payment *PaymentUsecase) AddPaymentToAdmin(req entities.AdminAccount) error {
	return payment.adapters.AddPaymentToAdmin(req)
}
