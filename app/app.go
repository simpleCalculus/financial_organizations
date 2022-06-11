package app

import (
	"context"
	"database/sql"
	"financial_organizations/config"
	"financial_organizations/pkg/repo"
	"financial_organizations/pkg/resp"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// App ...
type App struct {
	userRepo *repo.UserRepo
	Router   *mux.Router
}

// NewApp ...
func NewApp(conf *config.DB) (*App, error) {
	app := &App{}
	connection := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		conf.Username, conf.Password, conf.DBName)

	db, err := sql.Open(conf.Driver, connection)
	if err != nil {
		return nil, fmt.Errorf("could not connect database %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	app.userRepo = repo.NewUserRepo(db)
	app.userRepo.CreateTables()
	app.userRepo.AddUsers()

	app.Router = mux.NewRouter()
	app.setRouters()
	return app, nil
}

// Run ...
func (a *App) Run(port string) {
	log.Printf("Server started ... \n")
	shutDown := make(chan os.Signal, 1)
	signal.Notify(shutDown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	server := http.Server{
		Addr:    "localhost" + port,
		Handler: a.Router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	<-shutDown

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown failed")
	}
	log.Printf("Server stopped ZzzzZzz")
}

// Middleware ...
func (a *App) Middleware(handlerFunc http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		usrID, err := strconv.Atoi(r.Header.Get("X-UserID"))
		if err != nil {
			resp.ResponseByCode(w, resp.BadRequestHeader, http.StatusBadRequest)
		}
		hmac1 := r.Header.Get("X-Digest")

		err = a.userRepo.Authenticate(usrID, hmac1)
		if err != nil {
			resp.ResponseByCode(w, resp.BadRequest, http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "user-info", repo.UserInfo{
			Id:   usrID,
			Hmac: hmac1,
		}))
		start := time.Now()
		handlerFunc(w, r)
		log.Printf("method :%s %s was called - time handling %v msec", r.RequestURI, r.Method, time.Since(start).Milliseconds())
	}
}

// setRouters ...
func (a *App) setRouters() {
	a.Router.HandleFunc("/login", a.Login).Methods(http.MethodPost)
	a.Router.HandleFunc("/authentication", a.Middleware(a.Authentication)).Methods(http.MethodPost)
	a.Router.HandleFunc("/replenishment", a.Middleware(a.WalletReplenishment)).Methods(http.MethodPost)
	a.Router.HandleFunc("/transactions", a.Middleware(a.Transactions)).Methods(http.MethodPost)
	a.Router.HandleFunc("/balance", a.Middleware(a.WalletBalance)).Methods(http.MethodPost)
}
