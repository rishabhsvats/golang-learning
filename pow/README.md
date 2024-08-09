
# Blockchain Proof of Work

Start with the genesis block:
~~~
go run main.go
~~~

Create new block in the blockchain
~~~
curl -X POST -H "Content-Type: application/json" -d '{"name": "19"}' http://localhost:8080
~~~
