package weatherlocation

import (
	"context"
	"errors"
	"github.com/lclpedro/weather-location/pkg/clients/viacep"
	weatherClient "github.com/lclpedro/weather-location/pkg/clients/weather"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/gofiber/fiber/v2"
	weatherlocation "github.com/lclpedro/weather-location/internal/scaffold/services/weather_location"
	"github.com/lclpedro/weather-location/pkg/requester"
	"github.com/spf13/viper"
)

type view struct {
	WeatherLocationService weatherlocation.Service
	Configs                *viper.Viper
	Tracer                 trace.Tracer
}

type View interface {
	WeatherLocationHandler(*fiber.Ctx) error
}

func NewView(tracer trace.Tracer, weatherLocationService weatherlocation.Service) View {
	return &view{
		WeatherLocationService: weatherLocationService,
		Configs:                viper.GetViper(),
		Tracer:                 tracer,
	}
}

func (v *view) WeatherLocationHandler(c *fiber.Ctx) error {
	carrier := propagation.HeaderCarrier(c.GetReqHeaders())
	ctx := context.Background()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, spam := v.Tracer.Start(ctx, "weather-location-api")
	defer spam.End()

	requesterViaCep := requester.NewRequester(ctx, v.Tracer, v.Configs.GetInt(viacep.ViaCEPTimeout))
	requesterWeather := requester.NewRequester(ctx, v.Tracer, v.Configs.GetInt(weatherClient.WeatherAPITimeout))

	v.WeatherLocationService.SetClients(
		viacep.NewClient(requesterViaCep),
		weatherClient.NewClient(requesterWeather),
	)

	cep := c.Params("cep")
	weather, err := v.WeatherLocationService.GetWeatherLocation(ctx, cep)
	if errors.Is(err, viacep.ErrInvalidCep) {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": err.Error()})
	}

	if errors.Is(err, viacep.ErrNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(weather)
}
