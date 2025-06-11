package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	BaseURL string
}

var (
	db          *sqlx.DB
	GitTag      string = "dev"
	CFG         *Config
	DataBaseURL string
	Version     string
)

func open(dbsource string) (*sqlx.DB, error) {
	var err error
	db, err = sqlx.Open("postgres", dbsource)
	if err != nil {
		err = fmt.Errorf("error open db: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	db.SetMaxOpenConns(50)                 // Set maximum number of open connections to the database
	db.SetMaxIdleConns(30)                 // Set maximum number of connections in the idle connection pool
	db.SetConnMaxLifetime(5 * time.Minute) // Set the maximum lifetime of a connection to the database

	return db, nil
}

// get parameters from the database
func getParameters(appName string) (map[string]string, error) {
	var err error
	DataBaseURL = os.Getenv("DATABASE_URL")
	if DataBaseURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err = open(DataBaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

	SQLStatement := `
        SELECT "key", value, value_type
        FROM settings
        ORDER BY updated_at ASC
    `

	if appName != "" {
		SQLStatement = `
        SELECT "key", value, value_type
        FROM settings
        WHERE "key" LIKE %s
        ORDER BY updated_at ASC`

		SQLStatement = fmt.Sprintf(SQLStatement, "'"+appName+".%'")
	}

	rows, err := db.Queryx(SQLStatement)
	if err != nil {
		log.Fatalf("Error executing query: %s", err)
	}

	defer rows.Close()
	var configMap = make(map[string]string)

	for rows.Next() {
		var key, value, valueType string
		err := rows.Scan(&key, &value, &valueType)
		if err != nil {
			return nil, err
		}

		value = strings.TrimSpace(value)
		configMap[key] = value
	}

	return configMap, nil
}

func populateStruct(data map[string]string, out any) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return errors.New("out must be pointer to struct")
	}
	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		key := field.Name
		strVal, ok := data[key]
		if !ok {
			continue
		}

		fv := v.Field(i)
		if !fv.CanSet() {
			continue
		}

		switch fv.Kind() {
		case reflect.String:
			fv.SetString(strVal)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			parsed, err := strconv.ParseInt(strVal, 10, fv.Type().Bits())
			if err != nil {
				return fmt.Errorf("parse int %q for %s: %w", strVal, key, err)
			}
			fv.SetInt(parsed)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			parsed, err := strconv.ParseUint(strVal, 10, fv.Type().Bits())
			if err != nil {
				return fmt.Errorf("parse uint %q for %s: %w", strVal, key, err)
			}
			fv.SetUint(parsed)

		case reflect.Float32, reflect.Float64:
			parsed, err := strconv.ParseFloat(strVal, fv.Type().Bits())
			if err != nil {
				return fmt.Errorf("parse float %q for %s: %w", strVal, key, err)
			}
			fv.SetFloat(parsed)

		case reflect.Bool:
			parsed, err := strconv.ParseBool(strVal)
			if err != nil {
				return fmt.Errorf("parse bool %q for %s: %w", strVal, key, err)
			}
			fv.SetBool(parsed)
		case reflect.Slice:
			// Assuming the slice is of type []string
			// value separated by space " "

			parsed := strings.Split(strVal, " ")
			slice := reflect.MakeSlice(fv.Type(), len(parsed), len(parsed))
			for j, v := range parsed {
				slice.Index(j).Set(reflect.ValueOf(v))
			}
			fv.Set(slice)

		default:
			log.Fatalf("unsupported field type %s for field %s", fv.Kind(), field.Name)
		}
	}

	return nil
}

func Load() error {
	CFG = &Config{}

	execName, err := os.Executable()
	if err != nil {
		return err
	}
	execName = filepath.Base(execName)
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = execName
	}

	configMap := make(map[string]string)

	configMap, err = getParameters("")
	if err != nil {
		log.Fatalf("Error getting parameters: %s", err)
	}

	err = populateStruct(configMap, CFG)
	if err != nil {
		log.Fatalf("Error populating struct: %s", err)
	}

	configMap, err = getParameters(appName)
	if err != nil {
		log.Fatalf("Error getting parameters: %s", err)
	}

	for k, v := range configMap {
		k = strings.TrimPrefix(k, appName+".")
		configMap[k] = v
	}

	err = populateStruct(configMap, CFG)
	if err != nil {
		log.Fatalf("Error populating struct: %s", err)
	}

	return nil
}
