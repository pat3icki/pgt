SELECT
    n.nspname AS type_schema,
    t.typname AS type_name,
    a.attname AS attribute_name,
    pg_catalog.format_type(a.atttypid, a.atttypmod) AS attribute_type,
    a.attnum AS attribute_position,
    NOT a.attnotnull AS is_nullable,
    pg_catalog.col_description(t.oid, a.attnum) AS attribute_comment
FROM
    pg_catalog.pg_type t
    JOIN pg_catalog.pg_namespace n ON n.oid = t.typnamespace
    JOIN pg_catalog.pg_class c ON c.oid = t.typrelid
    JOIN pg_catalog.pg_attribute a ON a.attrelid = c.oid AND a.attnum > 0 AND NOT a.attisdropped
WHERE
    t.typtype = 'c'
    AND n.nspname NOT IN ('pg_catalog', 'information_schema', 'pg_toast')
    AND NOT EXISTS (
        SELECT 1 
        FROM pg_catalog.pg_class tbl 
        WHERE tbl.relname = t.typname 
          AND tbl.relnamespace = n.oid 
          AND tbl.relkind = 'r'
    )
ORDER BY
    type_schema,
    type_name,
    a.attnum;
