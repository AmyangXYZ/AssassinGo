FROM alpine

RUN apk --no-cache add ca-certificates

ADD AssassinGo /

CMD [ "./AssassinGo" ]
