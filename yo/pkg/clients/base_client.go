package clients

import (
	"log"
	"myapp/environment"
	"myapp/pkg/database"

	"github.com/mymmrac/telego"
)

type DefaultUser interface {
	ListAndService(telego.Update) error
	AddFriend(userName string) error
	DeleteFriend(userName string) error
	Yo(bot telego.Bot) error
}

type UserClient struct {
	chatId int
}

func NewUserClient(chatID int) *UserClient {
	return &UserClient{
		chatId: chatID,
	}
}

func (c *UserClient) Cmd(cmdFromMessage string, bot telego.Bot) error {
	switch cmdFromMessage {
	case environment.CmdStart:
		params := telego.SendMessageParams{ChatID: telego.ChatID{ID: int64(c.chatId)},
			Text: environment.MsgStart}
		_, err := bot.SendMessage(&params)
		if err != nil {
			log.Printf("[ERR] can`t cmd client -> %s ", err)
			return err
		}

	case environment.CmdHelp:
		params := telego.SendMessageParams{ChatID: telego.ChatID{ID: int64(c.chatId)},
			Text: environment.MsgHelp}
		_, err := bot.SendMessage(&params)
		if err != nil {
			log.Printf("[ERR] can`t cmd client -> %s ", err)
			return err
		}
	case environment.CmdAddfriend:

	case environment.CmdDelFriend:
	case environment.CmdYo:
	default:
		params := telego.SendMessageParams{ChatID: telego.ChatID{ID: int64(c.chatId)},
			Text: environment.MsgCmdNotFound}
		_, err := bot.SendMessage(&params)
		if err != nil {
			log.Printf("[ERR] can`t cmd client -> %s ", err)
			return err
		}
	}
	return nil
}

func (c *UserClient) AddFriend(userName string) error {
	db, err := database.NewDataBase(environment.MustPath())
	if err != nil {
		log.Printf("[ERR]user:%s can`t connect to database -> %s", userName, err)
		return err
	}
	usersFriend, err := db.TakeUserDataFromDbWithUserName(userName)
	if err != nil {
		return err
	}
	err = db.AddFrienToUserList(c.chatId, usersFriend.СhatId)
	return err
}

func (c *UserClient) DeleteFriend(userName string) error {
	db, err := database.NewDataBase(environment.MustPath())
	if err != nil {
		log.Printf("[ERR]user:%s can`t connect to database -> %s", userName, err)
		return err
	}
	usersFriend, err := db.TakeUserDataFromDbWithUserName(userName)
	if err != nil {
		return err
	}
	err = db.RemoveFriendFromUserList(c.chatId, usersFriend.UserName)
	if err != nil {
		return err
	}
	return nil
}

func (c *UserClient) Yo(bot telego.Bot) error {

	db, err := database.NewDataBase(environment.MustPath())
	if err != nil {
		log.Printf("[ERR]user can`t connect to database -> %s", err)
		return err
	}
	user, err := db.TakeUserDataFromDbWithChatId(c.chatId)
	if err != nil {
		return err
	}
	friends, err := db.GetUserFriends(c.chatId)
	if err != nil {
		return err
	}

	//
	//СДЕЛАТЬ ЦИКЛ ОТПРАВКИ ДРУЗЬЯМ ЙООООООУ
	//
	defer bot.StopLongPolling()

	for _, friend := range friends {
		params := telego.SendMessageParams{
			ChatID: telego.ChatID{ID: int64(friend.СhatId)},
			Text:   friend.UserName + " send Yo!",
		}
		_, err = bot.SendMessage(&params)
		if err != nil {
			log.Printf("[ERR]user: %s can`t send Yo to friens -> %s", user.UserName, err)
		}
	}
	return nil

}
