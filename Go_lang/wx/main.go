package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/eatmoreapple/openwechat"
	"gopkg.in/toast.v1"
)

func main() {
	bot := openwechat.DefaultBot(openwechat.Desktop)
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	defer reloadStorage.Close()

	ex, err := os.Executable()
	CheckError(err)
	exPath := filepath.Dir(ex)

	// 登陆
	if err := bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
		fmt.Println(err)
		return
	}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	} else {
		push("Wechat", self.NickName, "登录成功", exPath+"/icon.png")
	}
	// 获取所有的好友
	friends, err := self.Friends()
	CheckError(err)

	// // 获取所有的群组
	// groups, err := self.Groups()
	// fmt.Println(groups, err)

	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
		if msg.IsSystem() {
			system_msg(msg, exPath)
		} else if msg.IsSendByFriend() {
			friend_msg(&friends, msg, exPath)
		} else if msg.IsSendByGroup() {
			group_msg(&friends, msg, exPath, self)
		} else {
			unknown_msg(msg, exPath)
		}
	}
	// 注册登陆二维码回调
	// bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
func system_msg(msg *openwechat.Message, path string) {
	receiver, err := msg.Receiver()
	CheckError(err)
	if msg.IsText() {
		str_time := time.Now().String()
		fmt.Println(str_time[0:19] + " " + "System" + ": " + msg.Content)
		push("Wechat: "+receiver.NickName, "System", msg.Content, path+"/avatar.png")
	} else if msg.IsPicture() {
		resp, err := msg.GetPicture()
		CheckError(err)
		defer func() { _ = resp.Body.Close() }()
		file, err := os.Create(path + "/pic.png")
		CheckError(err)
		defer func() { _ = file.Close() }()
		_, err = io.Copy(file, resp.Body)
		CheckError(err)
		str_time := time.Now().String()
		fmt.Println(str_time[0:19] + " " + "System" + ": [图片]")
		push("Wechat: "+receiver.NickName, "System", "[图片]", path+"/pic.png")
	} else if msg.IsVideo() {
		str_time := time.Now().String()
		fmt.Println(str_time[0:19] + " " + "System" + ": [视频]")
		push("Wechat: "+receiver.NickName, "System", "[视频]", path+"/avatar.png")
	} else if msg.IsMedia() {
		str_time := time.Now().String()
		fmt.Println(str_time[0:19] + " " + "System" + ": [App]" + msg.FileName)
		push("Wechat: "+receiver.NickName, "System", "[App]"+msg.FileName, path+"/avatar.png")
	} else {
		str_time := time.Now().String()
		fmt.Println(str_time[0:19] + " " + "System" + ": " + msg.Content)
		push("Wechat: "+receiver.NickName, "System", msg.Content, path+"/avatar.png")
	}
}

func friend_msg(friends *openwechat.Friends, msg *openwechat.Message, path string) {
	receiver, err := msg.Receiver()
	CheckError(err)
	sender, err := msg.Sender()
	CheckError(err)
	sult := friends.Search(1, func(friend *openwechat.Friend) bool { return friend.NickName == sender.NickName })
	if len(sult) != 0 && !sender.IsMP() {
		sult[0].SaveAvatar(path + "/avatar.png")
		if msg.IsText() && msg.Content != "" {
			str_time := time.Now().String()
			fmt.Println(str_time[0:19] + " " + sender.NickName + ": " + msg.Content)
			push("Wechat: "+receiver.NickName, sender.NickName, msg.Content, path+"/avatar.png")
		} else if msg.IsPicture() {
			resp, err := msg.GetPicture()
			CheckError(err)
			defer func() { _ = resp.Body.Close() }()
			file, err := os.Create(path + "/pic.png")
			CheckError(err)
			defer func() { _ = file.Close() }()
			_, err = io.Copy(file, resp.Body)
			CheckError(err)
			str_time := time.Now().String()
			fmt.Println(str_time[0:19] + " " + sender.NickName + ": [图片]")
			push("Wechat: "+receiver.NickName, sender.NickName, "[图片]", path+"/pic.png")
		} else if msg.IsVoice() {
			str_time := time.Now().String()
			fmt.Println(str_time[0:19] + " " + sender.NickName + ": [语音]")
			push("Wechat: "+receiver.NickName, sender.NickName, "[语音]", path+"/avatar.png")
		} else if msg.IsMedia() {
			str_time := time.Now().String()
			fmt.Println(str_time[0:19] + " " + sender.NickName + ": [App]" + msg.FileName)
			push("Wechat: "+receiver.NickName, sender.NickName, "[App]"+msg.FileName, path+"/avatar.png")
		} else if msg.IsVideo() {
			str_time := time.Now().String()
			fmt.Println(str_time[0:19] + " " + sender.NickName + ": [视频]")
			push("Wechat: "+receiver.NickName, sender.NickName, "[视频]", path+"/avatar.png")
		} else if msg.IsEmoticon() {
			resp, err := msg.GetPicture()
			CheckError(err)
			defer func() { _ = resp.Body.Close() }()
			file, err := os.Create(path + "/pic.png")
			CheckError(err)
			defer func() { _ = file.Close() }()
			_, err = io.Copy(file, resp.Body)
			CheckError(err)
			str_time := time.Now().String()
			fmt.Println(str_time[0:19] + " " + sender.NickName + ": [Emoticon]")
			push("Wechat: "+receiver.NickName, sender.NickName, "[Emoticon]", path+"/pic.png")
		} else {
			str_time := time.Now().String()
			fmt.Println(str_time[0:19] + " " + sender.NickName + ": " + msg.Content)
			push("Wechat: "+receiver.NickName, sender.NickName, msg.Content, path+"/avatar.png")
		}
	}
	// if sender.IsMP() {
	// 	str_time := time.Now().String()
	// 	fmt.Println(str_time[0:19] + " " + sender.NickName + msg.FileName)
	// 	push("Wechat: "+receiver.NickName, sender.NickName, msg.FileName, path+"/avatar.png")
	// }
}

