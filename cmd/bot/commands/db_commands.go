package commands

import (
	"database/sql"
	"log"
)

var DB *sql.DB

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

	tx.Commit()

	return true
}