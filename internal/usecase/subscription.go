package usecase

import (
	"time"

	"github.com/Aiszhio/Task/internal/domain"
	errors_package "github.com/Aiszhio/Task/internal/errors"
	"github.com/Aiszhio/Task/internal/repository/postgres"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionUseCase interface {
	AcceptSubscription(ctx *gin.Context, Subscription *domain.Subscription) error
	GetSubscription(ctx *gin.Context, id uuid.UUID) (*domain.Subscription, error)
	RemoveSubscription(ctx *gin.Context, id uuid.UUID) error
	RefreshSubscription(ctx *gin.Context, Subscription *domain.Subscription) error
	GetListSubscriptions(ctx *gin.Context, Subscription *domain.SubscriptionSummary) (uint64, error)
}

type SubscriptionUseCaseImpl struct {
	db postgres.SubscriptionRepository
}

func NewSubscriptionUseCase(db postgres.SubscriptionRepository) *SubscriptionUseCaseImpl {
	return &SubscriptionUseCaseImpl{
		db: db,
	}
}

func validatesub(Subscription *domain.Subscription) error {
	if Subscription == nil {
		return errors_package.EmptySub
	}

	if Subscription.ServiceName == "" {
		return errors_package.EmptySubName
	}

	if Subscription.UserId == uuid.Nil {
		return errors_package.EmptyUser
	}

	if Subscription.Price <= 0 {
		return errors_package.InvalidPrice
	}

	start := Subscription.StartDate
	now := time.Now()

	if start.Year() < now.Year() ||
		(start.Year() == now.Year() && start.Month() < now.Month()) {
		return errors_package.SubscriptionInPast
	}

	return nil
}

func (uc *SubscriptionUseCaseImpl) AcceptSubscription(ctx *gin.Context, Subscription *domain.Subscription) error {
	err := validatesub(Subscription)
	if err != nil {
		return err
	}

	err = uc.db.CreateSubscription(ctx, Subscription)
	if err != nil {
		return err
	}

	return nil
}

func (uc *SubscriptionUseCaseImpl) GetSubscription(ctx *gin.Context, id uuid.UUID) (*domain.Subscription, error) {
	if id == uuid.Nil {
		return nil, errors_package.ErrEmptyId
	}

	subscr, err := uc.db.ReadSubscription(ctx, id)
	if err != nil {
		return nil, err
	}

	return subscr, nil
}

func (uc *SubscriptionUseCaseImpl) RemoveSubscription(ctx *gin.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors_package.ErrEmptyId
	}

	err := uc.db.DeleteSubscription(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *SubscriptionUseCaseImpl) RefreshSubscription(ctx *gin.Context, Subscription *domain.Subscription) error {
	err := validatesub(Subscription)
	if err != nil {
		return err
	}

	err = uc.db.UpdateSubscription(ctx, Subscription)
	if err != nil {
		return err
	}

	return nil
}

func (uc *SubscriptionUseCaseImpl) GetListSubscriptions(ctx *gin.Context, Subscription *domain.SubscriptionSummary) (uint64, error) {
	if Subscription.EndDate.Before(Subscription.StartDate) {
		return 0, errors_package.InvalidPeriod
	}

	return uc.db.ListSubscriptions(ctx, Subscription.StartDate, Subscription.EndDate, Subscription.UserID, Subscription.ServiceName)
}
