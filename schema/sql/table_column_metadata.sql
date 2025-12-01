
-- Go Comment: 
-- This Gets the metadatas from a given schema
-- This Query Collects 2 pararmeters 
-- Note: this query does not check if 
-- its system table and weather if its a Temp Table also 
-- Schema Name: string - $1
-- Table Name: string - $2  
-- Returns :

-- Go Variable: TableMetaData
SELECT 
    n.nspname AS table_schema,
    c.relname AS table_name,
    CASE c.relkind
        WHEN 'r' THEN 'BASE TABLE'
        WHEN 'v' THEN 'VIEW'
        WHEN 'm' THEN 'MATERIALIZED VIEW'
        WHEN 'f' THEN 'FOREIGN TABLE'
        WHEN 'p' THEN 'PARTITIONED TABLE'
        ELSE 'OTHER'
    END AS table_type,
    a.attname AS column_name,
    a.attnum AS ordinal_position,
    pg_catalog.format_type(a.atttypid, a.atttypmod) AS data_type,
    CASE 
        WHEN a.attnotnull THEN FALSE 
        ELSE TRUE 
    END AS is_nullable,
    -- CASE 
    --     WHEN pk.column_name IS NOT NULL THEN TRUE
    --     ELSE FALSE
    -- END AS is_primary_key,
    -- CASE 
    --     WHEN uq.column_name IS NOT NULL THEN TRUE
    --     ELSE FALSE
    -- END AS is_unique,
    -- CASE 
    --     WHEN fk.column_name IS NOT NULL THEN TRUE
    --     ELSE FALSE
    -- END AS is_foreign_key,
    COALESCE(pg_catalog.pg_get_expr(ad.adbin, ad.adrelid), '') AS column_default,
    col_description(c.oid, a.attnum) AS column_comment,
    -- CASE 
    --     WHEN n.nspname LIKE 'pg_temp%' THEN TRUE
    --     ELSE FALSE
    -- END AS is_temporary,
    --  check constraints using aggregation
    -- (
    --     SELECT string_agg(pg_catalog.pg_get_constraintdef(con.oid, true), ' | ')
    --     FROM pg_catalog.pg_constraint con
    --     WHERE con.conrelid = c.oid
    --       AND con.contype = 'c'
    --       AND a.attnum = ANY(con.conkey)
    -- ) AS check_constraint_definitions,

    -- Character maximum length for string types
    CASE 
        WHEN a.atttypid IN (1042, 1043) THEN -- char, varchar
            CASE 
                WHEN a.atttypmod = -1 THEN NULL
                ELSE a.atttypmod - 4
            END
        ELSE NULL
    END AS character_maximum_length,
    -- Numeric precision and scale
    CASE 
        WHEN a.atttypid IN (1700) THEN -- numeric
            ((a.atttypmod - 4) >> 16) & 65535
        ELSE NULL
    END AS numeric_precision,
    CASE 
        WHEN a.atttypid IN (1700) THEN -- numeric
            (a.atttypmod - 4) & 65535
        ELSE NULL
    END AS numeric_scale

FROM pg_catalog.pg_class c
JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
JOIN pg_catalog.pg_attribute a ON a.attrelid = c.oid
LEFT JOIN pg_catalog.pg_attrdef ad ON (a.attrelid = ad.adrelid AND a.attnum = ad.adnum)
-- primary key information
LEFT JOIN (
    SELECT 
        kcu.table_schema,
        kcu.table_name, 
        kcu.column_name
    FROM information_schema.key_column_usage kcu
    JOIN information_schema.table_constraints tc 
        ON kcu.constraint_name = tc.constraint_name 
        AND kcu.table_schema = tc.table_schema
    WHERE tc.constraint_type = 'PRIMARY KEY'
) pk ON pk.table_schema = n.nspname 
    AND pk.table_name = c.relname 
    AND pk.column_name = a.attname
-- Unique constraint information
LEFT JOIN (
    SELECT 
        kcu.table_schema,
        kcu.table_name, 
        kcu.column_name
    FROM information_schema.key_column_usage kcu
    JOIN information_schema.table_constraints tc 
        ON kcu.constraint_name = tc.constraint_name 
        AND kcu.table_schema = tc.table_schema
    WHERE tc.constraint_type = 'UNIQUE'
) uq ON uq.table_schema = n.nspname 
    AND uq.table_name = c.relname 
    AND uq.column_name = a.attname
-- foreign key information
LEFT JOIN (
    SELECT 
        kcu.table_schema,
        kcu.table_name, 
        kcu.column_name
    FROM information_schema.key_column_usage kcu
    JOIN information_schema.table_constraints tc 
        ON kcu.constraint_name = tc.constraint_name 
        AND kcu.table_schema = tc.table_schema
    WHERE tc.constraint_type = 'FOREIGN KEY'
) fk ON fk.table_schema = n.nspname 
    AND fk.table_name = c.relname 
    AND fk.column_name = a.attname

WHERE c.relkind IN ('r', 'p', 'f') -- Include regular, partitioned, and foreign tables
  AND n.nspname NOT IN ('information_schema', 'pg_catalog', 'pg_toast')
  AND n.nspname NOT LIKE 'pg_temp%'
  AND a.attnum > 0
  AND NOT a.attisdropped
  AND c.relpersistence IN ('p', 'u') -- 'p' = permanent, 'u' = unlogged
--   AND n.nspname = 'public' -- $1
--   AND c.relname = 'example' -- $2

ORDER BY n.nspname, c.relname, a.attnum;