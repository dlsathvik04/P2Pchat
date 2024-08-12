package db

import (
	"database/sql"
	"fmt"
	"log"
)

type P2PchatDB struct {
	db *sql.DB
}
type Peer struct {
	Username string
	Addr     string
}

func NewP2PchatDB(db *sql.DB) *P2PchatDB {
	// store := NewStore(db, )

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Db connection sucessful....")

	p2pdb := P2PchatDB{
		db: db,
	}
	p2pdb.setUp([]func(*sql.DB) error{CreateAddrMapTable})

	return &p2pdb
}

func (p2pdb *P2PchatDB) setUp(setupFuncs []func(*sql.DB) error) {
	for _, setupFunc := range setupFuncs {
		if err := setupFunc(p2pdb.db); err != nil {
			log.Fatal("failed in setup: ", err)
		}
	}
	fmt.Println("Db setup sucessful....")

}

func CreateAddrMapTable(db *sql.DB) error {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS AddrMap 
			(
				username VARCHAR PRIMARY KEY, 
				addr VARCHAR, password VARCHAR
			);`,
	)
	return err
}

func (p2pdb *P2PchatDB) CreateUser(username string, password string, addr string) (Peer, error) {
	var user Peer
	err := p2pdb.db.QueryRow(`INSERT INTO AddrMap (username, password, addr) VALUES ($1, $2, $3) RETURNING username,addr`, username, password, addr).Scan(&user.Username, &user.Addr)
	return user, err
}

func (p2pdb *P2PchatDB) UpdateUserIp(username string, addr string) (Peer, error) {
	var user Peer
	err := p2pdb.db.QueryRow(`UPDATE AddrMap SET addr=$1 WHERE username=$2 RETURNING username, addr`, addr, username).Scan(&user.Username, &user.Addr)
	return user, err
}

func (p2pdb *P2PchatDB) GetUserByUsername(username string) (Peer, string, error) {
	var user Peer
	var password string
	err := p2pdb.db.QueryRow(`SELECT username, addr, password FROM AddrMap where username=$1`, username).Scan(&user.Username, &user.Addr, &password)
	return user, password, err
}

func (p2pdb *P2PchatDB) CleanUp() error {
	fmt.Println("cleaning up .... ")
	_, err := p2pdb.db.Exec(`DROP TABLE IF EXISTS AddrMap;`)
	if err == nil {
		fmt.Println("Cleaning up successful....")
	}
	return err
}
