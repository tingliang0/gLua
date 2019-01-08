function permgen(a, n)
    n = n or #a
    if n <= 1 then
        coroutine.yield(a)
    else
        for i=1, n do
            a[n],  a[i] = a[i], a[n]
            permgen(a, n - 1)
            a[n], a[i] = a[i], a[n]
        end
    end
end

function permutations(a)
    local co = coroutine.create(function()
            permgen(a)
    end)
    return function()
        local _, res = coroutine.resume(co)
        return res
    end
end

for p in permutations{"a", "b", "c"} do
    print(table.concat(p, ","))
end
