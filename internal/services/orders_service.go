package services

import (
	"context"
	"errors"
	"fp_kata/common/utils"
	"fp_kata/internal/datasources"
	"fp_kata/internal/datasources/dsmodels"
	"fp_kata/internal/models"
)

const compOrdersService = "OrdersService"

type OrdersService interface {
	StoreOrder(ctx context.Context, userId int, order models.Order) (*models.Order, error)
	GetOrder(ctx context.Context, userId int, id int) (*models.Order, error)
	GetOrders(ctx context.Context, userId int) ([]*models.Order, error)
	GetOrdersWithFilter(ctx context.Context, userId int, filter func(order *models.Order) bool) ([]*models.Order, error)
}

type ordersService struct {
	storage              datasources.OrdersDatasource
	paymentService       PaymentsService
	authorizationService AuthorizationService
}

func NewOrdersService(storage datasources.OrdersDatasource, paymentService PaymentsService, authorizationService AuthorizationService) OrdersService {
	return &ordersService{storage: storage, paymentService: paymentService, authorizationService: authorizationService}
}

const errUserRequired = "user id is required"

func (service *ordersService) StoreOrder(ctx context.Context, userId int, order models.Order) (*models.Order, error) {
	utils.LogAction(ctx, compOrdersService, "StoreOrder")

	// Validate user
	if userId == 0 || order.User == nil || order.User.ID != userId {
		return nil, errors.New(errUserRequired)
	}

	isNewOrder := order.ID == 0
	// Generate new order ID if not present
	if isNewOrder {
		order.ID = utils.GenerateNewId()
	}

	// Process payments
	// payment Ids inside order will be updated <-- side effect
	storedPayments, err := service.processPayments(ctx, &order)
	if err != nil {
		return nil, err
	}

	// Store order in database
	var storedOrderModel *dsmodels.Order
	if isNewOrder {
		storedOrderModel, err = service.storage.InsertOrder(ctx, *order.ToDSModel())
	} else {
		storedOrderModel, err = service.storage.UpdateOrder(ctx, *order.ToDSModel())
	}
	if err != nil {
		return nil, err
	}

	// Map stored order to the response model
	newOrder := models.MapToOrder(*storedOrderModel, storedPayments, order.User)
	return newOrder, nil
}

// processPayments handles storing payments and updating payment IDs.
func (service *ordersService) processPayments(ctx context.Context, order *models.Order) ([]*models.Payment, error) {
	storedPayments := make([]*models.Payment, len(order.Payments))

	for i, payment := range order.Payments {
		payment.Order = order
		storedPayment, err := service.paymentService.StorePayment(ctx, *payment)
		if err != nil {
			return nil, err
		}
		order.Payments[i].Id = storedPayment.Id
		storedPayments[i] = storedPayment
	}

	return storedPayments, nil
}

func (service *ordersService) GetOrder(ctx context.Context, userId int, id int) (*models.Order, error) {
	utils.LogAction(ctx, compOrdersService, "GetOrder")

	if userId == 0 {
		return nil, errors.New("user id is required")
	}

	dsOrder, err := service.storage.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	if dsOrder.UserId != userId {
		return nil, errors.New("user is not authorized to access this order")
	}

	payments, err := service.paymentService.GetPaymentsByOrder(ctx, id)
	if err != nil {
		return nil, err
	}

	order := models.MapToOrder(*dsOrder, payments, &models.User{ID: dsOrder.UserId})
	return order, nil
}

func (service *ordersService) GetOrders(ctx context.Context, userId int) ([]*models.Order, error) {
	utils.LogAction(ctx, compOrdersService, "GetOrders")

	if userId == 0 {
		return nil, errors.New("user id is required")
	}
	dsOrders, err := service.storage.GetAllOrdersForUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	var orders []*models.Order
	for _, dsOrder := range dsOrders {
		payments, err := service.paymentService.GetPaymentsByOrder(ctx, dsOrder.ID)
		if err != nil {
			return nil, err
		}
		order := models.MapToOrder(dsOrder, payments, &models.User{ID: dsOrder.UserId})
		orders = append(orders, order)
	}

	return orders, nil
}

func (service *ordersService) GetOrdersWithFilter(ctx context.Context, userId int, filter func(order *models.Order) bool) ([]*models.Order, error) {
	utils.LogAction(ctx, compOrdersService, "GetOrdersWithFilter")

	if userId == 0 {
		return nil, errors.New("user id is required")
	}
	allDsOrders, err := service.storage.GetAllOrdersForUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	var filteredOrders []*models.Order
	for _, dsOrder := range allDsOrders {
		order := models.MapToOrder(dsOrder, []*models.Payment{}, &models.User{ID: dsOrder.UserId})
		if filter(order) {
			order.Payments, err = service.paymentService.GetPaymentsByOrder(ctx, order.ID)
			if err != nil {
				return nil, err
			}
			filteredOrders = append(filteredOrders, order)
		}
	}

	return filteredOrders, nil
}
