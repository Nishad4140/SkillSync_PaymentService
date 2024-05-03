package adapters

import "github.com/Nishad4140/SkillSync_PaymentService/entities"

type PaymentAdapterInterface interface {
	AddPayment(req entities.Payment) (string, error)
	AddPaymentToFreelancer(req entities.FreelancerAccount) error
	AddPaymentToAdmin(req entities.AdminAccount) error
}
