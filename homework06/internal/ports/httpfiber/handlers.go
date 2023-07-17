package httpfiber

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"homework6/internal/app"
)

// Метод для создания объявления (ad)
func createAd(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody createAdRequest
		err := c.BodyParser(&reqBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		ad, createErr := a.CreateAd(c.Context(), reqBody.Title, reqBody.Text, reqBody.UserID)

		switch createErr {
		case app.ValidationError:
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(createErr))
		case nil:
			return c.JSON(AdSuccessResponse(&ad))
		default:
			c.Status(http.StatusInternalServerError)
			return c.JSON(AdErrorResponse(createErr))
		}
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody changeAdStatusRequest
		if err := c.BodyParser(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		adID, err := c.ParamsInt("ad_id")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		ad, changeErr := a.ChangeAdStatus(c.Context(), int64(adID), reqBody.UserID, reqBody.Published)

		switch changeErr {
		case app.WrongUserError:
			c.Status(http.StatusForbidden)
			return c.JSON(AdErrorResponse(changeErr))
		case nil:
			return c.JSON(AdSuccessResponse(&ad))
		default:
			c.Status(http.StatusInternalServerError)
			return c.JSON(AdErrorResponse(changeErr))
		}
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody updateAdRequest
		if err := c.BodyParser(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		adID, err := c.ParamsInt("ad_id")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		ad, updateErr := a.UpdateAd(c.Context(), int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)

		switch updateErr {
		case app.WrongUserError:
			c.Status(http.StatusForbidden)
			return c.JSON(AdErrorResponse(updateErr))
		case app.ValidationError:
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(updateErr))
		case nil:
			return c.JSON(AdSuccessResponse(&ad))
		default:
			c.Status(http.StatusInternalServerError)
			return c.JSON(AdErrorResponse(updateErr))
		}
	}
}
