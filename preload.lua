function TestRequireModule(t)
    modules = {
        "argparse",
        "base64",
        "cert_util",
        "chef",
        "cloudwatch",
        "cmd",
        "crypto",
        "db",
        "filepath",
        "goos",
        "humanize",
        "inspect",
        "ioutil",
        "json",
        "log",
        "pb",
        "plugin",
        "pprof",
        "prometheus",
        "regexp",
        "runtime",
        "shellescape",
        "stats",
        "storage",
        "strings",
        "tac",
        "tcp",
        "telegram",
        "template",
        "time",
        "xmlpath",
        "yaml",
        "zabbix",
    }
    for _, module in ipairs(modules) do
        t:Run(module, function(t)
            require(module)
        end)
    end
end

