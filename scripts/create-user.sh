curl -v -X POST -H "Authorization: Bearer XXX" -H "Content-Type: application/json" -d '{
  "name": "John Snow"
}' http://localhost:8082/create

docker logs $(docker ps -qf "name=globalwebindex_user.1*") -f