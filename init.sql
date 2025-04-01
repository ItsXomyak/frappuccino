-- ENUM Types
CREATE TYPE order_status AS ENUM ('pending', 'preparing', 'ready', 'delivered', 'cancelled');
CREATE TYPE payment_method AS ENUM ('cash', 'card', 'online');
CREATE TYPE item_size AS ENUM ('small', 'medium', 'large');
CREATE TYPE staff_role AS ENUM ('admin', 'chef', 'waiter', 'cashier');

-- Tables
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    preferences JSONB DEFAULT '{}'::JSONB
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
    status order_status DEFAULT 'pending',
    total_amount DECIMAL(10,2) NOT NULL CHECK (total_amount >= 0),
    payment_method payment_method NOT NULL,
    special_instructions JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id INT NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    customizations JSONB DEFAULT '{}'::JSONB
);

CREATE TABLE menu_items (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    categories TEXT[] DEFAULT '{}',
    allergens TEXT[] DEFAULT '{}',
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    available BOOLEAN DEFAULT TRUE,
    size item_size NOT NULL
);

CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    stock DECIMAL(10,2) NOT NULL CHECK (stock >= 0),
    unit TEXT NOT NULL,
    reorder_threshold DECIMAL(10,2) NOT NULL CHECK (reorder_threshold >= 0),
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0)
);


CREATE TABLE menu_item_ingredients (
    id SERIAL PRIMARY KEY,
    menu_item_id INT REFERENCES menu_items(id) ON DELETE CASCADE,
    ingredient_id INT REFERENCES inventory(id) ON DELETE RESTRICT,
    quantity DECIMAL(10,2) NOT NULL CHECK (quantity > 0),
    unit TEXT NOT NULL
);

