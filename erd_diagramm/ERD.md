erDiagram
customers ||--o{ orders : places
orders ||--|{ order_items : contains
order_items ||--o{ menu_items : includes
menu_items ||--o{ menu_item_ingredients : requires
menu_item_ingredients ||--o{ inventory : uses
orders ||--o{ order_status_history : tracks
menu_items ||--o{ price_history : logs
inventory ||--o{ inventory_transactions : records

    customers {
        int id PK
        varchar name
        jsonb preferences
    }

    orders {
        int id PK
        int customer_id FK
        enum status "pending, preparing, ready, delivered, cancelled"
        decimal total_amount
        enum payment_method "cash, card, online"
        jsonb special_instructions
        timestamp_tz created_at
        timestamp_tz updated_at
    }

    order_items {
        int id PK
        int order_id FK
        int menu_item_id FK
        int quantity
        decimal price
        jsonb customizations
    }

    menu_items {
        int id PK
        varchar name
        varchar description
        text[] categories
        text[] allergens
        decimal price
        boolean available
        enum size "small, medium, large"
    }

    inventory {
        int id PK
        varchar name UK
        decimal stock
        varchar unit
        decimal reorder_threshold
        decimal price
    }

    menu_item_ingredients {
        int id PK
        int menu_item_id FK
        int ingredient_id FK
        decimal quantity
        varchar unit
    }

    order_status_history {
        int id PK
        int order_id FK
        enum status "pending, preparing, ready, delivered, cancelled"
        timestamp_tz changed_at
    }

    price_history {
        int id PK
        int menu_item_id FK
        decimal old_price
        decimal new_price
        timestamp_tz changed_at
    }

    inventory_transactions {
        int id PK
        int ingredient_id FK
        decimal change_amount
        varchar transaction_type "purchase, use"
        timestamp_tz occurred_at
    }
