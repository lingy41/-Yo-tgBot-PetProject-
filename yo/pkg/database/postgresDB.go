package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"sync"
	"time"
)

type UserTable struct {
	userId   int
	UserName string
	СhatId   int
	friends  int
}

type friendsTable struct {
	id int
	//user_id из userTable (userId)
	userId int
	//user_id из userTable (userId)
	friendId int
}

type Database struct {
	mu   sync.Mutex
	pool *pgxpool.Pool
}

type Storage interface {
	TakeUserDataFromDbWithUserName(userName string) (UserTable, error)
	TakeUserDataFromDbWithChatId(chatID int) (UserTable, error)
	Ident(chatID int) error
	AddUsr(chatID int, usrName string) (idUsr int, err error)
	PlusCountOfYo(chatID int)
	//exactly BanUser , in future user can`t register or login in bot
	RemoveUsr(usrName string) error
	AddFrienToUserList(chatId int, friend string) error
	RemoveFriendFromUserList(chatId int, friend string) error
	GetUserFriends(chatId int) ([]UserTable, error)
	//LookAtDatabase()( _ , err error)
	//LookAtUserHistory(usrName string)( _ , err error)

}

func NewDataBase(connPath string) (*Database, error) {
	pool, err := pgxpool.Connect(context.Background(), connPath)
	if err != nil {
		log.Fatal(fmt.Errorf("[ERR]cant connect to database -> %s", err))
		return nil, fmt.Errorf("[ERR]cant connect to database -> %s", err)
	}
	return &Database{
		mu:   sync.Mutex{},
		pool: pool,
	}, nil
}

// user identification
func (d *Database) Ident(chatID int, userName string) error {
	_, err := d.pool.Query(context.Background(), `SELECT user FROM user_table WHERE chat_id like $1 ;`, chatID)
	if err != nil {
		id, err := d.AddUsr(chatID, userName)
		if err != nil {
			return fmt.Errorf("[ERR]can`t add new user -> %s", err)
		}
		log.Printf("[NEW_USER] new user: %s with id:%v", userName, id)
	}

	return nil

}

func (d *Database) AddUsr(chatID int, usrName string) (idUsr int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	var id int

	_ = d.pool.QueryRow(context.Background(),
		`INSERT INTO user_table (user_name , chat_id , created)
		VALUES ($1 , $2 , $3) RETURNING user_id; `, usrName, chatID, time.Now()).Scan(&id)
	if id == 0 {
		log.Printf("[ERR]can`t add user:%s into databse", usrName)
		return 0, fmt.Errorf("can`t add user")
	}

	return id, nil
}

func (d *Database) AddFrienToUserList(chatId int, friendsChatId int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	user, err := d.TakeUserDataFromDbWithChatId(chatId)
	if err != nil {
		return fmt.Errorf("[ERR]from databse: can`t take data about user with chatId:%v -> %s ",
			chatId, err)
	}
	usersFriend, err := d.TakeUserDataFromDbWithChatId(friendsChatId)
	if err != nil {
		return fmt.Errorf("[ERR]from databse: can`t take data about user with chatId:%v -> %s ",
			chatId, err)
	}

	_ = d.pool.QueryRow(context.Background(),
		`SELECT user_id , user_name FROM user_table WHERE chat_id=$1;`,
		friendsChatId).Scan(
		&usersFriend.СhatId,
		&usersFriend.UserName,
	)

	var frineds_link int

	_ = d.pool.QueryRow(context.Background(),
		`SELECT id FROM friends_table 
		WHERE user_id=$1 AND friend_id=$2 RETURNING id; `,
		user.СhatId, usersFriend.userId,
	).Scan(
		&frineds_link,
	)
	if frineds_link == 0 {
		_, err := d.pool.Exec(context.Background(),
			`INSERT INTO friends_table (user_id , friend_id) 
			VALUES ($1 , $2);`, user.userId, usersFriend.userId)
		if err != nil {
			return fmt.Errorf("[ERR] can`t make user:%s and user:%s friends -> %s",
				user.UserName, usersFriend.UserName, err)
		}
		return nil
	}
	return nil
}

func (d *Database) RemoveFriendFromUserList(chatId int, friend string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	user, err := d.TakeUserDataFromDbWithChatId(chatId)
	if err != nil {
		return fmt.Errorf("[ERR]from databse: can`t take data about user with chatId:%v -> %s ",
			chatId, err)
	}
	usersFriend, err := d.TakeUserDataFromDbWithUserName(friend)
	if err != nil {
		return fmt.Errorf("[ERR]from databse: can`t take data about user with chatId:%v -> %s ",
			chatId, err)
	}

	var frineds_link int

	_ = d.pool.QueryRow(context.Background(),
		`SELECT id FROM friends_table 
		WHERE user_id=$1 AND friend_id=$2 RETURNING id; `,
		user.СhatId, usersFriend.userId,
	).Scan(
		&frineds_link,
	)
	if frineds_link == 0 {
		return fmt.Errorf("[ERR] can`t remove friend:%s from user:%s -> %s",
			usersFriend.UserName, user.UserName, err)
	}
	_, err = d.pool.Exec(context.Background(),
		`DELETE FROM friends_table WHERE friend_id=$1 AND user_id=$2`,
		usersFriend.СhatId, chatId)
	if err != nil {
		return err
	}
	return err
}

func (d *Database) TakeUserDataFromDbWithChatId(chatID int) (UserTable, error) {

	user := UserTable{
		СhatId: chatID,
	}
	_ = d.pool.QueryRow(context.Background(),
		`SELECT user_id , user_name , friends FROM user_table WHERE chat_id=$1;`,
		user.СhatId).Scan(
		&user.userId,
		&user.UserName,
		&user.friends,
	)

	if user.userId == 0 {
		return UserTable{}, fmt.Errorf("[ERR]no current user with chatID:%v", chatID)
	}
	return user, nil

}

func (d *Database) TakeUserDataFromDbWithUserName(userName string) (UserTable, error) {

	user := UserTable{
		UserName: userName,
	}
	_ = d.pool.QueryRow(context.Background(),
		`SELECT user_id , chat_id , friends FROM user_table WHERE user_name=$1;`,
		user.UserName).Scan(
		&user.userId,
		&user.СhatId,
		&user.friends,
	)

	if user.userId == 0 {
		return UserTable{}, fmt.Errorf("[ERR]no current user with username:%s", userName)
	}
	return user, nil

}

func (d *Database) GetUserFriends(chatId int) ([]UserTable, error) {
	sliceOfChatIds, err := d.pool.Query(context.Background(),
		`SELECT friend_id FROM friends_table where user_id=$1 `,
		chatId)
	if err != nil {
		return nil, err
	}
	defer sliceOfChatIds.Close()

	var sliceOfFriends []UserTable

	for sliceOfChatIds.Next() {
		var chatID int
		var friend UserTable
		sliceOfChatIds.Scan(
			&chatID,
		)
		friend, err = d.TakeUserDataFromDbWithChatId(chatID)
		if err != nil {
			return nil, err
		}
		sliceOfFriends = append(sliceOfFriends, friend)
	}

	return sliceOfFriends, nil

}
