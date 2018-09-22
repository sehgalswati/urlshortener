package postgres

import (
	"database/sql"
	"fmt"
	"time"

	// This loads the postgres drivers.
	_ "github.com/lib/pq"

	"github.com/sehgalswati/urlshortener/algorithm"
	"github.com/sehgalswati/urlshortener/backend"
)

// New returns a postgres backed backend service.
func New(host, port, user, password, dbName string) (backend.BackendService, error) {
	// Connect postgres
	connect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	fmt.Printf("Opening the database")
	db, err := sql.Open("postgres", connect)
	if err != nil {
		fmt.Printf("Error in opening data base")
	
		return nil, err
	}

	// Ping to connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	strQuery := "CREATE TABLE IF NOT EXISTS urlshortenertable (uid serial NOT NULL, url VARCHAR not NULL, " +"visited boolean DEFAULT FALSE, count INTEGER DEFAULT 0, entrytime TIMESTAMP, lastvisitedtime TIMESTAMP);"

	_, err = db.Exec(strQuery)
	if err != nil {
		return nil, err
	}
	return &postgres{db}, nil
}

type postgres struct{ db *sql.DB }

func (p *postgres) Save(url string) (string, error) {
	var id int64
	//Stale entries which haven't been accessed for more than a year are deleted everytime a new entry is created
	// for testing purpose change the interval time to '30 seconds'
	//TODO: Make deletion of stale entries a separate go routine
	err := p.db.QueryRow("DELETE FROM urlshortenertable where lastvisitedtime < CURRENT_TIMESTAMP - INTERVAL '1 year' RETURNING uid").Scan(&id)
	fmt.Printf("remove stale entries before entering entries to %v", id )
	
	err = p.db.QueryRow("INSERT INTO urlshortenertable(url,visited,count,entrytime,lastvisitedtime) VALUES($1,$2,$3,$4,$5) returning uid;",url,false,0,time.Now(),time.Now()).Scan(&id)
	if err != nil {
		return "", err
	}
	return algorithm.Encode(id-1), nil
}

func (p *postgres) Load(code string) (string, error) {
	id, err := algorithm.Decode(code)
	if err != nil {
		return "", err
	}

	var url string
	err = p.db.QueryRow("update urlshortenertable set visited=true, count = count + 1, lastvisitedtime = $1 where uid=$2 RETURNING url", time.Now(),id+1).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}
func (p *postgres) LoadInfo(code string) (*backend.Tuple, error) {
	id, err := algorithm.Decode(code)
	if err != nil {
		return nil, err
	}

	var tuple backend.Tuple
	err = p.db.QueryRow("SELECT url, visited, count, entrytime,lastvisitedtime FROM urlshortenertable where uid=$1 limit 1", id+1).
		Scan(&tuple.URL, &tuple.Visited, &tuple.Count,&tuple.EntryTime,&tuple.LastVisitedTime)
	if err != nil {
		return nil, err
	}

	return &tuple, nil
}

func (p *postgres) Close() error { return p.db.Close() }