func group_msg(friends *openwechat.Friends, msg *openwechat.Message, path string, self *openwechat.Self) {
	sender, err := msg.SenderInGroup()
	CheckError(err)
	receiver, err := msg.Receiver()
	CheckError(err)
	fmt.Println("group", receiver.NickName, sender.NickName)
	sult := friends.Search(1, func(friend *openwechat.Friend) bool { return friend.NickName == sender.NickName })
	if len(sult) != 0 && sender.NickName != self.NickName {
		sult[0].SaveAvatar("avatar.png")
		if msg.IsText() {
			str_time := time.Now().String()
			fmt.Println(str_time[0:19] + " " + sender.NickName + ": " + msg.Content)
			push("Wechat: "+receiver.NickName, sender.NickName, msg.Content, path+"/avatar.png")
		} else if msg.IsPicture() {
			str_time := time.Now().String()
			fmt.Println(str_time[0:19] + " " + sender.NickName + ": [图片]")
			push("Wechat: "+receiver.NickName, sender.NickName, "[图片]", path+".png")
		} else if msg.IsVoice() {
			str_time := time.Now().String()
			fmt.Println(str_time[0:19] + " " + sender.NickName + ": [语音]")
			push("Wechat: "+receiver.NickName, sender.NickName, "[语音]", path+"/avatar.png")
		} else if msg.IsVideo() {
			str_time := time.Now().String()
			fmt.Println(str_time[0:19] + " " + sender.NickName + ": [视频]")
			push("Wechat: "+receiver.NickName, sender.NickName, "[视频]", path+"/avatar.png")
		} else if msg.IsMedia() {
			str_time := time.Now().String()
			fmt.Println(str_time[0:19] + " " + sender.NickName + ": [App]" + msg.FileName)
			push("Wechat: "+receiver.NickName, sender.NickName, "[App]"+msg.FileName, path+"/avatar.png")
		} else {
			str_time := time.Now().String()
			fmt.Println(str_time[0:19] + " " + sender.NickName + ": " + msg.Content)
			push("Wechat: "+receiver.NickName, sender.NickName, msg.Content, path+"/avatar.png")
		}
	}
}

func unknown_msg(msg *openwechat.Message, path string) {
	receiver, err := msg.Receiver()
	CheckError(err)
	sender, err := msg.Sender()
	CheckError(err)
	fmt.Println("unknown", receiver.NickName, sender.NickName)
	if msg.IsText() {
		str_time := time.Now().String()
		fmt.Println(str_time[0:19] + " " + sender.NickName + ": " + msg.Content)
		push("Wechat: "+receiver.NickName, sender.NickName, msg.Content, path+"/avatar.png")
	} else if msg.IsPicture() {
		str_time := time.Now().String()
		fmt.Println(str_time[0:19] + " " + sender.NickName + ": [图片]")
		push("Wechat: "+receiver.NickName, sender.NickName, "[图片]", path+"/avatar.png")
	} else if msg.IsVoice() {
		str_time := time.Now().String()
		fmt.Println(str_time[0:19] + " " + sender.NickName + ": [语音]")
		push("Wechat: "+receiver.NickName, sender.NickName, "[语音]", path+"/avatar.png")
	} else if msg.IsVideo() {
		str_time := time.Now().String()
		fmt.Println(str_time[0:19] + " " + sender.NickName + ": [视频]")
		push("Wechat: "+receiver.NickName, sender.NickName, "[视频]", path+"/avatar.png")
	} else if msg.IsMedia() {
		str_time := time.Now().String()
		fmt.Println(str_time[0:19] + " " + sender.NickName + ": [App]" + msg.FileName)
		push("Wechat: "+receiver.NickName, sender.NickName, "[App]"+msg.FileName, path+"/avatar.png")
	} else {
		str_time := time.Now().String()
		fmt.Println(str_time[0:19] + " " + sender.NickName + ": " + msg.Content)
		push("Wechat: "+receiver.NickName, sender.NickName, msg.Content, path+"/avatar.png")
	}
}

func push(title string, sender string, message string, address string) {
	title = title + "\n"
	message = "\t" + message
	notification := toast.Notification{
		AppID:   title,
		Title:   sender,
		Message: message,
		Icon:    address, // This file must exist (remove this line if it doesn't)
		Actions: []toast.Action{
			{"protocol", "I'm a button", ""},
			{"protocol", "Me too!", ""},
		},
	}
	err := notification.Push()
	CheckError(err)
}

func CheckError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
