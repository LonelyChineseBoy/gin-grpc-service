package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"
	"user-srv/global"
	"user-srv/model"
	"user-srv/proto"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func ModelToUserResponse(user model.User) *proto.UserInfoResponse {
	userInfoRsp := proto.UserInfoResponse{
		Id:         uint32(user.ID),
		CreateTime: user.CreatedAt.Unix(),
		Userinfo: &proto.UserBase{
			Username: user.UserName,
			Nickname: user.NickName,
			Password: user.Password,
			Mobile:   user.Mobile,
			Email:    user.Email,
			Status:   user.Status,
			Usertype: user.UserType,
			Gender:   user.Gender,
		},
	}
	if &user.UpdatedAt != nil {
		userInfoRsp.UpdateTime = user.UpdatedAt.Unix()
	}
	if user.Birthday != nil {
		userInfoRsp.Userinfo.Birthday = user.Birthday.Unix()
	}
	return &userInfoRsp
}

func ModelToAddressResponse(address model.UserAddress) *proto.AddressResponse {
	addressInfo := proto.AddressResponse{
		Id:            uint32(address.ID),
		AddressDetail: address.AddressDetail,
	}
	return &addressInfo
}

func (u *UserServer) CreateUser(ctx context.Context, req *proto.UserInfoRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(model.User{Mobile: req.Userinfo.Mobile}).Or(model.User{UserName: req.Userinfo.Username}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "当前手机号或用户名已存在!")
	}
	password, err := HashPassword(req.Userinfo.Password)
	if err != nil {
		zap.S().Info("密码加密失败!")
		return nil, status.Errorf(codes.Internal, "内部服务错误!")
	}
	if req.Userinfo.Birthday == 0 {
		user.Birthday = nil
	} else {
		birthday := time.Unix(req.Userinfo.Birthday, 0)
		user.Birthday = &birthday
	}
	user.UserName = req.Userinfo.Username
	user.NickName = req.Userinfo.Nickname
	user.Mobile = req.Userinfo.Mobile
	user.Email = req.Userinfo.Email
	user.Status = req.Userinfo.Status
	user.UserType = req.Userinfo.Usertype
	user.Gender = req.Userinfo.Gender
	user.Password = password
	result = global.DB.Create(&user)
	if result.Error != nil {
		zap.S().Errorf("gorm调用Create创建用户失败，%v", result.Error.Error())
		return nil, status.Errorf(codes.Internal, "内部服务错误，请稍后重试!错误信息:%v", result.Error.Error())
	}
	return ModelToUserResponse(user), nil
}

func (u *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Model: gorm.Model{ID: uint(req.Id)}}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "未发现指定用户!")
	}
	user.NickName = req.Userinfo.Nickname
	user.Mobile = req.Userinfo.Mobile
	user.Email = req.Userinfo.Email
	user.Status = req.Userinfo.Status
	user.UserType = req.Userinfo.Usertype
	user.Gender = req.Userinfo.Gender
	if req.Userinfo.Birthday == 0 {
		user.Birthday = nil
	} else {
		birthday := time.Unix(req.Userinfo.Birthday, 0)
		user.Birthday = &birthday
	}
	result = global.DB.Save(&user)
	if result.RowsAffected == 1 {
		return ModelToUserResponse(user), nil
	}
	return nil, status.Errorf(codes.Internal, "保存用户信息时，发生内部错误，请稍后重试!")
}

func (u *UserServer) GetUserList(ctx context.Context, req *proto.UserListRequest) (*proto.UserListResponse, error) {
	var users []model.User
	var total int64
	result := global.DB.Model(&model.User{}).Find(&users)
	total = result.RowsAffected
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "查询用户数据失败，请稍后重试!")
	}
	data := &proto.UserListResponse{}
	data.Total = uint64(total)
	global.DB.Scopes(Paginate(int(req.Page), int(req.Size))).Find(&users)
	for _, user := range users {
		data.Users = append(data.Users, ModelToUserResponse(user))
	}
	return data, nil
}

func (u *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有找到用户信息!")
	}
	return ModelToUserResponse(user), nil
}

func (u *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Model: gorm.Model{ID: uint(req.Id)}}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有找到用户信息!")
	}
	return ModelToUserResponse(user), nil
}

func (u *UserServer) CheckPassword(ctx context.Context, req *proto.PasswordRequest) (*proto.CheckResultResponse, error) {
	password := req.Password
	encryptPassword := req.EncryptPassword
	result := CheckPasswordHash(password, encryptPassword)
	return &proto.CheckResultResponse{Result: result}, nil
}

func (u *UserServer) CreateUserAddress(ctx context.Context, req *proto.UserAddressInfoRequest) (*proto.AddressResponse, error) {
	var user model.User
	var address model.UserAddress
	if result := global.DB.Where(model.User{Model: gorm.Model{ID: uint(req.UserId)}}).First(&user); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "保存地址对应用户没有找到!")
	}
	address.UserId = int(req.UserId)
	address.AddressDetail = req.AddressDetail
	result := global.DB.Save(&address)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "数据保存失败，请稍后重试!")
	}
	return ModelToAddressResponse(address), nil
}

func (u *UserServer) GetUserAddressList(ctx context.Context, req *proto.IdRequest) (*proto.UserAddressListResponse, error) {
	var addresses []model.UserAddress
	result := global.DB.Where(model.UserAddress{UserId: int(req.Id)}).Find(&addresses)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取地址失败，请稍后重试!")
	}
	total := uint32(result.RowsAffected)
	addressList := proto.UserAddressListResponse{}.AddressList
	for _, address := range addresses {
		item := ModelToAddressResponse(address)
		addressList = append(addressList, item)
	}
	return &proto.UserAddressListResponse{
		Total:       total,
		AddressList: addressList,
	}, nil
}

func (u *UserServer) UpdateUserAddressInfo(ctx context.Context, req *proto.UserAddressDetailRequest) (*proto.AddressResponse, error) {
	var address model.UserAddress
	if result := global.DB.Where(model.UserAddress{Model: gorm.Model{ID: uint(req.Id)}}).First(&address); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有找到地址信息，无法修改!")
	}
	address.AddressDetail = req.AddressDetail
	result := global.DB.Save(&address)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "保存地址信息错误，请稍后重试!")
	}
	return ModelToAddressResponse(address), nil
}
