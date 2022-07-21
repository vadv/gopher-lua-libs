local crypto = require("crypto")

function Test_crypto(t)
    t:Run("md5", function(t)
        assert(crypto.md5("1\n") == "b026324c6904b2a9cb4b88d6d61c81d1")
    end)

    t:Run("sha256", function(t)
        assert(crypto.sha256("1\n") == "4355a46b19d348dc2f57c046f8ef63d4538ebb936000f3c9ee954a27460dd865")
    end)
end
