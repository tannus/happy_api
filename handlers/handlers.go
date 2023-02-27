package handlers

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
       	voucherProgramBuilder := model.Link{
			StartDate:					voucherProgram.StartDate,
			EndDate:					voucherProgram.EndDate,
			MaxProductsPerVoucher:		voucherProgram.MaxProductsPerVoucher,
			TotalVouchers:				voucherProgram.TotalVouchers,
			CreatedAt:					time.Now().Format(time.RFC3339),
			UpdatedAt:					time.Now().Format(time.RFC3339),
		}

		insertStmt := Link.INSERT(Link.StartDate, Link.EndDate, Link.MaxProductsPerVoucher, Link.TotalVouchers, Link.CreatedAt, Link.UpdatedAt).
			MODEL(voucherProgramBuilder)

        // Execute the SQL query and get the new voucher program ID
        result, err := insertStmt.Exec(db)
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
        voucherProgramBuilder := model.Link{
			StartDate:					voucherProgram.StartDate,
			EndDate:					voucherProgram.EndDate,
			MaxProductsPerVoucher:		voucherProgram.MaxProductsPerVoucher,
			TotalVouchers:				voucherProgram.TotalVouchers,
			UpdatedAt:					time.Now().Format(time.RFC3339),

		}
		updateStmt := Link.INSERT(Link.StartDate, Link.EndDate, Link.MaxProductsPerVoucher, Link.TotalVouchers, Link.CreatedAt, Link.UpdatedAt).
			MODEL(voucherProgramBuilder).
			WHERE(Link.VoucherProgramId.EQ(id))

        // Execute the SQL query
        result, err := updateStmt.Exec(query, args...)
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
