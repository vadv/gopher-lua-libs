strings = require("strings")
cmd = require("cmd")

hex_to_char = function(x)
  return string.char(tonumber(x, 16))
end

urldecode = function(url)
  if url == nil then
    return
  end
  url = url:gsub("+", " ")
  url = url:gsub("%%(%x%x)", hex_to_char)
  return url
end

format_time_with_offset = function(i)
  command = ("date --iso-8601=seconds --date '+"..i.." hour'")
  encoded, err = cmd.exec(command)
  if err then error(err) end
  return strings.trim((encoded.stdout), "\n")
end
