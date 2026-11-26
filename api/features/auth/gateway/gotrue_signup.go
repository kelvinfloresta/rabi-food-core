package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
)

type adminCreateUserRequest struct {
	Email        string         `json:"email"`
	Password     string         `json:"password"`
	EmailConfirm bool           `json:"email_confirm"`
	UserMetadata map[string]any `json:"user_metadata"`
	AppMetadata  map[string]any `json:"app_metadata"`
}

func (g *GoTrueGatewayAdapter) CreateUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	body := adminCreateUserRequest{
		Email:        input.Email,
		Password:     input.Password,
		EmailConfirm: true,
		UserMetadata: map[string]any{
			"name":         input.Name,
			"phone":        input.Phone,
			"tax_id":       input.TaxID,
			"social_id":    input.SocialID,
			"city":         input.City,
			"state":        input.State,
			"zip":          input.ZIP,
			"street":       input.Street,
			"complement":   input.Complement,
			"neighborhood": input.Neighborhood,
		},
		AppMetadata: map[string]any{
			"role": input.Role,
		},
	}

	raw, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrAuthUserCreateFailed, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/admin/users", bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrAuthUserCreateFailed, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.serviceKey)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrAuthUserCreateFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnprocessableEntity {
		body, _ := io.ReadAll(resp.Body)
		if strings.Contains(string(body), "email_exists") {
			return nil, errs.ErrUserAlreadyExists
		}

		return nil, fmt.Errorf("%w: status=422 body=%s", errs.ErrAuthUserCreateFailed, body)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, wrapResponseError(errs.ErrAuthUserCreateFailed, resp)
	}

	var result CreateUserOutput
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrAuthUserCreateFailed, err)
	}

	return &result, nil
}
