if os.getenv("CI") then
    -- travis
    function TestCI(t)
        t:Skip("CI")
    end
else
    dofile("./test/test_api_restream.lua")
end
