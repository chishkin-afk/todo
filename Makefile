

.PHONY: all quick localtls

all: quick

quick:
	echo "writing .env..."
	mv .env.example .env || true
	echo "building..."
	docker-compose up -d

localtls:
	mkcert --install

	echo "generating certs & keys..."
	mkcert localhost 127.0.0.1
	mv localhost+1.pem server.crt
	mv localhost+1-key.pem server.key

	mkdir certs
	echo "moving files to cert dir..."
	mv server.crt certs/
	mv server.key certs/