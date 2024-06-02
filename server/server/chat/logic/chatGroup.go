package logic

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/net"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/chat/model"
	"sync"
	"time"
)

// 聊天频道
type ChatGroup struct {
	userMutex sync.RWMutex
	msgMutex  sync.RWMutex
	//用户
	users map[int]*User
	//消息列表
	msgs ItemQueue
}

func (c *ChatGroup) Enter(user *User) {
	c.userMutex.Lock()
	defer c.userMutex.Unlock()
	c.users[user.rid] = user
}

func (c *ChatGroup) GetUser(rid int) (*User, bool) {
	c.userMutex.RLock()
	defer c.userMutex.RUnlock()
	u, ok := c.users[rid]
	return u, ok
}

func (c *ChatGroup) Exit(rid int) {
	c.userMutex.Lock()
	defer c.userMutex.Unlock()
	delete(c.users, rid)
}

func (c *ChatGroup) History() []model.ChatMsg {
	//消息列表
	c.msgMutex.RLock()
	defer c.msgMutex.RUnlock()
	msgs := c.msgs
	items := msgs.items
	chatMsgs := make([]model.ChatMsg, 0)
	for _, item := range items {
		msg := item.(*Msg)
		cm := model.ChatMsg{RId: msg.RId, NickName: msg.NickName, Time: msg.Time.Unix(), Msg: msg.Msg}
		chatMsgs = append(chatMsgs, cm)
	}
	return chatMsgs
}

func (c *ChatGroup) PushMsg(rid int, msg string, t int8) *model.ChatMsg {
	c.userMutex.RLock()
	u, ok := c.users[rid]
	if !ok {
		return nil
	}
	c.userMutex.RUnlock()

	m := &Msg{
		Msg:      msg,
		RId:      rid,
		NickName: u.nickName,
		Time:     time.Now(),
	}
	c.msgMutex.Lock()
	if c.msgs.Size() > 100 {
		c.msgs.Dequeue()
	}
	c.msgs.Enqueue(m)

	c.msgMutex.Unlock()

	chatMsg := &model.ChatMsg{
		RId:      rid,
		Msg:      msg,
		NickName: u.nickName,
		Type:     t,
		Time:     time.Now().Unix(),
	}
	//消息要广播出去 所有的频道用户广播出去
	for _, user := range c.users {
		net.Mgr.PushByRoleId(user.rid, "chat.push", chatMsg)
	}
	return chatMsg
}

func NewGroup() *ChatGroup {
	return &ChatGroup{users: map[int]*User{}}
}
