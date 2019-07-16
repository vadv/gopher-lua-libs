if os.getenv("TRAVIS") then
  -- travis
else
  dofile("./test/test_api_auth.lua")
end
