package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pintu-crypto/sre-playground/otel-testing/services/internal/database"
	"github.com/pintu-crypto/sre-playground/otel-testing/services/internal/utils"
	logger "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type ServiceHTTP struct {
	queries *database.Queries
}

func NewService(queries *database.Queries) *ServiceHTTP {
	return &ServiceHTTP{
		queries: queries,
	}
}

func (s *ServiceHTTP) RegisterHandlers(router *gin.Engine) {
	router.POST("/authors", s.CreateAuthor)
	router.GET("/authors/:id", s.GetAuthor)
	router.PUT("/authors/:id", s.FullUpdateAuthor)
	router.PATCH("/authors/:id", s.PartialUpdateAuthor)
	router.DELETE("/authors/:id", s.DeleteAuthor)
	router.GET("/authors", s.ListAuthors)
}

type apiAuthor struct {
	ID   int64
	Name string `json:"name,omitempty" binding:"required,max=32"`
	Bio  string `json:"bio,omitempty" binding:"required"`
}

type apiAuthorPartialUpdate struct {
	Name *string `json:"name,omitempty" binding:"omitempty,max=32"`
	Bio  *string `json:"bio,omitempty" binding:"omitempty"`
}

type pathParameters struct {
	ID int64 `uri:"id" binding:"required"`
}

func (s *ServiceHTTP) CreateAuthor(ctx *gin.Context) {
	var request apiAuthor

	_, span := otel.Tracer("CreateAuthor").Start(ctx.Request.Context(), "CreateAuthor Func()", oteltrace.WithAttributes(attribute.String("CreateAuthor() Func", "Excecuted")))
	defer span.End()

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := database.CreateAuthorParams{
		Name: request.Name,
		Bio:  request.Bio,
	}
	author, err := s.queries.CreateAuthor(ctx.Request.Context(), params)
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	response := fromDB(author)
	ctx.IndentedJSON(http.StatusCreated, response)

	logger.WithFields(logger.Fields{
		"dd.trace_id": utils.ConvertHexId(span.SpanContext().TraceID().String()),
		"dd.span_id":  utils.ConvertHexId(span.SpanContext().SpanID().String()),
	}).WithContext(ctx).Info("on CreateAuthorFunc() Excecuted ")
}

func (s *ServiceHTTP) GetAuthor(ctx *gin.Context) {
	var pathParameters pathParameters

	_, span := otel.Tracer("GetAuthor").Start(ctx.Request.Context(), "GetAuthor Func()", oteltrace.WithAttributes(attribute.String("GetAuthor() Func", "Excecuted")))
	defer span.End()

	if err := ctx.ShouldBindUri(&pathParameters); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	author, err := s.queries.GetAuthor(ctx.Request.Context(), pathParameters.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	response := fromDB(author)
	ctx.IndentedJSON(http.StatusOK, response)

	logger.WithFields(logger.Fields{
		"dd.trace_id": utils.ConvertHexId(span.SpanContext().TraceID().String()),
		"dd.span_id":  utils.ConvertHexId(span.SpanContext().SpanID().String()),
	}).WithContext(ctx).Info("on GetAuthorFunc() Excecuted ")
}

func (s *ServiceHTTP) FullUpdateAuthor(ctx *gin.Context) {
	var pathParameters pathParameters

	_, span := otel.Tracer("FullUpdateAuthor").Start(ctx.Request.Context(), "FullUpdateAuthor Func()", oteltrace.WithAttributes(attribute.String("FullUpdateAuthor() Func", "Excecuted")))
	defer span.End()

	if err := ctx.ShouldBindUri(&pathParameters); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var request apiAuthor
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parameter := database.UpdateAuthorParams{
		ID:   pathParameters.ID,
		Name: request.Name,
		Bio:  request.Bio,
	}

	author, err := s.queries.UpdateAuthor(ctx.Request.Context(), parameter)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
	}

	response := fromDB(author)
	ctx.IndentedJSON(http.StatusOK, response)

	logger.WithFields(logger.Fields{
		"dd.trace_id": utils.ConvertHexId(span.SpanContext().TraceID().String()),
		"dd.span_id":  utils.ConvertHexId(span.SpanContext().SpanID().String()),
	}).WithContext(ctx).Info("on FullUpdateAuthorFunc() Excecuted ")
}

func (s *ServiceHTTP) PartialUpdateAuthor(ctx *gin.Context) {
	var pathParameters pathParameters

	_, span := otel.Tracer("PartialUpdateAuthor").Start(ctx.Request.Context(), "PartialUpdateAuthor Func()", oteltrace.WithAttributes(attribute.String("PartialAuthor() Func", "Excecuted")))
	defer span.End()

	if err := ctx.ShouldBindUri(&pathParameters); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var request apiAuthorPartialUpdate
	if error := ctx.ShouldBindJSON(&request); error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": error})
		return
	}

	parameter := database.PartialUpdateAuthorParams{ID: pathParameters.ID}
	if request.Name != nil {
		parameter.UpdateName = true
		parameter.Name = *request.Name
	}
	if request.Bio != nil {
		parameter.UpdateBio = true
		parameter.Bio = *request.Bio
	}

	author, err := s.queries.PartialUpdateAuthor(ctx.Request.Context(), parameter)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	response := fromDB(author)
	ctx.IndentedJSON(http.StatusOK, response)

	logger.WithFields(logger.Fields{
		"dd.trace_id": utils.ConvertHexId(span.SpanContext().TraceID().String()),
		"dd.span_id":  utils.ConvertHexId(span.SpanContext().SpanID().String()),
	}).WithContext(ctx).Info("on PartialUpdateAuthorFunc() Excecuted ")
}

func (s *ServiceHTTP) ListAuthors(ctx *gin.Context) {
	_, span := otel.Tracer("ListAuthor").Start(ctx.Request.Context(), "ListAuthor Func()", oteltrace.WithAttributes(attribute.String("ListAuthor() Func", "Excecuted")))
	defer span.End()

	authors, err := s.queries.ListAuthors(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	if len(authors) == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	var response []*apiAuthor
	for _, author := range authors {
		response = append(response, fromDB(author))
	}

	ctx.IndentedJSON(http.StatusOK, authors)

	logger.WithFields(logger.Fields{
		"dd.trace_id": utils.ConvertHexId(span.SpanContext().TraceID().String()),
		"dd.span_id":  utils.ConvertHexId(span.SpanContext().SpanID().String()),
	}).WithContext(ctx).Info("on ListAuthorFunc() Excecuted ")

}

func (s *ServiceHTTP) DeleteAuthor(ctx *gin.Context) {
	var pathParameters pathParameters

	_, span := otel.Tracer("DeleteAuthor").Start(ctx.Request.Context(), "DeleteAuthor Func()", oteltrace.WithAttributes(attribute.String("DeleteAuthor() Func", "Excecuted")))
	defer span.End()

	if err := ctx.ShouldBindUri(&pathParameters); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.queries.DeleteAuthor(ctx.Request.Context(), pathParameters.ID); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusServiceUnavailable, gin.H{"errors": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)

	logger.WithFields(logger.Fields{
		"dd.trace_id": utils.ConvertHexId(span.SpanContext().TraceID().String()),
		"dd.span_id":  utils.ConvertHexId(span.SpanContext().SpanID().String()),
	}).WithContext(ctx).Info("on DeleteAuthorFunc() Excecuted ")
}

func fromDB(author database.Author) *apiAuthor {
	return &apiAuthor{
		ID:   author.ID,
		Name: author.Name,
		Bio:  author.Bio,
	}
}
