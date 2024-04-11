package events

import (
	"context"
	"time"
	"townwatch/base"
	"townwatch/models"
	"townwatch/services/events/arcgis"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
)

func fetchAndStoreEvents(b *base.Base, ctx context.Context, fromDate time.Time, toDate time.Time) (int64, *utils.CError) {
	tempEventParams := []models.CreateTempEventsParams{}

	if events, cerr := arcgis.FetchAndConverYorkEvents(b, ctx, fromDate, toDate); cerr == nil {
		tempEventParams = append(tempEventParams, *events...)
	}
	if events, cerr := arcgis.FetchAndConverHaltonEvents(b, ctx, fromDate, toDate); cerr == nil {
		tempEventParams = append(tempEventParams, *events...)
	}
	if events, cerr := arcgis.FetchAndConverPeelEvents(b, ctx, fromDate, toDate); cerr == nil {
		tempEventParams = append(tempEventParams, *events...)
	}
	if events, cerr := arcgis.FetchAndConverTorontoEvents(b, ctx, fromDate, toDate); cerr == nil {
		tempEventParams = append(tempEventParams, *events...)
	}

	return storeEvents(b, ctx, tempEventParams)
}

func storeEvents(b *base.Base, ctx context.Context, eventParams []models.CreateTempEventsParams) (int64, *utils.CError) {

	eventParams = removeEventParamsDuplicates(eventParams)

	tx, err := b.DB.Pool.Begin(ctx)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return 0, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	defer tx.Rollback(ctx)
	qtx := b.DB.Queries.WithTx(tx)

	if err := qtx.CreateTempEventsTable(ctx); err != nil {
		eventID := sentry.CaptureException(err)
		return 0, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	count, errInsert := qtx.CreateTempEvents(ctx, eventParams)
	if errInsert != nil {
		eventID := sentry.CaptureException(errInsert)
		return 0, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   errInsert,
		}
	}

	if err := qtx.MoveFromTempEventsToEvents(ctx); err != nil {
		eventID := sentry.CaptureException(err)
		return 0, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}

	if err := tx.Commit(ctx); err != nil {
		eventID := sentry.CaptureException(err)
		return 0, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}

	return count, nil
}

func removeEventParamsDuplicates(params []models.CreateTempEventsParams) []models.CreateTempEventsParams {
	uniqueMap := make(map[string]models.CreateTempEventsParams)
	for _, param := range params {
		// Only add to the map if ExternalID doesn't exist already
		if _, ok := uniqueMap[param.ExternalID]; !ok {
			uniqueMap[param.ExternalID] = param
		}
	}
	uniqueParams := make([]models.CreateTempEventsParams, 0, len(uniqueMap))
	for _, v := range uniqueMap {
		uniqueParams = append(uniqueParams, v)
	}

	return uniqueParams
}
