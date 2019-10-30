package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"awesomeProject1/db"
	"github.com/go-pg/pg"

	"awesomeProject1/account"
	"github.com/go-kit/kit/log"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fmt.Println(q.FormattedQuery())
	return nil
}

var (
	flagHttpAddr = flag.String("http_address", "0.0.0.0:8080", "Http address for web server running")

	flagDBAddr     = flag.String("db_address", "localhost:5432", "Address to connect to PostgreSQL server")
	flagDBUser     = flag.String("db_user", "postgres", "PostgreSQL connection user")
	flagDBPassword = flag.String("db_password", "", "PostgreSQL connection password")
	flagDBDatabase = flag.String("database", "payments", "PostgreSQL database name")
	flagDBAppName  = flag.String("app_name", "payments", "PostgreSQL application name (for logging)")
	flagDBPoolSize = flag.Int("pool_size", 10, "PostgreSQL connection pool size")
	flagDBLog = flag.Bool("db_log", false, "Switch for statements logging")
)

func main() {
	flag.Parse()

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	conn := setupDB(logger)
	defer func() {
		if err := conn.Close(); err != nil {
			_ = logger.Log("error", err)
		}
	}()

	var (
		accounts = db.NewAccountRepository(conn)
	)

	as := setupAccountService(accounts, logger)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	mux.Handle("/api/accounts/v1/", account.MakeHandler(as, httpLogger))


	http.Handle("/", accessControl(mux))


	errs := make(chan error, 2)
	go func() {
		_ = logger.Log("transport", "http", "address", *flagHttpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*flagHttpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	_ = logger.Log("terminated", <-errs)
}

func setupDB(logger log.Logger) *pg.DB {
	conn := pg.Connect(&pg.Options{
		Addr:            *flagDBAddr,
		User:            *flagDBUser,
		Password:        *flagDBPassword,
		Database:        *flagDBDatabase,
		ApplicationName: *flagDBAppName,
		PoolSize:        *flagDBPoolSize,
	})
	if *flagDBLog {
		conn.AddQueryHook(dbLogger{})
	}
	if err := db.CreateSchema(conn); err != nil {
		_ = logger.Log("transport", "DB", "address", *flagDBAddr, "msg", err)
	}
	return conn
}


func setupAccountService(accounts account.Repository, logger log.Logger) account.Service {

	as := account.NewService(accounts)
	as = account.NewLoggingService(log.With(logger, "component", "account"), as)

	return as
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
