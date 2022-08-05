package router

import "fmt"

type Business struct {
	//Get string `plage:"get:user/"`
}

type PriceReq struct {
	ProductID int    `json:"product_id" plate:"product_id"`
	UserName  string `json:"user_name" plate:"user_name,mid"`
}

type PriceResp struct {
	Price    int    `json:"out"`
	UserName string `json:"user_string"`
}

var (
	_Business *Business
)

func (b *Business) Price(req interface{}) (resp interface{}, err error) {
	args := req.(*PriceReq)
	res := &PriceResp{}
	res.Price = args.ProductID * 10
	res.UserName = args.UserName
	return res, nil
}

type UserReq struct {
	Uid uint32 `json:"uid" plate:"Uid,header"`
}

type UserResp struct {
	Uid      uint32 `json:"uid" plate:"uid,mid"`
	UserName string `json:"user_name" plate:"user_name,mid"`
}

func UserMiddleware(req interface{}) (resp interface{}, err error) {
	args := req.(*UserReq)

	if args.Uid == 111 || args.Uid == 0 {
		return nil, fmt.Errorf("forbiden user")
	}

	userResp := &UserResp{
		Uid:      args.Uid,
		UserName: fmt.Sprintf("user:%d", args.Uid),
	}

	return userResp, nil
}
