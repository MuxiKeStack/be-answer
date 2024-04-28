package grpc

import (
	"context"
	"github.com/MuxiKeStack/be-answer/domain"
	"github.com/MuxiKeStack/be-answer/service"
	answerv1 "github.com/MuxiKeStack/be-api/gen/proto/answer/v1"
	"github.com/ecodeclub/ekit/slice"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"math"
	"time"
)

type AnswerServiceSever struct {
	answerv1.UnimplementedAnswerServiceServer
	svc service.AnswerService
}

func NewAnswerServiceSever(svc service.AnswerService) *AnswerServiceSever {
	return &AnswerServiceSever{svc: svc}
}

func (a *AnswerServiceSever) Publish(ctx context.Context, request *answerv1.PublishRequest) (*answerv1.PublishResponse, error) {
	aid, err := a.svc.Publish(ctx, convertToDomain(request.GetAnswer()))
	return &answerv1.PublishResponse{
		AnswerId: aid,
	}, err
}

func (a *AnswerServiceSever) Detail(ctx context.Context, request *answerv1.DetailRequest) (*answerv1.DetailResponse, error) {
	answer, err := a.svc.GetDetailById(ctx, request.GetAnswerId())
	return &answerv1.DetailResponse{
		Answer: convertToV(answer),
	}, err
}

func (a *AnswerServiceSever) ListForQuestion(ctx context.Context, request *answerv1.ListForQuestionRequest) (*answerv1.ListForQuestionResponse, error) {
	curAnswerId := request.GetCurAnswerId()
	if curAnswerId == 0 {
		curAnswerId = math.MaxInt64
	}
	answers, err := a.svc.ListForQuestion(ctx, request.GetQuestionId(), curAnswerId, request.GetLimit())
	return &answerv1.ListForQuestionResponse{
		Answers: slice.Map(answers, func(idx int, src domain.Answer) *answerv1.Answer {
			return convertToV(src)
		}),
	}, err
}

func (a *AnswerServiceSever) ListForUser(ctx context.Context, request *answerv1.ListForUserRequest) (*answerv1.ListForUserResponse, error) {
	curAnswerId := request.GetCurAnswerId()
	if curAnswerId == 0 {
		curAnswerId = math.MaxInt64
	}
	answers, err := a.svc.ListForUser(ctx, request.GetUid(), curAnswerId, request.GetLimit())
	return &answerv1.ListForUserResponse{
		Answers: slice.Map(answers, func(idx int, src domain.Answer) *answerv1.Answer {
			return convertToV(src)
		}),
	}, err
}

func (a *AnswerServiceSever) CountForQuestion(ctx context.Context, request *answerv1.CountForQuestionRequest) (*answerv1.CountForQuestionResponse, error) {
	cnt, err := a.svc.CountForQuestion(ctx, request.GetQuestionId())
	return &answerv1.CountForQuestionResponse{
		Cnt: cnt,
	}, err
}

func (a *AnswerServiceSever) DelAnswerById(ctx context.Context, request *answerv1.DelAnswerByIdRequest) (*answerv1.DelAnswerByIdResponse, error) {
	err := a.svc.DelAnswerById(ctx, request.GetAnswerId(), request.GetUid())
	return &answerv1.DelAnswerByIdResponse{}, err
}

func (a *AnswerServiceSever) Register(server *grpc.Server) {
	answerv1.RegisterAnswerServiceServer(server, a)
}

func convertToDomain(a *answerv1.Answer) domain.Answer {
	return domain.Answer{
		Id:          a.Id,
		PublisherId: a.PublisherId,
		QuestionId:  a.QuestionId,
		Content:     a.Content,
		Utime:       time.UnixMilli(a.Utime),
		Ctime:       time.UnixMilli(a.Ctime),
	}
}

func convertToV(a domain.Answer) *answerv1.Answer {
	return &answerv1.Answer{
		Id:          a.Id,
		PublisherId: a.PublisherId,
		QuestionId:  a.QuestionId,
		Content:     a.Content,
		Utime:       a.Utime.UnixMilli(),
		Ctime:       a.Ctime.UnixMilli(),
	}
}
