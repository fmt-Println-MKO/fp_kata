package services

import (
	"context"
	"errors"
	"fp_kata/common/constants"
	"fp_kata/common/utils"
	"fp_kata/internal/datasources"
	"fp_kata/internal/datasources/dsmodels"
	"fp_kata/internal/models"
	"github.com/samber/mo"
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
	var storedOrderModelResult mo.Result[dsmodels.Order]
	if isNewOrder {
		storedOrderModelResult = service.storage.InsertOrder(ctx, *order.ToDSModel())
	} else {
		storedOrderModelResult = service.storage.UpdateOrder(ctx, *order.ToDSModel())
	}

	newOrder := mo.Fold[error, dsmodels.Order, mo.Result[*models.Order]](storedOrderModelResult,
		func(dsOrder dsmodels.Order) mo.Result[*models.Order] {
			// Map stored order to the response model
			newOrder := models.MapToOrder(dsOrder)
			newOrder.Payments = storedPayments
			return mo.Ok(newOrder)
		}, func(err error) mo.Result[*models.Order] {
			return mo.Err[*models.Order](err)
		})
	return newOrder.Get()
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

	orderResult := service.processDsOrder(ctx, userId, dsOrderResult.MustGet())
	if orderResult.IsError() {
		return nil, orderResult.Error()
	}

	return orderResult.Get()
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
		orderResult := service.processDsOrder(ctx, userId, dsOrder)
		if orderResult.IsError() {
			return nil, orderResult.Error()
		}
		orders = append(orders, orderResult.MustGet())
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
			filteredOrderResult := service.processOrder(ctx, userId, order)
			if filteredOrderResult.IsError() {
				return nil, filteredOrderResult.Error()
			}
			filteredOrders = append(filteredOrders, filteredOrderResult.MustGet())
		}
	}

	return filteredOrders, nil
}

func (service *ordersService) processDsOrder(ctx context.Context, userId int, storedOrder dsmodels.Order) mo.Result[*models.Order] {
	// Map the stored order
	order := models.MapToOrder(storedOrder)

	return service.processOrder(ctx, userId, order)
}

func (service *ordersService) processOrder(ctx context.Context, userId int, order *models.Order) mo.Result[*models.Order] {

	// Authorization check
	return service.verifyAuthorization(ctx, userId, order).
		// Add payments to the order
		FlatMap(func(value *models.Order) mo.Result[*models.Order] {
			return service.addPayments(ctx, value)
		}).
		// Add user to the order
		FlatMap(func(value *models.Order) mo.Result[*models.Order] {
			return service.addUser(ctx, value)
		})
}
func (service *ordersService) verifyAuthorization(ctx context.Context, userId int, order *models.Order) mo.Result[*models.Order] {
	utils.LogAction(ctx, compOrdersService, "verifyAuthorization")

	isAuthorized, err := service.authorizationService.IsAuthorized(ctx, userId, order)
	if err != nil {
		return mo.Err[*models.Order](err)
	}
	if !isAuthorized {
		return mo.Errf[*models.Order]("user is not authorized to access this order")
	}
	return mo.Ok(order)
}

func (service *ordersService) addPayments(ctx context.Context, order *models.Order) mo.Result[*models.Order] {
	utils.LogAction(ctx, compOrdersService, "addPayment")

	payments, err := service.paymentService.GetPaymentsByOrder(ctx, order.ID)
	if err != nil {
		return mo.Err[*models.Order](err)
	}
	order.Payments = payments
	return mo.Ok(order)
}

func (service *ordersService) addUser(ctx context.Context, order *models.Order) mo.Result[*models.Order] {
	utils.LogAction(ctx, compOrdersService, "addUser")

	user := ctx.Value(constants.AuthenticatedUserKey).(*models.User)
	if user == nil {
		return mo.Errf[*models.Order](errUserRequired)
	}
	order.User = user
	return mo.Ok(order)
}
