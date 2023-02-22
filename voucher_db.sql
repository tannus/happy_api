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

INSERT INTO VoucherProgram (voucher_program_id, start_date, end_date, max_products_per_voucher, total_vouchers, created_at, updated_at)
VALUES
	(1, '2022-01-01', '2022-01-31', 5, 10, '2022-01-01', '2022-01-01'),
	(2, '2022-02-01', '2022-02-28', 10, 20, '2022-02-01', '2022-02-01'),
	(3, '2022-03-01', '2022-03-31', 7, 15, '2022-03-01', '2022-03-01'),
	(4, '2022-04-01', '2022-04-30', 3, 5, '2022-04-01', '2022-04-01'),
	(5, '2022-05-01', '2022-05-31', 2, 10, '2022-05-01', '2022-05-01');

INSERT INTO Voucher (voucher_code_char, voucher_program_id, email_address, created_at, updated_at)
SELECT
    CONCAT('VC', LPAD(ROW_NUMBER() OVER(), 4, '0')),
    FLOOR(RAND() * 5) + 1,
    CONCAT('user', LPAD(ROW_NUMBER() OVER(), 2, '0'), '@example.com'),
    '2022-01-01',
    '2022-01-01'
FROM
    (SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5) a
    CROSS JOIN (SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5) b
ORDER BY
    RAND()
LIMIT 50;

INSERT INTO VoucherClaim (voucher_claim_id, voucher_code_char, product_quantity, recipient_email, recipient_name, address, created_at, updated_at)
SELECT 
	ROW_NUMBER() OVER (ORDER BY RAND()),
	voucher_code_char,
	FLOOR(RAND() * max_products_per_voucher) + 1,
	CONCAT('recipient', ROW_NUMBER() OVER (ORDER BY RAND()), '@example.com'),
	CONCAT('Recipient ', ROW_NUMBER() OVER (ORDER BY RAND())),
	CONCAT('Address ', ROW_NUMBER() OVER (ORDER BY RAND())),
	NOW(), NOW()
FROM Voucher v
INNER JOIN VoucherProgram vp ON v.voucher_program_id = vp.voucher_program_id
WHERE vp.total_vouchers >= 1
LIMIT 17;
