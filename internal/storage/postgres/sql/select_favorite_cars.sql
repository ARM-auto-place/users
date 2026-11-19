SELECT c.id, c.mark, c.model, c.build_year, p.price, c.kilometre,
       c.condition::text, c.color, c.images -> 0, links.created_at
FROM cars.favorite_cars AS links
JOIN cars.car AS c ON c.id = links.car_id
LEFT JOIN cars.car_price AS p ON p.car_id = c.id
WHERE links.user_id = $1 AND c.release = true
ORDER BY links.created_at DESC, c.id DESC;
