package services

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/Nishad4140/SkillSync_PaymentService/entities"
	helperstruct "github.com/Nishad4140/SkillSync_PaymentService/helperStruct"
	"github.com/Nishad4140/SkillSync_PaymentService/internal/usecase"
	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/razorpay/razorpay-go"
	"google.golang.org/grpc"
)

type PaymentService struct {
	usecases    usecase.PaymentUsecaseInterface
	UserConn    pb.UserServiceClient
	ProjectConn pb.ProjectServiceClient
}

func NewPaymentService(usecases usecase.PaymentUsecaseInterface) *PaymentService {
	userConn, _ := grpc.Dial("localhost:4001", grpc.WithInsecure())
	projectConn, _ := grpc.Dial("localhost:4002", grpc.WithInsecure())
	return &PaymentService{
		usecases:    usecases,
		UserConn:    pb.NewUserServiceClient(userConn),
		ProjectConn: pb.NewProjectServiceClient(projectConn),
	}
}

func (payment *PaymentService) userPayment(c *gin.Context) {
	fmt.Println(1)
	projectId := c.Query("project_id")
	userId := c.Query("user_id")
	fmt.Println(2)
	project, err := payment.ProjectConn.GetProject(context.Background(), &pb.GetProjectById{
		Id:     projectId,
		UserId: userId,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "error getting project info",
			Error:      err.Error(),
		})
		return
	}
	fmt.Println(project.Status)

	if project.Status != "completed" {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "project not yet completed",
			Error:      "project not yet completed",
		})
		return
	}
	fmt.Println(project.Status)

	if project.IsPaid {
		c.HTML(200, "alreadyPaid.html", gin.H{})
		return
	}
	fmt.Println(project.Price)
	client := razorpay.NewClient(os.Getenv("RAZORPAY_ID"), os.Getenv("RAZORPAY_SECRET"))
	data := map[string]interface{}{
		"amount":   project.Price * 100,
		"currency": "INR",
		"receipt":  "test_receipt_id",
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "error creating razorpay client",
			Error:      err.Error(),
		})
		return
	}
	fmt.Println(5)
	value := body["id"]
	razorPayId := value.(string)
	c.HTML(200, "payment.html", gin.H{
		"total_amount": project.Price,
		"total":        project.Price,
		"orderid":      razorPayId,
		"project_id":   projectId,
		"userId":       userId,
	})
	fmt.Println(6)
}

func (payment *PaymentService) verifyPayment(c *gin.Context) {
	paymentRef := c.Query("payment_ref")
	userId := c.Query("user_id")
	projectId := c.Query("project_id")

	userID, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "please chose a valid user to proceed",
			Error:      err.Error(),
		})
		return
	}

	projectID, err := uuid.Parse(projectId)
	if err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "please chose a valid project to proceed",
			Error:      err.Error(),
		})
		return
	}

	if paymentRef == "" {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "payment failed",
		})
		return
	}
	paymentId, err := payment.usecases.AddPayment(entities.Payment{
		UserId:     userID,
		PaymentRef: paymentRef,
		ProjectId:  projectID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "error while adding payment",
			Error:      err.Error(),
		})
		return
	}

	paymentID, err := uuid.Parse(paymentId)
	if err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "error while converting to uuid",
			Error:      err.Error(),
		})
		return
	}

	project, err := payment.ProjectConn.GetProject(context.Background(), &pb.GetProjectById{
		Id:     projectId,
		UserId: userId,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "error getting project info",
			Error:      err.Error(),
		})
		return
	}

	if err := payment.usecases.AddPaymentToAdmin(entities.AdminAccount{
		Amount:    float64(project.Price) / 10,
		PaymentId: paymentID,
	}); err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "error while adding payment",
			Error:      err.Error(),
		})
		return
	}

	freelancerID, err := uuid.Parse(project.FreelancerId)
	if err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "error while converting to uuid",
			Error:      err.Error(),
		})
		return
	}

	if err := payment.usecases.AddPaymentToFreelancer(entities.FreelancerAccount{
		Amount:       float64(project.Price) - (float64(project.Price) / 10),
		PaymentId:    paymentID,
		FreelancerId: freelancerID,
	}); err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "error while adding payment",
			Error:      err.Error(),
		})
		return
	}
	if _, err := payment.ProjectConn.PaymentStatusChange(context.Background(), &pb.PaymentStatusRequest{
		ProjectId: projectId,
		Status:    true,
	}); err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "error while changin payment status",
			Error:      err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, helperstruct.Response{
		StatusCode: 200,
		Message:    "payment verified",
		Data:       true,
	})

}
func (p *PaymentService) servePaymentSuccessPage(c *gin.Context) {
	c.HTML(200, "paymentVerified.html", gin.H{})
}
