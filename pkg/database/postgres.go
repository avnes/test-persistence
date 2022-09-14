package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/avnes/test-persistence/pkg/common"
	_ "github.com/lib/pq"
	"github.com/pborman/uuid"
)

const (
	postgresUser     string = "postgres"
	postgresPassword string = "demo1234"
	postgresDatabase string = "postgres"
	postgresHostname string = "localhost"
	postgresPort     string = "5432"
)

type randomData struct {
	UUID     string `json:"uuid"`
	Modified string `json:"last_modified"`
}

func (r *randomData) record(uuid string, modified string) randomData {
	r.UUID = uuid
	r.Modified = modified
	return *r
}

func getConnectString() string {
	username, exist := os.LookupEnv("POSTGRES_USER")
	if !exist {
		username = postgresUser
	}
	password, exist := os.LookupEnv("POSTGRES_PASSWORD")
	if !exist {
		password = postgresPassword
	}
	database, exist := os.LookupEnv("POSTGRES_DB")
	if !exist {
		database = postgresDatabase
	}
	hostname, exist := os.LookupEnv("POSTGRES_HOSTNAME")
	if !exist {
		hostname = postgresHostname
	}
	port, exist := os.LookupEnv("POSTGRES_PORT")
	if !exist {
		port = postgresPort
	}
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", username, password, hostname, port, database)
	return connStr
}

func GetConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", getConnectString())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createDatabase() error {
	db, err := GetConnection()
	if err != nil {
		return err
	}
	rs, err := db.Exec("CREATE TABLE IF NOT EXISTS random_data(uuid VARCHAR(128) NOT NULL, modified VARCHAR(64) NOT NULL)")
	_ = rs
	if err != nil {
		return err
	}
	return nil
}

func populateTable(numberOfRecords int) error {
	db, err := GetConnection()
	if err != nil {
		return err
	}
	for i := 0; i < numberOfRecords; i++ {
		uuid := uuid.NewRandom()
		now := time.Now()
		stmt := fmt.Sprintf("INSERT INTO random_data (uuid, modified) VALUES ('%s','%s')", uuid, now.Format(time.RFC3339Nano))
		_, err := db.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func truncateTable() error {
	db, err := GetConnection()
	if err != nil {
		return err
	}
	rs, err := db.Exec("TRUNCATE TABLE random_data")
	_ = rs
	if err != nil {
		return err
	}
	return nil
}

func SetupDatabase() error {
	err := createDatabase()
	if err != nil {
		return err
	}
	return nil
}

func httpGetRandomData() ([]randomData, error) {
	db, err := GetConnection()
	if err != nil {
		return nil, err
	}
	var jsonFragment []randomData
	rs, err := db.Query("SELECT uuid, modified FROM random_data")
	defer rs.Close()
	if err != nil {
		return nil, err
	}
	var result randomData
	for rs.Next() {
		rs.Scan(&result.UUID, &result.Modified)
		record := new(randomData).record(result.UUID, result.Modified)
		jsonFragment = append(jsonFragment, record)

	}
	return jsonFragment, nil
}

func httpGetRandomDataCount() ([]common.Counter, error) {
	db, err := GetConnection()
	if err != nil {
		return nil, err
	}
	var jsonFragment []common.Counter
	rs, err := db.Query("SELECT COUNT(*) FROM random_data")
	defer rs.Close()
	if err != nil {
		return nil, err
	}
	var result common.Counter
	for rs.Next() {
		rs.Scan(&result.Count)
		record := new(common.Counter).Record(result.Count)
		jsonFragment = append(jsonFragment, record)
	}
	fmt.Printf("DEBUG: %v", result)
	return jsonFragment, nil
}
