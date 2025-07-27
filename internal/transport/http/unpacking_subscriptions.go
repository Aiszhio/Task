package transport

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MonthYear time.Time

type SubscriptionRequest struct {
	Id          uuid.UUID  `json:"id"`
	ServiceName string     `json:"service_name" binding:"required"`
	Price       int        `json:"price"        binding:"required,min=0"`
	UserID      uuid.UUID  `json:"user_id"      binding:"required"`
	StartDate   MonthYear  `json:"start_date"   binding:"required"`
	EndDate     *MonthYear `json:"end_date,omitempty"`
}

type SubscriptionSummaryRequest struct {
	UserID      uuid.UUID `json:"user_id"       binding:"required"`
	ServiceName string    `json:"service_name"  binding:"required"`
	StartDate   MonthYear `json:"start_date"    binding:"required"`
	EndDate     MonthYear `json:"end_date"      binding:"required"`
}

func (m *MonthYear) UnmarshalJSON(data []byte) error {
	str := ""
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	formattedTime, err := time.Parse("01-2006", str)
	if err != nil {
		return fmt.Errorf("invalid date format %s: %s", str, err)
	}

	*m = MonthYear(formattedTime)

	return nil
}

func (m *MonthYear) MarshalJSON() ([]byte, error) {
	t := time.Time(*m)
	return []byte(t.Format("07-2025")), nil
}
