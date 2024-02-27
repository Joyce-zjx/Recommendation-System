package routes

import (
	api "Recommendation-System/external/api"
	"Recommendation-System/repository"
	"Recommendation-System/routes/entity"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler interface {
	recommendEvents(ctx *gin.Context)
	listEvents(ctx *gin.Context)
	likeEvent(ctx *gin.Context)
	dislikeEvent(ctx *gin.Context)
}

func NewUserHandler(logger *zap.Logger, repo repository.UserRepo, mapUtil api.MapUtilities, eventSearcher api.EventsSearcher) UserHandler {
	return &userHandler{
		logger:        logger,
		userRepo:      repo,
		mapUtilities:  mapUtil,
		eventSearcher: eventSearcher,
	}
}

type userHandler struct {
	logger        *zap.Logger
	userRepo      repository.UserRepo
	mapUtilities  api.MapUtilities
	eventSearcher api.EventsSearcher
}

type listEventsResp struct {
	Events []entity.Event `json:"events"`
}

func (u *userHandler) listEvents(ctx *gin.Context) {
	latString := ctx.Query("lat")
	if latString == "" {
		u.logger.Warn("lat is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "lat is required"})
		return
	}
	longString := ctx.Query("long")
	if longString == "" {
		u.logger.Warn("long is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "long is required"})
		return
	}
	lat, _ := strconv.ParseFloat(latString, 64)
	long, _ := strconv.ParseFloat(longString, 64)
	u.logger.Info("list events", zap.Float64("lat", lat), zap.Float64("long", long))
	eventsResp, err := u.eventSearcher.ListEvents(lat, long, 0, "")
	if err != nil {
		u.logger.Error("failed to get events", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get events"})
		return
	}

	// compose responses
	eventsToOccur := make(map[string]entity.Event)
	events := make([]entity.Event, 0)
	for _, event := range eventsResp.Embedded.Events {
		var imageURL string
		if len(event.Images) > 0 {
			imageURL = event.Images[0].URL
		}
		long, _ := strconv.ParseFloat(event.EmbeddedVenues.Venues[0].Location.Longitude, 64)
		lat, _ := strconv.ParseFloat(event.EmbeddedVenues.Venues[0].Location.Latitude, 64)
		event := entity.Event{
			ID:   event.ID,
			Name: event.Name,
			Date: entity.Date{
				LocalDate: event.Dates.Start.LocalDate,
				LocalTime: event.Dates.Start.LocalTime,
				Timezone:  event.Dates.Timezone,
				Status:    event.Dates.Status.Code,
			},
			Address: event.EmbeddedVenues.Venues[0].Address.Line1 + ", " + event.
				EmbeddedVenues.Venues[0].City.Name + ", " +
				event.EmbeddedVenues.Venues[0].State.Name + ", " + event.EmbeddedVenues.Venues[0].Country.Name,
			Genres:    []string{event.Classifications[0].Genre.Name, event.Classifications[0].SubGenre.Name},
			ImageURL:  imageURL,
			URL:       event.URL,
			Longitude: long,
			Latitude:  lat,
		}

		// remove duplicates
		if _, ok := eventsToOccur[event.Name]; !ok {
			eventsToOccur[event.Name] = event
			events = append(events, event)
		}
	}
	// calculate distance
	for i := range events {
		events[i].Distance, _ = u.mapUtilities.CalculateDistance(lat, long, events[i].Latitude, events[i].Longitude)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Distance < events[j].Distance
	})
	ctx.JSON(http.StatusOK, listEventsResp{Events: events})
}

type likeReq struct {
	Genre    string `json:"genre" binding:"required"`
	SubGenre string `json:"subgenre" binding:"required"`
}

type dislikeReq struct {
	Genre    string `json:"genre" binding:"required"`
	SubGenre string `json:"subgenre" binding:"required"`
}

func (u *userHandler) likeEvent(ctx *gin.Context) {
	userID := "1111"
	req := &likeReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		u.logger.Warn("invalid request body")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	genre := req.Genre
	subgenre := req.SubGenre
	u.userRepo.InsertUserLikeGenreAndSubGenre(userID, genre, subgenre)

	ctx.JSON(http.StatusCreated, gin.H{"message": "liked event"})
}

func (u *userHandler) dislikeEvent(ctx *gin.Context) {
	// userID, _ := ctx.Get(_ctxKey_UserID)
	userID := "1111"

	req := &dislikeReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		u.logger.Warn("invalid request body")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	genre := req.Genre
	subgenre := req.SubGenre
	u.userRepo.InsertUserDisLikeGenreAndSubGenre(userID, genre, subgenre)

	ctx.JSON(http.StatusOK, gin.H{"message": "disliked event"})
}

func (u *userHandler) recommendEvents(ctx *gin.Context) {
	userID := "1111"
	likedGenreAndSubgenre, err := u.userRepo.SelectUserLikedGenresAndSubGenre(userID)
	if err != nil {
		u.logger.Warn("failed to get liked genres", zap.Error(err))
		ctx.JSON(http.StatusOK, gin.H{"Warn": "failed to get disliked genres. No recommendations will be made"})
		return
	}
	dislikedGenreAndSubgenre, err := u.userRepo.SelectUserDisLikedGenreAndSubGenre(userID)
	if err != nil {
		u.logger.Warn("failed to get disliked genres", zap.Error(err))
		ctx.JSON(http.StatusOK, gin.H{"Warn": "failed to get disliked genres. No recommendations will be made"})
		return
	}
	likedGenreToCount := make(map[string]int)
	likedSubGenreToCount := make(map[string]int)

	var mostLikedGenre, mostLikedSubGenre string
	var mostLikedGenreCount, mostLikedSubGenreCount int
	for _, genre := range likedGenreAndSubgenre {
		likedGenreToCount[genre.Genre]++
		likedSubGenreToCount[genre.SubGenre]++
		if likedGenreToCount[genre.Genre] > mostLikedGenreCount {
			mostLikedGenre = genre.Genre
			mostLikedGenreCount = likedGenreToCount[genre.Genre]
		}
		if likedSubGenreToCount[genre.SubGenre] > mostLikedSubGenreCount {
			mostLikedSubGenre = genre.SubGenre
			mostLikedSubGenreCount = likedSubGenreToCount[genre.SubGenre]
		}
	}
	u.logger.Info("most liked genre", zap.String("genre", mostLikedGenre))
	dislikedGenreToCount := make(map[string]int)
	dislikedSubGenreToCount := make(map[string]int)

	var mostDislikedGenre, mostDislikedSubGenre string
	var mostDislikedGenreCount, mostDislikedSubGenreCount int
	for _, genre := range dislikedGenreAndSubgenre {
		dislikedGenreToCount[genre.Genre]++
		dislikedSubGenreToCount[genre.SubGenre]++
		if dislikedGenreToCount[genre.Genre] > mostDislikedGenreCount {
			mostDislikedGenre = genre.Genre
			mostDislikedGenreCount = dislikedGenreToCount[genre.Genre]
		}
		if dislikedSubGenreToCount[genre.SubGenre] > mostDislikedSubGenreCount {
			mostDislikedSubGenre = genre.SubGenre
			mostDislikedSubGenreCount = dislikedSubGenreToCount[genre.SubGenre]
		}
	}
	u.logger.Info("most disliked genre", zap.String("genre", mostDislikedGenre))

	// search with keyword equals to most liked genre or subgenre

	latString := ctx.Query("lat")
	longString := ctx.Query("long")
	lat, _ := strconv.ParseFloat(latString, 64)
	long, _ := strconv.ParseFloat(longString, 64)

	if err != nil {
		u.logger.Error("failed to get lat long", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to get lat long"})
		return
	}

	resp, err := u.eventSearcher.ListEvents(lat, long, 0, mostLikedSubGenre)
	if err != nil {
		u.logger.Error("failed to get events", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to get events"})
		return
	}
	u.logger.Info("events", zap.Any("events number", len(resp.Embedded.Events)))

	// compose responses
	eventsToOccur := make(map[string]*entity.Event)
	events := make([]entity.Event, 0)
	for _, event := range resp.Embedded.Events {
		// filter out events that are disliked
		if event.Classifications[0].Genre.Name == mostDislikedGenre || event.Classifications[0].SubGenre.Name == mostDislikedSubGenre {
			continue
		}
		var imageURL string
		if len(event.Images) > 0 {
			imageURL = event.Images[0].URL
		}
		long, _ := strconv.ParseFloat(event.EmbeddedVenues.Venues[0].Location.Longitude, 64)
		lat, _ := strconv.ParseFloat(event.EmbeddedVenues.Venues[0].Location.Latitude, 64)
		event := entity.Event{
			ID:   event.ID,
			Name: event.Name,
			Date: entity.Date{
				LocalDate: event.Dates.Start.LocalDate,
				LocalTime: event.Dates.Start.LocalTime,
				Timezone:  event.Dates.Timezone,
				Status:    event.Dates.Status.Code,
			},
			Address: event.EmbeddedVenues.Venues[0].Address.Line1 + ", " + event.
				EmbeddedVenues.Venues[0].City.Name + ", " +
				event.EmbeddedVenues.Venues[0].State.Name + ", " + event.EmbeddedVenues.Venues[0].Country.Name,
			Genres:    []string{event.Classifications[0].Genre.Name, event.Classifications[0].SubGenre.Name},
			ImageURL:  imageURL,
			URL:       event.URL,
			Longitude: long,
			Latitude:  lat,
		}

		// remove duplicates
		if _, ok := eventsToOccur[event.Name]; !ok {
			events = append(events, event)
			eventsToOccur[event.Name] = &event
		}
	}
	// calculate distance
	for i := range events {
		events[i].Distance, _ = u.mapUtilities.CalculateDistance(lat, long, events[i].Latitude, events[i].Longitude)
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].Distance < events[j].Distance
	})
	ctx.JSON(http.StatusOK, listEventsResp{Events: events})
}
