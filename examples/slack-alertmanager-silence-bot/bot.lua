local http = require("http")
local inspect = require("inspect")
local json = require("json")
local time = require("time")
local cmd = require("cmd")
local strings = require("strings")
local filepath = require("filepath")

package.path = filepath.dir(debug.getinfo(1).source) .. './?.lua;' .. package.path
functions = require("functions")

log_level = os.getenv("GLUA_SCRIPT_LOG_LEVEL") or 'info'
local server, err = http.server("127.0.0.1:3001")
if err then error(err) end

-- main block program
while true do
  local req, resp = server:accept() -- lock and wait request
  
  local raw_body, err = req.body()
  if err then error(err) end
  resp:code(200)
  resp:done()
  
  if req.method == 'POST' then
    
    body, err = json.decode(urldecode(raw_body:gsub("payload=", "")))
    if err then error(err) end
    
    actions = body['actions']
    response_url = body['response_url']
    user = body['user']['name']
    
    original_message = body['original_message']
    original_attach = original_message['attachments']
    original_fields = original_attach[1]['fields']
    
    if log_level == 'debug' then
      print("[DEBUG] RAW BODY: \n"..inspect(body))
      print("[DEBUG] ACTIONS: \n"..inspect(actions))
    end
    
    for _, v in pairs(actions) do --parce actions
      
      if v['value'] ~= "" then
        j, err = json.decode(v['value'])
        duration = j['duration']
        alertmanager_url = j['url']
        if err then error(err) end
        
        if j['labels'] then
          matchers = {}
          for n, v in pairs(j['labels']) do --making data for alertmanager
            for k, v in pairs(v) do
              matcher = {isRegex = false, name = k, value = v}
              table.insert(matchers, n, matcher)
            end
          end
          
        else
          print("[INFO] No labels in value container")
        end
        
        end_time = format_time_with_offset(duration)
        start_time = format_time_with_offset(0)
        silence = {comment = "Silenced by Slack!", createdBy = user, startsAt = start_time, endsAt = end_time, matchers = matchers}
        data, err = json.encode(silence)
        if err then error(err) end
        
        client = http.client()
        url = alertmanager_url.."/api/v2/silences" -- Bot uses api v2
        request, err = http.request("POST", url, data)
        request:header_set("Content-Type", "application/json")
        
        if err then error(err) end
        
        local result, err = client:do_request(request)
        if err then
          error(err)
        elseif result.code ~= 200 then
          error(result.body)
        else
          -- We need to create new message for slack with changed field

          silenced_text = ":thumbsup: successfully silenced by "..user.." for "..duration.." hours (ends "..end_time..")"
          for _, v in pairs(original_fields) do
            if v['title'] == 'Status' then
              v['value'] = silenced_text
            end
          end
          
          new_actions = {}
          
          for i, v in pairs(original_attach[1]['actions']) do
            if v['name'] and v['name'] == 'silence' then
              print('[INFO] Skiping action SILENCE because we have a new one')
            else
              table.insert(new_actions, v)
            end
          end
          new_attach = original_attach
          new_attach[1]['actions'] = new_actions
          
          new_message = {}
          new_message['text'] = ""
          new_message['attachments'] = new_attach
          new_message['replace_original'] = true
          new_message['response_type'] = "in_channel"
          
          if log_level == 'debug' then
            print("[DEBUG] NEW CREATED MESSAGE:\n"..inspect(new_message))
          end
          
          message, err = json.encode(new_message)
          if err then error(err) end
          
          request, err = http.request("POST", response_url, message)
          local result, err = client:do_request(request)
          
          if err then error(err) end
          print("[DEBUG] Slack responce: \n"..result.body)
        end
      else
        print("[INFO] Skipping action because we do not value in action")
      end
    end
  end
end
