postgresql_local:
	@go run . postgres --ip 127.0.0.1 --mode master --name master_1 \
	--dbport 5432 --dbpassword hui123456
	@go run . postgres --ip 127.0.0.1 --mode slave --name slave_1 \
	--dbport 5432 --dbpassword hui123456 --master_ip 192.168.110.113

# docker run --name autobot-postgres14 \
-e POSTGRES_PASSWORD=zx123456 -d -p 54302:5432 \
postgres:12-alpine

postgresql_remote_12:
	@go run . postgres \
	--ip 172.18.3.9 \
	--username root \
	--password root \
	--mode master \
	--name autobot-postgres14 \
	--dbport 54302 \
	--dbpassword zxpostgres \
	--update true \
	--version 12

	@go run . postgres \
	--ip 172.18.3.9 \
	--username root \
	--password root \
	--mode slave \
	--name db-slave-1 \
	--dbport 54033 \
	--dbpassword zxpostgres \
	--master_ip 172.18.3.9 \
	--master_port 54302 \
	--version 12

postgresql_local_clean:
	@docker rm -f master_1
	@docker rm -f slave_1
	@docker image prune