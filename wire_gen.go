// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/MuxiKeStack/be-answer/grpc"
	"github.com/MuxiKeStack/be-answer/ioc"
	"github.com/MuxiKeStack/be-answer/pkg/grpcx"
	"github.com/MuxiKeStack/be-answer/repository"
	"github.com/MuxiKeStack/be-answer/repository/dao"
	"github.com/MuxiKeStack/be-answer/service"
)

// Injectors from wire.go:

func InitGRPCServer() grpcx.Server {
	logger := ioc.InitLogger()
	db := ioc.InitDB(logger)
	answerDAO := dao.NewGORMAnswerDAO(db)
	answerRepository := repository.NewAnswerRepository(answerDAO)
	answerService := service.NewAnswerService(answerRepository)
	answerServiceSever := grpc.NewAnswerServiceSever(answerService)
	client := ioc.InitEtcdClient()
	server := ioc.InitGRPCxKratosServer(answerServiceSever, client, logger)
	return server
}
