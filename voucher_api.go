package main

import (
	"database/sql"
    "encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-jet/jet/v2/mysqljet"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
)

// VoucherProgram represents a voucher program in the database
type VoucherProgram struct {
	VoucherProgramID      int       `json:"voucher_program_id"`
	StartDate             time.Time `json:"start_date"`
	EndDate               time.Time `json:"end_date"`
	MaxProductsPerVoucher int       `json:"max_products_per_voucher"`
	TotalVouchers         int       `json:"total_vouchers"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// Voucher represents a voucher in the database
type Voucher struct {
	VoucherCodeChar string    `json:"voucher_code_char"`
	VoucherProgram  int       `json:"voucher_program"`
	EmailAddress    string    `json:"email_address"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// VoucherClaim represents a voucher claim in the database
type VoucherClaim struct {
	VoucherClaimID   int       `json:"voucher_claim_id"`
	VoucherCodeChar  string    `json:"voucher_code_char"`
	ProductQuantity  int       `json:"product_quantity"`
	RecipientEmail   string    `json:"recipient_email"`
	RecipientName    string    `json:"recipient_name"`
	Address          string    `json:"address"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
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
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

	r.Get("/", indexHandler)

    // Define the routes and handlers for voucher programs
	router.Get("/voucher_programs", ListVoucherPrograms(db))
	router.Post("/voucher_programs", CreateVoucherProgram(db))
	router.Get("/voucher_programs/{id}", GetVoucherProgram(db))
	router.Put("/voucher_programs/{id}", UpdateVoucherProgram(db))
	router.Delete("/voucher_programs/{id}", DeleteVoucherProgram(db))

	// Define the routes and handlers for vouchers
	router.Get("/voucher", ListVouchers(db))
	router.Post("/voucher", CreateVoucher(db))
	router.Get("/voucher/{code}", GetVoucher(db))
	router.Put("/voucher/{code}", UpdateVoucher(db))
	router.Delete("/voucher/{code}", DeleteVoucher(db))

	// Define the routes and handlers for voucher claims
	router.Get("/voucher_claims", ListVoucherClaims(db))
	router.Post("/voucher_claims", CreateVoucherClaim(db))
	router.Get("/voucher_claims/{id}", GetVoucherClaim(db))
	router.Put("/voucher_claims/{id}", UpdateVoucherClaim(db))
	router.Delete("/voucher_claims/{id}", DeleteVoucherClaim(db))


	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

//VoucherPrograms CRUD Operations
func CreateVoucherProgram(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        //TO-DO
    }
}

func UpdateVoucherProgram(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        //TO-DO
    }
}

func DeleteVoucherProgram(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        //TO-DO
    }
}

func GetAllVoucherProgram(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        //TO-DO
    }
}

func GetAllActiveVoucherProgram(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        //TO-DO
    }
}

func GetVoucherProgramById(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        //TO-DO
    }
}
//Voucher CRUD Operations
func CreateVoucherGivenVoucherProgram(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        //TO-DO
    }
}

func UpdateVoucher(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        //TO-DO
    }
}

func DeleteVoucher(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        //TO-DO
    }
}

func DeleteAllUnclaimedVoucherGivenVoucherProgram(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        //TO-DO
    }
}

func GetAllVoucherGivenVoucherProgram(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        //TO-DO
    }
}

//VoucherClaim CRUD Operations
func CreateVoucherClaim(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        //TO-DO
    }
}

func UpdateVoucherClaim(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        //TO-DO
    }
}

func DeleteVoucherClaim(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        //TO-DO
    }
}

func GetVoucherClaim(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        //TO-DO
    }
}
