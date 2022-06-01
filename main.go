package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/sensu/sensu-go/types"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	gelf "github.com/vin01/graylog-golang"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	PostgresURL   string
	GaylogHost    string
	GraylogPort   int
	EnableGraylog bool
}

type SLQuery struct {
	datname         string
	usename         string
	query_start     pgtype.Timestamp
	wait_event      sql.NullString
	wait_event_type sql.NullString
	state           sql.NullString
	backend_type    sql.NullString
	query           string
}

const SlowResponseQuery = `
	select datname,usename,query_start,wait_event,wait_event_type,state,backend_type,query 
		from pg_stat_activity 
		where 
			(now() - query_start) > interval '1 seconds' 
		AND 
			(now() - query_start) < interval '1 minute';
`

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-pg-metric-log",
			Short:    "A Hacked up tool, to find slow responding queries, provide a count of them, and log to graylog",
			Keyspace: "sensu.io/plugins/sensu-pg-metric-log/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[string]{
			Env:      "POSTGRESURL",
			Argument: "pgurl",
			Usage:    "URL to the postgres database",
			Value:    &plugin.PostgresURL,
		},
		&sensu.PluginConfigOption[string]{
			Env:      "GAYLOGHOST",
			Argument: "glhost",
			Usage:    "hostname of the graylog server",
			Value:    &plugin.GaylogHost,
		},
		&sensu.PluginConfigOption[int]{
			Env:      "GRAYLOGPORT",
			Argument: "glport",
			Default:  0,
			Usage:    "port of the graylog server",
			Value:    &plugin.GraylogPort,
		},
		&sensu.PluginConfigOption[bool]{
			Env:      "ENABLEGRAYLOG",
			Argument: "enable",
			Default:  false,
			Usage:    "Log results to graylog",
			Value:    &plugin.EnableGraylog,
		},
	}
)

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {

	if len(plugin.PostgresURL) == 0 {
		return sensu.CheckStateWarning, fmt.Errorf("Please specify a postgres url ( --pgurl postgres://user:password@127.0.0.1:5132/postgres")
	}

	if plugin.EnableGraylog == true {
		if len(plugin.GaylogHost) == 0 {
			return sensu.CheckStateWarning, fmt.Errorf("Please specify a graylog hostname ( -glhost graylog.example.com )")
		}
		if plugin.GraylogPort == 0 {
			return sensu.CheckStateWarning, fmt.Errorf("Please specify a graylog port ( -glport 12249 )")
		}
	}

	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {
	now := time.Now()
	offenders := getOffenders()
	metrics_ouput := fmt.Sprintf("query %d %d", len(offenders), now.Unix()) // type Graphite format
	fmt.Println(metrics_ouput)

	if plugin.EnableGraylog {

		g := gelf.New(gelf.Config{
			GraylogPort:     plugin.GraylogPort,
			GraylogHostname: plugin.GaylogHost,
		})

		message_templ := `{
				"timestamp": %d,
				"short_message": "%s",
				"_datname": "%s",
				"_usename": "%s",
				"_query_start": "%s",
				"_wait_event": "%s",
				"_wait_event_type": "%s",
				"_state": "%s",
				"_backend_type": "%s",
				"_query": "%s"
				}
				`

		for _, offender := range offenders {
			sqltxt := strings.Replace(offender.query, "\"", "\\\"", -1)
			logmessage := fmt.Sprintf(
				message_templ,
				now.Unix(),
				TruncateString(sqltxt, 150),
				offender.datname,
				offender.usename,
				offender.query_start.Time.String(),
				unnullify(offender.wait_event),
				unnullify(offender.wait_event_type),
				unnullify(offender.state),
				unnullify(offender.backend_type),
				sqltxt,
			)
			g.Log(logmessage)
		}
	}
	return sensu.CheckStateOK, nil
}

func getOffenders() []SLQuery {
	var offenders = []SLQuery{}
	conn, _ := pgx.Connect(context.Background(), plugin.PostgresURL)
	if rows, err := conn.Query(context.Background(), SlowResponseQuery); err != nil {
		fmt.Println("Unable to query due to: ", err)
		os.Exit(1)
	} else {
		defer rows.Close()
		var tmp SLQuery
		for rows.Next() {
			rows.Scan(&tmp.datname, &tmp.usename, &tmp.query_start, &tmp.wait_event, &tmp.wait_event_type, &tmp.state, &tmp.backend_type, &tmp.query)
			offenders = append(offenders, tmp)
		}
		if rows.Err() != nil {
			fmt.Println(rows.Err())
		}
	}
	return offenders
}

func unnullify(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

func TruncateString(str string, length int) string {
	if length <= 0 {
		return ""
	}
	truncated := ""
	count := 0
	for _, char := range str {
		truncated += string(char)
		count++
		if count >= length {
			break
		}
	}
	return truncated
}
