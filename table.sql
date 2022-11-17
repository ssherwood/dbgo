create table user_account(id bigserial primary key,
                   account_status varchar(255),
                   city varchar(255),
                   state varchar(255),
                   auto_renew bool,
                   about varchar(255),
                   days_active integer,
                   encryption_level smallint,
                   confidence real,
                   created_date timestamptz);