package db

import (
	"log"
	"time"

	"github.com/gocql/gocql"
)

//func GetUser(u_id gocql.UUID) {
//	session = getSessionObj()
//	//var id gocql.UUID
//	var name string
//	if err := session.Query(`SELECT name FROM experts WHERE exp_id = ? `,
//		u_id).Scan( &name); err != nil {
//		log.Fatal(err)
//	}
//}

//func getSessionObj () *gocql.Session{
//	if session != nil {
//		return session
//	}
//	return nil
//}

var session *gocql.Session

func ConnectToDB() (*gocql.Session, error) {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "chat_app"
	cluster.Consistency = gocql.One
	var err error
	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal("Trouble while connecting to db", err.Error())
	}
	//defer session.Close()
	//log.SetFlags(log.LstdFlags | log.Lshortfile)

	// if err := session.Query(`DROP TABLE IF EXISTS chat_app.experts;`).Exec(); err != nil {
	// 	log.Fatal(err)
	// }

	// if err := session.Query(`DROP TABLE IF EXISTS chat_app.patients;`).Exec(); err != nil {
	// 	log.Fatal(err)
	// }

	if err := session.Query(`DROP TABLE IF EXISTS chat_app.users;`).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := session.Query(`DROP TABLE IF EXISTS chat_app.chat_messages;`).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := session.Query(`DROP TABLE IF EXISTS chat_app.documents;`).Exec(); err != nil {
		log.Fatal(err)
	}
	// if err := session.Query(`CREATE TABLE IF NOT EXISTS chat_app.experts (exp_id uuid PRIMARY KEY, name text, last_update_timestamp timestamp );`).Exec(); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := session.Query(`CREATE TABLE IF NOT EXISTS chat_app.patients (patient_id uuid PRIMARY KEY, name text, last_update_timestamp timestamp );`).Exec(); err != nil {
	// 	log.Fatal(err)
	// }
	if err := session.Query(`CREATE TABLE IF NOT EXISTS chat_app.users (user_id VARCHAR PRIMARY KEY, name text,user_type text, last_update_timestamp timestamp );`).Exec(); err != nil {
		log.Fatal(err)
	}
	if err := session.Query(`CREATE TABLE IF NOT EXISTS chat_app.chat_messages (chat_id VARCHAR, from_user VARCHAR, to_user VARCHAR, body text, status VARCHAR, reply_for_chat_id VARCHAR, documents set<VARCHAR>, time timestamp, PRIMARY KEY ((from_user, to_user), time)) WITH CLUSTERING ORDER BY (time ASC);`).Exec(); err != nil {
		log.Fatal(err)
	}
	if err := session.Query(`CREATE TABLE IF NOT EXISTS chat_app.documents ( id VARCHAR, chat_id VARCHAR, document_content blob, document_name text, PRIMARY KEY (id, chat_id));`).Exec(); err != nil {
		log.Fatal(err)
	}

	// if err := session.Query(`INSERT INTO experts (name,exp_id,last_update_timestamp) VALUES (?,?,?)`,
	// 	"doc01",  gocql.TimeUUID().String(), time.Now()).Exec(); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := session.Query(`INSERT INTO experts (name,exp_id,last_update_timestamp) VALUES (?,?,?)`,
	// 	"doc02",  gocql.TimeUUID().String(), time.Now()).Exec(); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := session.Query(`INSERT INTO experts (name,exp_id,last_update_timestamp) VALUES (?,?,?)`,
	// 	"doc03",  gocql.TimeUUID().String(), time.Now()).Exec(); err != nil {
	// 	log.Fatal(err)
	// }

	// if err := session.Query(`INSERT INTO patients (name,patient_id,last_update_timestamp) VALUES (?,?,?)`,
	// 	"p01",  gocql.TimeUUID().String(), time.Now()).Exec(); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := session.Query(`INSERT INTO patients (name,patient_id,last_update_timestamp) VALUES (?,?,?)`,
	// 	"p02",  gocql.TimeUUID().String(), time.Now()).Exec(); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := session.Query(`INSERT INTO patients (name,patient_id,last_update_timestamp) VALUES (?,?,?)`,
	// 	"p03",  gocql.TimeUUID().String(), time.Now()).Exec(); err != nil {
	// 	log.Fatal(err)
	// }
	p01Id := gocql.TimeUUID().String()
	if err := session.Query(`INSERT INTO users (name,user_id,user_type,last_update_timestamp) VALUES (?,?,?,?)`,
		"p01", p01Id, "patient", time.Now()).Exec(); err != nil {
		log.Fatal(err)
	}
	p02Id := gocql.TimeUUID().String()
	if err := session.Query(`INSERT INTO users (name,user_id,user_type,last_update_timestamp) VALUES (?,?,?,?)`,
		"p02", p02Id, "patient", time.Now()).Exec(); err != nil {
		log.Fatal(err)
	}
	d01Id := gocql.TimeUUID().String()
	if err := session.Query(`INSERT INTO users (name,user_id,user_type,last_update_timestamp) VALUES (?,?,?,?)`,
		"d01", d01Id, "expert", time.Now()).Exec(); err != nil {
		log.Fatal(err)
	}
	d02Id := gocql.TimeUUID().String()
	if err := session.Query(`INSERT INTO users (name,user_id,user_type,last_update_timestamp) VALUES (?,?,?,?)`,
		"d02", d02Id, "expert", time.Now()).Exec(); err != nil {
		log.Fatal(err)
	}
	d03Id := gocql.TimeUUID().String()
	if err := session.Query(`INSERT INTO users (name,user_id,user_type,last_update_timestamp) VALUES (?,?,?,?)`,
		"d03", d03Id, "expert", time.Now()).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := session.Query(`INSERT INTO chat_app.chat_messages (chat_id, from_user, to_user, body, status, time) VALUES (?,?,?,?,?,?)`,
		gocql.TimeUUID().String(), p01Id, d02Id, "notDelivered msg from p01 to d02", "notDelivered", time.Now()).Exec(); err != nil {
		log.Fatal(err)
	}
	if err := session.Query(`INSERT INTO chat_app.chat_messages (chat_id, from_user, to_user, body, status, time) VALUES (?,?,?,?,?,?)`,
		gocql.TimeUUID().String(), p01Id, d02Id, "delivered msg from p01 to d02", "delivered", time.Now()).Exec(); err != nil {
		log.Fatal(err)
	}
	msgTime := time.Now()
	if err := session.Query(`INSERT INTO chat_app.chat_messages (chat_id, from_user, to_user, body, status, time) VALUES (?,?,?,?,?,?)`,
		gocql.TimeUUID().String(), p02Id, d02Id, "seen msg from p02 to d02", "seen", msgTime).Exec(); err != nil {
		log.Fatal(err)
	}
	msgTime2 := time.Now()
	if err := session.Query(`INSERT INTO chat_app.chat_messages (chat_id, from_user, to_user, body, status, time) VALUES (?,?,?,?,?,?)`,
		gocql.TimeUUID().String(), p02Id, d02Id, "delivered msg from p02 to d02", "delivered", msgTime2).Exec(); err != nil {
		log.Fatal(err)
	}
	// scan := session.Query(`SELECT chat_id, body, time FROM chat_app.chat_messages WHERE from_user = ? AND to_user = ? AND time <= ? LIMIT ? ALLOW FILTERING `, p02Id, d02Id, msgTime2, 10).Iter().Scanner()
	// for scan.Next() {
	// 	var (
	// 		body    string
	// 		chatId  string
	// 		msgTime time.Time
	// 	)
	// 	err2 := scan.Scan(&chatId, &body, &msgTime)
	// 	if err2 != nil {
	// 		fmt.Printf("error while querying; %v", err2.Error())
	// 		// return nil, err2
	// 	}
	// 	b, err2 := msgTime.MarshalText()
	// 	if err2 != nil {
	// 		fmt.Printf("error while marshaling time; %v\n", err2.Error())
	// 	}
	// 	fmt.Printf("chatId %v,\nbody %v,\nmsgTime %v\n", chatId, body, string(b))
	// }
	// scan2 := session.Query(`SELECT chat_id, body, time, from_user ,to_user FROM chat_app.chat_messages WHERE to_user IN ( ?, ?) AND from_user IN ( ?, ? ) AND time <= ? LIMIT 20 ALLOW FILTERING `, p02Id, d02Id, p02Id, d02Id, msgTime2).Iter().Scanner()
	// for scan2.Next() {
	// 	var (
	// 		body      string
	// 		chatId    string
	// 		msgTime   time.Time
	// 		from_user string
	// 		to_user   string
	// 	)
	// 	err2 := scan2.Scan(&chatId, &body, &msgTime, &from_user, &to_user)
	// 	if err2 != nil {
	// 		fmt.Printf("error while querying; %v", err2.Error())
	// 		// return nil, err2
	// 	}
	// 	b, err2 := msgTime.MarshalText()
	// 	if err2 != nil {
	// 		fmt.Printf("error while marshaling time; %v\n", err2.Error())
	// 	}
	// 	fmt.Printf("chatId %v,\nbody %v,\nmsgTime %v,\n from_user %v,\n to_user %v", chatId, body, string(b), from_user, to_user)
	// }

	// b, err2 := msgTime.MarshalText()
	// if err2 != nil {
	// 	fmt.Printf("error while marshaling time; %v\n", err2.Error())
	// } else {
	// 	// fmt.Printf("Marshaled time to string %v\n", string(b))
	// 	fmt.Printf("%v %v\n", msgTime, string(b))
	// }

	// var msgTimew time.Time
	// if err := msgTimew.UnmarshalText([]byte(string(b))); err != nil {
	// 	fmt.Printf("Trouble while parsing Time; %v\n", err.Error())
	// }

	// if err := session.Query(`UPDATE chat_app.chat_messages SET is_read = True WHERE from_user = ? AND to_user = ? AND time = ?`, u04Id.String(), u02Id.String(), msgTimew).Exec(); err != nil {
	// 	fmt.Printf("Trouble while updating db; %v\n", err.Error())
	// }

	// scan := session.Query(`SELECT user_id,name FROM chat_app.users;`).Iter().Scanner()

	// for scan.Next() {
	// 	var (
	// 		userID string
	// 		name   string
	// 	)
	// 	err := scan.Scan(&userID, &name)
	// 	if err != nil {
	// 		fmt.Printf("error while querying; %v", err.Error())
	// 	}
	// 	fmt.Printf("UserId : %v\n", userID)
	// 	fmt.Printf("Name : %v\n", name)
	// }

	return session, nil
}
