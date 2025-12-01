SELECT
    n.nspname AS type_schema,
    t.typname AS type_name,
    e.enumlabel AS enum_value,
    e.enumsortorder AS sort_order,
    pg_catalog.format_type(t.oid, NULL) AS type_definition
FROM
    pg_catalog.pg_type t
    JOIN pg_catalog.pg_namespace n ON n.oid = t.typnamespace
    JOIN pg_catalog.pg_enum e ON e.enumtypid = t.oid
WHERE
    t.typtype = 'e'  -- Enum types
    AND n.nspname NOT IN ('pg_catalog', 'information_schema', 'pg_toast')
ORDER BY
    type_schema,
    type_name,
    e.enumsortorder;