var reserveStockLua = `
    local stockKey = KEYS[1]
    local amount = tonumber(ARGV[1])

    local current = redis.call('get', stockKey)
    if not current then
        return -1
    end

    local currentStock = tonumber(current)
    if currentStock < amount then
        return -2
    end

    local remaining = redis.call('decrby', stockKey, amount)
    return remaining
`