tab1 = { key1 = "val1", key2 = "val2", "val3" }
for k, v in pairs(tab1) do
    print(k .. " - " .. v)
end
print("tab1 len:", #tab1)
print("=============================")

tab1.key1 = nil
for k, v in pairs(tab1) do
    print(k .. " - " .. v)
end
print("tab1 len:", #tab1)
print("=============================")

tab2 = {"v1", "v2", "v3", nil, "v4"}
for k, v in ipairs(tab2) do
    print(k .. " - " .. v)
end
print("----")
for k, v in pairs(tab2) do
    print(k .. " - " .. v)
end
print("tab2 len:", #tab2)
print("=============================")