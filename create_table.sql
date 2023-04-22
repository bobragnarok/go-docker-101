create table log_date (
	id serial PRIMARY KEY,
	ping varchar(50) NOT NULL,
	created_date TIMESTAMP NOT NULL
)