package transport

import (
	"log"

	"github.com/Aiszhio/Task/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create() gin.HandlerFunc
	GetByID() gin.HandlerFunc
	Update() gin.HandlerFunc
	DeleteByID() gin.HandlerFunc
	ListByPeriod() gin.HandlerFunc
}

type SubscriptionHandler struct {
	Repository usecase.SubscriptionUseCase
}

func NewSubscriptionHandler(repo usecase.SubscriptionUseCase) *SubscriptionHandler {
	return &SubscriptionHandler{
		Repository: repo,
	}
}

// CreateSubscription godoc
// @Summary     Создать подписку
// @Description Создаёт новую запись о подписке
// @Tags        subscriptions
// @Accept      json
// @Produce     json
// @Param       body  body  SubscriptionRequest true "JSON"
// @Success     201   {object} domain.Subscription
// @Failure     400   {object} ErrorResponse
// @Router      /subscriptions [post]
func (handler *SubscriptionHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req SubscriptionRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		sub, err := TransportToDomain(&req)
		if err != nil {
			log.Printf("error in converting data to domain: %v\n", err)
			ctx.JSON(400, gin.H{"error": err.Error()})
		}

		err = handler.Repository.AcceptSubscription(ctx, sub)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"message": "Subscription created"})
	}
}

// GetSubscription godoc
// @Summary     Получить подписку по ID
// @Tags        subscriptions
// @Produce     json
// @Param       id   path  string true "Subscription ID"
// @Success     200  {object} domain.Subscription
// @Failure     400  {object} ErrorResponse
// @Router      /subscriptions/{id} [get]
func (handler *SubscriptionHandler) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		subscription, err := handler.Repository.GetSubscription(ctx, id)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"message": subscription})
	}
}

// UpdateSubscription godoc
// @Summary     Полное обновление подписки
// @Tags        subscriptions
// @Accept      json
// @Produce     json
// @Param       id   path  string true "Subscription ID"
// @Param       body body  SubscriptionRequest true "JSON"
// @Success     200  {object} domain.Subscription
// @Failure     400  {object} ErrorResponse
// @Router      /subscriptions/{id} [put]
func (handler *SubscriptionHandler) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req SubscriptionRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		sub, err := TransportToDomain(&req)
		if err != nil {
			log.Printf("error in converting data to domain: %v\n", err)
			ctx.JSON(400, gin.H{"error": err.Error()})
		}

		err = handler.Repository.RefreshSubscription(ctx, sub)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"message": "Subscription was successfully updated"})
	}
}

// DeleteSubscription godoc
// @Summary     Удалить подписку
// @Tags        subscriptions
// @Produce     json
// @Param       id   path  string true "Subscription ID"
// @Success     200  {object} SuccessResponse
// @Failure     400  {object} ErrorResponse
// @Router      /subscriptions/{id} [delete]
func (handler *SubscriptionHandler) DeleteByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err = handler.Repository.RemoveSubscription(ctx, id)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"message": "Subscription deleted"})
	}
}

// ListSubscriptions godoc
// @Summary     Сумма подписок за период
// @Description Возвращает суммарную стоимость подписок по фильтру
// @Tags        subscriptions
// @Accept      json
// @Produce     json
// @Param       body body SubscriptionSummaryRequest true "JSON"
// @Success     200  {object} SuccessResponse  "total_cost"
// @Failure     400  {object} ErrorResponse
// @Router      /subscriptions/list [post]
func (handler *SubscriptionHandler) ListByPeriod() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req SubscriptionSummaryRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		list, err := ListToDomain(&req)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		subscriptions, err := handler.Repository.GetListSubscriptions(ctx, list)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"total_cost": subscriptions})
	}
}
