package repository

import (
	"context"
	"github.com/MuxiKeStack/be-answer/domain"
	"github.com/MuxiKeStack/be-answer/repository/dao"
	"github.com/ecodeclub/ekit/slice"
	"time"
)

type AnswerRepository interface {
	Create(ctx context.Context, answer domain.Answer) (int64, error)
	FindById(ctx context.Context, answerId int64) (domain.Answer, error)
	FindByQuestionId(ctx context.Context, questionId int64, curAnswerId int64, limit int64) ([]domain.Answer, error)
	FindByUid(ctx context.Context, uid int64, curAnswerId int64, limit int64) ([]domain.Answer, error)
	CountForQuestion(ctx context.Context, questionId int64) (int64, error)
	DelById(ctx context.Context, answerId int64, uid int64) error
}

type answerRepository struct {
	dao dao.AnswerDAO
}

func NewAnswerRepository(dao dao.AnswerDAO) AnswerRepository {
	return &answerRepository{dao: dao}
}

func (repo *answerRepository) Create(ctx context.Context, answer domain.Answer) (int64, error) {
	return repo.dao.Insert(ctx, repo.toEntity(answer))
}

func (repo *answerRepository) FindById(ctx context.Context, answerId int64) (domain.Answer, error) {
	answer, err := repo.dao.FindById(ctx, answerId)
	return repo.toDomain(answer), err
}

func (repo *answerRepository) FindByQuestionId(ctx context.Context, questionId int64, curAnswerId int64, limit int64) ([]domain.Answer, error) {
	answers, err := repo.dao.FindByQuestionId(ctx, questionId, curAnswerId, limit)
	return slice.Map(answers, func(idx int, src dao.Answer) domain.Answer {
		return repo.toDomain(src)
	}), err
}

func (repo *answerRepository) FindByUid(ctx context.Context, uid int64, curAnswerId int64, limit int64) ([]domain.Answer, error) {
	answers, err := repo.dao.FindByUid(ctx, uid, curAnswerId, limit)
	return slice.Map(answers, func(idx int, src dao.Answer) domain.Answer {
		return repo.toDomain(src)
	}), err
}

func (repo *answerRepository) CountForQuestion(ctx context.Context, questionId int64) (int64, error) {
	return repo.dao.CountForQuestion(ctx, questionId)
}

func (repo *answerRepository) DelById(ctx context.Context, answerId int64, uid int64) error {
	return repo.dao.DelById(ctx, answerId, uid)
}

func (repo *answerRepository) toDomain(answer dao.Answer) domain.Answer {
	return domain.Answer{
		Id:          answer.Id,
		PublisherId: answer.PublisherId,
		QuestionId:  answer.QuestionId,
		Content:     answer.Content,
		Utime:       time.UnixMilli(answer.Utime),
		Ctime:       time.UnixMilli(answer.Ctime),
	}
}

func (repo *answerRepository) toEntity(answer domain.Answer) dao.Answer {
	return dao.Answer{
		Id:          answer.Id,
		PublisherId: answer.PublisherId,
		QuestionId:  answer.QuestionId,
		Content:     answer.Content,
		Utime:       answer.Utime.UnixMilli(),
		Ctime:       answer.Ctime.UnixMilli(),
	}
}
