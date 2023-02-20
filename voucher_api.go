package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

type VoucherProgram struct {
	ID                   int
	StartDate, EndDate   time.Time
	MaxProductsPerVoucher int
	TotalVouchers         int
	CreatedAt, UpdatedAt  time.Time
}

type Voucher struct {
	CodeChar       string
	ProgramID      int
	EmailAddress   string
	CreatedAt, UpdatedAt time.Time
}

type VoucherClaim struct {
	ID               int
	VoucherCodeChar  string
	ProductQuantity int
	RecipientEmail  string
	RecipientName   string
	Address         string
	CreatedAt, UpdatedAt time.Time
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/voucher_db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	r := chi.NewRouter()

	r.Get("/", indexHandler)
	r.Get("/voucher-programs", getVoucherProgramsHandler)
	r.Post("/voucher-programs", createVoucherProgramHandler)
	r.Get("/vouchers", getVouchersHandler)
	r.Post("/vouchers", createVoucherHandler)
	r.Get("/voucher-claims", getVoucherClaimsHandler)
	r.Post("/voucher-claims", createVoucherClaimHandler)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
