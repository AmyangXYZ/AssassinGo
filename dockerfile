FROM alpine

ADD web/templates/ /web/templates
ADD web/static/ /web/static
ADD gatherer/dict/ /gatherer/dict/
ADD AssassinGo /

CMD [ "./AssassinGo" ]
