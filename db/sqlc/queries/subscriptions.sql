-- name: GetSubscription :one
SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
FROM subscriptions
WHERE id = $1;

-- name: CreateSubscription :one
INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at;

-- name: UpdateSubscription :one
UPDATE subscriptions
SET
    service_name = COALESCE($2, service_name),
    price = COALESCE($3, price),
    start_date = COALESCE($4, start_date),
    end_date = COALESCE($5, end_date),
    updated_at = NOW()
WHERE id = $1
RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at;

-- name: DeleteSubscription :exec
DELETE FROM subscriptions
WHERE id = $1;

-- name: ListSubscriptions :many
SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
FROM subscriptions
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListSubscriptionsFiltered :many
SELECT *
FROM subscriptions
WHERE
    ($1::uuid IS NULL OR user_id = $1)
    AND ($2::text IS NULL OR service_name = $2)
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;
