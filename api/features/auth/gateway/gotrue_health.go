package gateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

func (g *GoTrueGatewayAdapter) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.baseURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrAuthServiceFailure, err)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrAuthServiceFailure, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return wrapResponseError(errs.ErrAuthServiceFailure, resp)
	}

	return nil
}
