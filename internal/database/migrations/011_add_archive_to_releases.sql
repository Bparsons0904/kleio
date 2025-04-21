SELECT CASE 
    WHEN COUNT(*) = 0 THEN
        'ALTER TABLE releases ADD COLUMN archive BOOLEAN NOT NULL DEFAULT FALSE;'
    ELSE
        'SELECT 1;' 
END AS sql_to_execute
FROM pragma_table_info('releases') 
WHERE name = 'archive';
