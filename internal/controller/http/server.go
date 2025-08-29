package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"subscription-service/internal/app/entity"
	"subscription-service/internal/app/usecase"
	"subscription-service/internal/controller/http/gen"
	pkg "subscription-service/internal/pkg/utils"
)

var _ gen.StrictServerInterface = (*Server)(nil)

const (
	defaultReadTimeout     = time.Second * 10
	defaultHeadReadTimeout = time.Second * 5
	defaultWriteTimeout    = time.Second * 15
	defaultIdleTimeout     = time.Minute * 2
)

type Server struct {
	address    string
	subUsecase usecase.SubscriptionUseCase
	pool       *pgxpool.Pool
	logger     *zap.Logger
}

func NewServer(address string, subUsecase usecase.SubscriptionUseCase, pool *pgxpool.Pool, logger *zap.Logger) *Server {
	return &Server{
		address:    address,
		subUsecase: subUsecase,
		pool:       pool,
		logger:     logger,
	}
}

func (r *Server) GetSubscriptions(
	ctx context.Context,
	request gen.GetSubscriptionsRequestObject,
) (gen.GetSubscriptionsResponseObject, error) {
	filter := new(entity.ListSubscriptionFilter)
	if request.Params.ServiceName != nil {
		filter.Title = request.Params.ServiceName
	}

	if request.Params.UserId != nil {
		filter.UserID = pkg.PointerTo(request.Params.UserId.String())
	}

	if request.Params.Price != nil {
		price := *request.Params.Price
		price64 := int64(price)
		filter.Price = &price64
	}

	t_s, err := time.Parse("01-2006", *request.Params.StartDate)
	if err != nil {
		if errors.Is(err, &time.ParseError{}) {
			return gen.GetSubscriptions400JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
		}
		return gen.GetSubscriptions500JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
	}

	startDateWithDay := time.Date(t_s.Year(), t_s.Month(), 1, 0, 0, 0, 0, time.UTC)

	t_e, err := time.Parse("01-2006", *request.Params.EndDate)
	if err != nil {
		if errors.Is(err, &time.ParseError{}) {
			return gen.GetSubscriptions400JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
		}
		return gen.GetSubscriptions500JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
	}

	endDateWithDay := time.Date(t_e.Year(), t_e.Month(), 1, 0, 0, 0, 0, time.UTC)

	filter.StartDate = &startDateWithDay
	filter.EndDate = &endDateWithDay

	filter.Limit = request.Params.Limit
	filter.Offset = request.Params.Offset

	subs, err := r.subUsecase.List(ctx, *filter)
	if err != nil {
		return gen.GetSubscriptions500JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
	}

	resp := make([]gen.Subscription, len(subs))

	for i, s := range subs {
		resp[i] = gen.Subscription{
			Id:          pkg.UUID(s.ID),
			ServiceName: s.Title,
			Price:       int(s.Price),
			UserId:      *pkg.UUID(s.UserID),
			StartDate:   s.StartDate.String(),
			EndDate:     pkg.PointerTo(s.EndDate.String()),
			CreatedAt:   pkg.PointerTo(time.UnixMilli(s.CreatedAt)),
			UpdatedAt:   pkg.PointerTo(time.UnixMilli(s.UpdatedAt)),
		}
	}
	return gen.GetSubscriptions200JSONResponse(resp), nil
}

