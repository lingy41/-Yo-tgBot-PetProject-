package processor

import (
	"log"
	"myapp/pkg/clients"
	"myapp/pkg/database"

	"github.com/mymmrac/telego"
)

func ListAndService(upd telego.Update, bot telego.Bot, db database.Database) error {
	chatId := upd.Message.Chat.ID
	user, err := db.TakeUserDataFromDbWithChatId(int(chatId))
	if err != nil {
		log.Printf("[ERR]main:can`t take data about user with chatId:%v -> %s ", chatId, err)
	}
	userClient := clients.NewUserClient(int(chatId))
	err = userClient.Cmd(upd.Message.Text, bot)
	if err != nil {
		log.Printf("[ERR]can`t to do CMD in processor with user :%s -> %s ", user.UserName, err)
	}
	return nil
}
