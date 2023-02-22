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
        // Extract voucher program ID from the URL parameter
        voucherProgramIDStr := chi.URLParam(r, "id")
        voucherProgramID, err := strconv.Atoi(voucherProgramIDStr)
        if err != nil {
            http.Error(w, fmt.Sprintf("Invalid voucher program ID: %v", err), http.StatusBadRequest)
            return
        }

        // Parse request body to get the updated voucher program fields
        var updatedVoucherProgram VoucherProgram
        err = json.NewDecoder(r.Body).Decode(&updatedVoucherProgram)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
            return
        }

        // Construct SQL query using Jet to update the voucher program
        query := mysql.Update("voucher_program").
            Set(mysql.L("start_date"), updatedVoucherProgram.StartDate).
            Set(mysql.L("end_date"), updatedVoucherProgram.EndDate).
            Set(mysql.L("max_products_per_voucher"), updatedVoucherProgram.MaxProductsPerVoucher).
            Set(mysql.L("total_vouchers"), updatedVoucherProgram.TotalVouchers).
            Set(mysql.L("updated_at"), time.Now().Format("2006-01-02 15:04:05")).
            Where(mysql.L("voucher_program_id").EQ(voucherProgramID))

        // Execute the query
        _, err = query.ExecContext(context.Background(), db)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error updating voucher program: %v", err), http.StatusInternalServerError)
            return
        }

        // Return success response
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Voucher program updated successfully"))
    }
}
func DeleteVoucherProgram(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        voucherProgramID, err := strconv.Atoi(chi.URLParam(r, "id"))
        if err != nil {
            http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
            return
        }

        // Delete the voucher program from the database
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        query := mysql.
            Delete("voucher_program").
            Where(mysql.ConditionRaw("voucher_program_id = ?", voucherProgramID))

        result, err := query.RunWith(db).ExecContext(ctx)
        if err != nil {
            http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
            return
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
            http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
            return
        }

        if rowsAffected == 0 {
            http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
            return
        }

        // Write a success response
        w.WriteHeader(http.StatusOK)
    }
}

func GetAllVoucherProgram(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        rows, err := mysql.
            Select("voucher_program_id", "start_date", "end_date", "max_products_per_voucher", "total_vouchers", "created_at", "updated_at").
            From("voucher_program").
            Exec(db)

        if err != nil {
            http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        voucherPrograms := []VoucherProgram{}
        for rows.Next() {
            voucherProgram := VoucherProgram{}
            err := rows.Scan(&voucherProgram.VoucherProgramID, &voucherProgram.StartDate, &voucherProgram.EndDate, &voucherProgram.MaxProductsPerVoucher, &voucherProgram.TotalVouchers, &voucherProgram.CreatedAt, &voucherProgram.UpdatedAt)
            if err != nil {
                http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
                return
            }
            voucherPrograms = append(voucherPrograms, voucherProgram)
        }

        if err = rows.Err(); err != nil {
            http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(voucherPrograms)
    }
}

func GetAllActiveVoucherProgram(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get current date and time
		now := time.Now().Format("2006-01-02 15:04:05")

		// Select all voucher programs that are active
		query := mysql.
			Select(voucherProgramTable.AllColumns).
			From(voucherProgramTable).
			Where(voucherProgramTable.StartDate.LessOrEqual(now)).
			Where(voucherProgramTable.EndDate.GreaterOrEqual(now)).
			OrderBy(voucherProgramTable.StartDate)

		rows, err := query.QueryContext(r.Context(), db)
		if err != nil {
			log.Printf("Error selecting voucher programs: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Iterate over the rows and build a slice of VoucherProgram objects
		voucherPrograms := make([]VoucherProgram, 0)
		for rows.Next() {
			var voucherProgram VoucherProgram
			err := rows.Scan(
				&voucherProgram.VoucherProgramID,
				&voucherProgram.StartDate,
				&voucherProgram.EndDate,
				&voucherProgram.MaxProductsPerVoucher,
				&voucherProgram.TotalVouchers,
				&voucherProgram.CreatedAt,
				&voucherProgram.UpdatedAt,
			)
			if err != nil {
				log.Printf("Error scanning voucher program row: %v\n", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			voucherPrograms = append(voucherPrograms, voucherProgram)
		}

		// Convert slice of VoucherProgram objects to JSON and write to response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(voucherPrograms)
		if err != nil {
			log.Printf("Error encoding JSON: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func GetVoucherProgramById(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract the ID from the request URL
        id := chi.URLParam(r, "id")

        // Query the database to get the voucher program with the specified ID
        row := db.QueryRow("SELECT id, name, start_date, end_date, is_active FROM voucher_programs WHERE id = $1", id)

        // Initialize a VoucherProgram struct to hold the query result
        var program VoucherProgram

        // Scan the query result into the VoucherProgram struct
        err := row.Scan(&program.Id, &program.Name, &program.StartDate, &program.EndDate, &program.IsActive)

        if err != nil {
            if err == sql.ErrNoRows {
                http.NotFound(w, r)
            } else {
                http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
            }
            return
        }

        // Encode the VoucherProgram struct as JSON and send it in the response
        w.Header().Set("Content-Type", "application/json")
        err = json.NewEncoder(w).Encode(program)
        if err != nil {
            http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        }
    }
}

//Voucher CRUD Operations
func CreateVoucherGivenVoucherProgram(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body to get voucher details
		var newVoucher Voucher
		err := json.NewDecoder(r.Body).Decode(&newVoucher)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if voucher program exists
		voucherProgramID := newVoucher.VoucherProgramID
		var voucherProgram VoucherProgram
		err = db.QueryRow("SELECT * FROM voucher_programs WHERE id=?", voucherProgramID).Scan(&voucherProgram.ID, &voucherProgram.Name, &voucherProgram.DiscountPercentage, &voucherProgram.StartDate, &voucherProgram.EndDate, &voucherProgram.IsActive)
		if err != nil {
			http.Error(w, "Voucher program does not exist", http.StatusBadRequest)
			return
		}

		// Insert new voucher
		result, err := db.Exec("INSERT INTO vouchers (code, expiry_date, voucher_program_id) VALUES (?, ?, ?)", newVoucher.Code, newVoucher.ExpiryDate, newVoucher.VoucherProgramID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get ID of newly inserted voucher
		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return ID of newly inserted voucher
		response := map[string]int64{"id": id}
		json.NewEncoder(w).Encode(response)
	}
}

func UpdateVoucher(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        voucherID := chi.URLParam(r, "voucherID")
        if voucherID == "" {
            http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
            return
        }

        var voucher Voucher
        err := json.NewDecoder(r.Body).Decode(&voucher)
        if err != nil {
            http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
            return
        }

        // Check if the voucher exists
        var count int
        err = db.QueryRow("SELECT COUNT(*) FROM vouchers WHERE id = $1", voucherID).Scan(&count)
        if err != nil {
            http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
            return
        }

        if count == 0 {
            http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
            return
        }

        // Update the voucher
        _, err = db.Exec("UPDATE vouchers SET code = $1, start_date = $2, end_date = $3, is_active = $4 WHERE id = $5",
            voucher.Code, voucher.StartDate, voucher.EndDate, voucher.IsActive, voucherID)

        if err != nil {
            http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
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
    return func(w http.ResponseWriter, r *http.Request){
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
