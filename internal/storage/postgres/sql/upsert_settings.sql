WITH updated AS (
    UPDATE users.user_settings
    SET dark_mode = COALESCE($2::boolean, dark_mode),
        private_data = COALESCE($3::boolean, private_data)
    WHERE user_id = $1
    RETURNING dark_mode, private_data
), inserted AS (
    INSERT INTO users.user_settings (user_id, dark_mode, private_data)
    SELECT $1, COALESCE($2::boolean, false), COALESCE($3::boolean, false)
    WHERE NOT EXISTS (SELECT 1 FROM updated)
    RETURNING dark_mode, private_data
)
SELECT dark_mode, private_data FROM updated
UNION ALL
SELECT dark_mode, private_data FROM inserted;
