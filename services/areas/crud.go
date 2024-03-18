package areas

import (
	"townwatch/base"
	"townwatch/models"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func CreateArea(b *base.Base, ctx *gin.Context, params *models.CreateAreaParams) *base.CError {
	err := b.DB.Queries.CreateArea(ctx, *params)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return &base.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	return nil
}

func UpdateArea(b *base.Base, ctx *gin.Context, params *models.UpdateAreaParams) *base.CError {
	err := b.DB.Queries.UpdateArea(ctx, *params)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return &base.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	return nil
}

func ReadArea(b *base.Base, ctx *gin.Context, params *models.GetAreaParams) (*models.Area, *base.CError) {
	area, err := b.DB.Queries.GetArea(ctx, *params)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &base.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	return &area, nil
}

func DeleteArea(b *base.Base, ctx *gin.Context, params *models.DeleteAreaParams) *base.CError {
	err := b.DB.Queries.DeleteArea(ctx, *params)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return &base.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	return nil
}

type GetAreasByUserParams struct {
	UserID string
}

func ReadAreasByUser(b *base.Base, ctx *gin.Context, params *GetAreasByUserParams) (*[]models.Area, *base.CError) {
	areas, err := b.DB.Queries.GetAreasByUser(ctx, params.UserID)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &base.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	return &areas, nil
}
