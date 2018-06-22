FROM node:8

WORKDIR /
RUN git clone https://github.com/wolfcw/libfaketime.git
WORKDIR /libfaketime/src
RUN make install

# Create app directory
WORKDIR /usr/src/app



# Install app dependencies
# A wildcard is used to ensure both package.json AND package-lock.json are copied
# where available (npm@5+)
COPY package*.json ./

RUN npm install
# If you are building your code for production
# RUN npm install --only=production

# Bundle app source
COPY . .

#EXPOSE 9999

CMD ["/bin/sh", "-c", "LD_PRELOAD=/usr/local/lib/faketime/libfaketime.so.1 FAKETIME_NO_CACHE=1 npm start"]
