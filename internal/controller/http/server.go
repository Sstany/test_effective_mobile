package http

import (
	"context"
	"subscrioption-service/internal/app/usecase"
	"subscrioption-service/internal/controller/http/gen"

	"go.uber.org/zap"
)

var _ gen.StrictServerInterface = (*Server)(nil)

type Server struct {
	address    string
	subUsecase usecase.SubscriptionUseCase
	logger     *zap.Logger
}

func NewServer(address string, subUsecase usecase.SubscriptionUseCase, logger *zap.Logger) *Server {
	return &Server{
		address:    address,
		subUsecase: subUsecase,
		logger:     logger,
	}
}

func (r *Server) GetSubscriptions(ctx context.Context, request gen.GetSubscriptionsRequestObject) (gen.GetSubscriptionsResponseObject, error) {

	return nil, nil
}

func (r *Server) PostSubscriptions(ctx context.Context, request gen.PostSubscriptionsRequestObject) (gen.PostSubscriptionsResponseObject, error) {
	return nil, nil
}

func (r *Server) GetSubscriptionsSum(ctx context.Context, request gen.GetSubscriptionsSumRequestObject) (gen.GetSubscriptionsSumResponseObject, error) {
	return nil, nil
}

func (r *Server) DeleteSubscriptionsId(ctx context.Context, request gen.DeleteSubscriptionsIdRequestObject) (gen.DeleteSubscriptionsIdResponseObject, error) {
	return nil, nil
}

func (r *Server) GetSubscriptionsId(ctx context.Context, request gen.GetSubscriptionsIdRequestObject) (gen.GetSubscriptionsIdResponseObject, error) {
	return nil, nil

}

func (r *Server) PutSubscriptionsId(ctx context.Context, request gen.PutSubscriptionsIdRequestObject) (gen.PutSubscriptionsIdResponseObject, error) {
	return nil, nil
}
