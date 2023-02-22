package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-jet/jet/v2/mysql"
	_ "github.com/go-sql-driver/mysql"
    
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	VoucherCodeChar string      `json:"voucher_code_char"`
	VoucherProgram  int         `json:"voucher_program"`
	EmailAddress    string      `json:"email_address"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// VoucherClaim represents a voucher claim in the database
type VoucherClaim struct {
	VoucherClaimID   int        `json:"voucher_claim_id"`
	VoucherCodeChar  string     `json:"voucher_code_char"`
	ProductQuantity  int        `json:"product_quantity"`
	RecipientEmail   string     `json:"recipient_email"`
	RecipientName    string     `json:"recipient_name"`
	Address          string     `json:"address"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
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

//VoucherPrograms CRUD Operations
func CreateVoucherProgram(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Parse the request body into a VoucherProgram struct
        var voucherProgram VoucherProgram
        err := json.NewDecoder(r.Body).Decode(&voucherProgram)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Build the SQL query using Jet SQL builder
        insertBuilder := mysql.
            Insert("voucher_program").
            Columns(
                "start_date",
                "end_date",
                "max_products_per_voucher",
                "total_vouchers",
                "created_at",
                "updated_at",
            ).
            Values(
                voucherProgram.StartDate,
                voucherProgram.EndDate,
                voucherProgram.MaxProductsPerVoucher,
                voucherProgram.TotalVouchers,
                time.Now().Format(time.RFC3339),
                time.Now().Format(time.RFC3339),
            )
        query, args := insertBuilder.Build()

        // Execute the SQL query and get the new voucher program ID
        result, err := db.Exec(query, args...)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        voucherProgramID, err := result.LastInsertId()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Set the voucher program ID in the struct and return it as JSON
        voucherProgram.VoucherProgramID = int(voucherProgramID)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(voucherProgram)
    }
}

func UpdateVoucherProgram(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Parse the request body into a VoucherProgram struct
        var voucherProgram VoucherProgram
        err := json.NewDecoder(r.Body).Decode(&voucherProgram)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Get the voucher program ID from the URL parameter
        id, err := strconv.Atoi(chi.URLParam(r, "id"))
        if err != nil {
            http.Error(w, "Invalid voucher program ID", http.StatusBadRequest)
            return
        }

        // Build the SQL query using Jet SQL builder
        updateBuilder := mysql.
            Update("voucher_program").
            Set(
                mysql.Assign("start_date", voucherProgram.StartDate),
                mysql.Assign("end_date", voucherProgram.EndDate),
                mysql.Assign("max_products_per_voucher", voucherProgram.MaxProductsPerVoucher),
                mysql.Assign("total_vouchers", voucherProgram.TotalVouchers),
                mysql.Assign("updated_at", time.Now().Format(time.RFC3339)),
            ).
            Where(
                mysql.Condition("voucher_program_id", "=", id),
            )
        query, args := updateBuilder.Build()

        // Execute the SQL query
        result, err := db.Exec(query, args...)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Check if the voucher program was updated
        rowsAffected, err := result.RowsAffected()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        if rowsAffected == 0 {
            http.Error(w, "Voucher program not found", http.StatusNotFound)
            return
        }

        // Return a success message
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, "Voucher program updated successfully")
    }
}

func DeleteVoucherProgram(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
    }
}

func GetAllVoucherProgram(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
    }
}

func GetAllActiveVoucherProgram(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func GetVoucherProgramById(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
    }
}

//Voucher CRUD Operations
func CreateVoucherGivenVoucherProgram(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func UpdateVoucher(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
    }
}

func DeleteVoucher(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
    }
}

func DeleteAllUnclaimedVoucherGivenVoucherProgram(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
    }
}

func GetAllVoucherGivenVoucherProgram(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
    }
}

//VoucherClaim CRUD Operations
func CreateVoucherClaim(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
    }
}

func UpdateVoucherClaim(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
    }
}

func DeleteVoucherClaim(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
    }
}

func GetVoucherClaim(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
    }
}
