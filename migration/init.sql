-- DDL 
-- create database
CREATE DATABASE go-ecommerce-cli;

-- create table users
CREATE TABLE users(id INT SERIAL PRIMARY KEY, email VARCHAR(100) UNIQUE NOT NULL, password VARCHAR(100) NOT NULL, name VARCHAR(100) NOT NULL, role_id INT REFERENCES roles(id) ON DELETE CASCADE) NOT NULL;

-- create table roles
CREATE TABLE roles(id INT SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL);

-- create table orders
CREATE TABLE orders(id INT SERIAL PRIMARY KEY, user_id INT REFERENCES users(id) NOT NULL, order_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL, status VARCHAR(100) (REFERENCE BELOM), total_amount decimal(12, 2) NOT NULL, update_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP);

-- order item (on delete casacde ensure if we delete the parent id, the child id on reference will automatically be deleted also)
CREATE TABLE order_item(order_id INT REFERENCES orders(id) ON DELETE CASCADE, product_id INT REFEREMCES products(id) ON DELETE CASCADE, quantity INT, price_at_order DECIMAL(12, 2), PRIMARY KEY (order_id, product_id), FOREIGN KEY (order_id) REFERENCES orders(id), FOREIGN KEY (product_id) REFERENCES products(id));

-- create table products
CREATE TABLE products(id INT SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL, description TEXT, created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP);