CREATE TABLE order_status_history (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    status order_status NOT NULL,
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE price_history (
    id SERIAL PRIMARY KEY,
    menu_item_id INT REFERENCES menu_items(id) ON DELETE CASCADE,
    old_price DECIMAL(10,2) NOT NULL CHECK (old_price >= 0),
    new_price DECIMAL(10,2) NOT NULL CHECK (new_price >= 0),
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE inventory_transactions (
    id SERIAL PRIMARY KEY,
    ingredient_id INT REFERENCES inventory(id) ON DELETE CASCADE,
    change_amount DECIMAL(10,2) NOT NULL,
    transaction_type TEXT NOT NULL CHECK (transaction_type IN ('purchase', 'use')),
    occurred_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Trigger for updated_at
CREATE FUNCTION update_timestamp() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_orders_timestamp
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();

-- Indexes
CREATE INDEX idx_menu_items_name ON menu_items USING GIN (to_tsvector('english', name));
CREATE INDEX idx_customers_name ON customers USING GIN (to_tsvector('english', name));
CREATE INDEX idx_order_items_order_id ON order_items (order_id);
CREATE INDEX idx_orders_created_at ON orders (created_at);

-- Мок дата
-- Customers 
INSERT INTO customers (name, preferences) VALUES
('Alice Smith', '{"coffee": "latte", "size": "large"}'),
('Bob Johnson', '{"no_sugar": true}'),
('Charlie Brown', '{}');

-- Inventory 
INSERT INTO inventory (name, stock, unit, reorder_threshold) VALUES
('Coffee Beans', 1000, 'g', 200),
('Milk', 5000, 'ml', 1000),
('Sugar', 2000, 'g', 500),
('Vanilla Syrup', 300, 'ml', 100),
('Flour', 10000, 'g', 2000),
('Butter', 5000, 'g', 1000),
('Eggs', 200, 'pcs', 50),
('Chocolate', 1500, 'g', 300),
('Cream', 2000, 'ml', 500),
('Cinnamon', 500, 'g', 100),
('Almonds', 1000, 'g', 200),
('Hazelnuts', 800, 'g', 150),
('Salt', 3000, 'g', 500),
('Oat Milk', 3000, 'ml', 600),
('Caramel Syrup', 400, 'ml', 100),
('Blueberries', 2000, 'g', 400),
('Yeast', 1000, 'g', 200),
('Honey', 500, 'ml', 100),
('Cocoa Powder', 800, 'g', 150),
('Whipped Cream', 1000, 'ml', 200);

-- Menu Items 
INSERT INTO menu_items (name, description, categories, allergens, price, available, size) VALUES
    ('Latte', 'Espresso with steamed milk', ARRAY['coffee', 'hot']::text[], ARRAY['milk']::text[], 4.50, TRUE, 'medium'),
    ('Espresso', 'Strong black coffee', ARRAY['coffee', 'hot']::text[], ARRAY[]::text[], 2.50, TRUE, 'small'),
    ('Muffin', 'Blueberry muffin', ARRAY['pastry']::text[], ARRAY['gluten']::text[], 3.00, TRUE, 'medium'),
    ('Cappuccino', 'Espresso with foamed milk', ARRAY['coffee', 'hot']::text[], ARRAY['milk']::text[], 4.00, TRUE, 'medium'),
    ('Croissant', 'Buttery pastry', ARRAY['pastry']::text[], ARRAY['gluten']::text[], 2.80, TRUE, 'medium'),
    ('Mocha', 'Coffee with chocolate', ARRAY['coffee', 'hot']::text[], ARRAY['milk']::text[], 5.00, TRUE, 'large'),
    ('Americano', 'Espresso with hot water', ARRAY['coffee', 'hot']::text[], ARRAY[]::text[], 3.00, TRUE, 'medium'),
    ('Chocolate Cake', 'Rich chocolate cake', ARRAY['dessert']::text[], ARRAY['gluten', 'milk']::text[], 6.50, TRUE, 'medium'),
    ('Oat Latte', 'Latte with oat milk', ARRAY['coffee', 'hot']::text[], ARRAY[]::text[], 5.50, TRUE, 'large'),
    ('Cinnamon Roll', 'Sweet roll with cinnamon', ARRAY['pastry']::text[], ARRAY['gluten']::text[], 3.50, TRUE, 'medium');


-- Menu Item Ingredients
INSERT INTO menu_item_ingredients (menu_item_id, ingredient_id, quantity, unit) VALUES
(1, 1, 30, 'g'),   -- Latte: Coffee Beans
(1, 2, 200, 'ml'),  -- Latte: Milk
(2, 1, 20, 'g'),    -- Espresso: Coffee Beans
(3, 5, 150, 'g'),   -- Muffin: Flour
(3, 16, 50, 'g'),   -- Muffin: Blueberries
(4, 1, 25, 'g'),    -- Cappuccino: Coffee Beans
(4, 2, 150, 'ml'),  -- Cappuccino: Milk
(5, 5, 100, 'g'),   -- Croissant: Flour
(5, 6, 50, 'g'),    -- Croissant: Butter
(6, 1, 30, 'g'),    -- Mocha: Coffee Beans
(6, 2, 200, 'ml'),  -- Mocha: Milk
(6, 8, 20, 'g'),    -- Mocha: Chocolate
(7, 1, 20, 'g'),    -- Americano: Coffee Beans
(8, 5, 200, 'g'),   -- Chocolate Cake: Flour
(8, 8, 100, 'g'),   -- Chocolate Cake: Chocolate
(9, 1, 30, 'g'),    -- Oat Latte: Coffee Beans
(9, 14, 200, 'ml'), -- Oat Latte: Oat Milk
(10, 5, 150, 'g'),  -- Cinnamon Roll: Flour
(10, 10, 5, 'g');   -- Cinnamon Roll: Cinnamon

-- Orders 
INSERT INTO orders (customer_id, status, total_amount, payment_method, special_instructions, created_at) VALUES
(1, 'delivered', 4.50, 'card', '{"extra_shot": true}', '2025-03-20 10:00:00+00'),
(2, 'pending', 2.50, 'cash', '{}', '2025-03-24 15:00:00+00'),
(3, 'preparing', 3.00, 'online', '{"no_wrapper": true}', '2025-03-25 09:00:00+00'),
(1, 'ready', 5.00, 'card', '{}', '2025-03-22 12:00:00+00'),
(2, 'cancelled', 3.50, 'cash', '{"urgent": true}', '2025-03-23 14:00:00+00'),
(3, 'delivered', 6.50, 'online', '{}', '2025-02-15 08:00:00+00'),
(1, 'delivered', 4.00, 'card', '{"less_milk": true}', '2025-01-10 11:00:00+00'),
(2, 'preparing', 2.80, 'cash', '{}', '2025-03-25 10:00:00+00'),
(3, 'pending', 5.50, 'online', '{"extra_syrup": true}', '2025-03-24 16:00:00+00'),
(1, 'delivered', 3.00, 'card', '{}', '2025-02-20 13:00:00+00'),
(2, 'ready', 4.50, 'cash', '{}', '2025-03-21 09:00:00+00'),
(3, 'delivered', 2.50, 'online', '{}', '2025-03-15 17:00:00+00'),
(1, 'preparing', 6.50, 'card', '{"no_cream": true}', '2025-03-25 11:00:00+00'),
(2, 'pending', 3.50, 'cash', '{}', '2025-03-23 18:00:00+00'),
(3, 'delivered', 5.00, 'online', '{}', '2025-02-25 10:00:00+00'),
(1, 'ready', 4.00, 'card', '{}', '2025-03-20 15:00:00+00'),
(2, 'cancelled', 2.80, 'cash', '{}', '2025-03-22 16:00:00+00'),
(3, 'delivered', 5.50, 'online', '{"extra_hot": true}', '2025-03-10 12:00:00+00'),
(1, 'preparing', 3.00, 'card', '{}', '2025-03-24 14:00:00+00'),
(2, 'pending', 6.50, 'cash', '{}', '2025-03-25 08:00:00+00'),
(3, 'delivered', 4.50, 'online', '{}', '2025-01-15 09:00:00+00'),
(1, 'ready', 2.50, 'card', '{}', '2025-03-23 11:00:00+00'),
(2, 'delivered', 3.50, 'cash', '{"no_sugar": true}', '2025-03-20 17:00:00+00'),
(3, 'preparing', 5.00, 'online', '{}', '2025-03-24 13:00:00+00'),
(1, 'pending', 4.00, 'card', '{}', '2025-03-25 12:00:00+00'),
(2, 'delivered', 2.80, 'cash', '{}', '2025-03-21 10:00:00+00'),
(3, 'ready', 5.50, 'online', '{"extra_shot": true}', '2025-03-22 14:00:00+00'),
(1, 'delivered', 3.00, 'card', '{}', '2025-02-10 16:00:00+00'),
(2, 'preparing', 6.50, 'cash', '{}', '2025-03-23 15:00:00+00'),
(3, 'pending', 4.50, 'online', '{}', '2025-03-25 07:00:00+00');

-- Order Items
INSERT INTO order_items (order_id, menu_item_id, quantity, price, customizations) VALUES
(1, 1, 1, 4.50, '{"milk": "oat"}'),
(2, 2, 1, 2.50, '{}'),
(3, 3, 1, 3.00, '{}'),
(4, 6, 1, 5.00, '{}'),
(5, 10, 1, 3.50, '{}'),
(6, 8, 1, 6.50, '{}'),
(7, 4, 1, 4.00, '{}'),
(8, 5, 1, 2.80, '{}'),
(9, 9, 1, 5.50, '{}'),
(10, 7, 1, 3.00, '{}'),
(11, 1, 1, 4.50, '{}'),
(12, 2, 1, 2.50, '{}'),
(13, 8, 1, 6.50, '{}'),
(14, 10, 1, 3.50, '{}'),
(15, 6, 1, 5.00, '{}'),
(16, 4, 1, 4.00, '{}'),
(17, 5, 1, 2.80, '{}'),
(18, 9, 1, 5.50, '{}'),
(19, 7, 1, 3.00, '{}'),
(20, 8, 1, 6.50, '{}'),
(21, 1, 1, 4.50, '{}'),
(22, 2, 1, 2.50, '{}'),
(23, 10, 1, 3.50, '{}'),
(24, 6, 1, 5.00, '{}'),
(25, 4, 1, 4.00, '{}'),
(26, 5, 1, 2.80, '{}'),
(27, 9, 1, 5.50, '{}'),
(28, 7, 1, 3.00, '{}'),
(29, 8, 1, 6.50, '{}'),
(30, 1, 1, 4.50, '{}');

-- Order Status History
INSERT INTO order_status_history (order_id, status, changed_at) VALUES
(1, 'pending', '2025-03-20 10:00:00+00'),
(1, 'preparing', '2025-03-20 10:05:00+00'),
(1, 'delivered', '2025-03-20 10:30:00+00'),
(4, 'pending', '2025-03-22 12:00:00+00'),
(4, 'preparing', '2025-03-22 12:10:00+00'),
(4, 'ready', '2025-03-22 12:20:00+00'),
(5, 'pending', '2025-03-23 14:00:00+00'),
(5, 'cancelled', '2025-03-23 14:15:00+00'),
(6, 'pending', '2025-02-15 08:00:00+00'),
(6, 'delivered', '2025-02-15 08:30:00+00');

-- Price History
INSERT INTO price_history (menu_item_id, old_price, new_price, changed_at) VALUES
(1, 4.00, 4.50, '2025-01-01 00:00:00+00'),
(2, 2.00, 2.50, '2025-02-01 00:00:00+00'),
(3, 2.50, 3.00, '2024-12-01 00:00:00+00'),
(4, 3.50, 4.00, '2025-01-15 00:00:00+00'),
(8, 6.00, 6.50, '2025-02-15 00:00:00+00');

-- Inventory Transactions
INSERT INTO inventory_transactions (ingredient_id, change_amount, transaction_type, occurred_at) VALUES
(1, 1000, 'purchase', '2025-03-01 08:00:00+00'),
(1, -30, 'use', '2025-03-20 10:00:00+00'),
(2, 5000, 'purchase', '2025-03-01 08:00:00+00'),
(2, -200, 'use', '2025-03-20 10:00:00+00'),
(5, 10000, 'purchase', '2025-03-01 08:00:00+00'),
(5, -150, 'use', '2025-03-25 09:00:00+00');