package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	//mcrRdsStore "github.com/go-macaron/session/redis"

	"github.com/astaxie/beego/session"
	_ "github.com/astaxie/beego/session/redis"
	ke "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/metadata"
	"gl.king.im/king-lib/framework"
	"gl.king.im/king-lib/framework/config"
)

var globalAuthToken string
var sessionManagerMap sync.Map

type SessionManager struct {
	m    *session.Manager
	conf SessionManagerConf
}

type SessionManagerConf struct {
	ProvideName string
	Cf          *session.ManagerConfig
	Prefix      string
}

func NewSessionManager(key string, conf SessionManagerConf) (*SessionManager, error) {

	if val, ok := sessionManagerMap.Load(key); ok {
		return val.(*SessionManager), nil
	}

	serviceConf := config.GetServiceConf()
	if serviceConf.Data.Redis.Addr == "" {
		panic("Redis Addr配置为空！")
	}

	if conf.Cf.ProviderConfig == "" {
		//配置依次为地址、连接个数、密码、db编号、连接空闲超时
		providerConfArr := []string{serviceConf.Data.Redis.Addr, "100"}
		providerConfArr = append(providerConfArr, serviceConf.Data.Redis.Password)
		providerConfArr = append(providerConfArr, "0")
		providerConfArr = append(providerConfArr, fmt.Sprintf("%d", serviceConf.Data.Redis.ReadTimeout.Seconds))
		providerConf := strings.Join(providerConfArr, ",")
		//fmt.Println(providerConf)
		conf.Cf.ProviderConfig = providerConf
	}

	innersm, err := session.NewManager("redis", conf.Cf)
	if err != nil {
		return nil, err
	}
	go innersm.GC()

	sm := &SessionManager{
		m:    innersm,
		conf: conf,
	}

	sessionManagerMap.Store(key, sm)

	return sm, nil
}

func GetSessionManager(key string) (*SessionManager, bool) {
	if val, ok := sessionManagerMap.Load(key); ok {
		return val.(*SessionManager), true
	} else {
		return nil, false
	}
}

func (sm *SessionManager) GenerateSessionId() string {
	return GenerateSessionId()
}

func (sm *SessionManager) addPrefix(token string) string {
	return sm.conf.Prefix + token
}

func (sm *SessionManager) SessionGet(token string, key interface{}) (interface{}, error) {
	token = sm.addPrefix(token)
	if !sm.m.GetProvider().SessionExist(token) {
		return nil, ke.New(401, "SessionExist_ERROR", "用户会话信息不存在！")
	}

	sesssionStore, err := sm.m.GetProvider().SessionRead(token)
	if err != nil {
		return nil, err
	}
	beegoval := sesssionStore.Get(key)
	return beegoval, nil

	//	macronStore, err := providerIns.Read(token)
	//	if err != nil {
	//		return nil, err
	//	}
	//	val := macronStore.Get(key)
	//	return val, nil
}

func (sm *SessionManager) SessionGetSilence(token string, key interface{}) (interface{}, error) {
	token = sm.addPrefix(token)
	if !sm.m.GetProvider().SessionExist(token) {
		return nil, nil
	}

	sesssionStore, err := sm.m.GetProvider().SessionRead(token)
	if err != nil {
		return nil, err
	}
	beegoval := sesssionStore.Get(key)
	return beegoval, nil

	//	macronStore, err := providerIns.Read(token)
	//	if err != nil {
	//		return nil, err
	//	}
	//	val := macronStore.Get(key)
	//	return val, nil
}

func (sm *SessionManager) SessionSet(token string, key, value interface{}) error {
	token = sm.addPrefix(token)
	if !sm.m.GetProvider().SessionExist(token) {
		return ke.New(401, "SessionExist_ERROR", "用户会话信息不存在！")
	}
	sesssionStore, err := sm.m.GetProvider().SessionRead(token)
	if err != nil {
		return err
	}
	sesssionStore.Set(key, value)
	sesssionStore.SessionRelease(NewNoneResponseWriter())
	return nil

	//	macronStore, err := providerIns.Read(token)
	//	if err != nil {
	//		return err
	//	}
	//
	//	err = macronStore.Set(key, value)
	//	if err != nil {
	//		return err
	//	}
	//	macronStore.Release()
	//
	//	return nil
}

