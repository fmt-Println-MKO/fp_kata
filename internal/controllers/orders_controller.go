package controllers

import (
	"context"
	"fp_kata/common/constants"
	"fp_kata/common/utils"
	"fp_kata/internal/models"
	"fp_kata/internal/services"
	"fp_kata/pkg/log"
	"fp_kata/pkg/transports"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"

	"strconv"
)

const compOrdersController = "OrdersController"

type OrdersController struct {
	orderService services.OrdersService
}

func NewOrdersController(orderService services.OrdersService) OrdersController {
	return OrdersController{
		orderService: orderService,
	}
}

func (c *OrdersController) RegisterOrderRoutes(app *fiber.App, authMiddleware fiber.Handler) {
	app.Post("/orders", c.CreateOrder, authMiddleware)
	app.Get("/orders", c.GetOrders, authMiddleware)
	app.Get("/orders/:id", c.GetOrder, authMiddleware)
}

func (c *OrdersController) CreateOrder(ctx fiber.Ctx) error {
	logger := log.GetFiberLogger(ctx)
	context := log.NewBackgroundContext(logger)
	utils.LogAction(context, compOrdersController, "CreateOrder")

	userID := ctx.Locals(constants.AuthenticatedUserIdKey).(int)
	user := ctx.Locals(constants.AuthenticatedUserKey).(models.User)

	var orderRequest = new(transports.OrderCreateRequest)

	if err := ctx.Bind().Body(orderRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	validate := validator.New()
	if err := validate.Struct(orderRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	order := orderRequest.ToOrder(user)

	newOrder, err := c.orderService.StoreOrder(context, userID, *order)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to create the order",
		})
	}

	orderResponse := transports.MapToOrderResponse(*newOrder)
	return ctx.Status(fiber.StatusCreated).JSON(orderResponse)
}

// GetOrders handles "/orders" with method "GET"
func (c *OrdersController) GetOrders(requestCtx fiber.Ctx) error {
	logger := log.GetFiberLogger(requestCtx)
	backgroundCtx := log.NewBackgroundContext(logger)
	utils.LogAction(backgroundCtx, compOrdersController, "GetOrders")

	user := requestCtx.Locals(constants.AuthenticatedUserKey).(models.User)

	backgroundCtx = context.WithValue(backgroundCtx, constants.AuthenticatedUserKey, &user)
	backgroundCtx = context.WithValue(backgroundCtx, constants.AuthenticatedUserIdKey, user.ID)

	var orders []*models.Order

	price := requestCtx.Query("price")
	if price != "" {
		priceInt, err := strconv.ParseFloat(price, 64)
		if err != nil {
			return requestCtx.Status(fiber.StatusBadRequest).SendString("Invalid price value")
		}

		filter := func(order *models.Order) bool {
			return order.Price > priceInt
		}

		orders, err = c.orderService.GetOrdersWithFilter(backgroundCtx, user.ID, filter)
		if err != nil {
			return requestCtx.Status(fiber.StatusInternalServerError).SendString("Error filtering orders")
		}
	} else {
		var err error
		orders, err = c.orderService.GetOrders(backgroundCtx, user.ID)
		if err != nil {
			return requestCtx.Status(fiber.StatusInternalServerError).SendString("Error loading orders")
		}
	}

	orderResponses := make([]*transports.OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = transports.MapToOrderResponse(*order)
	}
	return requestCtx.Status(fiber.StatusOK).JSON(orderResponses)

}

// GetOrder handles "/orders/{id}" with method "GET"
func (c *OrdersController) GetOrder(requestCtx fiber.Ctx) error {

	orderId := requestCtx.Params("id")
	logger := log.GetFiberLogger(requestCtx).With().Str("orderId", orderId).Logger()
	log.SetFiberLogger(requestCtx, &logger)
	backgroundCtx := log.NewBackgroundContext(&logger)
	utils.LogAction(backgroundCtx, compOrdersController, "GetOrder")

	oid, err := strconv.Atoi(orderId)
	if err != nil {
		return requestCtx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	backgroundCtx = context.WithValue(backgroundCtx, "orderId", oid)

	user := requestCtx.Locals(constants.AuthenticatedUserKey).(models.User)
	backgroundCtx = context.WithValue(backgroundCtx, constants.AuthenticatedUserKey, &user)
	backgroundCtx = context.WithValue(backgroundCtx, constants.AuthenticatedUserIdKey, user.ID)

	order, err := c.orderService.GetOrder(backgroundCtx, user.ID, oid)
	if err != nil {
		return requestCtx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	orderResponse := transports.MapToOrderResponse(*order)

	return requestCtx.Status(fiber.StatusOK).JSON(orderResponse)
}
