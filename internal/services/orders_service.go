package services

import (
	"context"
	"errors"
	"fp_kata/common/constants"
	"fp_kata/common/monads"
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
	userService          UsersService
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
	var storedOrderModelResult monads.Result[dsmodels.Order]
	if isNewOrder {
		storedOrderModelResult = service.storage.InsertOrder(ctx, *order.ToDSModel())
	} else {
		storedOrderModelResult = service.storage.UpdateOrder(ctx, *order.ToDSModel())
	}
	if storedOrderModelResult.IsError() {
		return nil, storedOrderModelResult.Error()
	}

	// Map stored order to the response model
	newOrder := models.MapToOrder(storedOrderModelResult.MustGet())
	newOrder.Payments = storedPayments

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

	dsOrderResult := service.storage.GetOrder(ctx, id)
	if dsOrderResult.IsError() {
		return nil, dsOrderResult.Error()
	}

	order, err := service.processDsOrder(ctx, userId, dsOrderResult.MustGet())
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (service *ordersService) GetOrders(ctx context.Context, userId int) ([]*models.Order, error) {
	utils.LogAction(ctx, compOrdersService, "GetOrders")

	if userId == 0 {
		return nil, errors.New("user id is required")
	}
	dsOrdersResult := service.storage.GetAllOrdersForUser(ctx, userId)
	if dsOrdersResult.IsError() {
		return nil, dsOrdersResult.Error()
	}

	var orders []*models.Order
	for _, dsOrder := range dsOrdersResult.MustGet() {
		order, err := service.processDsOrder(ctx, userId, dsOrder)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (service *ordersService) GetOrdersWithFilter(ctx context.Context, userId int, filter func(order *models.Order) bool) ([]*models.Order, error) {
	utils.LogAction(ctx, compOrdersService, "GetOrdersWithFilter")

	if userId == 0 {
		return nil, errors.New("user id is required")
	}
	dsOrdersResult := service.storage.GetAllOrdersForUser(ctx, userId)
	if dsOrdersResult.IsError() {
		return nil, dsOrdersResult.Error()
	}

	var filteredOrders []*models.Order
	for _, dsOrder := range dsOrdersResult.MustGet() {
		order := models.MapToOrder(dsOrder)
		if filter(order) {
			order, err := service.processOrder(ctx, userId, order)
			if err != nil {
				return nil, err
			}
			filteredOrders = append(filteredOrders, order)
		}
	}

	return filteredOrders, nil
}

func (service *ordersService) processDsOrder(ctx context.Context, userId int, storedOrder dsmodels.Order) (*models.Order, error) {
	// Map the stored order
	order := models.MapToOrder(storedOrder)

	order, err := service.processOrder(ctx, userId, order)

	if err != nil {
		return nil, err
	}
	return order, nil
}

func (service *ordersService) processOrder(ctx context.Context, userId int, order *models.Order) (*models.Order, error) {

	// Authorization check
	isAuthorized, err := service.authorizationService.IsAuthorized(ctx, userId, order)
	if err != nil {
		return nil, err
	}
	if !isAuthorized {
		return nil, errors.New("user is not authorized to access this order")
	}

	// Add payments to the order
	order, err = service.addPayments(ctx, order)
	if err != nil {
		return nil, err
	}

	// Add user to the order
	order, err = service.addUser(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (service *ordersService) addPayments(ctx context.Context, order *models.Order) (*models.Order, error) {
	utils.LogAction(ctx, compOrdersService, "addPayment")

	payments, err := service.paymentService.GetPaymentsByOrder(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	order.Payments = payments
	return order, nil
}

func (service *ordersService) addUser(ctx context.Context, order *models.Order) (*models.Order, error) {
	utils.LogAction(ctx, compOrdersService, "addPayment")

	user := ctx.Value(constants.AuthenticatedUserKey).(*models.User)

	if user == nil {
		return nil, errors.New(errUserRequired)
	}
	order.User = user
	return order, nil
}
