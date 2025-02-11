package rbd

// EnvClusterMap 环境到集群ID的映射表
var EnvClusterMap = map[string]string{
	"store-test":      "da2fad090772456ab58a49c3c35cd2ee",
	"ent3-store":      "f793e5c13b33484fa8f08fe7bc6483c9",
	"ent4-store":      "e2691a8c566d46769c0a07742ccb8c6d",
	"ent3-enterprise": "408d8b5529ed444787d4ddcb1e8fd092",
	"ent4-enterprise": "440e1333bffc43109cb977a2fb2b2a77",
}

// EnvTokenMap 环境token映射表
var EnvTokenMap = map[string]string{
	"store-test":      "6a1db07534f1d4df8fff95deff376ea1d49f069b",
	"ent3-store":      "679fe72254d8da4cc22c253a2c123d186b5b407e",
	"ent4-store":      "925282dd80799de8870b64fbf3275dfd8ea46bc2",
	"ent3-enterprise": "89197686aa66a43b8e92ff3a07647f6af99e211c",
	"ent4-enterprise": "52e449875a850330c4808df92a1626bf2f2ad9d6",
}

// EnvApiURLMap 环境apiURL映射表
var EnvApiURLMap = map[string]string{
	"store-test":      "http://172.21.14.149:7070",
	"ent4-store":      "https://rbdwg.hwwt2.com",
	"ent3-store":      "https://rbdent3.hwwt2.com",
	"ent3-enterprise": "https://etrbd-prd-tc.hwwt2.com",
	"ent4-enterprise": "https://etrbd-prd.hwwt2.com",
}

// USER APIS

// ApiPathMap 存储不同API功能对应的路径
var ApiPathMap = map[string]string{
	"user_create": "/openapi/v1/users",
}
