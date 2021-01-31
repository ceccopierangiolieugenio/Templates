FROM python:3.9.1-alpine3.12

RUN apk add git
RUN pip install git+https://github.com/behave/behave@v1.2.6
RUN pip install requests

COPY functest .

CMD sleep 10 ; behave
