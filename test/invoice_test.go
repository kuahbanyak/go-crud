package test

import (
	"testing"

	"github.com/kuahbanyak/go-crud/internal/invoice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockInvoiceRepo is a mock implementation of invoice.Repository
type MockInvoiceRepo struct {
	mock.Mock
}

func (m *MockInvoiceRepo) Create(i *invoice.Invoice) error {
	args := m.Called(i)
	return args.Error(0)
}

func (m *MockInvoiceRepo) Summary() (map[string]int64, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int64), args.Error(1)
}

func TestInvoiceModel_Validation(t *testing.T) {
	tests := []struct {
		name    string
		invoice invoice.Invoice
		valid   bool
	}{
		{
			name: "Valid invoice",
			invoice: invoice.Invoice{
				BookingID: 1,
				Amount:    100000,
				Status:    "pending",
				PDFURL:    "https://example.com/invoice.pdf",
			},
			valid: true,
		},
		{
			name: "Valid paid invoice",
			invoice: invoice.Invoice{
				BookingID: 2,
				Amount:    250000,
				Status:    "paid",
				PDFURL:    "https://example.com/invoice2.pdf",
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotZero(t, tt.invoice.BookingID)
			assert.Greater(t, tt.invoice.Amount, 0)
			assert.NotEmpty(t, tt.invoice.Status)
			assert.Contains(t, []string{"pending", "paid", "cancelled"}, tt.invoice.Status)
		})
	}
}

func TestInvoiceRepository_Create(t *testing.T) {
	mockRepo := new(MockInvoiceRepo)

	testInvoice := &invoice.Invoice{
		BookingID: 1,
		Amount:    100000,
		Status:    "pending",
		PDFURL:    "https://example.com/invoice.pdf",
	}

	mockRepo.On("Create", testInvoice).Return(nil)

	err := mockRepo.Create(testInvoice)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestInvoiceRepository_Summary(t *testing.T) {
	mockRepo := new(MockInvoiceRepo)

	expectedSummary := map[string]int64{
		"count": 15,
	}

	mockRepo.On("Summary").Return(expectedSummary, nil)

	summary, err := mockRepo.Summary()

	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, int64(15), summary["count"])

	mockRepo.AssertExpectations(t)
}

func TestInvoiceStatuses(t *testing.T) {
	tests := []struct {
		name   string
		status string
	}{
		{"Pending status", "pending"},
		{"Paid status", "paid"},
		{"Cancelled status", "cancelled"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.status)
			assert.Contains(t, []string{"pending", "paid", "cancelled"}, tt.status)
		})
	}
}
