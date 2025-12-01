CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" VARCHAR(20),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "country_code" VARCHAR(6),
  "country" VARCHAR(100)
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigserial,
  "amount" bigint NOT NULL,
  "description" VARCHAR(250),
  "payment_type" VARCHAR(200),
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigserial,
  "to_account_id" bigserial,
  "amount" bigint NOT NULL,
  "description" VARCHAR(250),
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "payout_transactions" (
  "id" bigint PRIMARY KEY,
  "reference" varchar(50),
  "account_id" bigint,
  "recipient_number" bigint,
  "recipient_name" VARCHAR(250),
  "recipient_bank_code" int,
  "amount" bigint,
  "status" smallint,
  "description" VARCHAR(200),
  "created_at" timestamptz,
  "paid_at" timestamptz,
  "merchant_name" varchar(50),
  "merchant_fee" smallint,
  "customer_fee" smallint
);

CREATE TABLE "fund_transactions" (
  "id" bigint PRIMARY KEY,
  "reference" VARCHAR(200),
  "account_id" bigint,
  "amount" bigint,
  "status" smallint,
  "target_id" varchar(50),
  "target_type" varchar(50),
  "currency" smallint,
  "payment_type" smallint,
  "ip_address" VARCHAR(200),
  "created_at" timestamptz,
  "paid_at" timestamptz,
  "merchant_name" varchar(50),
  "merchant_fee" smallint,
  "customer_fee" smallint
);

CREATE TABLE "cards_authorization" (
  "id" bigint PRIMARY KEY,
  "owner_id" bigint,
  "authorization_code" varchar(250),
  "card_type" varchar(50),
  "exp_mounth" smallint,
  "exp_year" smallint,
  "last4" int
);

CREATE TABLE "safeboxes" (
  "id" bigint PRIMARY KEY,
  "owner_id" bigint,
  "amount" bigint,
  "withdraw_date" timestamptz
);

CREATE TABLE "safebox_transactions" (
  "id" bigint PRIMARY KEY,
  "safebox_id" bigint,
  "amount" bigint,
  "trans_type" varchar(20),
  "inital_balance" bigint,
  "finial_balance" bigint
);

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "payout_transactions" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "fund_transactions" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "cards_authorization" ADD FOREIGN KEY ("owner_id") REFERENCES "accounts" ("id");

ALTER TABLE "safeboxes" ADD FOREIGN KEY ("owner_id") REFERENCES "accounts" ("id");

ALTER TABLE "safebox_transactions" ADD FOREIGN KEY ("safebox_id") REFERENCES "safeboxes" ("id");

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "accounts" ("id");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

CREATE INDEX ON "payout_transactions" ("id");

CREATE INDEX ON "payout_transactions" ("account_id");

CREATE INDEX ON "payout_transactions" ("reference");

CREATE INDEX ON "fund_transactions" ("id");

CREATE INDEX ON "fund_transactions" ("account_id");

CREATE INDEX ON "fund_transactions" ("reference");

CREATE INDEX ON "cards_authorization" ("id");

CREATE INDEX ON "cards_authorization" ("owner_id");

CREATE INDEX ON "safeboxes" ("id");

CREATE INDEX ON "safeboxes" ("owner_id");

CREATE INDEX ON "safebox_transactions" ("id");

CREATE INDEX ON "safebox_transactions" ("safebox_id");

CREATE INDEX ON "safebox_transactions" ("id", "safebox_id");

COMMENT ON COLUMN "entries"."amount" IS 'can  be positive or negative';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';


-- First, define a custom ENUM type for priority, required by the main table
CREATE TYPE task_priority AS ENUM ('Low', 'Medium', 'High', 'Urgent');

-- Create the single, complex 'ProjectTasks' table
CREATE TABLE ProjectTasks (
    -- Primary Key: Modern auto-incrementing ID
    task_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    -- Standard Text Fields
    title VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Date & Time Fields with Defaults and Constraints
    created_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    due_date DATE,
    
    -- Custom ENUM Type with Default
    priority task_priority NOT NULL DEFAULT 'Medium'::task_priority,

    -- Numeric Types
    estimated_hours NUMERIC(5, 2) CHECK (estimated_hours > 0),
    progress_percentage SMALLINT NOT NULL DEFAULT 0 CHECK (progress_percentage BETWEEN 0 AND 100),

    -- Boolean Type with Default
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,

    -- JSONB for unstructured data (e.g., tags, sub-tasks)
    metadata JSONB,

    -- Email Field with Format Validation using a CHECK constraint and regex
    assigned_email VARCHAR(255) CHECK (assigned_email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,4}$'),

    -- Table-level CHECK constraint: ensures due date is after creation date
    CONSTRAINT due_date_after_created_date CHECK (due_date IS NULL OR due_date >= DATE(created_date)),

    -- Table-level UNIQUE constraint: prevents duplicate titles for the same assigned user (email)
    CONSTRAINT unique_task_for_user UNIQUE (title, assigned_email)
);

-- Add an index for faster lookups by assignee's email
CREATE INDEX idx_projecttasks_email ON ProjectTasks(assigned_email);



CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    department_id INTEGER,
    salary DECIMAL(10,2) CHECK (salary > 0),
    hire_date DATE DEFAULT CURRENT_DATE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);


CREATE TEMP TABLE session_logs (
    log_id SERIAL PRIMARY KEY,
    username VARCHAR(100),
    action TEXT,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


-- When you backup, constraint definitions include deferrability
CREATE TABLE example (
    id SERIAL PRIMARY KEY,
    email VARCHAR UNIQUE DEFERRABLE INITIALLY DEFERRED
);

-- Ensure no circles overlap within a certain area
CREATE TABLE circles (
    id serial PRIMARY KEY,
    position point,
    radius float,
    EXCLUDE USING gist (circle(position, radius) WITH &&)
);


CREATE DOMAIN valid_email AS text NOT NULL CHECK (value ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$');


CREATE TABLE mailing_list (
    id SERIAL PRIMARY KEY,
    email VALID_EMAIL
);



CREATE TYPE bug_status AS ENUM ('new', 'open', 'closed');

CREATE TABLE bug (
    id serial,
    description text,
    status bug_status
);

CREATE TYPE full_address AS (
    street text,
    city text,
    zipcode varchar(10)
);

-- Use the composite type in a table
CREATE TABLE company_branches (
    branch_id serial PRIMARY KEY,
    location full_address
);

CREATE TYPE address_type AS (
    street text,
    city text,
    zip_code integer,
    floor smallint
    country text
);


CREATE TABLE contacts (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address address_type NOT NULL
);


CREATE TYPE public.accounts AS (id bigint, owner character varying, balance bigint, currency character varying(20), created_at timestamp with time zone, country_code character varying(6), country character varying(100));