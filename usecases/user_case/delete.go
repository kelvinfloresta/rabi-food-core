package user_case

import "context"

func (c *UserCase) Delete(ctx context.Context, id string) (bool, error) {
	return c.gateway.Delete(id)
}
