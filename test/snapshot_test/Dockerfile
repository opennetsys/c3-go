FROM node:8

WORKDIR /
RUN git clone https://github.com/wolfcw/libfaketime.git
WORKDIR /libfaketime/src
RUN make install

WORKDIR /usr/src/app

COPY package*.json ./

RUN npm install

COPY . .

CMD ["/bin/sh", "-c", "LD_PRELOAD=/usr/local/lib/faketime/libfaketime.so.1 FAKETIME_NO_CACHE=1 faketime -f '@2017-01-01 00:00:00' npm start"]
