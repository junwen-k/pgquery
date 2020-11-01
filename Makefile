connect:
	docker run --name postgresqldb -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -p 5433:5433 -d postgres
	docker exec -it postgresqldb psql -U postgres
	
clean-testcache:
	go clean -testcache