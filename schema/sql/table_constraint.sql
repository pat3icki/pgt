
SELECT
    n.nspname AS table_schema,
    c.relname AS table_name,
    con.conname AS constraint_name,
    CASE con.contype
        WHEN 'c' THEN 'CHECK'
        WHEN 'f' THEN 'FOREIGN KEY'
        WHEN 'p' THEN 'PRIMARY KEY'
        WHEN 'u' THEN 'UNIQUE'
        WHEN 't' THEN 'TRIGGER'
        WHEN 'x' THEN 'EXCLUSION'
        ELSE 'OTHER'
    END AS constraint_type,
    pg_catalog.pg_get_constraintdef(con.oid, true) AS constraint_definition,
    -- Get the column names for the constraint
    (
        SELECT string_agg(a.attname, ', ' ORDER BY u.ordinality)
        FROM unnest(con.conkey) WITH ORDINALITY AS u(attnum, ordinality)
        JOIN pg_catalog.pg_attribute a ON a.attrelid = con.conrelid AND a.attnum = u.attnum
    ) AS column_names,

    -- Additional foreign key specific columns
    CASE 
        WHEN con.contype = 'f' THEN 
            CASE con.confupdtype
                WHEN 'a' THEN 'NO ACTION'
                WHEN 'r' THEN 'RESTRICT'
                WHEN 'c' THEN 'CASCADE'
                WHEN 'n' THEN 'SET NULL'
                WHEN 'd' THEN 'SET DEFAULT'
                ELSE 'UNKNOWN'
            END
        ELSE NULL
    END AS update_rule, 
    CASE 
        WHEN con.contype = 'f' THEN 
            CASE con.confdeltype
                WHEN 'a' THEN 'NO ACTION'
                WHEN 'r' THEN 'RESTRICT'
                WHEN 'c' THEN 'CASCADE'
                WHEN 'n' THEN 'SET NULL'
                WHEN 'd' THEN 'SET DEFAULT'
                ELSE 'UNKNOWN'
            END
        ELSE NULL
    END AS delete_rule,


    con.condeferrable AS is_deferrable,
    con.condeferred AS initially_deferred,

        -- Foreign key specific information
    CASE WHEN con.contype = 'f' THEN
        (SELECT nspname FROM pg_catalog.pg_namespace WHERE oid = con.connamespace)
    END AS f_referenced_table_schema,

    CASE WHEN con.contype = 'f' THEN
        (SELECT relname FROM pg_catalog.pg_class WHERE oid = con.confrelid)
    END AS f_referenced_table_name,

    CASE WHEN con.contype = 'f' THEN
        (
            SELECT string_agg(a.attname, ', ' ORDER BY u.ordinality)
            FROM unnest(con.confkey) WITH ORDINALITY AS u(attnum, ordinality)
            JOIN pg_catalog.pg_attribute a ON a.attrelid = con.confrelid AND a.attnum = u.attnum
        )
    END AS f_referenced_columns



FROM
    pg_catalog.pg_constraint con
    INNER JOIN pg_catalog.pg_class c ON c.oid = con.conrelid
    INNER JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
WHERE
    n.nspname NOT IN ('information_schema', 'pg_catalog', 'pg_toast')
ORDER BY
    table_schema,
    table_name,
    constraint_type,
    constraint_name;



