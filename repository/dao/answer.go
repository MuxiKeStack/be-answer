package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type AnswerDAO interface {
	Insert(ctx context.Context, answer Answer) (int64, error)
	FindById(ctx context.Context, answerId int64) (Answer, error)
	FindByQuestionId(ctx context.Context, questionId int64, curAnswerId int64, limit int64) ([]Answer, error)
	FindByUid(ctx context.Context, uid int64, curAnswerId int64, limit int64) ([]Answer, error)
	CountForQuestion(ctx context.Context, questionId int64) (int64, error)
	DelById(ctx context.Context, answerId int64, uid int64) error
}

type Answer struct {
	Id          int64 `gorm:"primaryKey,autoIncrement"`
	QuestionId  int64 `gorm:"index"`
	PublisherId int64
	Content     string `gorm:"varchar(200)"` // 产品给的长度
	Utime       int64
	Ctime       int64
}

type GORMAnswerDAO struct {
	db *gorm.DB
}

func NewGORMAnswerDAO(db *gorm.DB) AnswerDAO {
	return &GORMAnswerDAO{db: db}
}

func (dao *GORMAnswerDAO) Insert(ctx context.Context, answer Answer) (int64, error) {
	now := time.Now().UnixMilli()
	answer.Utime = now
	answer.Ctime = now
	err := dao.db.WithContext(ctx).Create(&answer).Error
	return answer.Id, err
}

func (dao *GORMAnswerDAO) FindById(ctx context.Context, answerId int64) (Answer, error) {
	var a Answer
	err := dao.db.WithContext(ctx).
		Where("id = ?", answerId).
		First(&a).Error
	return a, err
}

func (dao *GORMAnswerDAO) FindByQuestionId(ctx context.Context, questionId int64, curAnswerId int64, limit int64) ([]Answer, error) {
	var answers []Answer
	err := dao.db.WithContext(ctx).
		Where("question_id = ? and id < ?", questionId, curAnswerId).
		Order("id desc").
		Limit(int(limit)).
		Find(&answers).Error
	return answers, err
}

func (dao *GORMAnswerDAO) FindByUid(ctx context.Context, uid int64, curAnswerId int64, limit int64) ([]Answer, error) {
	var answers []Answer
	err := dao.db.WithContext(ctx).
		Where("publisher_id = ? and id < ?", uid, curAnswerId).
		Order("id desc").
		Limit(int(limit)).
		Find(&answers).Error
	return answers, err
}

func (dao *GORMAnswerDAO) CountForQuestion(ctx context.Context, questionId int64) (int64, error) {
	var cnt int64
	err := dao.db.WithContext(ctx).
		Model(&Answer{}).
		Where("question_id = ?", questionId).
		Count(&cnt).Error
	return cnt, err
}

func (dao *GORMAnswerDAO) DelById(ctx context.Context, answerId int64, uid int64) error {
	return dao.db.WithContext(ctx).
		Where("id = ? and publisher_id = ?", answerId, uid).
		Delete(&Answer{}).Error
}
