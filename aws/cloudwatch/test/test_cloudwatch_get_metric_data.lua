local cloudwatch = require("cloudwatch")
local inspect = require("inspect")

local clw_client, err = cloudwatch.new()
if err then error(err) end

local query1 = {
  namespace = "AWS/RDS",
  metric = "CPUUtilization",
  dimension_name = "DBInstanceIdentifier",
  dimension_value = os.getenv("DBINSTANCE"),
  stat = "Average",
  period = 60,
}
local query2 = {
    namespace = "AWS/RDS",
    metric = "ReadIOPS",
    dimension_name = "DBInstanceIdentifier",
    dimension_value = os.getenv("DBINSTANCE"),
    stat = "Average",
    period = 60,
}
local result, err = clw_client:get_metric_data({queries={cpu=query1, iops=query2}})
if err then error(err) end
print(inspect(result))
