CREATE USER samurai WITH PASSWORD '0000';
CREATE DATABASE l0 WITH ENCODING 'UTF8';
GRANT ALL PRIVILEGES ON DATABASE l0 TO postgres;

\connect l0 postgres

-- таблица orders хранит данные структуры Order
CREATE TABLE "orders"
(
    order_uid          text not null unique primary key,
    track_number       text,
    entry              text,
    locale             text,
    internal_signature text,
    customer_id        text,
    delivery_service   text,
    shardkey           text,
    sm_id              integer,
    date_created       timestamp(0),
    oof_shard          text
);

-- таблица deliveries хранит данные структуры Delivery
CREATE TABLE "deliveries"
(
    order_uid text not null unique primary key,
    phone     text,
    name      text,
    zip       text,
    city      text,
    address   text,
    region    text,
    email     text
);

-- таблица payments хранит данные структуры Payment
CREATE TABLE "payments"
(
    order_uid     text not null unique primary key,
    transaction   text,
    request_id    text,
    currency      text,
    provider      text,
    amount        integer,
    payment_dt    integer,
    bank          text,
    delivery_cost integer,
    goods_total   integer,
    custom_fee    integer
);

-- таблица items хранит данные структуры Item
CREATE TABLE "items"
(
    chrt_id      integer not null unique primary key,
    track_number text,
    price        integer,
    rid          text,
    name         text,
    sale         integer,
    size         text,
    total_price  integer,
    nm_id        integer,
    brand        text,
    status       integer
);

-- таблица ordersitems хранит связи между таблицами orders и items
CREATE TABLE "ordersitems"
(
    order_uid text    not null,
    chrt_id   integer not null
);



-- функция check_order проверяет наличие сохраненной структуры Order
CREATE OR REPLACE FUNCTION check_order(ordr_uid text)
    RETURNS boolean
    LANGUAGE sql
AS
$$
-- запрос, определяющий наличие order_uid в таблице orders
SELECT count(order_uid) > 0
FROM orders
WHERE order_uid = ordr_uid;
$$;

-- функция get_orders выводит данные всех сохраненных структур Order
CREATE OR REPLACE FUNCTION get_orders()
-- формирование таблицы с возвращаемыми данными
    RETURNS TABLE
            (
                ordr_uid          text,
                trck_number       text,
                entr              text,
                lcale             text,
                intrnal_signature text,
                cstomer_id        text,
                dlivery_service   text,
                shrdkey           text,
                smid              integer,
                dte_created       timestamp(0),
                oof_shrd          text
            )
    LANGUAGE plpgsql
AS
$$
BEGIN
    -- запрос, выводящий все сохраненные записи структуры Order
    RETURN QUERY
        SELECT order_uid,
               track_number,
               entry,
               locale,
               internal_signature,
               customer_id,
               delivery_service,
               shardkey,
               sm_id,
               date_created,
               oof_shard
        FROM orders;
END;
$$;

-- функция get_orders_delivery выводит данные структуры Delivery,
-- соответствующей заданной структуре Order
CREATE OR REPLACE FUNCTION get_orders_delivery(order_id text)
-- формирование таблицы с возвращаемыми данными
    RETURNS TABLE
            (
                nme    text,
                phne   text,
                zp     text,
                cty    text,
                addrss text,
                rgion  text,
                emal   text
            )
    LANGUAGE sql
AS
$$
    -- запрос, выводящий записи структуры Delivery,
-- в которых ключ равен значению входного параметра
-- (связь один к одному)
SELECT name,
       phone,
       zip,
       city,
       address,
       region,
       email
FROM deliveries
WHERE order_uid = $1;
$$;

-- функция get_orders_payment выводит данные структуры Payment,
-- соответствующей заданной структуре Order
CREATE OR REPLACE FUNCTION get_orders_payment(order_id text)
-- формирование таблицы с возвращаемыми данными
    RETURNS TABLE
            (
                trnsaction   text,
                rquest_id    text,
                crrency      text,
                prvider      text,
                amnt         integer,
                pyment_dt    integer,
                bnk          text,
                dlivery_cost integer,
                gds_total    integer,
                cstom_fee    integer
            )
    LANGUAGE sql
AS
$$
    -- запрос, выводящий записи структуры Payment,
-- в которых ключ равен значению входного параметра
-- (связь один к одному)
SELECT transaction,
       request_id,
       currency,
       provider,
       amount,
       payment_dt,
       bank,
       delivery_cost,
       goods_total,
       custom_fee
FROM payments
WHERE order_uid = $1;
$$;

-- функция get_orders_items выводит данные всех структур Item,
-- соответствующих заданной структуре Order*/
CREATE OR REPLACE FUNCTION get_orders_items(order_id text)
-- формирование таблицы с возвращаемыми данными
    RETURNS TABLE
            (
                chrtid      integer,
                trck_number text,
                prce        integer,
                rd          text,
                nme         text,
                sle         integer,
                sze         text,
                ttal_price  integer,
                nmid        integer,
                brnd        text,
                sttus       integer
            )
    LANGUAGE sql
