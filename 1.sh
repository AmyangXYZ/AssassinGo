sudo docker-compose stop 
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o AssassinGo
sudo docker-compose up --build -d