func (r *Server) PostSubscriptions(
	ctx context.Context,
	request gen.PostSubscriptionsRequestObject,
) (gen.PostSubscriptionsResponseObject, error) {
	filter := new(entity.CreateSubscriptionRequest)
	filter.Title = request.Body.ServiceName

	filter.UserID = request.Body.UserId.String()
	filter.Price = int64(request.Body.Price)

	t_s, err := time.Parse("01-2006", request.Body.StartDate)
	if err != nil {
		if errors.Is(err, &time.ParseError{}) {
			return gen.PostSubscriptions400JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
		}
		return gen.PostSubscriptions500JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
	}

	startDateWithDay := time.Date(t_s.Year(), t_s.Month(), 1, 0, 0, 0, 0, time.UTC)

	t_e, err := time.Parse("01-2006", *request.Body.EndDate)
	if err != nil {
		if errors.Is(err, &time.ParseError{}) {
			return gen.PostSubscriptions400JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
		}
		return gen.PostSubscriptions500JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
	}

	endDateWithDay := time.Date(t_e.Year(), t_e.Month(), 1, 0, 0, 0, 0, time.UTC)

	filter.StartDate = startDateWithDay
	filter.EndDate = &endDateWithDay
	filter.CreatedAt = time.Now().UnixMilli()
	filter.UpdatedAt = time.Now().UnixMilli()

	s, err := r.subUsecase.Create(ctx, *filter)
	if err != nil {
		if errors.Is(err, usecase.ErrSubscriptionAlreadyExists) {
			return gen.PostSubscriptions400JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
		}
		return gen.PostSubscriptions500JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
	}
	created := time.UnixMilli(s.CreatedAt)
	updated := time.UnixMilli(s.UpdatedAt)

	return gen.PostSubscriptions201JSONResponse{
		Id:          pkg.PointerTo(*pkg.UUID(s.ID)),
		ServiceName: s.Title,
		Price:       int(s.Price),
		StartDate:   s.StartDate.Format("01-2006"),
		EndDate: func() *string {
			if s.EndDate != nil {
				str := s.EndDate.Format("01-2006")
				return &str
			}
			return nil
		}(),
		UserId:    *pkg.UUID(s.UserID),
		CreatedAt: &created,
		UpdatedAt: &updated,
	}, nil
}

func (r *Server) GetSubscriptionsSum(
	ctx context.Context,
	request gen.GetSubscriptionsSumRequestObject,
) (gen.GetSubscriptionsSumResponseObject, error) {
	filter := new(entity.ListSubscriptionFilter)
	filter.Title = request.Params.ServiceName
	if request.Params.UserId != nil {
		filter.UserID = pkg.PointerTo(request.Params.UserId.String())
	}

	t_s, err := time.Parse("01-2006", request.Params.StartDate)
	if err != nil {
		if errors.Is(err, &time.ParseError{}) {
			return gen.GetSubscriptionsSum400JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
		}
		return gen.GetSubscriptionsSum500JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
	}

	startDateWithDay := time.Date(t_s.Year(), t_s.Month(), 1, 0, 0, 0, 0, time.UTC)

	t_e, err := time.Parse("01-2006", request.Params.EndDate)
	if err != nil {
		if errors.Is(err, &time.ParseError{}) {
			return gen.GetSubscriptionsSum400JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
		}
		return gen.GetSubscriptionsSum500JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
	}

	endDateWithDay := time.Date(t_e.Year(), t_e.Month(), 1, 0, 0, 0, 0, time.UTC)

	filter.StartDate = &startDateWithDay
	filter.EndDate = &endDateWithDay

	sum, err := r.subUsecase.Sum(ctx, *filter)
	if err != nil {
		return gen.GetSubscriptionsSum500JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
	}

	return gen.GetSubscriptionsSum200JSONResponse(gen.AggregationResult{TotalCost: int(sum)}), nil
}

func (r *Server) DeleteSubscriptionsId(ctx context.Context, request gen.DeleteSubscriptionsIdRequestObject) (gen.DeleteSubscriptionsIdResponseObject, error) {
	err := r.subUsecase.Delete(ctx, request.Id.String())
	if err != nil {
		r.logger.Error("delete subscription", zap.Error(err))

		if errors.Is(err, usecase.ErrNotFound) {
			return gen.DeleteSubscriptionsId404JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
		}
		return gen.DeleteSubscriptionsId500JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
	}
	return gen.DeleteSubscriptionsId204Response{}, nil
}

