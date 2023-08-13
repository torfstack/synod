FROM node:18

COPY build build
COPY package.json .
COPY package-lock.json .

CMD ["npm", "ci", "--omit", "dev"]

EXPOSE 3000
CMD ["node", "build"]