# Heroku link
```
https://murmuring-scrubland-25884.herokuapp.com/
```

# Local Setup

## Pre-requisites
Create config.env at root directory
```
LIMIT=60
LIMIT_DURATION=60
REDIS_ENDPOINT=127.0.0.1:6379
REDIS_PASSWORD=
```

## Run Redis through docker-compose
docker compose version 3.7+
```
docker-compose up -d
```

# In-Memory Database
This Request limiter use Redis as its IMDB to store the request counter. It use the native atomic increment operation to deal with concurrency. We may use built-in data structure such as sync.atomic / sync.mutex on native golang, but it only worked for a single services. May need to use a sticky sessions to let the request know which server this request should go to. So we choose Redis as a its storage since easier to use when it's coming to scaling. Other option like Memcached may also worked, but in terms of data persistence Redis would be better option. 