package net

import (
	"sync"
)

// 用户缓存
var Mgr = &WsMgr{
	userCache: make(map[int]WSConn),
	connCache: make(map[int64]WSConn),
	roleCache: make(map[int]WSConn),
}

type WsMgr struct {
	uc sync.RWMutex
	cc sync.RWMutex
	rc sync.RWMutex

	userCache map[int]WSConn
	connCache map[int64]WSConn
	roleCache map[int]WSConn
}

func (m *WsMgr) UserLogin(conn WSConn, uid int, token string) {
	m.uc.Lock()
	defer m.uc.Unlock()
	oldConn := m.userCache[uid]
	if oldConn != nil {
		//有用户登录着呢
		if conn != oldConn {
			//通过旧客户端 有用户抢登录了
			oldConn.Push("robLogin", nil)
		}
	}
	m.userCache[uid] = conn
	conn.SetProperty("uid", uid)
	conn.SetProperty("token", token)
}

func (w *WsMgr) UserLogout(wsConn WSConn) {
	w.RemoveUser(wsConn)
}
func (w *WsMgr) RemoveUser(conn WSConn) {
	w.uc.Lock()
	uid, err := conn.GetProperty("uid")
	if err == nil {
		//只删除自己的conn
		id := uid.(int)
		c, ok := w.userCache[id]
		if ok && c == conn {
			delete(w.userCache, id)
		}
	}
	w.uc.Unlock()

	w.rc.Lock()
	rid, err := conn.GetProperty("rid")
	if err == nil {
		//只删除自己的conn
		id := rid.(int)
		c, ok := w.roleCache[id]
		if ok && c == conn {
			delete(w.roleCache, id)
		}
	}
	w.rc.Unlock()

	conn.RemoveProperty("session")
	conn.RemoveProperty("uid")
	conn.RemoveProperty("role")
	conn.RemoveProperty("rid")
}

func (w *WsMgr) PushByRoleId(rid int, msgName string, data interface{}) bool {
	if rid <= 0 {
		return false
	}
	w.rc.Lock()
	defer w.rc.Unlock()
	c, ok := w.roleCache[rid]
	if ok {
		c.Push(msgName, data)
		return true
	} else {
		return false
	}
}

func (m *WsMgr) RoleEnter(conn WSConn, rid int) {
	m.rc.Lock()
	defer m.rc.Unlock()
	conn.SetProperty("rid", rid)
	m.roleCache[rid] = conn
}
