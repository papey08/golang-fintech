package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework10/internal/app"
	"homework10/internal/model/ads"
	"homework10/internal/model/errs"
	"homework10/internal/model/filter"
	"net/http"
)

// Метод для создания объявления (ad)
func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(err))
		}

		newAd, createErr := a.CreateAd(c, reqBody.Title, reqBody.Text, reqBody.UserID)

		switch createErr {
		case errs.ValidationError:
			c.JSON(http.StatusBadRequest, ErrorResponse(createErr))
		case errs.UserNotExist:
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse(createErr))
		case nil:
			c.JSON(http.StatusOK, AdSuccessResponse(&newAd))
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse(createErr))
		}
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(err))
		}

		adID, err := getID(c, "ad_id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(err))
		}

		ad, changeErr := a.ChangeAdStatus(c, adID, reqBody.UserID, reqBody.Published)

		switch changeErr {
		case errs.WrongUserError:
			c.JSON(http.StatusForbidden, ErrorResponse(changeErr))
		case errs.AdNotExist:
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse(changeErr))
		case errs.UserNotExist:
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse(changeErr))
		case nil:
			c.JSON(http.StatusOK, AdSuccessResponse(&ad))
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse(changeErr))
		}
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(err))
		}

		adID, err := getID(c, "ad_id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(err))
		}

		ad, updateErr := a.UpdateAd(c, adID, reqBody.UserID, reqBody.Title, reqBody.Text)

		switch updateErr {
		case errs.WrongUserError:
			c.JSON(http.StatusForbidden, ErrorResponse(updateErr))
		case errs.ValidationError:
			c.JSON(http.StatusBadRequest, ErrorResponse(updateErr))
		case errs.AdNotExist:
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse(updateErr))
		case errs.UserNotExist:
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse(updateErr))
		case nil:
			c.JSON(http.StatusOK, AdSuccessResponse(&ad))
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse(updateErr))
		}
	}
}

// Метод для получения объявления по его ID
func getAdByID(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		adID, err := getID(c, "ad_id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(err))
		}

		ad, getErr := a.GetAdByID(c, adID)

		switch getErr {
		case errs.AdNotExist:
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse(getErr))
		case nil:
			c.JSON(http.StatusOK, AdSuccessResponse(&ad))
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse(getErr))
		}
	}
}

// Метод для получения списка объявлений с фильтрами
func getAdsList(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody filterRequest
		var f filter.Filter

		if c.Request.ContentLength == 0 {
			f = filter.Filter{
				PublishedBy: false,
				AuthorBy:    false,
				AuthorID:    0,
				DateBy:      false,
				Date:        ads.Date{},
			}
		} else if err := c.BindJSON(&reqBody); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(err))
		} else {
			f = filter.Filter{
				PublishedBy: reqBody.PublishedBy,
				AuthorBy:    reqBody.AuthorBy,
				AuthorID:    reqBody.AuthorID,
				DateBy:      reqBody.DateBy,
				Date:        reqBody.Date,
			}
		}

		adsList, listErr := a.GetAdsList(c, f)

		if listErr != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		} else {
			c.JSON(http.StatusOK, AdsSuccessResponse(adsList))
		}
	}
}

func searchAds(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody searchAdRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(err))
		}

		adsList, searchErr := a.SearchAds(c, reqBody.Pattern)

		if searchErr != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		} else {
			c.JSON(http.StatusOK, AdsSuccessResponse(adsList))
		}
	}
}

func deleteAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody deleteAdRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(err))
		}

		adID, err := getID(c, "ad_id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(err))
		}

		deleteErr := a.DeleteAd(c, reqBody.UserID, adID)

		switch deleteErr {
		case errs.WrongUserError:
			c.JSON(http.StatusForbidden, ErrorResponse(deleteErr))
		case errs.AdNotExist:
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse(deleteErr))
		case errs.UserNotExist:
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse(deleteErr))
		case nil:
			c.JSON(http.StatusOK, nil)
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse(deleteErr))
		}
	}
}
