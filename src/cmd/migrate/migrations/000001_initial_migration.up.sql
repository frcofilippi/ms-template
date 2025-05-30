CREATE TABLE
    IF NOT EXISTS customers (
        id SERIAL PRIMARY KEY,
        name VARCHAR(250) NOT NULL,
        created_at TIME
        WITH
            TIME ZONE NOT NULL
    );


create table if not exists products (
	id serial primary key,
	customer_id int not null,
	name varchar(250) not null,
	cost float not null
);

CREATE TABLE IF NOT EXISTS outbox_messages (
    id serial primary key,
    aggregate_type TEXT NOT NULL,          -- e.g. "Order", "Inventory", "Payment"
    aggregate_id INT NOT NULL,            -- ID of the entity that triggered the event
    event_type TEXT NOT NULL,              -- e.g. "OrderSubmitted", "InventoryReserved"
    payload JSONB NOT NULL,                -- Event data
    occurred_at TIMESTAMP WITH TIME ZONE DEFAULT now(), -- when the event was created
    processed BOOLEAN DEFAULT FALSE,       -- Has it been sent to the message broker?
    processed_at TIMESTAMP WITH TIME ZONE  -- Optional: when it was published
);

alter table products enable row level security;

create policy customer_isolation_policy on products
for all
using (customer_id = current_setting('app.customer_id')::int);

GRANT USAGE ON SCHEMA public TO pedimeapp;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO pedimeapp;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO pedimeapp;