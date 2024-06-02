package middleware

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/net"
	"log"
)

func CheckRoleId() net.MiddlewareFunc {
	return func(next net.HandlerFunc) net.HandlerFunc {
		return func(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
			log.Println("进入到角色检测....")
			_, err := req.Conn.GetProperty("rid")
			if err != nil {
				rsp.Body.Code = constant.SessionInvalid
				return
			}
			next(req, rsp)
		}
	}
}
