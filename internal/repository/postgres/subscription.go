package postgres

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/Aiszhio/Task/internal/db"
	"github.com/Aiszhio/Task/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	insertQuery = `INSERT INTO subscriptions
    				(id, service_name, price, user_id, start_date, end_date)
					VALUES ($1, $2, $3, $4, $5, $6);`

	selectUserQuery = `SELECT id, service_name, price, user_id, start_date, end_date
						FROM subscriptions
						WHERE id = $1;`

	deleteQuery = `DELETE FROM subscriptions WHERE id = $1;`

	updateQuery = `UPDATE subscriptions
					SET service_name = $1,
						price        = $2,
						start_date   = $3,
						end_date     = $4
					WHERE user_id = $5;`

	selectFilterQuery = `SELECT COALESCE(SUM(price), 0) AS total_cost
						FROM   subscriptions
						WHERE  user_id = $1
						  AND  service_name = $2
						  AND  start_date  <= $4
						  AND  (end_date IS NULL
								OR end_date >= $3);`

	uniqueFilterQuery = `SELECT COUNT(*) FROM subscriptions
							WHERE user_id = $1 and service_name = $2 and start_date = $3`
)

type SubscriptionRepository interface {
	CreateSubscription(ctx *gin.Context, sub *domain.Subscription) error
	ReadSubscription(ctx *gin.Context, id uuid.UUID) (*domain.Subscription, error)
	DeleteSubscription(ctx *gin.Context, id uuid.UUID) error
	UpdateSubscription(ctx *gin.Context, sub *domain.Subscription) error
	ListSubscriptions(ctx *gin.Context, startDate, endDate time.Time, userId uuid.UUID, serviceName string) (uint64, error)
}

type PGSubscriptionRepository struct {
	db     *db.Pool
	logger *slog.Logger
}

func NewRepository(pool *db.Pool, logger *slog.Logger) *PGSubscriptionRepository {
	return &PGSubscriptionRepository{
		db:     pool,
		logger: logger,
	}
}

func (repo *PGSubscriptionRepository) CreateSubscription(ctx *gin.Context, sub *domain.Subscription) error {
	var count int
	err := repo.db.Client.QueryRow(
		ctx, uniqueFilterQuery,
		sub.UserId, sub.ServiceName, sub.StartDate,
	).Scan(&count)
	if err != nil {
		repo.logger.Error("failed to check uniqueness", "err", err)
		return err
	}

	if count > 0 {
		repo.logger.Info("subscription already exists",
			"user_id", sub.UserId,
			"service", sub.ServiceName,
			"start_date", sub.StartDate,
		)
		return fmt.Errorf("subscription already exists for user %s service %s at %s",
			sub.UserId, sub.ServiceName, sub.StartDate.Format("2006-01-02"),
		)
	}

	_, err = repo.db.Client.Exec(ctx, insertQuery, sub.Id,
		sub.ServiceName, sub.Price, sub.UserId, sub.StartDate, sub.EndDate)
	if err != nil {
		log.Printf("error creating subscription: %v", err)
		return err
	}

	repo.logger.Info("successfully created subscription")
	return nil
}

func (repo *PGSubscriptionRepository) ReadSubscription(ctx *gin.Context, id uuid.UUID) (*domain.Subscription, error) {
	subscr := &domain.Subscription{}

	err := repo.db.Client.QueryRow(ctx, selectUserQuery, id).Scan(&subscr.Id, &subscr.ServiceName, &subscr.Price,
		&subscr.UserId, &subscr.StartDate, &subscr.EndDate)
	if err != nil {
		log.Printf("error reading subscription: %v", err)
		return nil, err
	}

	repo.logger.Info("successfully read subscription")
	return subscr, nil
}

func (repo *PGSubscriptionRepository) DeleteSubscription(ctx *gin.Context, id uuid.UUID) error {
	cmd, err := repo.db.Client.Exec(ctx, deleteQuery, id)
	if err != nil {
		repo.logger.Error("delete failed", "err", err)
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("subscription %s not found", id)
	}

	repo.logger.Info("successfully deleted subscription", "id", id)
	return nil

}

func (repo *PGSubscriptionRepository) UpdateSubscription(ctx *gin.Context, sub *domain.Subscription) error {
	_, err := repo.db.Client.Exec(ctx, updateQuery, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate, sub.UserId)
	if err != nil {
		log.Printf("error updating subscription: %v", err)
		return err
	}

	repo.logger.Info("successfully updated subscription")
	return nil
}

func (repo *PGSubscriptionRepository) ListSubscriptions(ctx *gin.Context,
	startDate, endDate time.Time,
	userId uuid.UUID,
	serviceName string) (uint64, error) {

	var total uint64
	err := repo.db.Client.QueryRow(
		ctx, selectFilterQuery,
		userId, serviceName, startDate, endDate,
	).Scan(&total)

	if err != nil {
		log.Printf("error calculating total cost: %v", err)
		return 0, err
	}

	repo.logger.Info("successfully list subscriptions")
	return total, nil
}
