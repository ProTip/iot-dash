run: ssl
	docker-compose up --build

ssl:
	openssl req -x509 -nodes -days 5 -newkey rsa:2048 -keyout ./server/ssl/key.pem -out ./server/ssl/cert.pem -config ./server/ssl/req.conf -extensions 'v3_req'


.PHONY: ssl