package yo_api

import (
	"myapp/pkg/database"

	"github.com/gorilla/mux"
	"github.com/mymmrac/telego"
)

const (
	cmdStart = "/start"
	cmdHelp  = "/help"
	cmdYo    = "/Yo"
	cmdAddYo = "/NowYouMyYo"
	cmdDelYo = "/NowYouNotMyYo"
)

func CMD(upd telego.Update) error {

	switch upd.Message.Text {
	case cmdStart:
	case cmdHelp:
	case cmdAddYo:
	case cmdDelYo:
	case cmdYo:
	}
}

func start(upd telego.Update) error {
	chatId := upd.Message.Chat.ID

}
