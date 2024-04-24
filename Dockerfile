FROM golang:1.17-alpine AS builder

COPY . /home/fire/

RUN cd /home/fire/ && CGO_ENABLED=0 GOOS=linux GO111MODULE="on" GOPROXY=http://goproxy.cn,direct go build -o /home/fire/bin/fire /home/fire/cmd/app/server.go

FROM alpine

WORKDIR kubeedge

COPY --from=builder /home/fire/bin/fire kubeedge/
COPY --from=builder /home/fire/config/config.yaml kubeedge/etc/
COPY --from=builder /home/fire/config/deviceProfile.json /opt/kubeedge/

ENTRYPOINT ["kubeedge/fire", "--v", "3"]

CMD ["fire"]
