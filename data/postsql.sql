CREATE TABLE orders (
    order_uid           character(19) primary key NOT NULL,
    track_number        character(14),
    entry               character(4),
    locale              character(2),
    internal_signature  character varying(20),
    customer_id         character varying(20),
    delivery_service    character varying(20),
    shardkey            character varying(20),
    sm_id               integer,
    date_created        timestamptz,
    oof_shard           character varying(20)
);

CREATE TABLE delivery (
    id                  serial primary key NOT NULL,
    order_uid           character(19) references orders(order_uid),
    name                character varying(20),
    phone               character varying(20),
    zip                 character varying(20),
    city                character varying(20),
    address             character varying(20),
    region              character varying(20),
    email               character varying(20)
);

CREATE TABLE payment (
    transaction         character(19) primary key NOT NULL,
    order_uid           character(19) references orders(order_uid),
    request_id          character varying(20),
    currency            character(3),
    provider            character varying(20),
    amount              integer,
    payment_dt          integer,
    bank                character varying(20),
    delivery_cost       integer,
    goods_total         integer,
    custom_fee          integer
);

CREATE TABLE items (
    track_number        character(14) primary key NOT NULL, 
    order_uid           character(19) references orders(order_uid),
    chrt_id             integer, 
    price               integer, 
    rid                 character(21), 
    name                character varying(20), 
    sale                integer, 
    size                character varying(20), 
    total_price         integer, 
    nm_id               integer, 
    brand               character varying(20), 
    status              integer 
);
