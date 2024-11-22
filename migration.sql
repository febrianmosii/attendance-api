CREATE TABLE m_users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE m_events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE m_classes (
    id SERIAL PRIMARY KEY,
    event_id INT NOT NULL REFERENCES m_events(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE m_operators (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE t_transactions (
    id SERIAL PRIMARY KEY,
    amount NUMERIC(10, 2) NOT NULL,
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    invoice_code VARCHAR(50) UNIQUE NOT NULL,
    transaction_status VARCHAR(50) NOT NULL,
    user_id INT NOT NULL REFERENCES m_users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE t_transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id INT NOT NULL REFERENCES t_transactions(id) ON DELETE CASCADE,
    class_id INT NOT NULL REFERENCES m_classes(id) ON DELETE CASCADE,
    ticket_code VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    attend_time TIMESTAMP DEFAULT NULL,
    attend_status BOOLEAN DEFAULT FALSE,
    attend_operator_id INT DEFAULT NULL REFERENCES m_operators(id) ON DELETE CASCADE,
    latest_sync_at TIMESTAMP DEFAULT NULL,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE UNIQUE INDEX idx_users_email ON m_users(email);
CREATE UNIQUE INDEX idx_users_phone ON m_users(phone);

CREATE INDEX idx_classes_event_id ON m_classes(event_id);
CREATE UNIQUE INDEX idx_operators_email ON m_operators(email);
CREATE UNIQUE INDEX idx_operators_phone ON m_operators(phone);
CREATE INDEX idx_transactions_user_id ON t_transactions(user_id);
CREATE INDEX idx_transactions_invoice_code ON t_transactions(invoice_code);
CREATE INDEX idx_transactions_status_date ON t_transactions(transaction_status, transaction_date);
CREATE INDEX idx_transaction_details_transaction_id ON t_transaction_details(transaction_id);
CREATE INDEX idx_transaction_details_class_id ON t_transaction_details(class_id);
CREATE INDEX idx_operators_classes_operator_id ON t_operators_classes(operator_id);
CREATE INDEX idx_operators_classes_class_id ON t_operators_classes(class_id);
CREATE UNIQUE INDEX idx_operators_classes_composite ON t_operators_classes(operator_id, class_id);
