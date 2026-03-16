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

-- name: DeleteSubscription :execrows
DELETE FROM subscriptions
WHERE id = $1;

-- name: ListSubscriptions :many
SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
FROM subscriptions
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: TotalSubscriptions :one
SELECT COUNT(*)::int AS total
FROM subscriptions;

-- name: SubscriptionsTotalCost :one
SELECT COALESCE(SUM(
    price * 
    CASE 
        WHEN effective_end < effective_start THEN 0
        ELSE (
            (EXTRACT(YEAR FROM effective_end) - EXTRACT(YEAR FROM effective_start)) * 12 +
            (EXTRACT(MONTH FROM effective_end) - EXTRACT(MONTH FROM effective_start)) +
            1
        )
    END
), 0)::bigint AS total_cost
FROM (
    SELECT 
        price,
        -- Логика старта: если фильтр задан, берем максимум, иначе старт подписки
        CASE 
            WHEN sqlc.narg(start_date)::date IS NOT NULL 
                THEN GREATEST(start_date, sqlc.narg(start_date)::date)
            ELSE start_date
        END AS effective_start,

        -- Логика конца:
        -- 1. Определяем реальный конец подписки (если NULL, то считаем до сегодня или бесконечно?)
        -- 2. Определяем границу фильтра (если NULL, то берем реальный конец подписки)
        LEAST(
            COALESCE(end_date, CURRENT_DATE), -- Конец подписки (или сегодня)
            COALESCE(sqlc.narg(end_date)::date, COALESCE(end_date, CURRENT_DATE)) -- Конец фильтра (или конец подписки/сегодня)
        ) AS effective_end

    FROM subscriptions
    WHERE 
        (sqlc.narg(user_id)::uuid IS NULL OR user_id = sqlc.narg(user_id))
        AND (sqlc.narg(service_name)::text IS NULL OR service_name = sqlc.narg(service_name))
        -- Фильтры работают только если переданы соответствующие даты
        AND (sqlc.narg(end_date)::date IS NULL OR start_date <= sqlc.narg(end_date)::date)
        AND (sqlc.narg(start_date)::date IS NULL OR end_date IS NULL OR end_date >= sqlc.narg(start_date)::date)
) AS calculated_periods;
