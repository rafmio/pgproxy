curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"query":"INSERT INTO team_schedule (team_name, clock_in_time, clock_out_time) VALUES ($1, $2, $3)", "params":["Random Team", "12:00:00", "17:00:00"]}' \
     http://0.0.0.0:9080/create