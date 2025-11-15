package order_case

import (
	"context"
	"fmt"
	"rabi-food-core/app_context"
	"rabi-food-core/domain/order"
	g "rabi-food-core/libs/database/gateways/order_gateway"
	"rabi-food-core/libs/database/gateways/product_gateway"
	"rabi-food-core/libs/errs"
	"rabi-food-core/libs/logger"

	"github.com/google/uuid"
)

type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  uint   `json:"quantity"`
}

type CreateInput struct {
	Items []OrderItem `json:"items"`
	Notes string      `json:"notes"`
}

func (c *OrderCase) Create(ctx context.Context, input CreateInput) (string, error) {
	productIds := make([]string, 0, len(input.Items))
	for _, item := range input.Items {
		productIds = append(productIds, item.ProductID)
	}

	products, err := c.productCase.List(ctx, productIds)
	if err != nil {
		return "", fmt.Errorf("failed to fetch products: %w", err)
	}

	if len(products) == 0 {
		return "", errs.ErrProductNotFound
	}

	if len(products) != len(input.Items) {
		logger.Get(ctx).Warn().Msgf("some products not found for the given IDs: %v, found: %v", productIds, products)
		return "", errs.ErrProductNotFound
	}

	productMap := make(map[string]product_gateway.ListOutput)
	for _, product := range products {
		productMap[product.ID] = product
	}

	orderItems := make([]g.OrderItem, 0, len(input.Items))
	totalPrice := uint(0)
	for _, item := range input.Items {
		product, exists := productMap[item.ProductID]
		if !exists {
			return "", fmt.Errorf("product with ID %s not found", item.ProductID)
		}

		itemTotal := product.Price * item.Quantity
		orderItems = append(orderItems, g.OrderItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			UnitPrice:   product.Price,
			Total:       itemTotal,
		})
		totalPrice += itemTotal
	}

	session := app_context.GetSession(ctx)
	id, err := c.gateway.Create(g.CreateInput{
		UserID:            session.UserID,
		TenantID:          session.TenantID,
		Code:              uuid.NewString(),
		PaymentStatus:     order.PaymentPending,
		FulfillmentStatus: order.FulfillmentPending,
		DeliveryStatus:    order.DeliveryPending,
		Notes:             input.Notes,
		TotalPrice:        totalPrice,
		Items:             orderItems,
	})

	if err != nil {
		return "", err
	}

	logger.L().Info().Str("tenant", session.TenantID).Str("order", id).Msg("order created")

	return id, nil
}
