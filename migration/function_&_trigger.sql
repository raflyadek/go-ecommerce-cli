-- =====================================================
-- Function global untuk auto-update kolom updated_at
-- =====================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- Trigger untuk tabel orders dan products
-- =====================================================
CREATE TRIGGER set_timestamp_orders
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER set_timestamp_products
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- VIEW: Produk Paling Laku
-- =====================================================
CREATE OR REPLACE VIEW view_best_selling_products AS
WITH product_sales AS (
    SELECT 
        p.id AS product_id,                    
        p.name AS product_name,                
        SUM(oi.quantity) AS total_sold,       
        SUM(oi.quantity * oi.price_at_order) AS total_revenue 
    FROM order_items oi
    JOIN products p ON p.id = oi.product_id   
    JOIN orders o ON o.id = oi.order_id       
    JOIN status_order so ON so.id = o.status_id 
    WHERE LOWER(so.status_name) = 'completed' 
    GROUP BY p.id, p.name
)
SELECT 
    product_id,
    product_name,
    total_sold,
    total_revenue
FROM product_sales
ORDER BY total_sold DESC;

-- =====================================================
-- VIEW: Order Aktif (User)
-- =====================================================
CREATE OR REPLACE VIEW orders_active AS
SELECT 
    o.id, 
    o.user_id,
    s.status_name, 
    o.total_amount, 
    COUNT(oi.product_id) AS item_count, 
    o.created_at
FROM orders o
JOIN status_order s ON o.status_id = s.id
LEFT JOIN order_items oi ON o.id = oi.order_id
WHERE o.status_id IN (1, 2, 3)  -- misal: pending, processing, shipped
GROUP BY o.id, o.user_id, s.status_name, o.total_amount, o.created_at;

-- =====================================================
-- VIEW: Riwayat Order (User)
-- =====================================================
CREATE OR REPLACE VIEW orders_history AS
SELECT 
    o.id, 
    o.user_id,
    s.status_name, 
    o.total_amount, 
    COUNT(oi.product_id) AS item_count, 
    o.created_at
FROM orders o
JOIN status_order s ON o.status_id = s.id
LEFT JOIN order_items oi ON o.id = oi.order_id
WHERE o.status_id IN (4, 5)  -- misal: completed, cancelled
GROUP BY o.id, o.user_id, s.status_name, o.total_amount, o.created_at;
