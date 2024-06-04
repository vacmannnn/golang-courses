make

./xkcd-server & sleep 5

response=$(curl -d '{"username":"admin", "password":"admin"}' "http://localhost:8080/login")

curl -X POST --header "Authorization:"$response "http://localhost:8080/update"

echo ""

curl -X GET --header "Authorization:"$response "http://localhost:8080/pics?search='apple,doctor'"
