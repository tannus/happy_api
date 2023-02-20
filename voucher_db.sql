CREATE DATABASE voucher_db;

USE voucher_db;

CREATE TABLE VoucherProgram (
	voucher_program_id INT NOT NULL,
	start_date DATE NOT NULL,
	end_date DATE NOT NULL,
	max_products_per_voucher INT NOT NULL,
	total_vouchers INT NOT NULL,
	created_at DATE NOT NULL,
	updated_at DATE NOT NULL,
	PRIMARY KEY(voucher_program_id)
);

CREATE TABLE Voucher (
	voucher_code_char VARCHAR(16) NOT NULL,
	voucher_program_id INT NOT NULL,
	email_address VARCHAR(50),
	created_at DATE NOT NULL,
	updated_at DATE NOT NULL,
	PRIMARY KEY (voucher_code_char),
	FOREIGN KEY (voucher_program_id) REFERENCES VoucherProgram(voucher_program_id)
);

CREATE TABLE VoucherClaim (
	voucher_claim_id INT NOT NULL
	voucher_code_char VARCHAR(16) NOT NULL,
	product_quantity INT NOT NULL,
	recipient_email VARCHAR(50) NOT NULL,
	recipient_name VARCHAR(256),
	address VARCHAR(256) NOT NULL,
	created_at DATE NOT NULL,
	updated_at DATE NOT NULL,
	PRIMARY KEY (voucher_claim_id),
	FOREIGN KEY (voucher_code_char) REFERENCES Voucher(voucher_code_char) ON DELETE CASCADE ON UPDATE CASCADE
);
