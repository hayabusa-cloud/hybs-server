package samplegame

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
	"github.com/hayabusa-cloud/hybs-server/application/model/samplegame/master"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

// a simple example API for receiving request parameters
func v1ExampleRequestParameters(ctx hybs.Ctx) {
	// example for receive query args like this:
	// GET http://localhost:8087/v1/example/request-parameters/?int_value=100&string_value=hello
	var queryArgs = ctx.QueryArgs().String()         // int_value=100&string_value=hello
	var intValue = ctx.FormInt("int_value")          // 100
	var stringValue = ctx.FormString("string_value") // hello
	// see parameters value in response
	ctx.SetResponseValue("queryArgs", queryArgs)
	ctx.SetResponseValue("intValue", intValue)
	ctx.SetResponseValue("stringValue", stringValue)
	// example for receiving body form values like this:
	// POST http://localhost:8087/v1/example/request-parameters/
	// form body(post form args):
	// post_int_value=200&post_string_value=hello2
	var postArgs = ctx.PostArgs().String()                    // post_int_value=200&post_string_value=hello2
	var postIntValue = ctx.FormInt("post_int_value")          // 200
	var postStringValue = ctx.FormString("post_string_value") // hello2
	// see parameters value in response
	ctx.SetResponseValue("postQueryArgs", postArgs)
	ctx.SetResponseValue("postIntValue", postIntValue)
	ctx.SetResponseValue("postStringValue", postStringValue)
	// set response status code
	ctx.StatusOk() // can be omitted
	// ctx.StatusBadRequest()
	// ctx.StatusForbidden()
	// ctx.StatusInternalServerError()
}

// a simple example API for receiving route parameters
func v1ExampleRouteParameters(ctx hybs.Ctx) {
	// for example, a route pattern defined in controller files like this:
	// http://localhost:8087/v1/example/route-parameters/:id/:name
	// and then, send a request like this:
	// http://localhost:8087/v1/example/route-parameters/100/hello
	var id = ctx.RouteParamInt("id")        // 100
	var name = ctx.RouteParamString("name") // hello
	// see parameters value in response
	ctx.SetResponseValue("intValue", id)
	ctx.SetResponseValue("stringValue", name)
	// set response status code
	// ctx.StatusNotFound()
	ctx.StatusOk()
}

// a simple example API to get master record data
func v1ExampleMasterData(ctx hybs.Ctx) {
	// get request parameter "level"
	var lv = ctx.FormInt("level")
	// find master data record
	var val, ok = ctx.MasterTable("LevelExp").Record(lv)
	if !ok {
		ctx.StatusBadRequest("level")
		return
	}
	// output to response
	ctx.SetResponseValue("exp", val.(*master.LevelExp).Exp)
}

