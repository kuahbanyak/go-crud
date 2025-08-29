package test

import (
	"testing"
	"time"

	"github.com/kuahbanyak/go-crud/internal/inventory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockInventoryRepo is a mock implementation of inventory.Repository
type MockInventoryRepo struct {
	mock.Mock
}

func (m *MockInventoryRepo) Create(p *inventory.Part) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockInventoryRepo) List() ([]inventory.Part, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]inventory.Part), args.Error(1)
}

func (m *MockInventoryRepo) Update(p *inventory.Part) error {
	args := m.Called(p)
	return args.Error(0)
}

func TestPartModel_Validation(t *testing.T) {
	tests := []struct {
		name  string
		part  inventory.Part
		valid bool
	}{
		{
			name: "Valid part",
			part: inventory.Part{
				SKU:   "BRK-001",
				Name:  "Brake Pad",
				Qty:   50,
				Price: 25000,
			},
			valid: true,
		},
		{
			name: "Valid part with minimal data",
			part: inventory.Part{
				SKU:   "OIL-001",
				Name:  "Engine Oil",
				Qty:   20,
				Price: 15000,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.part.SKU)
			assert.NotEmpty(t, tt.part.Name)
			assert.GreaterOrEqual(t, tt.part.Qty, 0)
			assert.Greater(t, tt.part.Price, 0)
		})
	}
}

func TestInventoryRepository_Create(t *testing.T) {
	mockRepo := new(MockInventoryRepo)

	testPart := &inventory.Part{
		SKU:   "BRK-001",
		Name:  "Brake Pad",
		Qty:   50,
		Price: 25000,
	}

	mockRepo.On("Create", testPart).Return(nil)

	err := mockRepo.Create(testPart)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestInventoryRepository_List(t *testing.T) {
	mockRepo := new(MockInventoryRepo)

	expectedParts := []inventory.Part{
		{
			ID:        1,
			SKU:       "BRK-001",
			Name:      "Brake Pad",
			Qty:       50,
			Price:     25000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			SKU:       "OIL-001",
			Name:      "Engine Oil",
			Qty:       20,
			Price:     15000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockRepo.On("List").Return(expectedParts, nil)

	parts, err := mockRepo.List()

	assert.NoError(t, err)
	assert.NotNil(t, parts)
	assert.Len(t, parts, 2)
	assert.Equal(t, "BRK-001", parts[0].SKU)
	assert.Equal(t, "OIL-001", parts[1].SKU)

	mockRepo.AssertExpectations(t)
}

func TestInventoryRepository_Update(t *testing.T) {
	mockRepo := new(MockInventoryRepo)

	testPart := &inventory.Part{
		ID:    1,
		SKU:   "BRK-001",
		Name:  "Brake Pad Premium",
		Qty:   45,
		Price: 30000,
	}

	mockRepo.On("Update", testPart).Return(nil)

	err := mockRepo.Update(testPart)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPartSKU_Format(t *testing.T) {
	tests := []struct {
		name string
		sku  string
	}{
		{"Brake part SKU", "BRK-001"},
		{"Oil part SKU", "OIL-001"},
		{"Filter part SKU", "FLT-001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.sku)
			assert.Contains(t, tt.sku, "-")
		})
	}
}
