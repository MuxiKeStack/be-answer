//go:build wireinject

package main

import (
	"github.com/MuxiKeStack/be-answer/grpc"
	"github.com/MuxiKeStack/be-answer/ioc"
	"github.com/MuxiKeStack/be-answer/pkg/grpcx"
	"github.com/MuxiKeStack/be-answer/repository"
	"github.com/MuxiKeStack/be-answer/repository/dao"
	"github.com/MuxiKeStack/be-answer/service"
	"github.com/google/wire"
)

func InitGRPCServer() grpcx.Server {
	wire.Build(
		ioc.InitGRPCxKratosServer,
		grpc.NewAnswerServiceSever,
		service.NewAnswerService,
		// producer
		ioc.InitProducer,
		// rpc client
		ioc.InitQuestionClient,
		repository.NewAnswerRepository,
		dao.NewGORMAnswerDAO,
		// 第三方
		ioc.InitKafka,
		ioc.InitEtcdClient,
		ioc.InitLogger,
		ioc.InitDB,
	)
	return grpcx.Server(nil)
}
