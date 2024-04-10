package events

import (
	"context"
	"townwatch/base"
	"townwatch/models"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
)

func storeEvents(b *base.Base, ctx context.Context, eventParams *[]models.CreateTempEventsParams) *utils.CError {
	tx, err := b.DB.Pool.Begin(ctx)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	defer tx.Rollback(ctx)
	qtx := b.DB.Queries.WithTx(tx)

	if err := qtx.CreateTempEventsTable(ctx); err != nil {
		eventID := sentry.CaptureException(err)
		return &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	_, errInsert := qtx.CreateTempEvents(ctx, *eventParams)
	if errInsert != nil {
		eventID := sentry.CaptureException(errInsert)
		return &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   errInsert,
		}
	}

	if err := qtx.MoveFromTempEventsToEvents(ctx); err != nil {
		eventID := sentry.CaptureException(err)
		return &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}

	if err := tx.Commit(ctx); err != nil {
		eventID := sentry.CaptureException(err)
		return &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}

	return nil
}

func removeEventParamsDuplicates(params *[]models.CreateTempEventsParams) *[]models.CreateTempEventsParams {
	uniqueMap := make(map[string]models.CreateTempEventsParams)
	for _, param := range *params {
		// Only add to the map if ExternalID doesn't exist already
		if _, ok := uniqueMap[param.ExternalID]; !ok {
			uniqueMap[param.ExternalID] = param
		}
	}
	uniqueParams := make([]models.CreateTempEventsParams, 0, len(uniqueMap))
	for _, v := range uniqueMap {
		uniqueParams = append(uniqueParams, v)
	}
	return &uniqueParams
}
