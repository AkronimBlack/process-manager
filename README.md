## Process Manager

This is a POC project to develop a bit dumbed down version of a bpm engine that works with json instead of the 
standardized xml notation

# How to start

Development
```
    go run main.go server:start -f <<file_with_actions.json>>
```

Build and run
```
    go build . -o process-manager
    ./process-manager -f <<file_with_actions.json>>
```

Docker 
```
    docker compose up -d
```