package category_case

import "context"

func (c *CategoryCase) Delete(ctx context.Context, id string) (bool, error) {
	return c.gateway.Delete(id)
}
