package commands

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/The1Dani/faf_bot_go/cmd/bot/messages"
)

var DB *sql.DB

const (
	pidor = "currentpidor"
	nice = "currentnice"
	pidor_stats = "pidorstats"
	nice_stats = "stats"
)

func GetUser(member_id, chat_id int64) (user, error) {
	
	var u user
	
	log.Printf("[DEBUG] Selecting user where member_id = %d and chat_id = %d\n", member_id, chat_id)
	row := DB.QueryRow(`SELECT full_name, nick_name FROM members WHERE member_id = $1 AND chat_id = $2`, member_id, chat_id)
	err := row.Scan(&u.full_name, &u.nick_name)
	u.member_id = member_id
	
	if err == sql.ErrNoRows {
		return user{}, err
	} else if err != nil {
		log.Println("[ERROR] ", err)
	}
	
	return u, nil
}

func GetAllMembers(users []user, chat_id int64) []user {

	rows, err := DB.Query(`
		SELECT full_name, 
		       nick_name, 
		       member_id, 
			   coefficient, 
			   pidor_coefficient 
		FROM members WHERE chat_id = $1`, chat_id)
	
	if err != nil {
		log.Println("[ERROR] ", err)
	}
	
	defer rows.Close()
	
	for rows.Next() {
		u := user{}
		err = rows.Scan(&u.full_name, &u.nick_name, &u.member_id, &u.coefficient, &u.pidor_coefficient)
		if err != nil {
			log.Println("[ERROR] ", err)
		}
		users = append(users, u)
	}
	
	return users
	
}

func CreateUser(chat_id, reg_member_id int64, full_name, user_name string) bool {
	

	is_user_in_chat := false

	tx, err := DB.Begin()

	if err != nil {
		log.Println("[ERROR] ", err)
		return false
	}

	defer tx.Rollback()

	err = tx.QueryRow(`
    SELECT EXISTS(SELECT 1 FROM members WHERE chat_id=$1 AND member_id=$2)
    `, chat_id, reg_member_id).Scan(&is_user_in_chat)

	if err != nil {
		log.Println("[ERROR]:_exists: ", err)
		return false
	}

	if is_user_in_chat {
		return false
	}

	_, err = tx.Exec(` --sql
		INSERT INTO members (chat_id, member_id, coefficient, pidor_coefficient, full_name, nick_name) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`, chat_id, reg_member_id, 10, 10, full_name, user_name)

	if err != nil {
		log.Println("[ERROR] ", err)
		return false
	}

	stats_of_user, pidor_stats_of_user := false, false

	resp := tx.QueryRow(`--sql
    SELECT EXISTS(SELECT 1 FROM stats WHERE chat_id=$1 AND member_id=$2)
	`, chat_id, reg_member_id);
		
	resp.Scan(&stats_of_user)

	resp = tx.QueryRow(`--sql
    SELECT EXISTS(SELECT 1 FROM pidorstats WHERE chat_id=$1 AND member_id=$2)
	`, chat_id, reg_member_id);
	
	resp.Scan(&pidor_stats_of_user)


	if !stats_of_user {
		_, err = tx.Exec(`
		INSERT INTO stats (chat_id, member_id, count) 
		VALUES ($1, $2, $3)
		`, chat_id, reg_member_id, 0)
		if err != nil {
			log.Println("[ERROR] ", err)
			return false
		}
	}

	if !pidor_stats_of_user {
		_, err = tx.Exec(`
		INSERT INTO pidorstats (chat_id, member_id, count) 
		VALUES ($1, $2, $3)
		`, chat_id, reg_member_id, 0)
		if err != nil {
			log.Println("[ERROR] ", err)
			return false
		}
	}

	exists_currPidor, exists_currNice := false, false

	resp = tx.QueryRow(`--sql
    SELECT EXISTS(SELECT 1 FROM currentpidor WHERE chat_id=$1)
	`, chat_id);

	resp.Scan(&exists_currPidor)
	
	resp = tx.QueryRow(`--sql
    SELECT EXISTS(SELECT 1 FROM currentnice WHERE chat_id=$1)
	`, chat_id);

	resp.Scan(&exists_currNice)

	if !exists_currPidor {
		_, err = tx.Exec(`--sql
		INSERT INTO currentpidor (chat_id, member_id, timestamp)
		VALUES ($1, $2, $3)
		`, chat_id, 0, 0)
		if err != nil {
			log.Println("[ERROR] ", err)
			return false
		}
	}

	if !exists_currNice {
		_, err = tx.Exec(`--sql
		INSERT INTO currentnice (chat_id, member_id, timestamp)
		VALUES ($1, $2, $3)
		`, chat_id, 0, 0)
		if err != nil {
			log.Println("[ERROR] ", err)
			return false
		}
	}
	// sql.Result.RowsAffected()

	tx.Commit()

	return true
}

func DeleteUser(chat_id, member_id int64) (bool, error) {

	tx, err := DB.Begin()
	defer tx.Commit()

	if err != nil {
		tx.Rollback()
		return false, err
	}

	res, err := tx.Exec("DELETE FROM members WHERE chat_id = $1 AND member_id = $2", chat_id, member_id)

	if err != nil {
		tx.Rollback()
		return false, err
	}

	if num , _ := res.RowsAffected(); num == 0 {
		return false, nil
	} else {
		return true, nil
	}
}


func TimeNotExpired(chat_id int64, mode string) (bool, user) {
	/*
		time.Now().Unix() ! This is the same as python time int
 	*/

	var db_timestamp int64
	var curr_user user
	
	var querry_string string
	
	querry_string = fmt.Sprintf(`SELECT timestamp, member_id FROM %s WHERE chat_id = $1`, mode)
	row := DB.QueryRow(querry_string, chat_id)
	err := row.Scan(&db_timestamp, &curr_user.member_id)
	
	if err != nil {
		log.Println("[ERROR]", err)
	}
	
	curr_user, err = GetUser(curr_user.member_id, chat_id)
	
	if err == sql.ErrNoRows {
		log.Println("Pidor is not a user")
	}
	
	curr_hour := int64(time.Now().Hour())
	day_start := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location()).Unix()
	
	curr_timeframe := (curr_hour / messages.EVERY_N) * messages.EVERY_N * 60 * 60
	
	return db_timestamp > day_start + curr_timeframe, curr_user
}

func CarmicDicesEnabled(chat_id int64) bool {
	
	var enabled bool = false
	
	err := DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM carmadicesenabled WHERE chat_id=$1)`, chat_id).Scan(&enabled)
	
	if err != nil {
		log.Println("[ERROR]", err)
	}
	
	return enabled
}

func UpdateStats(chat_id, member_id int64, mode string) int32 {
	
	tx, err := DB.Begin()
	
	if err != nil {
		log.Println("[ERROR]", err)
	}
	
	update_string := fmt.Sprintf(`UPDATE %s SET count = count + 1 WHERE member_id = $1 AND chat_id = $2`, mode)
	_, err = tx.Exec(update_string, member_id, chat_id)
	
	if err != nil {
		log.Println("[ERROR]", err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
	
	var count int32
	
	query_string := fmt.Sprintf(`SELECT count FROM %s WHERE member_id = $2 AND chat_id = $3`, mode)
	err = DB.QueryRow(query_string, member_id, chat_id).Scan(&count)
	
	if err != nil {
		log.Println("[ERROR]", err)
	}
	
	return count
}