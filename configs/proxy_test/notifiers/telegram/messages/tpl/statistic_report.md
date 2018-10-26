Date: {{.Date}}

_Authorization statistics_

Amount of token requests {{.AuthStats.AmountOfExecutions}}
Amount of errors {{.AuthStats.AmountOfErrors}}
Amount of timeouts {{.AuthStats.AmountOfTimeouts}}
Average response time {{.AuthStats.AverageResponseTime}}

_Tests statistics_
{{with .Tests}}
{{range $k, $v := .}}
        
{{$k.Definition.ID}} {{$k.Definition.HTTPMethod}}
[{{$k.Definition.URL}}]
Timeouts: {{$k.Definition.TimeOut}}
Run period: {{$k.Definition.RunPeriod}}
Amount of executions: {{$v.AmountOfExecutions}}
Amount of errors: {{$v.AmountOfErrors}}
Average response time: {{$v.AverageResponseTime}}

---------------------------------------------------

{{end}}
{{end}}
