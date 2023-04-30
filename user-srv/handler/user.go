package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"time"
	"user-srv/global"
	"user-srv/model"
	"user-srv/proto"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func (u UserServer) CreateUser(ctx context.Context, req *proto.UserInfoRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Userinfo.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "当前手机号绑定用户已存在!")
	}
	result = global.DB.Where(&model.User{UserName: req.Userinfo.Username}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "当前用户名已存在!")
	}
	password, err := HashPassword(req.Userinfo.Password)
	if err != nil {
		zap.S().Info("密码加密失败!")
		return nil, status.Errorf(codes.Internal, "内部服务错误!")
	}
	birthday := time.Unix(req.Userinfo.Birthday, 0)
	user.UserName = req.Userinfo.Username
	user.NickName = req.Userinfo.Nickname
	user.Mobile = req.Userinfo.Mobile
	user.Email = req.Userinfo.Email
	user.Status = req.Userinfo.Status
	user.UserType = req.Userinfo.Usertype
	user.Gender = req.Userinfo.Gender
	user.Birthday = &birthday
	user.Password = password
	result = global.DB.Create(&user)
	if result.Error != nil {
		zap.S().Errorf("gorm调用Create创建用户失败，%v", result.Error.Error())
		return nil, status.Errorf(codes.Internal, "内部服务错误，请稍后重试!错误信息:%v", result.Error.Error())
	}
	return &proto.UserInfoResponse{
		Id:         uint32(user.ID),
		CreateTime: user.CreatedAt.Unix(),
		UpdateTime: user.UpdatedAt.Unix(),
		Userinfo: &proto.UserBase{
			Username: user.UserName,
			Nickname: user.NickName,
			Password: user.Password,
			Mobile:   user.Mobile,
			Email:    user.Email,
			Status:   user.Status,
			Usertype: user.UserType,
			Gender:   user.Gender,
			Birthday: user.Birthday.Unix(),
		},
	}, nil
}

func (u UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UserInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}

func (u UserServer) GetUserList(ctx context.Context, req *proto.UserListRequest) (*proto.UserListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserList not implemented")
}
func (u UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByMobile not implemented")
}
func (u UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserById not implemented")
}
func (u UserServer) CheckPassword(ctx context.Context, req *proto.PasswordRequest) (*proto.CheckResultResponse, error) {
	password := req.Password
	encryptPassword := req.EncryptPassword
	result := CheckPasswordHash(password, encryptPassword)
	return &proto.CheckResultResponse{Result: result}, nil
}
