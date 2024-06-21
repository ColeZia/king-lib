package coder

import (
	"encoding/json"
	"io/ioutil"
	netHttp "net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/encoding"

	//_ "github.com/go-kratos/kratos/v2/encoding/form"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	"gl.king.im/king-lib/framework/filters"
	"gl.king.im/king-lib/framework/internal/tracing"
	"gl.king.im/king-lib/protobuf/api/common"

	//"github.com/golang/protobuf/jsonpb"
	//"github.com/golang/protobuf/proto"

	commonPb "gl.king.im/king-lib/protobuf/api/common"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

//var DebugInfo interface{}

type SerializedJsonDebugInfo string

const RESPONSE_FORMAT_JSON = "application/json"
const RESPONSE_FORMAT_HTML = "text/html"
const RESPONSE_FORMAT_PLAIN = "text/plain"

var responseFormat string
var lock = &sync.Mutex{}

func init() {
	SetResponseFormat(RESPONSE_FORMAT_JSON)
}

func SetResponseFormat(format string) {
	lock.Lock()
	defer lock.Unlock()
	responseFormat = format
}

func recuirseMessage(path []int, ir protoreflect.Message, jsonReplace map[string]string) map[string]string {
	replaceFlagWrapper := func(s string) string { return "xxxxx" + s + "xxxxx" }

	index := 0
	ir.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		path = append(path, index)
		//fmt.Println("range::", path, fd.FullName(), v)

		optValue := proto.GetExtension(fd.Options(), common.E_BaseOpts).(*common.BaseOpts)

		if optValue == nil {

		} else {
			if optValue.JsonSerialized {
				//fmt.Println("optValue.JsonSerialize::", path, fd.FullName())
				fval := ir.Get(fd)

				pathStr := []string{}
				for _, v := range path {
					pathStr = append(pathStr, strconv.Itoa(v))
				}
				replaceKey := strings.Join(pathStr, "-")
				replaceKey = replaceFlagWrapper(replaceKey)
				ir.Set(fd, protoreflect.ValueOf(replaceKey))
				jsonReplace[replaceKey] = fval.String()
			}
		}

		if fd.Message() != nil {
			//fmt.Println("fd.Message() != nil::", path, fd.FullName())
			if fd.IsList() {
				fval := ir.Get(fd)
				flist := fval.List()

				for i := 0; i < flist.Len(); i++ {
					path = append(path, i)
					recuirseMessage(path, flist.Get(i).Message(), jsonReplace)
					path = path[:len(path)-1]
				}
			} else if fd.IsMap() {
				fval := ir.Get(fd)
				fmap := fval.Map()

				mapIndex := 0
				fmap.Range(func(mk protoreflect.MapKey, v2 protoreflect.Value) bool {
					switch v2.Interface().(type) {
					case protoreflect.Message:
						path = append(path, mapIndex)
						recuirseMessage(path, v2.Message(), jsonReplace)
						path = path[:len(path)-1]
					}

					mapIndex++

					return true
				})
			} else {
				recuirseMessage(path, ir.Get(fd).Message(), jsonReplace)
			}
		}
		path = path[:len(path)-1]
		index++
		return true
	})

	return jsonReplace

}

//var (
//	traceDebugInfo     sync.Map
//	traceDebugInfoOpen bool
//)
//
//func OpenTraceDebugInfo() {
//	traceDebugInfoOpen = true
//}

//func AppendTraceDebugInfoBAK(ctx context.Context, data interface{}) {
//	if !traceDebugInfoOpen {
//		return
//	}
//
//	traceIdValuer := tracing.TraceID()
//	traceId := traceIdValuer(ctx).(string)
//	if traceId == "" {
//		return
//	}
//
//	storeValI, ok := traceDebugInfo.Load(traceId)
//	if ok {
//		storeValI := storeValI.([]DebugInfo)
//		storeValI = append(storeValI, DebugInfo{
//			Time: time.Now(),
//			Data: data,
//		})
//		traceDebugInfo.Store(traceId, storeValI)
//	} else {
//		traceDebugInfo.Store(traceId, []DebugInfo{{
//			Time: time.Now(),
//			Data: data,
//		}})
//	}
//}