func (sm *SessionManager) SessionDelete(token string, key interface{}) (bool, error) {
	token = sm.addPrefix(token)
	if !sm.m.GetProvider().SessionExist(token) {
		return false, ke.New(401, "SessionExist_ERROR", "用户会话信息不存在！")
	}

	sesssionStore, err := sm.m.GetProvider().SessionRead(token)
	if err != nil {
		return false, err
	}

	err = sesssionStore.Delete(key)

	if err != nil {
		return false, err
	}

	sesssionStore.SessionRelease(NewNoneResponseWriter())

	return true, nil
}

func (sm *SessionManager) SessionDestroy(token string) error {
	token = sm.addPrefix(token)
	return sm.m.GetProvider().SessionDestroy(token)

	//	macronStore, err := providerIns.Read(token)
	//	if err != nil {
	//		return nil, err
	//	}
	//	val := macronStore.Get(key)
	//	return val, nil
}

func (sm *SessionManager) SessionCreate(sessionID string) error {
	sessionID = sm.addPrefix(sessionID)

	sesssionStore, err := sm.m.GetProvider().SessionRead(sessionID)
	if err != nil {
		return err
	}

	sesssionStore.SessionRelease(NewNoneResponseWriter())

	return nil
}

var (
	//providerIns    mcrRdsStore.RedisProvider
	//globalSessions *session.Manager
	globalSessions *SessionManager
)

func init() {

	//	providerIns = mcrRdsStore.RedisProvider{}
	//	prefix := "BOSS-AUTH-SESSION-"
	//	configs := fmt.Sprintf("network=%s,addr=%s,password=%s,db=%s,pool_size=%s,idle_timeout=%s,prefix=%s", "tcp", ":6379", "", "0", "100", "180", prefix)
	//	err := providerIns.Init(int64(time.Minute*60/time.Second), configs)
	//	if err != nil {
	//		panic(err)
	//	}
}

func NewGlobalSessionManager() {
	if globalSessions != nil {
		return
	}
	serviceConf := config.GetServiceConf()
	if serviceConf.Data.Redis.Addr == "" {
		panic("Redis Addr配置为空！")
	}

	//配置依次为地址、连接个数、密码、db编号、连接空闲超时
	providerConfArr := []string{serviceConf.Data.Redis.Addr, "100"}
	providerConfArr = append(providerConfArr, serviceConf.Data.Redis.Password)
	providerConfArr = append(providerConfArr, "0")
	providerConfArr = append(providerConfArr, fmt.Sprintf("%d", serviceConf.Data.Redis.ReadTimeout.Seconds))
	providerConf := strings.Join(providerConfArr, ",")
	//fmt.Println(providerConf)
	smc := SessionManagerConf{
		ProvideName: "redis",
		Cf:          &session.ManagerConfig{CookieName: "sessionid", Gclifetime: int64(time.Hour*48) / int64(time.Second), ProviderConfig: providerConf},
		Prefix:      "BOSS-AUTH-SESSION-",
	}

	innergs, _ := session.NewManager("redis", smc.Cf)
	go innergs.GC()

	globalSessions = &SessionManager{
		m:    innergs,
		conf: smc,
	}

}

func GenerateSessionId() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func SetGlobalAuthToken(token string) {
	globalAuthToken = token
}

func GetGlobalAuthToken() string {
	return globalAuthToken
}

func GetContextSessionId(ctx context.Context) (string, error) {
	md, ok := metadata.FromServerContext(ctx)
	log.Println("GetContextSessionId...md, ok:", md, ok)
	if !ok {
		return "", ke.New(400, "METADATA_INFO_ERROR", "元信息获取失败！")
	}

	authToken := md.Get(framework.METADATA_KEY_AUTH_TOKEN)

	return authToken, nil
}

