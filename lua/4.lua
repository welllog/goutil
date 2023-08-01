local http = require("http")

print(http.name)

resp, err = http.get("https://www.baidu.com")
if resp then
    print(resp.status)
    print(resp.body)
    for k, v in pairs(resp.headers) do
        print(k .. " - " .. v)
    end
else
    print(err)
end