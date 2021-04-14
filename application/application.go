package application

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	hybs "github.com/hayabusa-cloud/hayabusa"
	"github.com/hayabusa-cloud/hybs-server/application/common"
	"github.com/hayabusa-cloud/hybs-server/application/middleware/realtime"
	"github.com/hayabusa-cloud/hybs-server/application/model/central"
	"github.com/hayabusa-cloud/hybs-server/application/model/games"
	"github.com/hayabusa-cloud/hybs-server/application/model/platform"
	"github.com/jinzhu/configor"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/mongo"
)

// Up init somethings at server startup timing
func Up(eng hybs.Engine) {
	loadAppConfig(eng)
	ensureMongoIndices(eng)
	loadRedisScripts(eng)
	realtimeDebugToken(eng)
}

// Down stops service before exit process
func Down(eng hybs.Engine) {
	eng.LogWarnf("hayabusa engine stopped")
}

func loadAppConfig(eng hybs.Engine) {
	var appConf = eng.AppConfig()
	var configLoader = configor.New(&configor.Config{
		AutoReload:         true,
		AutoReloadInterval: common.ConfigReload,
	})
	var err = configLoader.Load(&common.Config, appConf)
	if err != nil {
		fmt.Printf("load application config failed:\n%s\n", err)
		eng.Stop()
	}
}

func ensureMongoIndices(eng hybs.Engine) {
	var ensureIndexFn = func(dbID string, collection string, index mongo.IndexModel) (err error) {
		err = eng.MongoDo(dbID, func(db *mongo.Database) (err error) {
			// try create collection
			db.CreateCollection(eng.Context(), collection)
			var idxView, _ = db.Collection(collection).Indexes(), &mongo.Cursor{}
			// create new index
			_, err = idxView.CreateOne(eng.Context(), index)
			return
		})
		return
	}
	var ensureAppIndexFn = func(collection string, index mongo.IndexModel) (err error) {
		for _, app := range common.Config.Apps {
			ensureIndexFn(app.DatabaseID, collection, index)
		}
		return
	}
	// platform tables
	ensureIndexFn("MongoPlatform", "hayabusa_id", platform.HayabusaIDIndex)
	// game common tables
	ensureAppIndexFn("onetime_token", games.OnetimeTokenIndex)
	ensureAppIndexFn("onetime_token", games.OnetimeTokenHayabusaIndex)
	ensureAppIndexFn("onetime_token", games.OnetimeTokenExpireIndex)
	ensureAppIndexFn("player", games.PlayerIndex)
	// central tables
	ensureIndexFn("MongoCentral", "sys_realtime_server", central.SysRealtimeServerAppIDIndex)
	ensureIndexFn("MongoCentral", "sys_realtime_server", central.SysRealtimeServerOwnerIndex)
	ensureIndexFn("MongoCentral", "sys_realtime_server", central.SysRealtimeServerExpireIndex)
}

func loadRedisScripts(eng hybs.Engine) {
	const tokenExpireSeconds = 60*60*6 + 60
	var onetimeTokenGetSetScript = fmt.Sprintf(`
	local exists = redis.call('EXISTS', KEYS[1])
	local result = {false, nil, nil, nil, nil}
	if exists ~= 1 then
		return result
	end
	local hayabusaID = redis.call('HGET', KEYS[1], 'HayabusaID')
	local serverID = redis.call('HGET', KEYS[1], 'ServerID')
	local permission = redis.call('HGET', KEYS[1], 'Permission')
	local generatedAt = redis.call('HGET', KEYS[1], 'GeneratedAt')
	redis.call('DEL', KEYS[1])
	redis.call('HSET', KEYS[2], 'HayabusaID', hayabusaID)
	redis.call('HSET', KEYS[2], 'ServerID', serverID)
	redis.call('HSET', KEYS[2], 'Permission', permission)
	redis.call('HSET', KEYS[2], 'GeneratedAt', ARGV[2])
	redis.call('EXPIRE', KEYS[2], %d)
	result[0] = true
	result[1] = hayabusaID
    result[2] = serverID
	result[3] = permission
    result[4] = generatedAt
	return result
`, tokenExpireSeconds)
	var onetimeTokenGetOnlyScript = fmt.Sprintf(`
	local exists = redis.call('EXISTS', KEYS[1])
	local result = {false, nil, nil, nil, nil}
	if exists ~= 1 then
		return result
	end
	local hayabusaID = redis.call('HGET', KEYS[1], 'HayabusaID')
	local serverID = redis.call('HGET', KEYS[1], 'ServerID')
	local permission = redis.call('HGET', KEYS[1], 'Permission')
	local generatedAt = redis.call('HGET', KEYS[1], 'GeneratedAt')
	redis.call('HSET', KEYS[1], 'GeneratedAt', ARGV[1])
	redis.call('EXPIRE', KEYS[1], %d)
	result[0] = true
	result[1] = hayabusaID
    result[2] = serverID
	result[3] = permission
    result[4] = generatedAt
	return result
`, tokenExpireSeconds)
	var err = eng.RedisDo("Redis", func(r redis.UniversalClient) error {
		// onetime token get set script
		var cmd = r.ScriptLoad(eng.Context(), onetimeTokenGetSetScript)
		if err := cmd.Err(); err != nil {
			return fmt.Errorf("redis script load failed:%s", err)
		}
		// store SHA1 of scripts
		eng.SetUserValue("/Redis/Scripts/OnetimeTokenGetSet", cmd.Val())

		// onetime token get only script
		cmd = r.ScriptLoad(eng.Context(), onetimeTokenGetOnlyScript)
		if err := cmd.Err(); err != nil {
			return fmt.Errorf("redis script load failed:%s", err)
		}
		// store SHA1 of scripts
		eng.SetUserValue("/Redis/Scripts/OnetimeTokenGetOnly", cmd.Val())

		return nil
	})
	if err != nil {
		eng.LogErrorf("redis error:%s", err)
	}
}

func realtimeDebugToken(eng hybs.Engine) {
	eng.CacheDo("Cache", func(c *cache.Cache) error {
		var debugToken = &realtime.AccessToken{
			AppID:      "jp:playground:hayabusa-cloud:debug",
			HayabusaID: []byte("hayabusa-cloud"),
			Mux:        "default",
			Token:      []byte("anA6cGxheWdyb3VuZDpoYXlhYnVzYS1jbG91ZDpkZWJ1ZyNoYXlhYnVzYS1jbG91ZCNta1pTWXlJd3VZR3dKdWJD"),
			ValidUntil: hybs.Time{},
		}
		var cacheKey = fmt.Sprintf("/Realtime/Tokens/%s", debugToken.Token)
		c.Set(cacheKey, debugToken, -1)
		return nil
	})
}

func init() {
	hybs.RegisterRealtimeHandler("EchoTest", func(ctx hybs.RealtimeCtx) {
		ctx.SetHeader(hybs.RealtimeHeader0SCOriginal)
		ctx.SetEventCode(0x80f0)
		ctx.Response()
	})
}