func GetContextScenario(ctx context.Context) (string, error) {
	md, ok := metadata.FromServerContext(ctx)
	if !ok {
		return "", ke.New(400, "METADATA_INFO_ERROR", "元信息获取失败！")
	}

	scenario := md.Get(framework.METADATA_KEY_CALL_SCENARIO)

	return scenario, nil
}

// 暂未使用，不同的session manager存储的用户结构可能不一样
type UserIdentity struct {
	Id            uint32
	Username      string
	Email         string
	Phone         string
	Password      string
	LastLoginTime string
	LdapDn        string
	LdapCn        string
}

func GetIdentityByToken(authToken string, identity interface{}) error {

	identityJson, err := globalSessions.SessionGet(authToken, "user")
	if err != nil {
		return ke.New(400, "FRAMEWORK_SESSION_GET_ERROR", "Session信息获取失败！"+err.Error())
	}

	if identityJson == nil {
		return ke.New(401, "SESSION_INFO_NIL", "获取Session信息为空！session已过期或sessionid错误！")
	}

	err = json.Unmarshal(identityJson.([]byte), &identity)
	if err != nil {
		return ke.New(400, "SESSION_UNMARSHAL_ERROR", "Session信息解析失败")
	}

	return nil
}

func SessionGet(token string, key interface{}) (interface{}, error) {
	return globalSessions.SessionGet(token, key)
}

func SessionGetSilence(token string, key interface{}) (interface{}, error) {
	return globalSessions.SessionGetSilence(token, key)
}

func SessionSet(token string, key, value interface{}) error {
	return globalSessions.SessionSet(token, key, value)
}

func SessionDelete(token string, key interface{}) (bool, error) {
	return globalSessions.SessionDelete(token, key)
}

func SessionDestroy(token string) error {
	return globalSessions.SessionDestroy(token)
}

func SessionCreate(token string) error {
	return globalSessions.SessionCreate(token)
}

type NoneResponseWriter struct {
	body       []byte
	statusCode int
	header     http.Header
}

func NewNoneResponseWriter() *NoneResponseWriter {
	return &NoneResponseWriter{
		header: http.Header{},
	}
}

func (w *NoneResponseWriter) Header() http.Header {
	return w.header
}

func (w *NoneResponseWriter) Write(b []byte) (int, error) {
	w.body = b
	// implement it as per your requirement
	return 0, nil
}

func (w *NoneResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

//	//gorilla/sessions ---SET
//
//	// Note: Don't store your key in your source code. Pass it via an
//	// environmental variable, or flag (or both), and don't accidentally commit it
//	// alongside your code. Ensure your key is sufficiently random - i.e. use Go's
//	// crypto/rand or securecookie.GenerateRandomKey(32) and persist the result.
//
//	// Fetch new store.
//	store, err := redistore.NewRediStore(10, "tcp", ":6379", "", []byte("SESSION-SECRET-KEY"))
//	if err != nil {
//		panic(err)
//	}
//	defer store.Close()
//
//	// Get a session. We're ignoring the error resulted from decoding an
//	// existing session: Get() always returns a session, even if empty.
//
//	token := "gorilla-sessions-token"
//	session, _ := store.Get(&http.Request{}, token)
//	// Set some session values.
//	session.Values["gorilla-sessions-is_login"] = true
//	session.Values["gorilla-sessions-id"] = 43
//	session.Values["gorilla-sessions-name"] = "gorilla-sessions-testname"
//	// Save it before we write to the response/return from the handler.
//
//	err = session.Save(&http.Request{}, NewCustomResponseWriter())
//	store.Save(&http.Request{}, NewCustomResponseWriter(), session)
//
//	log.Println("err::", err)

//	//gorilla/sessions  ----GET
//	redistore.RediStore{}
//	store, err := redistore.NewRediStore(10, "tcp", ":6379", "", []byte("SESSION-SECRET-KEY"))
//	if err != nil {
//		panic(err)
//	}
//	defer store.Close()
//
//	token := "gorilla-sessions-token"
//	session, _ := store.Get(&http.Request{}, token)
//
//	log.Println("session::", session, session.Values)
