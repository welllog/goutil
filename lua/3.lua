Counter = {value = 0}

function Counter:new()
    local o = {}
    setmetatable(o, self)
    self.__index = self
    return o
end

function Counter:count()
    return self.value
end

function Counter:incr(v)
    self.value = self.value + (v or 0)
end

local c1 = Counter:new()
c1:incr(1)
print(c1:count())
c1:incr(2)
print(c1:count())

local c2 = Counter:new()
print(c2:count())
c2:incr(1)
print(c2:count())