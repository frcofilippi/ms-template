create table if not exists products (
	id serial primary key,
	customer_id int not null,
	name varchar(250) not null,
	cost float not null
);

alter table products enable row level security;

create policy customer_isolation_policy on products
for all
using (customer_id = current_setting('app.customer_id')::int);

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO pedimeapp;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO pedimeapp;