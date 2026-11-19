WITH updated AS (
    UPDATE users."user"
    SET name        = COALESCE($2::text, name),
        middle_name = CASE WHEN $3::boolean THEN $4::text ELSE middle_name END,
        sure_name   = CASE WHEN $5::boolean THEN $6::text ELSE sure_name END,
        phone       = COALESCE($7::text, phone),
        photo_path  = CASE WHEN $8::boolean THEN $9::text ELSE photo_path END,
        age         = CASE WHEN $10::boolean THEN $11::integer ELSE age END,
        city        = CASE WHEN $12::boolean THEN $13::text ELSE city END,
        address     = CASE WHEN $14::boolean THEN $15::text ELSE address END
    WHERE id = $1
    RETURNING id, email, name, middle_name, sure_name, phone, photo_path, age, city, address
)
SELECT u.id, u.email, u.name, u.middle_name, u.sure_name, u.phone,
       u.photo_path, u.age, u.city, u.address,
       COALESCE(s.dark_mode, false), COALESCE(s.private_data, false)
FROM updated AS u
LEFT JOIN users.user_settings AS s ON s.user_id = u.id;
