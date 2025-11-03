package payment
import (
	"fmt"
)
type StripeClient struct {
	apiKey  string
	baseURL string
}
type PaymentService interface {
	ProcessPayment(amount float64, currency, paymentMethod string) (*PaymentResult, error)
	RefundPayment(paymentID string, amount float64) error
}
type PaymentResult struct {
	ID     string
	Amount float64
	Status string
}
func NewStripeClient(apiKey string) *StripeClient {
	return &StripeClient{
		apiKey:  apiKey,
		baseURL: "https://api.stripe.com/v1",
	}
}
func (s *StripeClient) ProcessPayment(amount float64, currency, paymentMethod string) (*PaymentResult, error) {
	return &PaymentResult{
		ID:     fmt.Sprintf("pi_%d", int(amount*100)),
		Amount: amount,
		Status: "succeeded",
	}, nil
}
func (s *StripeClient) RefundPayment(paymentID string, amount float64) error {
	return nil
}

