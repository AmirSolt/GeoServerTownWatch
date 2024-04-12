package areas

import (
	"fmt"
	"strings"
	"townwatch/base"
	"townwatch/models"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func CreateArea(b *base.Base, ctx *gin.Context, params *models.CreateAreaParams) (*models.Area, *utils.CError) {

	count, errc := b.Queries.CountAreasByUser(ctx, params.UserID)
	if errc != nil {
		eventID := sentry.CaptureException(errc)
		cerr := &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   errc,
		}
		return nil, cerr
	}

	if count >= int64(b.MaxAreasByUser) {
		err := fmt.Errorf("user has reached maximum area count")
		eventID := sentry.CaptureException(err)
		cerr := &utils.CError{
			EventID: eventID,
			Message: "you have reached maximum area count",
			Error:   err,
		}
		return nil, cerr
	}

	params.Address = removeSpaceAndCapitalize(params.Address)

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
	area.Point = nil
	return area
}

func censorPostalCode(str string) string {
	newStr := ""
	numUncensored := 3
	for i := 0; i < len(str); i++ {
		if i < numUncensored {
			newStr = newStr + string(str[i])
		} else {
			newStr = newStr + "*"
		}
	}
	return newStr
}

// func roundCoordinates(num float64, decimals int) float64 {
// 	multiplier := math.Pow10(decimals)
// 	rounded := math.Round(num * multiplier)
// 	return rounded / multiplier
// }

func removeSpaceAndCapitalize(input string) string {
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ToUpper(input)
	return input
}
