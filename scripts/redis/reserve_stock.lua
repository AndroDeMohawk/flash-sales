local stockKey = KEYS[1]
local amount = tonumber(ARGV[1])

local current = redis.call('get', stockKey)
if not current then
    return -1 -- Ключ не найден в Redis
end

local currentStock = tonumber(current)
if currentStock < amount then
    return -2 -- Недостаточно билетов (Sold Out)
end

-- decrby автоматически возвращает НОВЫЙ остаток после списания (например, 99, 50, или 0)
local remaining = redis.call('decrby', stockKey, amount)
return remaining