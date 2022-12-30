package storage

const sqlDbSchema = `
CREATE TABLE IF NOT EXISTS orders (
	order_uid			varchar PRIMARY KEY,
	track_number		varchar,
	entry				varchar,
	
	name				varchar,
	phone				varchar,
	zip					varchar,
	city				varchar,
	address				varchar,
	region				varchar,
	email				varchar,
	
	transaction			varchar,
	request_id			varchar,
	currency			varchar,
	provider			varchar,
	amount				bigint,
	payment_dt			bigint,
	bank				varchar,
	delivery_cost		bigint,
	goods_total			bigint,
	custom_fee			bigint,
	
	locale				varchar,
	internal_signature	varchar,
	customer_id			varchar,
	delivery_service	varchar,
	shardkey			varchar,
	sm_id				bigint,
	date_created		varchar,
	oof_shard			varchar
);

CREATE TABLE IF NOT EXISTS items (
	order_uid			varchar NOT NULL REFERENCES orders ON DELETE CASCADE,
	item_id				serial,
	chrt_id				bigint,
	track_number		varchar,
	price				bigint,
	rid					varchar,
	name				varchar,
	sale				bigint,
	size				varchar,
	total_price			bigint,
	nm_id				bigint,
	brand				varchar,
	status				bigint,
	PRIMARY KEY (order_uid, item_id)
);
`
const sqlOrderInsert = `INSERT INTO orders (
	order_uid, track_number, entry, 
	name, phone, zip, city, address, region, email, 
	transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee, 
	locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
) VALUES (
	:order_uid, :track_number, :entry, 
	:name, :phone, :zip, :city, :address, :region, :email, 
	:transaction, :request_id, :currency, :provider, :amount, :payment_dt, :bank, :delivery_cost, :goods_total, :custom_fee, 
	:locale, :internal_signature, :customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard
)
`
const sqlItemInsert = `INSERT INTO items(
	order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
) VALUES (
	:order_uid, :chrt_id, :track_number, :price, :rid, :name, :sale, :size, :total_price, :nm_id, :brand, :status
)
`
const sqlOrderSelect = "SELECT * FROM orders WHERE order_uid = $1"

const sqlItemSelect = "SELECT * FROM items WHERE order_uid = $1"

const sqlOrderListSelect = "SELECT order_uid FROM orders"
