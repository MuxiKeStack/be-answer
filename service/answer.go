package service

import (
	"context"
	"github.com/MuxiKeStack/be-answer/domain"
	"github.com/MuxiKeStack/be-answer/repository"
)

type AnswerService interface {
	Publish(ctx context.Context, answer domain.Answer) (int64, error)
	GetDetailById(ctx context.Context, answerId int64) (domain.Answer, error)
	ListForQuestion(ctx context.Context, questionId int64, curAnswerId int64, limit int64) ([]domain.Answer, error)
	ListForUser(ctx context.Context, uid int64, curAnswerId int64, limit int64) ([]domain.Answer, error)
	CountForQuestion(ctx context.Context, questionId int64) (int64, error)
	DelAnswerById(ctx context.Context, answerId int64, uid int64) error
}

type answerService struct {
	repo repository.AnswerRepository
}

func NewAnswerService(repo repository.AnswerRepository) AnswerService {
	return &answerService{repo: repo}
}

func (a *answerService) Publish(ctx context.Context, answer domain.Answer) (int64, error) {
	return a.repo.Create(ctx, answer)
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
