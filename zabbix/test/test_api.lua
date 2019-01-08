if os.getenv("TRAVIS") then
  -- travis
else
  dofile("./test/test_api_restream.lua")
end
