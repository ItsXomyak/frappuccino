#!/bin/bash

# Base URL for the API
BASE_URL="http://localhost:9090"

# Function to post menu and get IDs
post_menu_and_get_ids() {
    echo "Adding menu items and collecting their IDs..."
    
    # Posting Latte and capturing response
    LATTE_RESP=$(curl -s -X POST "${BASE_URL}/menu" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Caffe Latte",
            "description": "Espresso with steamed milk",
            "price": 3.50,
            "ingredients": [
                {
                    "ingredient_id": "espresso_shot",
                    "quantity": 1
                },
                {
                    "ingredient_id": "milk",
                    "quantity": 200
                }
            ]
        }')
    LATTE_ID=$(echo "$LATTE_RESP" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    echo "Latte ID: $LATTE_ID"
    sleep 2

    # Posting Muffin and capturing response
    MUFFIN_RESP=$(curl -s -X POST "${BASE_URL}/menu" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Blueberry Muffin",
            "description": "Freshly baked muffin with blueberries",
            "price": 2.00,
            "ingredients": [
                {
                    "ingredient_id": "flour",
                    "quantity": 100
                },
                {
                    "ingredient_id": "blueberries",
                    "quantity": 20
                },
                {
                    "ingredient_id": "sugar",
                    "quantity": 30
                }
            ]
        }')
    MUFFIN_ID=$(echo "$MUFFIN_RESP" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    echo "Muffin ID: $MUFFIN_ID"
    sleep 2

    # Posting Espresso and capturing response
    ESPRESSO_RESP=$(curl -s -X POST "${BASE_URL}/menu" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Espresso",
            "description": "Strong and bold coffee",
            "price": 2.50,
            "ingredients": [
                {
                    "ingredient_id": "espresso_shot",
                    "quantity": 1
                }
            ]
        }')
    ESPRESSO_ID=$(echo "$ESPRESSO_RESP" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    echo "Espresso ID: $ESPRESSO_ID"
    sleep 2

    # Store IDs in temporary files
    echo "$LATTE_ID" > /tmp/latte_id
    echo "$MUFFIN_ID" > /tmp/muffin_id
    echo "$ESPRESSO_ID" > /tmp/espresso_id
}

# Function to post inventory items
post_inventory() {
    echo "Adding inventory items..."
    
    # Posting Espresso Shot
    curl -X POST "${BASE_URL}/inventory" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Espresso Shot",
            "quantity": 500,
            "unit": "shots"
        }'
    echo -e "\n"
    sleep 2

    # Posting Milk
    curl -X POST "${BASE_URL}/inventory" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Milk",
            "quantity": 5000,
            "unit": "ml"
        }'
    echo -e "\n"
    sleep 2

    # Posting Flour
    curl -X POST "${BASE_URL}/inventory" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Flour",
            "quantity": 10000,
            "unit": "g"
        }'
    echo -e "\n"
    sleep 2

    # Posting Blueberries
    curl -X POST "${BASE_URL}/inventory" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Blueberries",
            "quantity": 2000,
            "unit": "g"
        }'
    echo -e "\n"
    sleep 2

    # Posting Sugar
    curl -X POST "${BASE_URL}/inventory" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Sugar",
            "quantity": 5000,
            "unit": "g"
        }'
    echo -e "\n"
    sleep 2
}

# Function to post orders using stored menu item IDs
post_orders() {
    # Read IDs from temporary files
    LATTE_ID=$(cat /tmp/latte_id)
    MUFFIN_ID=$(cat /tmp/muffin_id)
    ESPRESSO_ID=$(cat /tmp/espresso_id)
    
    echo "Posting orders using menu IDs:"
    echo "Latte ID: $LATTE_ID"
    echo "Muffin ID: $MUFFIN_ID"
    echo "Espresso ID: $ESPRESSO_ID"
    
    # Order 1
    curl -X POST "${BASE_URL}/order" \
        -H "Content-Type: application/json" \
        -d '{
            "customer_name": "Alice Johnson Smith",
            "items": [
                {
                    "product_id": "'"$LATTE_ID"'",
                    "quantity": 2
                },
                {
                    "product_id": "'"$MUFFIN_ID"'",
                    "quantity": 1
                }
            ]
        }'
    echo -e "\n"
    sleep 2

    # Order 2
    curl -X POST "${BASE_URL}/order" \
        -H "Content-Type: application/json" \
        -d '{
            "customer_name": "Bob Anderson",
            "items": [
                {
                    "product_id": "'"$ESPRESSO_ID"'",
                    "quantity": 1
                },
                {
                    "product_id": "'"$MUFFIN_ID"'",
                    "quantity": 2
                }
            ]
        }'
    echo -e "\n"
    sleep 2

    # Order 3
    curl -X POST "${BASE_URL}/order" \
        -H "Content-Type: application/json" \
        -d '{
            "customer_name": "Carol Davis Wilson",
            "items": [
                {
                    "product_id": "'"$LATTE_ID"'",
                    "quantity": 1
                },
                {
                    "product_id": "'"$ESPRESSO_ID"'",
                    "quantity": 1
                }
            ]
        }'
    echo -e "\n"
    
    # Clean up temporary files
    rm -f /tmp/latte_id /tmp/muffin_id /tmp/espresso_id
}

# Handle command line arguments
process_args() {
    local has_orders=false
    local has_menu=false
    local has_inventory=false
    
    # First check what operations are requested
    for arg in "$@"; do
        case $arg in
            --orders)
                has_orders=true
                ;;
            --menu)
                has_menu=true
                ;;
            --inventory)
                has_inventory=true
                ;;
            *)
                echo "Unknown option: $arg"
                echo "Usage: $0 [--orders] [--menu] [--inventory]"
                exit 1
                ;;
        esac
    done
    
    # Execute in correct order
    if [ "$has_inventory" = true ]; then
        post_inventory
        sleep 2
    fi
    
    if [ "$has_menu" = true ]; then
        post_menu_and_get_ids
        sleep 2
    fi
    
    if [ "$has_orders" = true ]; then
        post_orders
    fi
}

# If no arguments provided, show usage
if [ $# -eq 0 ]; then
    echo "Usage: $0 [--orders] [--menu] [--inventory]"
    exit 1
fi

# Process arguments
process_args "$@"