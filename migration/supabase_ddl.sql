-- =====================================
-- TABLE: roles
-- =====================================
CREATE TABLE roles (
id SERIAL PRIMARY KEY,
name VARCHAR(100) NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =====================================
-- TABLE: users
-- =====================================
CREATE TABLE users (
id SERIAL PRIMARY KEY,
email VARCHAR(100) NOT NULL UNIQUE,
password VARCHAR(100) NOT NULL,
name VARCHAR(100) NOT NULL,
role_id INT NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY (role_id) REFERENCES roles(id)
);

-- =====================================
-- TABLE: products
-- =====================================
CREATE TABLE products (
id SERIAL PRIMARY KEY,
name VARCHAR(100) NOT NULL,
description TEXT,
price DECIMAL(12,2) NOT NULL,
stock INT DEFAULT 0,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =====================================
-- TABLE: status_order
-- =====================================
CREATE TABLE status_order (
id SERIAL PRIMARY KEY,
status_name VARCHAR(50) NOT NULL UNIQUE,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =====================================
-- TABLE: orders
-- =====================================
CREATE TABLE orders (
id SERIAL PRIMARY KEY,
user_id INT NOT NULL,
order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
status_id INT NOT NULL,
total_amount DECIMAL(12,2) NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY (user_id) REFERENCES users(id),
FOREIGN KEY (status_id) REFERENCES status_order(id)
);

-- =====================================
-- TABLE: order_items
-- =====================================
CREATE TABLE order_items (
order_id INT NOT NULL,
product_id INT NOT NULL,
quantity INT NOT NULL,
price_at_order DECIMAL(12,2),
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY (order_id, product_id),
FOREIGN KEY (order_id) REFERENCES orders(id),
FOREIGN KEY (product_id) REFERENCES products(id)
);
