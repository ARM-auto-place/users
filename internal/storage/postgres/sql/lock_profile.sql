SELECT id FROM users."user"
WHERE id = $1
FOR NO KEY UPDATE;
