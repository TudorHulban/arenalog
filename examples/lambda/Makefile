build:
	docker build -t go-lambda-bench .

run:
	docker run --rm --memory="128m" -p 8080:8080 go-lambda-bench

request:
	curl -XPOST "http://localhost:8080/2015-03-31/functions/function/invocations" -d '{}'
