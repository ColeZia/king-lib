由于和具体微服务的`internal/conf/conf.proto`一样，会导致如下报错，故特意将`conf.proto`修改为`fwconf.proto`
```
panic: proto: file "internal/conf/conf.proto" is already registered
        previously from: "boss-pay/internal/conf"
        currently from:  "gl.king.im/king-lib/framework/internal/conf"
See https://developers.google.com/protocol-buffers/docs/reference/go/faq#namespace-conflict
```