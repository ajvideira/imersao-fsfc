package grpc

import (
	"context"

	"github.com/ajvideira/imersao-fullstack-fullcycle/codepix/application/grpc/pb"
	"github.com/ajvideira/imersao-fullstack-fullcycle/codepix/application/usecase"
)

//PixGrpcService represents a service for pix keys operations
type PixGrpcService struct {
	PixUseCase usecase.PixKeyUseCase
	pb.UnimplementedPixKeyServiceServer
}

//RegisterPixKey registers a new key
func (service *PixGrpcService) RegisterPixKey(ctx context.Context, in *pb.PixKeyRegistration) (*pb.PixKeyCreatedResult, error) {
	pixKey, err := service.PixUseCase.RegisterKey(in.Kind, in.Key, in.AccountID);
	if err != nil {
		return &pb.PixKeyCreatedResult{
			Status: "not created",
			Error: err.Error(),
		}, nil
	}
	return &pb.PixKeyCreatedResult{
		Id: pixKey.ID,
		Status: "created",
	}, nil
}

//Find searches for a existing key
func (service *PixGrpcService) Find(ctx context.Context, in *pb.PixKeyFind) (*pb.Pixkey, error) {
	pixKey, err := service.PixUseCase.FindKey(in.Kind, in.Key)
	if err != nil {
		return &pb.Pixkey{}, err
	}
	return &pb.Pixkey{
		Id: pixKey.ID,
		Kind: pixKey.Kind,
		Key: pixKey.Key,
		Account: &pb.Account{
			Id: pixKey.Account.ID,
			Number: pixKey.Account.Number,
			OwnerName: pixKey.Account.OwnerName,
			CreatedAt: pixKey.Account.CreatedAt.String(),
			Bank: &pb.Bank{
				Id: pixKey.Account.Bank.ID,
				Code: pixKey.Account.Bank.Code,
				Name: pixKey.Account.Bank.Name,
			},
		},
	}, err
}