SELECT u.id, u.email, u.name, u.middle_name, u.sure_name, u.phone,
       u.photo_path, u.age, u.city, u.address,
       COALESCE(s.dark_mode, false), COALESCE(s.private_data, false)
FROM users."user" AS u
LEFT JOIN users.user_settings AS s ON s.user_id = u.id
WHERE u.id = $1;