func (r *Server) GetSubscriptionsId(
	ctx context.Context,
	request gen.GetSubscriptionsIdRequestObject,
) (gen.GetSubscriptionsIdResponseObject, error) {
	sub, err := r.subUsecase.Read(ctx, request.Id.String())
	if err != nil {
		r.logger.Error("get subscription", zap.Error(err))

		if errors.Is(err, usecase.ErrNotFound) {
			return gen.GetSubscriptionsId404JSONResponse{}, nil
		}
		return nil, err
	}

	userID := pkg.UUID(sub.UserID)

	created := time.UnixMilli(sub.CreatedAt)
	updated := time.UnixMilli(sub.UpdatedAt)

	return gen.GetSubscriptionsId200JSONResponse{
		Id:          pkg.UUID(sub.ID),
		ServiceName: sub.Title,
		Price:       int(sub.Price),
		StartDate:   sub.StartDate.String(),
		EndDate:     pkg.PointerTo(sub.EndDate.String()),
		UserId:      *userID,
		CreatedAt:   &created,
		UpdatedAt:   &updated,
	}, nil
}

func (r *Server) PutSubscriptionsId(
	ctx context.Context,
	request gen.PutSubscriptionsIdRequestObject,
) (gen.PutSubscriptionsIdResponseObject, error) {
	sub := new(entity.UpdateSubscriptionRequest)
	sub.Title = request.Body.ServiceName
	price := int64(request.Body.Price)
	sub.Price = price
	t_s, err := time.Parse("01-2006", request.Body.StartDate)
	if err != nil {
		r.logger.Error("parse start date", zap.Error(err))

		if errors.Is(err, &time.ParseError{}) {
			return gen.PutSubscriptionsId400JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
		}
		return gen.PutSubscriptionsId400JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
	}
	sub.StartDate = time.Date(t_s.Year(), t_s.Month(), 1, 0, 0, 0, 0, time.UTC)
	if request.Body.EndDate != nil {
		t_e, err := time.Parse("01-2006", *request.Body.EndDate)
		if err != nil {
			r.logger.Error("parse end date", zap.Error(err))

			if errors.Is(err, usecase.ErrInvalidSubscriptionData) {
				return gen.PutSubscriptionsId404JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
			}
		}
		endDateWithDay := time.Date(t_e.Year(), t_e.Month(), 1,
			0, 0, 0, 0, time.UTC)
		sub.EndDate = &endDateWithDay
	}

	sub.ID = request.Id.String()
	sub.UpdatedAt = time.Now().UnixMilli()

	err = r.subUsecase.Update(ctx, *sub)
	if err != nil {
		r.logger.Error("update subscription", zap.Error(err))

		if errors.Is(err, usecase.ErrSubscriptionAlreadyExists) {
			return gen.PutSubscriptionsId400JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
		}
		if errors.Is(err, usecase.ErrNotFound) {
			return gen.PutSubscriptionsId404JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
		}

		return nil, err
	}

	return gen.PutSubscriptionsId204Response{}, nil
}

func requestErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	responseErr(w, err.Error(), http.StatusInternalServerError)
}

//nolint:revive // ok.
func responseErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	responseErr(w, err.Error(), http.StatusInternalServerError)
}

func responseErr(w http.ResponseWriter, errStr string, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	resp := gen.ErrorResponse{
		Errors: &errStr,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	_, err = w.Write(jsonResp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
}

func (r *Server) Start() {
	srv := gen.NewStrictHandlerWithOptions(
		r,
		[]gen.StrictMiddlewareFunc{},
		gen.StrictHTTPServerOptions{
			RequestErrorHandlerFunc:  requestErrorHandler,
			ResponseErrorHandlerFunc: responseErrorHandler,
		},
	)
	handler := gen.Handler(srv)

	router := chi.NewRouter()
	router.Mount("/", handler)

	s := http.Server{
		Addr:              r.address,
		Handler:           router,
		ReadTimeout:       defaultReadTimeout,
		ReadHeaderTimeout: defaultHeadReadTimeout,
		WriteTimeout:      defaultWriteTimeout,
		IdleTimeout:       defaultIdleTimeout,
	}

	log.Fatal(s.ListenAndServe())
}
