WITH patch AS (SELECT $2::jsonb AS fields),
updated AS (
    UPDATE users."user" AS u
    SET name        = COALESCE(p.fields ->> 'name', u.name),
        middle_name = CASE WHEN p.fields ? 'middleName' THEN p.fields ->> 'middleName' ELSE u.middle_name END,
        sure_name   = CASE WHEN p.fields ? 'surname' THEN p.fields ->> 'surname' ELSE u.sure_name END,
        phone       = COALESCE(p.fields ->> 'phone', u.phone),
        photo_path  = CASE WHEN p.fields ? 'photoPath' THEN p.fields ->> 'photoPath' ELSE u.photo_path END,
        age         = CASE WHEN p.fields ? 'age' THEN (p.fields ->> 'age')::integer ELSE u.age END,
        city        = CASE WHEN p.fields ? 'city' THEN p.fields ->> 'city' ELSE u.city END,
        address     = CASE WHEN p.fields ? 'address' THEN p.fields ->> 'address' ELSE u.address END
    FROM patch AS p
    WHERE u.id = $1
    RETURNING u.id, u.email, u.name, u.middle_name, u.sure_name, u.phone, u.photo_path, u.age, u.city, u.address
)
SELECT u.id, u.email, u.name, u.middle_name, u.sure_name, u.phone,
       u.photo_path, u.age, u.city, u.address,
       COALESCE(s.dark_mode, false), COALESCE(s.private_data, false)
FROM updated AS u
LEFT JOIN users.user_settings AS s ON s.user_id = u.id;
