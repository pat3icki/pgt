



SELECT
    n.nspname AS schema_name,
    t.relname AS table_name,
    c.relname AS index_name,
    array_agg(a.attname ORDER BY k.ord) AS indexed_columns,
    pg_get_indexdef(i.indexrelid) AS index_definition,
    i.indisunique AS is_unique,
    i.indisprimary AS is_primary_key,
    i.indisexclusion AS is_exclusion,
    i.indimmediate AS is_immediate,
    i.indisclustered AS is_clustered,
    am.amname AS index_method,
    pg_relation_size(c.oid) AS index_size_bytes
FROM pg_index i
JOIN pg_class c ON c.oid = i.indexrelid
JOIN pg_class t ON t.oid = i.indrelid
JOIN pg_namespace n ON n.oid = t.relnamespace
JOIN pg_am am ON c.relam = am.oid
JOIN LATERAL unnest(i.indkey) WITH ORDINALITY AS k(attnum, ord) ON TRUE
JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = k.attnum
WHERE
   c.relkind = 'i'  AND
   i.indisvalid = true 
--    AND
    --  n.nspname = 'public' AND -- $1
    -- t.relname = 'example' -- $2

GROUP BY 
    n.nspname, t.relname, c.relname, i.indexrelid, 
    i.indisunique, i.indisprimary, i.indisexclusion, 
    i.indimmediate, i.indisvalid, i.indisready, 
    i.indisclustered, am.amname, c.oid
ORDER BY table_name, index_name;