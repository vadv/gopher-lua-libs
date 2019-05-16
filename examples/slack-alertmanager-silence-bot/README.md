# Usage

Let's assume you have such configuration in alertmanager for slack.
```
    slack_configs:
      - channel: "<you_channel>"
        api_url: "<you_apr_url>"
        callback_id: 'silence' 
        title: "[ {{ .GroupLabels.alertname }}: {{ .CommonAnnotations.summary }} ]"
        text: '{{ range .Alerts }} {{ .Annotations.description  }} {{ "\n" }} {{ end }}'

        actions:
          - type: 'button'
            text: '{{ if eq .Status "firing" }} Silence for 24 hours {{ else }}{{ end }}'
            name: 'silence' # We used that value in bot
            style: '{{ if eq .Status "firing" }}danger{{ else }}good{{ end }}'
            value: '{{ if eq .Status "firing" }}{"labels": [ {{ range .GroupLabels.SortedPairs }} { "{{ .Name }}": "{{ .Value }}" }, {{ end }} { "alertname": "{{ .GroupLabels.alertname }}" } ], "duration": "24", "url": "{{ .ExternalURL }}" } {{end}}'

```

Alertmanager uses value as container for pass data to this bot. Field `duration` indicates silence time period in hours.

Start bot:
```
$ GLUA_LOG_SCRIPT_LEVEL=debug glua-libs ./bot.lua 
```
