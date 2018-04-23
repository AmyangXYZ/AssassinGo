FROM alpine

ADD web/templates/ /web/templates
ADD gatherer/dict/ /gatherer/dict/
ADD AssassinGo /

CMD [ "./AssassinGo" ]
