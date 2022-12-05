go build
MONGODB_CONNSTRING_ENV="mongodb://root:example@localhost:27017/" MONGODB_DBNAME_ENV="tema_sprc_02" go run service