// a simple example API for rw sqlite tables
func v1ExampleSqlite(ctx hybs.Ctx) {
	type playerExampleData struct {
		HayabusaID  string `sql:"hayabusa_id" json:"hayabusaId"`
		IntField    int    `sql:"int_field" json:"intField"`
		StringField string `sql:"string_field" json:"stringField"`
	}
	var exampleData = playerExampleData{
		HayabusaID:  ctx.User().ID(),
		IntField:    1000,
		StringField: "hello hayabusa",
	}

	var tx = ctx.Sqlite().Table("player_example_data")
	// insert data
	if err := tx.Create(&exampleData).Error; err != nil {
		ctx.SysLogf("sqlite insert player_example_data failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
	// update data
	if err := tx.Where("hayabusa_id=?", ctx.User().ID()).
		Update("int_field", 1001).Error; err != nil {
		ctx.SysLogf("sqlite update player_example_data failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
	// select data
	if err := tx.Where("hayabusa_id=?", ctx.User().ID()).First(&exampleData).Error; err != nil {
		ctx.SysLogf("sqlite update player_example_data failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
	// delete data
	if err := tx.Where("hayabusa_id=?", ctx.User().ID()).Delete(&exampleData).Error; err != nil {
		ctx.SysLogf("mysql delete player_example_data failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
	ctx.SetResponseValue("exampleData", exampleData)
}

// a simple example API for rw mysql tables
func v1ExampleMysql(ctx hybs.Ctx) {
	type playerExampleData struct {
		HayabusaID  string `sql:"hayabusa_id" json:"hayabusaId"`
		IntField    int    `sql:"int_field" json:"intField"`
		StringField string `sql:"string_field" json:"stringField"`
	}
	var exampleData = playerExampleData{
		HayabusaID:  ctx.User().ID(),
		IntField:    1000,
		StringField: "hello hayabusa",
	}

	var tx = ctx.MySQL().Table("player_example_data")
	// insert data
	if err := tx.Create(&exampleData).Error; err != nil {
		ctx.SysLogf("mysql insert player_example_data failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
	// update data
	if err := tx.Where("hayabusa_id=?", ctx.User().ID()).
		Update("int_field", 1001).Error; err != nil {
		ctx.SysLogf("mysql update player_example_data failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
	// select data
	if err := tx.Where("hayabusa_id=?", ctx.User().ID()).First(&exampleData).Error; err != nil {
		ctx.SysLogf("mysql update player_example_data failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
	// delete data
	if err := tx.Where("hayabusa_id=?", ctx.User().ID()).Delete(&exampleData).Error; err != nil {
		ctx.SysLogf("mysql delete player_example_data failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
	ctx.SetResponseValue("exampleData", exampleData)
}

// a simple example API for rw mongodb collection
func v1ExampleMongodb(ctx hybs.Ctx) {
	type PlayerExampleData struct {
		IntValue      int       `bson:"int_value" json:"intValue"`
		StringValue   string    `bson:"string_value" json:"stringValue"`
		DatetimeValue hybs.Time `bson:"datetime_value" json:"datetimeValue"`
		SliceValue    []int     `bson:"slice_value" json:"sliceValue"`
	}
	var exampleData = &PlayerExampleData{
		IntValue:      1000,
		StringValue:   "hello hayabusa",
		DatetimeValue: ctx.Now(),
		SliceValue:    []int{50, 100, 150},
	}
	var mgoCollection = ctx.Mongo().Collection("player_example_data")
	// insert data like this
	if _, err := mgoCollection.InsertOne(ctx.Context(), exampleData); err != nil {
		ctx.SysLogf("mongo insert example data failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
	// update data like this
	if _, err := mgoCollection.UpdateOne(
		ctx.Context(),
		bson.M{
			"int_value": 1000,
		},
		bson.M{
			"$set": bson.M{
				"string_value": "updated",
			},
			"$push": bson.M{
				"slice_value": 200,
			},
		}); err != nil {
		ctx.SysLogf("mongo update example data failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
	// select data like this
	if err := mgoCollection.FindOne(
		ctx.Context(),
		bson.M{
			"int_value": 1000,
		}).
		Decode(&exampleData); err != nil {
		ctx.SysLogf("mongo find example data failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
	// delete data like this
	if _, err := mgoCollection.DeleteOne(ctx.Context(), bson.M{"int_value": 1000}); err != nil {
		ctx.SysLogf("mongo delete example data failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
	// return example data
	ctx.SetResponseValue("exampleData", exampleData)
}

// an example API for rw redis cache
func v1ExampleRedis(ctx hybs.Ctx) {
	var redisClient = ctx.Redis("Redis")
	// default redis id is defined by controller files, can be omitted as follow:
	// var redisClient = ctx.Redis()
	redisClient.Set(ctx.Context(), "TestSpace:TestKey", ctx.Now().Unix(), time.Hour)
	var val, err = redisClient.Get(ctx.Context(), "TestSpace:TestKey").Int64()
	if err != nil {
		ctx.SysLogf("get redis value failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
	ctx.SetResponseValue("testValue", val)
}

func v1ExampleDistributeLock(ctx hybs.Ctx) {
	var redisClient = ctx.Redis("Redis")
	// example global distribute lock/unlock via redis
	if err := redisClient.Lock("SampleMutex"); err != nil {
		ctx.SysLogf("redis distribute lock failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
	// do somethings here
	// ...
	time.Sleep(time.Second)
	// release lock
	if _, err := redisClient.Unlock("SampleMutex"); err != nil {
		ctx.SysLogf("redis distribute unlock failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
}

// an example API for rw local memory cache
func v1ExampleCache(ctx hybs.Ctx) {
	// cache somethings in local memory
	ctx.Cache("Cache").Set("SampleCacheKey1", 2000, time.Hour)
	// cache id can be omitted
	ctx.Cache().Set("SampleCacheKey2", "hello hayabusa", time.Hour)
	_ = ctx.Cache().Add("SampleCacheKey3", 100, time.Hour)
	_, _ = ctx.Cache().IncrementInt("SampleCacheKey3", 1)
	var iVal, _ = ctx.Cache().Get("SampleCacheKey3")
	// set response
	ctx.SetResponseValue("sampleCacheKey3", iVal)
}

func init() {
	hybs.RegisterService("ExampleRequestParametersV1", v1ExampleRequestParameters)
	hybs.RegisterService("ExampleRouteParametersV1", v1ExampleRouteParameters)
	hybs.RegisterService("ExampleMasterDataV1", v1ExampleMasterData)
	hybs.RegisterService("ExampleSqliteV1", v1ExampleSqlite)
	hybs.RegisterService("ExampleMySQLV1", v1ExampleMysql)
	hybs.RegisterService("ExampleMongodbV1", v1ExampleMongodb)
	hybs.RegisterService("ExampleRedisV1", v1ExampleRedis)
	hybs.RegisterService("ExampleDistributeLock", v1ExampleDistributeLock)
	hybs.RegisterService("ExampleCacheV1", v1ExampleCache)
}
