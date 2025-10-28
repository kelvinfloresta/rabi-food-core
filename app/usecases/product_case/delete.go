package product_case

import "context"

func (c *ProductCase) Delete(ctx context.Context, id string) (bool, error) {
	return c.gateway.Delete(id)
}
