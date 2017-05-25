docker run -d -p 2181:2181 -p 9092:9092 --env ADVERTISED_HOST=127.0.0.1 --env ADVERTISED_PORT=9092 --name=kafka spotify/kafka
docker run -d -p 27017:27017 --name mongo -v /var/mongo/datadir:/data/db mongo
docker run -d -p 6379:6379 --name redis -v /data/redis:/data redis redis-server --appendonly yes 