func HttpResponseEncoder() http.EncodeResponseFunc {
	return func(w netHttp.ResponseWriter, r *netHttp.Request, i interface{}) error {
		//低版本有bug，获取不到kratos的context，暂时注释
		ctx := r.Context()

		var ctxDebugInfo *[]*tracing.DebugInfo
		if filters.UseDebugInfoFilter {
			ctxDebugInfo = tracing.DebugInfoFromServerContext(ctx)
		}

		traceId := w.Header().Get("X-TRACE-ID")

		resp := []byte{}
		switch responseFormat {
		case RESPONSE_FORMAT_JSON:
			switch typeVal := i.(type) {
			case *commonPb.HttpResponseReplyWrapperJSON:
				resp = []byte(typeVal.Content)
			default:
				type responseInterface interface{}

				type response struct {
					Code      int         `json:"code"`
					Data      interface{} `json:"data"`
					Ts        string      `json:"ts"`
					Message   string      `json:"message"`
					RequestId string      `json:"request_id"`
				}

				type debugInfoRsp struct {
					response
					DebugInfo interface{} `json:"debug_info"`
				}

				ma := protojson.MarshalOptions{
					EmitUnpopulated: true,
					UseProtoNames:   true,
				}

				ir := i.(proto.Message).ProtoReflect()

				jsonReplace := map[string]string{}
				jsonReplace = recuirseMessage([]int{}, ir, jsonReplace)

				data, err := ma.Marshal(i.(proto.Message))
				//fmt.Println("jsonReplace::", jsonReplace)

				if err != nil {
					return err
				}

				replaceFlag := "||||flag_for_replace||||"

				var reply responseInterface
				baseRsp := response{
					Code:      200,
					Data:      replaceFlag,
					Message:   "",
					Ts:        time.Now().String(),
					RequestId: traceId, //time.Now().UnixNano(),
				}

				debugInfoReplaceFlag := "||||debugInfoReplaceFlag||||"

				if filters.UseDebugInfoFilter {
					debugInfoRsp := &debugInfoRsp{
						baseRsp,
						ctxDebugInfo,
					}

					debugInfoRsp.DebugInfo = debugInfoReplaceFlag
					//switch DebugInfo.(type) {
					//case SerializedJsonDebugInfo:
					//	reply.DebugInfo = debugInfoReplaceFlag
					//default:
					//
					//}

					//if traceDebugInfoOpen {
					//	storeValI, ok := traceDebugInfo.Load(traceId)
					//	if ok {
					//		debugInfoRsp.DebugInfo = storeValI
					//	}
					//}

					reply = debugInfoRsp
				} else {
					reply = baseRsp
				}

				codec := encoding.GetCodec("json")
				//goLog.Println("framework http.ResponseEncoder codec::", codec)
				resp, err = codec.Marshal(reply)
				//data, err := json.Marshal(reply)
				if err != nil {
					return err
				}

				respStr := strings.Replace(string(resp), `"`+replaceFlag+`"`, string(data), 1)

				for k, v := range jsonReplace {
					respStr = strings.Replace(respStr, `"`+k+`"`, v, 1)
				}

				//switch typeval := DebugInfo.(type) {
				//case SerializedJsonDebugInfo:
				//	respStr = strings.Replace(respStr, `"`+debugInfoReplaceFlag+`"`, string(typeval), 1)
				//default:
				//
				//}

				debuginfoSerializedJsonBytes, err := json.Marshal(ctxDebugInfo)
				if err != nil {
					return err
				}

				if filters.UseDebugInfoFilter {
					respStr = strings.Replace(respStr, `"`+debugInfoReplaceFlag+`"`, string(debuginfoSerializedJsonBytes), 1)
				}

				resp = []byte(respStr)
			}

		case RESPONSE_FORMAT_HTML:
			resp = []byte(i.(*commonPb.HttpResponseReplyWrapperHTML).Content)
		}

		w.Header().Set("Content-Type", responseFormat+";charset=utf-8")
		//w.Header().Add("Content-Type", "charset=utf-8")
		w.Write(resp)

		SetResponseFormat(RESPONSE_FORMAT_JSON)
		return nil
	}
}

// DefaultRequestDecoder decodes the request body to object.
func HttpRequestDecoder() http.DecodeRequestFunc {
	return func(r *netHttp.Request, v interface{}) error {

		codec, ok := CodecForRequest(r, "Content-Type")
		if !ok {
			return errors.BadRequest("CODEC", r.Header.Get("Content-Type"))
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return errors.BadRequest("CODEC", err.Error())
		}
		if err = codec.Unmarshal(data, v); err != nil {
			return errors.BadRequest("CODEC", err.Error())
		}
		return nil
	}
}

// CodecForRequest get encoding.Codec via http.Request
func CodecForRequest(r *netHttp.Request, name string) (encoding.Codec, bool) {
	for _, accept := range r.Header[name] {
		codec := encoding.GetCodec(ContentSubtype(accept))
		if codec != nil {
			return codec, true
		}
	}
	return encoding.GetCodec("json"), false
}

func ContentSubtype(contentType string) string {
	left := strings.Index(contentType, "/")
	if left == -1 {
		return ""
	}
	right := strings.Index(contentType, ";")
	if right == -1 {
		right = len(contentType)
	}
	if right < left {
		return ""
	}
	return contentType[left+1 : right]
}
