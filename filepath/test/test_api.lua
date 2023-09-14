local filepath = require('filepath')
local assert = require('assert')
local require = require('require')
local inspect = require('inspect')

function Test_join_and_separator(t)
    local path = "1"
    local need_path = path .. filepath.separator() .. "2" .. filepath.separator() .. "3"
    path = filepath.join(path, "2", "3")
    assert(path == need_path, string.format("expected %s; got %s", need_path, path))
end

function Test_join(t)
    tests = {
        {
            name = "1/2/3",
            args = { "1", "2", "3" },
        },
        {
            name = "foo",
            args = { "foo" },
        },
        {
            name = "",
            args = {},
        },
        {
            name = "/a/b/c",
            args = { "/", "a", "b", "c" },
        },
    }
    for _, tt in ipairs(tests) do
        t:Run(tt.name, function(t)
            assert:Equal(t, tt.name, filepath.join(unpack(tt.args)))
        end)
    end
end

function Test_glob(t)
    local results = filepath.glob("test" .. filepath.separator() .. "*")
    assert(#results == 1, string.format("expected one glob result; got %d", #results))
    for k, v in pairs(results) do
        if k == 1 then
            assert(v == "test" .. filepath.separator() .. "test_api.lua", v)
        end
    end
end

function Test_abs(t)
    pwd = os.getenv('PWD')
    path, err = filepath.abs('foo')
    require:NoError(t, err)
    assert:Equal(t, filepath.join(pwd, 'foo'), path)
end

function Test_clean(t)
    cleaned = filepath.clean('/foo//bar')
    assert:Equal(t, '/foo/bar', cleaned)
end

function Test_eval_symlinks(t)
    temp_dir = t:TempDir()
    temp_dir, err = filepath.eval_symlinks(temp_dir)
    require:NoError(t, err)

    test_path = filepath.join(temp_dir, 'foo')
    os.execute('ln -s . ' .. test_path)
    path, err = filepath.eval_symlinks(test_path)
    require:NoError(t, err)
    assert:Equal(t, temp_dir, path)
end

function Test_from_slash(t)
    t:Skip('TODO')
end

function Test_ext(t)
    ext = filepath.ext('foo.bar')
    assert:Equal(t, '.bar', ext)
end

function Test_is_abs(t)
    tests = {
        {
            name = 'is absolute',
            path = '/foo/bar',
            expected = true
        },
        {
            name = 'is not absolute',
            path = 'foo/bar',
            expected = false
        },
    }

    for _, tt in ipairs(tests) do
        t:Run(tt.name, function(t)
            assert:Equal(t, tt.expected, filepath.is_abs(tt.path))
        end)
    end
end

function Test_match(t)
    tests = {
        {
            name = 'Should match',
            path = '/foo/bar',
            pattern = '/*/bar',
            expected = true,
        },
        {
            name = 'Should NOT match',
            path = '/yada/yada',
            pattern = '/*/bar',
            expected = false,
        },
    }
    for _, tt in ipairs(tests) do
        t:Run(tt.name, function(t)
            assert:Equal(t, tt.expected, filepath.match(tt.pattern, tt.path))
        end)
    end
end

function Test_rel(t)
    tests = {
        {
            name = 'targpath is under basepath',
            basepath = '/foo/bar',
            targpath = '/foo/bar/baz',
            expected = 'baz',
            want_err = false,
        },
        {
            name = 'targpath is NOT under basepath',
            basepath = '/foo/bar',
            targpath = 'yada/yada/baz',
            want_err = true,
        },
    }
    for _, tt in ipairs(tests) do
        t:Run(tt.name, function(t)
            path, err = filepath.rel(tt.basepath, tt.targpath)
            if tt.want_err then
                require:Errorf(t, err, "basepath = %s, targpath = %s, path = %s", tt.basepath, tt.targpath, path)
                return
            end
            require:NoError(t, err)
            assert:Equal(t, tt.expected, path)
        end)
    end
end

function Test_split(t)
    tests = {
        {
            name = '/foo/bar/baz',
            expected_dir = '/foo/bar/',
            expected_file = 'baz',
        },
        {
            name = 'bar/baz',
            expected_dir = 'bar/',
            expected_file = 'baz',
        },
        {
            name = 'nodir',
            expected_dir = '',
            expected_file = 'nodir',
        },
    }
    for _, tt in ipairs(tests) do
        t:Run(tt.name, function(t)
            dir, file = filepath.split(tt.name)
            assert:Equal(t, tt.expected_dir, dir)
            assert:Equal(t, tt.expected_file, file)
        end)
    end
end

function Test_split_list(t)
    tests = {
        {
            name = 'foo/bar/baz:/yada/yada',
            expected = { 'foo/bar/baz', '/yada/yada' },
        },
        {
            name = 'a/b/c',
            expected = { 'a/b/c' },
        },
    }
    for _, tt in ipairs(tests) do
        t:Run(tt.name, function(t)
            assert:Equal(t, inspect(tt.expected), inspect(filepath.split_list(tt.name)))
        end)
    end
end

function Test_to_slash(t)
    t:Skip('TODO')
end

function Test_volume_name(t)
    t:Skip('TODO')
end
