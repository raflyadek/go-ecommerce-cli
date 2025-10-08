-- Role sesuai dokumen: Admin, Staff, User
INSERT INTO roles (name) VALUES
('admin'),
('staff'),
('user');

-- Status order sesuai alur use case: pending → paid → shipped → completed → cancelled
INSERT INTO status_order (status_name) VALUES
('pending'),
('paid'),
('shipped'),
('completed'),
('cancelled');

-- Dummy akun admin & staff (untuk testing CLI)
-- Password nanti di-hash pakai bcrypt di aplikasi Go
INSERT INTO users (email, password, name, role_id) VALUES
('admin@mail.com', 'admin123', 'Admin One', 1),
('staff@mail.com', 'staff123', 'Staff One', 2),
('user@mail.com', 'user123', 'User One', 3);

-- Produk contoh (pakaian)
INSERT INTO products (name, description, price, stock) VALUES
('T-Shirt Oversized', 'Kaos oversized bahan premium', 120000, 50),
('Hoodie Unisex', 'Hoodie hangat cocok untuk semua musim', 250000, 30),
('Celana Cargo', 'Celana cargo bahan tebal dan stylish', 180000, 40),
('Jaket Bomber', 'Jaket bomber ringan, cocok untuk hangout', 300000, 25),
('Sweater Rajut', 'Sweater rajut lembut, nyaman dipakai', 200000, 35);

