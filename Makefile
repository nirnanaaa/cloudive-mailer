start:
	docker-compose up -d --remove-orphans
stop:
	docker-compose down --remove-orphans
build:
	./build.py
clean:
	rm -rf build