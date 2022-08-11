CREATE TYPE gender AS ENUM ('male', 'female');
CREATE TYPE user_type AS ENUM ('client', 'driver');
CREATE TYPE place_type AS ENUM ('home', 'work', 'custom');
CREATE TYPE ride_status AS ENUM ('new', 'on_the_way', 'cancelled', 'done');
CREATE TYPE order_type AS ENUM ('city', 'interregional');
CREATE TYPE order_status AS ENUM ('client_cancelled', 'driver_cancelled', 'new', 'driver_accepted', 'driver_arrived', 'client_going_out', 'trip_started', 'order_completed');
CREATE TYPE payment_type AS ENUM ('cash', 'card');
CREATE TYPE cargo_type AS ENUM ('no', 'small', 'medium', 'large');
CREATE TYPE message_type AS ENUM ('message', 'audio', 'file');
CREATE TYPE driver_sts AS ENUM ('online', 'offline', 'on_the_way');

CREATE TABLE IF NOT EXISTS users
(
    id serial not null unique,
    driver_id int,
    client_id int,
    user_type user_type not null default 'client',
    login varchar(255) not null,
    password_hash varchar(255) not null,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS clients
(
    id serial not null unique,
    name varchar(255) not null,
    surname varchar(255) not null,
    birthdate date not null,
    gender gender not null,
    avatar varchar(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS saved_addresses
(
    id serial not null unique,
    user_id int not null,
    name varchar(255) not null,
    place_type place_type not null,
    location varchar(255) not null,
    address varchar(255) not null,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS credit_cards
(
    id serial not null unique,
    user_id int not null,
    card_info varchar(255) not null,
    is_active boolean not null default false,
    is_main boolean not null default false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rides
(
    id serial not null unique,
    driver_id int not null,
    from_district_id int not null,
    to_district_id int not null,
    departure_date timestamp not null,
    price varchar(255) not null,
    passenger_count integer not null,
    comments varchar(500),
    status ride_status not null,
    view_count integer not null default 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders
(
    id serial not null unique,
    driver_id int,
    client_id int not null,
    order_id int not null,
    order_type order_type not null default 'city',
    order_status order_status not null default 'new',
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS interregional_orders
(
    id serial not null unique,
    ride_id int not null,
    from_district_id int not null,
    to_district_id int not null,
    price numeric(12,2) not null,
    passenger_count integer not null,
    departure_date timestamp not null,
    comments varchar(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS city_orders
(
    id serial not null unique,
    points jsonb not null default '{}'::jsonb,
    tariff_id int not null,
    cargo_type cargo_type not null default 'no',
    payment_type payment_type not null default 'cash',
    card_id int,
    has_conditioner boolean not null default false,
    for_another boolean not null default false,
    for_another_phone varchar(100),
    receiver_comments varchar(255),
    receiver_phone varchar(100),
    price numeric(12,2) not null,
    comments varchar(500),
    ride_info jsonb not null default '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS canceled_orders
(
    id serial not null unique,
    order_type order_type not null default 'city',
    user_type user_type not null default 'client',
    user_id int not null,
    order_id int not null,
    comments varchar(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS canceled_order_reasons
(
    id serial not null unique,
    canceled_order_id int not null,
    reason_id int not null
);

CREATE TABLE IF NOT EXISTS rated_orders
(
    id serial not null unique,
    order_type order_type not null default 'city',
    user_type user_type not null default 'client',
    rate int not null,
    user_id int not null,
    order_id int not null,
    comments varchar(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rated_order_reasons
(
    id serial not null unique,
    rated_order_id int not null,
    reason_id int not null
);

CREATE TABLE IF NOT EXISTS chat_messages
(
    id serial not null unique,
    user_type user_type not null default 'client',
    driver_id int not null default 0,
    client_id int not null default 0,
    ride_id int not null default 0,
    order_id int not null default 0,
    message_type message_type not null default 'message',
    content text not null,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ride_view_counts
(
    id serial not null unique,
    user_id int not null,
    ride_id int not null,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS driver_enabled_tariffs
(
    id serial not null unique,
    user_id int not null,
    tariff_id int not null,
    is_active boolean not null default true
);

CREATE TABLE IF NOT EXISTS driver_statuses
(
    id serial not null unique,
    user_id int not null,
    driver_status driver_sts not null default 'offline'
);