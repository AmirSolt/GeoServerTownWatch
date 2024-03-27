package areas

import (
	"math"
	"strings"
	"townwatch/base"
	"townwatch/models"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func CreateArea(b *base.Base, ctx *gin.Context, params *models.CreateAreaParams) (*models.Area, *utils.CError) {
	area, err := b.DB.Queries.CreateArea(ctx, *params)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}

	cenArea := CensorArea(area)
	return &cenArea, nil
}

func UpdateArea(b *base.Base, ctx *gin.Context, params *models.UpdateAreaParams) (*models.Area, *utils.CError) {
	area, err := b.DB.Queries.UpdateArea(ctx, *params)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}

	cenArea := CensorArea(area)
	return &cenArea, nil
}

func ReadArea(b *base.Base, ctx *gin.Context, params *models.GetAreaParams) (*models.Area, *utils.CError) {
	area, err := b.DB.Queries.GetArea(ctx, *params)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}

	cenArea := CensorArea(area)
	return &cenArea, nil
}

func DeleteArea(b *base.Base, ctx *gin.Context, params *models.DeleteAreaParams) (*models.Area, *utils.CError) {
	area, err := b.DB.Queries.DeleteArea(ctx, *params)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}

	cenArea := CensorArea(area)
	return &cenArea, nil
}

type GetAreasByUserParams struct {
	UserID string `json:"user_id"`
}

func ReadAreasByUser(b *base.Base, ctx *gin.Context, params *GetAreasByUserParams) (*[]models.Area, *utils.CError) {
	areas, err := b.DB.Queries.GetAreasByUser(ctx, params.UserID)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}

	cenAreas := CensorAreas(areas)
	return &cenAreas, nil
}

func CensorAreas(areas []models.Area) []models.Area {
	cenAreas := []models.Area{}
	for _, area := range areas {
		cenAreas = append(cenAreas, CensorArea(area))
	}
	return cenAreas
}
func CensorArea(area models.Area) models.Area {
	area.Address = censorPostalCode(area.Address)
	area.Lat = roundCoordinates(area.Lat, 3)
	area.Long = roundCoordinates(area.Long, 3)
	area.Point = nil
	return area
}

func censorPostalCode(str string) string {
	if len(str) == 0 {
		return ""
	}
	trimmedStr := strings.ReplaceAll(str, " ", "")
	numUncensored := 3
	censored := make([]byte, len(trimmedStr))
	copy(censored, trimmedStr[:numUncensored])
	for i := numUncensored; i < len(trimmedStr); i++ {
		censored[i] = '#'
	}
	return string(censored)
}

func roundCoordinates(num float64, decimals int) float64 {
	multiplier := math.Pow10(decimals)
	rounded := math.Round(num * multiplier)
	return rounded / multiplier
}
