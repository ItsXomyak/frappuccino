erDiagram
    substations ||--|{ power_lines : connects
    substations ||--|{ transformers : contains
    transformers ||--|{ voltage_levels : supports
    transformers ||--|{ protection_equipment : secures
    power_lines ||--o{ line_sections : segments
    voltage_levels ||--o{ electrical_loads : supplies
    electrical_loads ||--|{ meters : measures
    protection_equipment ||--o{ fault_records : logs

    substations {
        int id PK
        varchar name
        decimal capacity_MVA
        varchar location
        timestamp_tz commissioned_at
        timestamp_tz updated_at
    }

    power_lines {
        int id PK
        int substation_id FK
        varchar voltage_class
        decimal length_km
        boolean is_operational
        timestamp_tz installed_at
    }

    transformers {
        int id PK
        int substation_id FK
        decimal rated_power_MVA
        varchar primary_voltage
        varchar secondary_voltage
        boolean is_in_service
        timestamp_tz installed_at
    }

    voltage_levels {
        int id PK
        int transformer_id FK
        varchar voltage_class
        decimal max_load_MW
    }

    electrical_loads {
        int id PK
        int voltage_level_id FK
        varchar consumer_type
        decimal demand_kW
    }

    meters {
        int id PK
        int electrical_load_id FK
        varchar meter_type
        decimal accuracy_class
        timestamp_tz last_calibrated_at
    }

    protection_equipment {
        int id PK
        int transformer_id FK
        varchar protection_type
        boolean active
        timestamp_tz last_tested_at
    }

    fault_records {
        int id PK
        int protection_equipment_id FK
        varchar fault_type
        text description
        timestamp_tz occurred_at
    }

    line_sections {
        int id PK
        int power_line_id FK
        decimal section_length_km
        varchar conductor_type
        timestamp_tz maintained_at
    }
