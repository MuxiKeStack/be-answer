package service

import (
	"context"
	"github.com/MuxiKeStack/be-answer/domain"
	"github.com/MuxiKeStack/be-answer/events"
	"github.com/MuxiKeStack/be-answer/pkg/logger"
	"github.com/MuxiKeStack/be-answer/repository"
	feedv1 "github.com/MuxiKeStack/be-api/gen/proto/feed/v1"
	questionv1 "github.com/MuxiKeStack/be-api/gen/proto/question/v1"
	"strconv"
	"time"
)

var ErrAnswerNotFound = repository.ErrAnswerNotFound

type AnswerService interface {
	Publish(ctx context.Context, answer domain.Answer) (int64, error)
	GetDetailById(ctx context.Context, answerId int64) (domain.Answer, error)
	ListForQuestion(ctx context.Context, questionId int64, curAnswerId int64, limit int64) ([]domain.Answer, error)
	ListForUser(ctx context.Context, uid int64, curAnswerId int64, limit int64) ([]domain.Answer, error)
	CountForQuestion(ctx context.Context, questionId int64) (int64, error)
	DelAnswerById(ctx context.Context, answerId int64, uid int64) error
}

type answerService struct {
	repo           repository.AnswerRepository
	producer       events.Producer
	questionClient questionv1.QuestionServiceClient
	l              logger.Logger
}

func NewAnswerService(repo repository.AnswerRepository, producer events.Producer, questionClient questionv1.QuestionServiceClient,
	l logger.Logger) AnswerService {
	return &answerService{repo: repo, producer: producer, questionClient: questionClient, l: l}
}

func (a *answerService) Publish(ctx context.Context, answer domain.Answer) (int64, error) {
	answerId, err := a.repo.Create(ctx, answer)
	if err != nil {
		return 0, err
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		// 获取接受者
		res, er := a.questionClient.GetDetailById(ctx, &questionv1.GetDetailByIdRequest{
			QuestionId: answer.QuestionId,
		})
		if er != nil {
			a.l.Error("获取问题发布者失败", logger.Int64("questionId", answer.QuestionId))
			return
		}
		// 发送回答消息
		er = a.producer.ProduceFeedEvent(ctx, events.FeedEvent{
			Type: feedv1.EventType_Answer,
			Metadata: map[string]string{
				"answerer":   strconv.FormatInt(answer.PublisherId, 10),
				"questioner": strconv.FormatInt(res.GetQuestion().GetQuestionerId(), 10),
				"questionId": strconv.FormatInt(answer.QuestionId, 10),
				"answerId":   strconv.FormatInt(answer.Id, 10),
			},
		})
		if er != nil {
			a.l.Error("发送回答消息失败",
				logger.Error(er),
				logger.Int64("answerer", answer.PublisherId),
				logger.Int64("questioner", res.GetQuestion().GetQuestionerId()),
				logger.Int64("questionId", answer.QuestionId),
				logger.Int64("answerId", answer.Id),
			)
		}
	}()
	return answerId, nil
}

func (a *answerService) GetDetailById(ctx context.Context, answerId int64) (domain.Answer, error) {
	return a.repo.FindById(ctx, answerId)
}

func (a *answerService) ListForQuestion(ctx context.Context, questionId int64, curAnswerId int64, limit int64) ([]domain.Answer, error) {
	return a.repo.FindByQuestionId(ctx, questionId, curAnswerId, limit)
}

func (a *answerService) ListForUser(ctx context.Context, uid int64, curAnswerId int64, limit int64) ([]domain.Answer, error) {
	return a.repo.FindByUid(ctx, uid, curAnswerId, limit)
}

func (a *answerService) CountForQuestion(ctx context.Context, questionId int64) (int64, error) {
	return a.repo.CountForQuestion(ctx, questionId)
}

func (a *answerService) DelAnswerById(ctx context.Context, answerId int64, uid int64) error {
	return a.repo.DelById(ctx, answerId, uid)
}
