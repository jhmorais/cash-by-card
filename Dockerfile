FROM ubuntu:18.04

RUN apt-get update -y
RUN apt-get install -y build-essential
RUN apt-get install -y wget git
RUN apt-get install -y curl
RUN apt-get install -y zip

RUN cd /tmp
RUN wget https://go.dev/dl/go1.20.linux-amd64.tar.gz
RUN tar -C /usr/lib -xzf go1.20.linux-amd64.tar.gz
RUN rm -rf go1.20.linux-amd64.tar.gz

ENV APP_HOME /cashbycard

ENV PATH=/usr/lib/go/bin:$PATH

WORKDIR $APP_HOME

ADD . $APP_HOME

RUN echo $PATH
RUN go mod download

RUN go build -o cashbycard ./cmd/restserver
RUN chmod +x cashbycard

COPY . /cashbycard

CMD ["./cashbycard"]