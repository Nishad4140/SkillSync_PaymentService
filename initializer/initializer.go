package initializer

import (
	"github.com/Nishad4140/SkillSync_PaymentService/internal/adapters"
	"github.com/Nishad4140/SkillSync_PaymentService/internal/services"
	"github.com/Nishad4140/SkillSync_PaymentService/internal/usecase"
	"gorm.io/gorm"
)

func Initializer(db *gorm.DB) *services.PaymentEngine {
	adapter := adapters.NewPaymentAdapter(db)
	usecases := usecase.NewPaymentUsecase(adapter)
	service := services.NewPaymentService(usecases)
	return services.NewPaymentEngine(service)
}