AS
$$

    -- запрос, выводящий все записи структуры Item из таблицы items,
-- в которых ключ таблицы items, в таблице ordersitems, соответствует ключу таблицы orders,
-- значение которого равно значению входного параметра
-- (связь многие ко многим)
SELECT chrt_id,
       track_number,
       price,
       rid,
       name,
       sale,
       size,
       total_price,
       nm_id,
       brand,
       status
FROM items
WHERE chrt_id IN (
    SELECT oi.chrt_id
    FROM ordersitems oi
    WHERE oi.order_uid = $1
);
$$;

-- процедура добавления данных структуры Order
CREATE OR REPLACE PROCEDURE insert_order(
    ordr_uid text,
    trck_number text,
    entr text,
    locle text,
    intrnl_signature text,
    cstmr_id text,
    dlvry_service text,
    shrdkey text,
    smid integer,
    dte_created timestamp(0),
    oof_shrd text
)
    LANGUAGE sql
AS
$$
INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id,
                    delivery_service, shardkey, sm_id, date_created, oof_shard)
VALUES (ordr_uid, trck_number, entr, locle, intrnl_signature, cstmr_id,
        dlvry_service, shrdkey, smid, dte_created, oof_shrd)
$$;

-- процедура добавления данных структуры Delivery
CREATE OR REPLACE PROCEDURE insert_orders_delivery(
    nme text,
    phne text,
    zp text,
    cty text,
    addrss text,
    rgon text,
    emal text,
    ordr_uid text
)
    LANGUAGE sql
AS
$$
INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
VALUES (ordr_uid, nme, phne, zp, cty, addrss, rgon, emal)
$$;

-- процедура добавления данных структуры Payment
CREATE OR REPLACE PROCEDURE insert_orders_payment(
    trnsaction text,
    rquest_id text,
    crrency text,
    prvider text,
    amnt integer,
    pyment_dt integer,
    bnk text,
    dlivery_cost integer,
    gds_total integer,
    cstom_fee integer,
    ordr_uid text
)
    LANGUAGE sql
AS
$$
INSERT INTO payments (order_uid, transaction, request_id, currency, provider, amount,
                      payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES (ordr_uid, trnsaction, rquest_id, crrency, prvider, amnt,
        pyment_dt, bnk, dlivery_cost, gds_total, cstom_fee)
$$;

-- процедура добавления данных структуры Item
CREATE OR REPLACE PROCEDURE insert_orders_item(
    ch_id integer,
    trck_number text,
    prc integer,
    r text,
    nme text,
    sle integer,
    sze text,
    ttl_price integer,
    nmid integer,
    brnd text,
    stts integer,
    ordr_uid text
)
    LANGUAGE plpgsql
AS
$$
DECLARE
    duplicate BOOLEAN;
BEGIN
    -- проверка, есть ли идентичная ранее добавленая запись в таблице items
    SELECT count(i.chrt_id) > 0 INTO duplicate FROM items i WHERE i.chrt_id = ch_id;
    IF duplicate THEN
        -- если дубликат существует, то просто добавляет связь в таблицу ordersitems
        INSERT INTO ordersitems (order_uid, chrt_id)
        VALUES (ordr_uid, ch_id);
    ELSE
        -- если дубликат не найден, то добавляем данные в таблицу items
        INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
        VALUES (ch_id, trck_number, prc, r, nme, sle, sze, ttl_price, nmid, brnd, stts);
        -- добавление связи в таблицу ordersitems
        INSERT INTO ordersitems (order_uid, chrt_id)
        VALUES (ordr_uid, ch_id);
    END IF;
END;
$$;

-- выдача пользователю привелегий для работы
/*GRANT ALL PRIVILEGES ON FUNCTION get_orders() TO samurai;
GRANT ALL PRIVILEGES ON FUNCTION check_order(text) TO samurai;
GRANT ALL PRIVILEGES ON FUNCTION get_orders_delivery(text) TO samurai;
GRANT ALL PRIVILEGES ON FUNCTION get_orders_items(text) TO samurai;
GRANT ALL PRIVILEGES ON FUNCTION get_orders_payment(text) TO samurai;
GRANT ALL PRIVILEGES ON PROCEDURE insert_order(text,text,text,text,text,text,text,integer,timestamp(0),text) TO samurai;
GRANT ALL PRIVILEGES ON PROCEDURE insert_orders_delivery(text,text,text,text,text,text,text) TO samurai;
GRANT ALL PRIVILEGES ON PROCEDURE insert_orders_item(integer,text,integer,text,text,
    integer,text,integer,integer,text,integer,text) TO samurai;
GRANT ALL PRIVILEGES ON PROCEDURE insert_orders_payment(text,text,text,text,integer,
    integer,text,integer,integer,integer,text) TO samurai;*/
GRANT ALL PRIVILEGES ON table orders TO samurai;
GRANT ALL PRIVILEGES ON table deliveries TO samurai;
GRANT ALL PRIVILEGES ON table payments TO samurai;
GRANT ALL PRIVILEGES ON table items TO samurai;
GRANT ALL PRIVILEGES ON table ordersitems TO samurai;
