FROM node AS npmbuild
	COPY interface /build
	WORKDIR /build

	RUN npm install
	RUN npm run build

FROM golang AS gobuild
	WORKDIR /build

	RUN go env -w CGO_ENABLED=0

	COPY web.go web.go
	RUN go mod init web
	RUN go mod tidy -v

	COPY --from=npmbuild /build/public interface/public

	RUN go build -v -o web

# Not scratch because of SSL certs.
FROM alpine
	EXPOSE 8080

	COPY --from=gobuild /build/web web

	ENTRYPOINT ["/web"]
