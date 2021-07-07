CREATE TABLE IF NOT EXISTS products
(
    id SERIAL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    sku varchar(100),
    CONSTRAINT products_pkey PRIMARY KEY (id)
);

/* Table 'users' */
CREATE TABLE public.users
(
    user_id            serial                      NOT NULL,
    username           character varying(100)      NOT NULL,
    description        character varying(500),
    email              character varying(200),
    api_key            character varying(64),
    api_key_updated    timestamp without time zone,
    api_key_expiration timestamp without time zone,
    user_type          character varying(10),
    created            timestamp without time zone NOT NULL,
    updated            timestamp without time zone,
    disabled           boolean                     NOT NULL,
    PRIMARY KEY (user_id)
);

