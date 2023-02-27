package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	table "./gen/voucher_db/table"
	model "./gen/voucher_db/model"
	jet "github.com/go-jet/jet/v2/mysql"

	entities "happy_api/entities"

	handlers "happy_api/handlers"

	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/voucher_db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	r := chi.NewRouter()
	//Setting base middleware stack
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponswWriter, r *http.Request){
        w.Write([]byte(""))
	})

	// Define the routes and handlers for voucher programs
	r.Route("/voucher_programs", func(r chi.Router) {
		//CreateVoucherProgram
		r.Post("/", CreateVoucherProgram(db)) 
		//UpdateVoucherProgram
		r.Put("/{id}", UpdateVoucherProgram(db))
		//DeleteVoucherProgram
		r.Delete("/{id}", DeleteVoucherProgram(db))
		//GetAllVoucherProgram
		r.Get("/", GetAllVoucherProgram(db))
		//GetAllActiveVoucherProgram
		r.Get("/", GetAllActiveVoucherProgram(db))
		//GetVoucherProgramById
		r.Get("/{id}", GetVoucherProgramById(db))
	})

	// Define the routes and handlers for vouchers
	r.Route("/voucher", func(r chi.Router) {    
		//CreateVoucherGivenVoucherProgram
		r.Post("/",CreateVoucherGivenVoucherProgram(db))
		//UpdateVoucher
		r.Put("/{voucherCodeChar}", UpdateVoucher(db))
		//DeleteVoucher
		r.Delete("/{voucherCodeChar}", DeleteVoucher(db))
		//DeleteAllUnclaimedVoucherGivenVoucherProgram
		r.Delete("/unclaimed/{voucherProgramId}", DeleteAllUnclaimedVoucherGivenVoucherProgram(db))
		//GetAllVoucherGivenVoucherProgram
		r.Get("/{voucherProgramId}", GetAllVoucherGivenVoucherProgram(db))
	})

	// Define the routes and handlers for voucher claims
	r.Route("/voucher_claim", func(r chi.Router) {
		//CreateVoucherClaim
		r.Post("/", CreateVoucherClaim(db))
		//UpdateVoucherClaim
		r.Put("/{id}", UpdateVoucherClaim(db))
		//DeleteVoucherClaim
		r.Delete("/{id}", DeleteVoucherClaim(db))
		//GetVoucherClaim
		r.Get("/{id}", GetVoucherClaim(db))
	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
