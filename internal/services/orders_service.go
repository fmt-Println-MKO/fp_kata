package services

import (
	"context"
	"fp_kata/common/constants"
	"fp_kata/common/utils"
	"fp_kata/internal/datasources"
	"fp_kata/internal/datasources/dsmodels"
	"fp_kata/internal/models"
	"github.com/samber/mo"
)

const compOrdersService = "OrdersService"

type OrdersService interface {
	StoreOrder(ctx context.Context, userId int, order models.Order) mo.Result[*models.Order]
	GetOrder(ctx context.Context, userId int, id int) mo.Result[*models.Order]
	GetOrders(ctx context.Context, userId int) mo.Result[[]*models.Order]
	GetOrdersWithFilter(ctx context.Context, userId int, filter func(order *models.Order) bool) mo.Result[[]*models.Order]
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

func (service *ordersService) StoreOrder(ctx context.Context, userId int, order models.Order) mo.Result[*models.Order] {
	utils.LogAction(ctx, compOrdersService, "StoreOrder")

	// Validate user
	if userId == 0 || order.User == nil || order.User.ID != userId {
		return mo.Errf[*models.Order](errUserRequired)
	}

	isNewOrder := order.ID == 0
	// Generate new order ID if not present
	if isNewOrder {
		order.ID = utils.GenerateNewId()
	}

	// Process payments
	// payment Ids inside order will be updated <-- side effect
	return mo.Fold(service.processPayments(ctx, &order),
		func(storedPayments []*models.Payment) mo.Result[*models.Order] {

			return mo.Fold(
				service.storeOrUpdateOrder(isNewOrder, ctx, order),
				func(dsOrder dsmodels.Order) mo.Result[*models.Order] {
					// Map stored order to the response model
					newOrder := models.MapToOrder(dsOrder)
					newOrder.Payments = storedPayments
					return mo.Ok(newOrder)
				}, func(err error) mo.Result[*models.Order] {
					return mo.Err[*models.Order](err)
				})

		},
		func(err error) mo.Result[*models.Order] {
			return mo.Err[*models.Order](err)
		},
	)
}

func (service *ordersService) storeOrUpdateOrder(isNewOrder bool, ctx context.Context, order models.Order) mo.Result[dsmodels.Order] {
	var storedOrderModelResult mo.Result[dsmodels.Order]
	if isNewOrder {
		storedOrderModelResult = service.storage.InsertOrder(ctx, *order.ToDSModel())
	} else {
		storedOrderModelResult = service.storage.UpdateOrder(ctx, *order.ToDSModel())
	}
	return storedOrderModelResult
}

// processPayments handles storing payments and updating payment IDs.
func (service *ordersService) processPayments(ctx context.Context, order *models.Order) mo.Result[[]*models.Payment] {
	storedPayments := make([]*models.Payment, len(order.Payments))

	for i, payment := range order.Payments {
		payment.Order = order
		storedPayment, err := service.paymentService.StorePayment(ctx, *payment)
		if err != nil {
			return mo.Err[[]*models.Payment](err)
		}
		order.Payments[i].Id = storedPayment.Id
		storedPayments[i] = storedPayment
	}

	return mo.Ok(storedPayments)
}

func (service *ordersService) GetOrder(ctx context.Context, userId int, id int) mo.Result[*models.Order] {
	utils.LogAction(ctx, compOrdersService, "GetOrder")

	if userId == 0 {
		return mo.Errf[*models.Order]("user id is required")
	}

	return mo.Fold(
		service.storage.GetOrder(ctx, id),
		func(dsOrder dsmodels.Order) mo.Result[*models.Order] {
			return service.processDsOrder(ctx, userId, dsOrder)
		},
		func(err error) mo.Result[*models.Order] {
			return mo.Err[*models.Order](err)
		})
}

func (service *ordersService) GetOrders(ctx context.Context, userId int) mo.Result[[]*models.Order] {
	utils.LogAction(ctx, compOrdersService, "GetOrders")

	if userId == 0 {
		return mo.Errf[[]*models.Order]("user id is required")
	}
	return mo.Fold(
		service.storage.GetAllOrdersForUser(ctx, userId),
		func(dsOrders []dsmodels.Order) mo.Result[[]*models.Order] {
			orders := make([]*models.Order, 0)
			for _, dsOrder := range dsOrders {
				orderResult := service.processDsOrder(ctx, userId, dsOrder)
				if orderResult.IsError() {
					return mo.Err[[]*models.Order](orderResult.Error())
				}
				orders = append(orders, orderResult.MustGet())
			}
			return mo.Ok(orders)
		},
		func(err error) mo.Result[[]*models.Order] {
			return mo.Err[[]*models.Order](err)
		},
	)
}

func (service *ordersService) GetOrdersWithFilter(ctx context.Context, userId int, filter func(order *models.Order) bool) mo.Result[[]*models.Order] {
	utils.LogAction(ctx, compOrdersService, "GetOrdersWithFilter")

	if userId == 0 {
		return mo.Errf[[]*models.Order]("user id is required")
	}
	return mo.Fold(
		service.storage.GetAllOrdersForUser(ctx, userId),
		func(dsOrders []dsmodels.Order) mo.Result[[]*models.Order] {
			filteredOrders := make([]*models.Order, 0)
			for _, dsOrder := range dsOrders {
				order := models.MapToOrder(dsOrder)
				if filter(order) {
					filteredOrderResult := service.processOrder(ctx, userId, order)
					if filteredOrderResult.IsError() {
						return mo.Err[[]*models.Order](filteredOrderResult.Error())
					}
					filteredOrders = append(filteredOrders, filteredOrderResult.MustGet())
				}
			}
			return mo.Ok(filteredOrders)
		},
		func(err error) mo.Result[[]*models.Order] {
			return mo.Err[[]*models.Order](err)
		},
	)
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
