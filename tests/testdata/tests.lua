local goos = require 'goos'

function TestTempDir(t)
    -- Ensures that a tempdir created in subtests doesn't exist any longer after the test is run
    local tempDir = ''

    t:Run('createTmpDir', function(t)
        tempDir = t:TempDir()
        stat = goos.stat(tempDir)
        assert(stat)
        assert(stat.is_dir)
    end)
    assert(tempDir ~= '', tempDir)
    assert(not goos.stat(tempDir))
end