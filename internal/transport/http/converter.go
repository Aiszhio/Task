package transport

import (
	"time"

	"github.com/Aiszhio/Task/internal/domain"
	errors_package "github.com/Aiszhio/Task/internal/errors"
)

func makesubscription(Subscription *SubscriptionRequest) (*domain.Subscription, error) {
	start := time.Time(Subscription.StartDate)
	var end time.Time

	if Subscription.EndDate != nil {
		end = time.Time(*Subscription.EndDate)
	} else {
		end = start.AddDate(0, 1, 0)
	}

	subscr := &domain.Subscription{
		Id:          Subscription.Id,
		ServiceName: Subscription.ServiceName,
		Price:       Subscription.Price,
		StartDate:   start,
		EndDate:     end,
		UserId:      Subscription.UserID,
	}

	return subscr, nil
}

func TransportToDomain(Subscription *SubscriptionRequest) (*domain.Subscription, error) {
	sub, err := makesubscription(Subscription)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func makelistsubscriptions(req *SubscriptionSummaryRequest) (*domain.SubscriptionSummary, error) {
	start := time.Time(req.StartDate)
	end := time.Time(req.EndDate)
	var zeroDate MonthYear

	if req.EndDate == zeroDate {
		return nil, errors_package.ErrDateRequired
	}

	return &domain.SubscriptionSummary{
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		StartDate:   start,
		EndDate:     end,
	}, nil
}

func ListToDomain(req *SubscriptionSummaryRequest) (*domain.SubscriptionSummary, error) {
	list, err := makelistsubscriptions(req)
	if err != nil {
		return nil, err
	}

	return list, nil
}
