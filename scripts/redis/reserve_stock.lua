-- KEYS[1]: Ключ остатка билетов, например "ticket:101:stock"
-- ARGV[1]: Количество запрашиваемых билетов (например, 1)

local stock_key = KEYS[1]
local requested = tonumber(ARGV[1])

-- 1. Получаем текущее количество билетов
local current_stock = tonumber(redis.call('GET', stock_key))

-- Если ключа нет в Redis (билет не инициализирован)
if not current_stock then
    return -1
end

-- 2. Проверяем, хватает ли билетов
if current_stock >= requested then
    -- Уменьшаем остаток
    local new_stock = redis.call('DECRBY', stock_key, requested)
    return new_stock -- Возвращаем оставшееся количество (>= 0)
else
    return -2 -- Билеты закончились (Sold Out)
end