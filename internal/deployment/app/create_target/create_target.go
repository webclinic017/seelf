package create_target

import (
	"context"

	auth "github.com/YuukanOO/seelf/internal/auth/domain"
	"github.com/YuukanOO/seelf/internal/deployment/domain"
	"github.com/YuukanOO/seelf/pkg/bus"
	"github.com/YuukanOO/seelf/pkg/monad"
	"github.com/YuukanOO/seelf/pkg/validate"
	"github.com/YuukanOO/seelf/pkg/validate/strings"
)

type Command struct {
	bus.Command[string]

	Name     string              `json:"name"`
	Url      monad.Maybe[string] `json:"url"`
	Provider any                 `json:"-"`
}

func (Command) Name_() string { return "deployment.command.create_target" }

func Handler(
	reader domain.TargetsReader,
	writer domain.TargetsWriter,
	provider domain.Provider,
) bus.RequestHandler[string, Command] {
	return func(ctx context.Context, cmd Command) (string, error) {
		var targetUrl domain.Url

		if err := validate.Struct(validate.Of{
			"name": validate.Field(cmd.Name, strings.Required),
			"url": validate.Maybe(cmd.Url, func(url string) error {
				return validate.Value(url, &targetUrl, domain.UrlFrom)
			}),
		}); err != nil {
			return "", err
		}

		config, err := provider.Prepare(ctx, cmd.Provider)

		if err != nil {
			return "", err
		}

		// Validate availability of both the target domain and the config
		var urlRequirement domain.TargetUrlRequirement

		if cmd.Url.HasValue() {
			urlRequirement, err = reader.CheckUrlAvailability(ctx, targetUrl)

			if err != nil {
				return "", err
			}
		}

		configRequirement, err := reader.CheckConfigAvailability(ctx, config)

		if err != nil {
			return "", err
		}

		if err = validate.Struct(validate.Of{
			"url":         validate.If(cmd.Url.HasValue(), urlRequirement.Error),
			config.Kind(): configRequirement.Error(),
		}); err != nil {
			return "", err
		}

		target, err := domain.NewTarget(
			cmd.Name,
			configRequirement,
			auth.CurrentUser(ctx).MustGet(),
		)

		if err != nil {
			return "", err
		}

		if cmd.Url.HasValue() {
			if err = target.ExposeServicesAutomatically(urlRequirement); err != nil {
				return "", err
			}
		}

		if err = writer.Write(ctx, &target); err != nil {
			return "", err
		}

		return string(target.ID()), nil
	}
}
