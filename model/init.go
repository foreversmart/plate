package model

var (
	mongo      *Model
	adminMongo *Model
)

// SetupModel 设置db
func SetupModel(model *Model) {
	mongo = model
}

// MongoModel 返回db
func MongoModel() *Model {
	return mongo
}

