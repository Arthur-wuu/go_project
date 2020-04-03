package apibackend

import (
	"fmt"
	"os"
)

/////////////新错误版本，兼容老版本///////////
/*
错误编码规则说明
---------------------------------------------------
公用错误：（基础组件错误，系统错误）
0             成功
100-199： 请求相关错误，一般是由于请求格式的输入错误导致。
200-299： 授权加密安全 账户 相关错误等
300-499： 预先定义的服务内部错误（对外显示）
500-999： 预先定义的服务内部错误(不对外显示错误信息，一些比较low的错误信息请写在这里)
555:      这个是程序Bug错误，必须修复
1000：    系统内部错误
---------------------------------------------------
>1000：开发者自定义错误，能用上面就用上面的，不能用的自行定义，注意不要分配太多错误号，规定一个服务最多分配50个

有些函数返回两种以上可能的错误，就写上出现概率最大的那个错误号
*/

type EnumBasErr int

const BasErrBegin EnumBasErr = 1000000                //区别老版本
const BasMoreErrBegin EnumBasErr = BasErrBegin + 1000 //用户自定义开始

const (
	//系统错误(公用)大家觉得需要自行添加
	BASERR_SUCCESS EnumBasErr = 0 //成功（不可加上起始值）
	//请求相关错误，一般是由于请求格式的输入错误导致。
	BASERR_INVALID_PARAMETER            EnumBasErr = BasErrBegin + 100 //请求参数无效。    一般用在参数缺失或者格式不对
	BASERR_UNSUPPORTED_METHOD           EnumBasErr = BasErrBegin + 101 //未知的方法。     请求的url的path未定义或者其它无法识别的方法
	BASERR_URLBODYTOOLARGE_METHOD       EnumBasErr = BasErrBegin + 102 //请求body过长。   数据长度过长
	BASERR_INVALID_CONTENTLENGTH_METHOD EnumBasErr = BasErrBegin + 103 //无效的ContentLength。这个是http的header中的
	BASERR_QID_ALREADY_EXISTS           EnumBasErr = BasErrBegin + 104 //qid已存在。      这个用在异步请求中，标识每次请求，理解成流水号。

	//授权 加密 账户 相关安全错误等, 200不要了
	BASERR_TOKEN_INVALID          EnumBasErr = BasErrBegin + 201 //无效的访问令牌。  token无效，可与过期混用，貌似目前就admin喝bkadmin用到，单点登录使用。
	BASERR_TOKEN_EXPIRED          EnumBasErr = BasErrBegin + 202 //访问令牌过期。    同上
	BASERR_INCORRECT_SIGNATURE    EnumBasErr = BasErrBegin + 203 //无效的签名。      签名无效
	BASERR_UNAUTHORIZED_METHOD    EnumBasErr = BasErrBegin + 204 //未授权的方法。    该用户方法无权限
	BASERR_UNAUTHORIZED_PARAMETER EnumBasErr = BasErrBegin + 205 //未授权的参数。    同上
	BASERR_ILLEGAL_DATA           EnumBasErr = BasErrBegin + 206 //非法数据，       解密失败或者无法解析

	BASERR_INCORRECT_ACCOUNT_PWD EnumBasErr = BasErrBegin + 220 //账号密码不对
	BASERR_INCORRECT_GA_PWD      EnumBasErr = BasErrBegin + 221 //GA密码错误
	BASERR_INCORRECT_PWD         EnumBasErr = BasErrBegin + 222 //密码错误
	//
	BASERR_WHITELIST_OUTSIDE EnumBasErr = BasErrBegin + 251 //白名单之外
	BASERR_BLACKLIST_INSIDE  EnumBasErr = BasErrBegin + 252 //黑名单内
	BASERR_FROZEN_ACCOUNT    EnumBasErr = BasErrBegin + 253 //账户未激活        账户处于初始状态
	BASERR_FROZEN_METHOD     EnumBasErr = BasErrBegin + 254 //功能未激活        部分功能未开启
	BASERR_BLOCK_ACCOUNT     EnumBasErr = BasErrBegin + 255 //账户被阻，        无法登陆，可能是长期不登录或者注销
	BASERR_DATAFROM_INVALID  EnumBasErr = BasErrBegin + 256 //数据源不对        数据的来源不合法

	//300-399： 预先定义的服务内部错误（对外显示错误信息）
	BASERR_OBJECT_NOT_FOUND      EnumBasErr = BasErrBegin + 301 //指定的对象不存在
	BASERR_OBJECT_EXISTS         EnumBasErr = BasErrBegin + 302 //指定的对象已存在
	BASERR_OBJECT_DATA_NOT_FOUND EnumBasErr = BasErrBegin + 303 //指定的对象数据不存在
	BASERR_ACCOUNT_NOT_FOUND     EnumBasErr = BasErrBegin + 304 //账号不存在
	BASERR_OBJECT_DATA_SAME      EnumBasErr = BasErrBegin + 305 //指定对象数据和上次一致。   这个一般在更新的时候，没必要重复更新。
	BASERR_OBJECT_DATA_NOT_SAME  EnumBasErr = BasErrBegin + 306 //指定对象数据不相同。   这
	BASERR_OBJECT_DATA_NULL      EnumBasErr = BasErrBegin + 307 //指定对象数据为空。   这
	BASERR_OBJECT_ZERO           EnumBasErr = BasErrBegin + 308 //指定对象为0。
	BASERR_OPERATE_FREQUENT      EnumBasErr = BasErrBegin + 311 //操作太频繁，超过了限制。   操作频率限制
	BASERR_INCORRECT_FORMAT      EnumBasErr = BasErrBegin + 312 //格式出错。
	BASERR_INCORRECT_PUBKEY      EnumBasErr = BasErrBegin + 313 //无效的公钥格式。  这个不知是否与上个重复
	BASERR_INCORRECT_FREQUENT    EnumBasErr = BasErrBegin + 314 //错误太频繁。
	BASERR_NOT_ALLOW_STATE       EnumBasErr = BasErrBegin + 315 //当前状态不允许

	//500-999： 预先定义的服务内部错误(不对外显示错误信息，错误信息 统一表示为系统未知错误)
	BASERR_SERVICE_UNKNOWN_ERROR         EnumBasErr = BasErrBegin + 500 //服务未知错误，     一般500到1000之间的错误。$实在不知道写什么错误就用这个$
	BASERR_INTERNAL_SERVICE_ACCESS_ERROR EnumBasErr = BasErrBegin + 501 //内部服务访问错误。  调用另一个服务出错, 或者返回的数据 不合理或者有错
	BASERR_INVALID_OPERATION             EnumBasErr = BasErrBegin + 502 //无效的操作方法。     可能这个方法未开发。 这个绝对是多余了

	BASERR_DATA_PACK_ERROR                 EnumBasErr = BasErrBegin + 510 //数据打包失败，     json或者pb等打包失败
	BASERR_DATA_UNPACK_ERROR               EnumBasErr = BasErrBegin + 511 //数据解包失败，     区别无效参数，这个一般用在请求其他服务的返回数据无法解包
	BASERR_DATABASE_ERROR                  EnumBasErr = BasErrBegin + 512 //数据库操作出错，请重试。 包括redis、mysql等
	BASERR_SERVICE_TEMPORARILY_UNAVAILABLE EnumBasErr = BasErrBegin + 513 //服务暂不可用      $这个应该也是你困惑的时候该选的$
	BASERR_SEND_TIMEOUT                    EnumBasErr = BasErrBegin + 514 //发送超时
	BASERR_STATE_NOT_ALLOW                 EnumBasErr = BasErrBegin + 515 //当前状态不允许

	BASERR_UNKNOWN_BUG EnumBasErr = BasErrBegin + 555 //未知bug请修复，这个明显是bug

	BASERR_INTERNAL_CONFIG_ERROR EnumBasErr = BasErrBegin + 556 //恭喜只是配置错误，不是bug。可以联系运维了。

	BASERR_SYSTEM_INTERNAL_ERROR EnumBasErr = BasErrBegin + 1000 //系统内部错误。      一些内存不足报错啊，什么系统错误啊都可以用这个表示。$这个应该也是你困惑的时候该选的$

	//开发者自定义错误BasMoreErrBegin, 每类应用使用区间50，并且BASERR_ServerName_开头
	//admin和bkadmin 1--49，原则上这些错误的描述信息表现为内部服务错误， 0不可用，与BASERR_SYSTEM_INTERNAL_ERROR 重复
	BASERR_ADMIN_INVALID_VERIFY_STATUS EnumBasErr = BasMoreErrBegin + 1 //无效的验证状态, 上次已验证的东西超过使用次数，或者 已验证的验证码、手机号、邮箱号等不匹配
	BASERR_ADMIN_INCORRECT_VERIFYCODE  EnumBasErr = BasMoreErrBegin + 2 //错误的验证码

	//50-99
	BASERR_BASNOTIFY_AWS_ERR                        EnumBasErr = BasMoreErrBegin + 50
	BASERR_BASNOTIFY_LANCHUANG_ERR                  EnumBasErr = BasMoreErrBegin + 51
	BASERR_BASNOTIFY_TWL_ERR                        EnumBasErr = BasMoreErrBegin + 52
	BASERR_BASNOTIFY_TEMPLATE_DEAD                  EnumBasErr = BasMoreErrBegin + 53
	BASERR_BASNOTIFY_RECIPIENT_EMPTY                EnumBasErr = BasMoreErrBegin + 54
	BASERR_BASNOTIFY_TEMPLATE_PARSE_FAIL            EnumBasErr = BasMoreErrBegin + 55
	BASERR_BASNOTIFY_Nexmo_ERR                      EnumBasErr = BasMoreErrBegin + 56
	BASERR_BASNOTIFY_RongLianYun_ERR                EnumBasErr = BasMoreErrBegin + 57
	BASERR_BASNOTIFY_DingDing_ERR                   EnumBasErr = BasMoreErrBegin + 58
	BASERR_BASNOTIFY_DingDing_QunName_NotSet_ERR    EnumBasErr = BasMoreErrBegin + 59
	BASERR_BASNOTIFY_DingDing_QunName_NotConfig_ERR EnumBasErr = BasMoreErrBegin + 60

	//1000-1049 红包裂变
	BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_NO_ON                   EnumBasErr = BasMoreErrBegin + 1000 //活动未上线
	BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_HAS_OFF                 EnumBasErr = BasMoreErrBegin + 1001 //活动已下线
	BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_NOFOUND                 EnumBasErr = BasMoreErrBegin + 1002 //活动未上线
	BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_RED_ZERO                EnumBasErr = BasMoreErrBegin + 1003 //活动红包抢完
	BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_RED_EXISTS              EnumBasErr = BasMoreErrBegin + 1004 //活动红包已存在
	BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_RED_NOFOUND             EnumBasErr = BasMoreErrBegin + 1005 //活动红包已存在
	BASERR_ACTIVITY_FISSIONSHARE_RED_ROB_ZERO                     EnumBasErr = BasMoreErrBegin + 1006 //红包名额抢完
	BASERR_ACTIVITY_FISSIONSHARE_RED_TIMEOUT                      EnumBasErr = BasMoreErrBegin + 1007 //红包过期
	BASERR_ACTIVITY_FISSIONSHARE_RED_ROBBER_EXISTS                EnumBasErr = BasMoreErrBegin + 1008 //红包名额已存在
	BASERR_ACTIVITY_FISSIONSHARE_ROBBER_TIMEOUT                   EnumBasErr = BasMoreErrBegin + 1009 //红包抢到已过期
	BASERR_ACTIVITY_FISSIONSHARE_ROBBER_NOFOUND                   EnumBasErr = BasMoreErrBegin + 1010 //红包抢到已过期
	BASERR_ACTIVITY_FISSIONSHARE_INCORRECT_PRECISION              EnumBasErr = BasMoreErrBegin + 1011 //活动精度设置有误
	BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_INVALID                 EnumBasErr = BasMoreErrBegin + 1012 //活动无效
	BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_RED_INVALID             EnumBasErr = BasMoreErrBegin + 1013 //活动红包无效
	BASERR_ACTIVITY_FISSIONSHARE_SMS_INCORRECT_VERIFYCODE         EnumBasErr = BasMoreErrBegin + 1014 //错误的验证码
	BASERR_ACTIVITY_FISSIONSHARE_SMS_VERIFYID_INVALID             EnumBasErr = BasMoreErrBegin + 1015 //无效验证ID
	BASERR_ACTIVITY_FISSIONSHARE_ILLEGAL_APIKEY                   EnumBasErr = BasMoreErrBegin + 1016 //非法APIKEY
	BASERR_ACTIVITY_FISSIONSHARE_ROBBER_TRANSFERFLAG_NOT_AFFECTED EnumBasErr = BasMoreErrBegin + 1017 //转账状态未改变
	BASERR_ACTIVITY_FISSIONSHARE_RED_ROB_ZERO_AND_ROBED           EnumBasErr = BasMoreErrBegin + 1018 //红包抢完，并且本人抢到
	BASERR_ACTIVITY_FISSIONSHARE_SPONSOR_NOFOUND                  EnumBasErr = BasMoreErrBegin + 1019 //未找到营销号
	BASERR_ACTIVITY_FISSIONSHARE_ILLEGAL_IP                       EnumBasErr = BasMoreErrBegin + 1020 //未找到营销号

	//1050-1099 //抽奖
	BASERR_MARKETING_LUCKDRAW_SPONSOR_NOFOUND                  EnumBasErr = BasMoreErrBegin + 1050 //未找到营销号
	BASERR_MARKETING_LUCKDRAW_SPONSOR_INVALID                  EnumBasErr = BasMoreErrBegin + 1051 //未找到营销号
	BASERR_MARKETING_LUCKDRAW_ILLEGAL_IP                       EnumBasErr = BasMoreErrBegin + 1052 //非法IP
	BASERR_MARKETING_LUCKDRAW_ILLEGAL_APIKEY                   EnumBasErr = BasMoreErrBegin + 1053 //非法APIKEY
	BASERR_MARKETING_LUCKDRAW_SMS_INCORRECT_VERIFYCODE         EnumBasErr = BasMoreErrBegin + 1054 //错误的验证码
	BASERR_MARKETING_LUCKDRAW_SMS_VERIFYID_INVALID             EnumBasErr = BasMoreErrBegin + 1055 //无效验证ID
	BASERR_MARKETING_LUCKDRAW_ACTIVITY_NO_ON                   EnumBasErr = BasMoreErrBegin + 1060 //活动未上线
	BASERR_MARKETING_LUCKDRAW_ACTIVITY_HAS_OFF                 EnumBasErr = BasMoreErrBegin + 1061 //活动已下线
	BASERR_MARKETING_LUCKDRAW_ACTIVITY_NOFOUND                 EnumBasErr = BasMoreErrBegin + 1062 //活动未找到
	BASERR_MARKETING_LUCKDRAW_ACTIVITY_INVALID                 EnumBasErr = BasMoreErrBegin + 1063 //活动无效
	BASERR_MARKETING_LUCKDRAW_INCORRECT_PRECISION              EnumBasErr = BasMoreErrBegin + 1064 //活动精度设置有误
	BASERR_MARKETING_LUCKDRAW_AWARD_EXISTS                     EnumBasErr = BasMoreErrBegin + 1070 //活动奖品已存在
	BASERR_MARKETING_LUCKDRAW_AWARD_NOFOUND                    EnumBasErr = BasMoreErrBegin + 1071 //活动奖品未找到
	BASERR_MARKETING_LUCKDRAW_AWARD_INVALID                    EnumBasErr = BasMoreErrBegin + 1072 //活动精度设置有误
	BASERR_MARKETING_LUCKDRAW_DRAWER_ZERO                      EnumBasErr = BasMoreErrBegin + 1075 //名额抢完
	BASERR_MARKETING_LUCKDRAW_DRAWER_EXISTS                    EnumBasErr = BasMoreErrBegin + 1076 //已抽过了
	BASERR_MARKETING_LUCKDRAW_DRAWER_TIMEOUT                   EnumBasErr = BasMoreErrBegin + 1077 //抽到的奖品过期
	BASERR_MARKETING_LUCKDRAW_DRAWER_NOFOUND                   EnumBasErr = BasMoreErrBegin + 1078 //未找到找到的奖品
	BASERR_MARKETING_LUCKDRAW_DRAWER_INVALID                   EnumBasErr = BasMoreErrBegin + 1079 //奖品无效
	BASERR_MARKETING_LUCKDRAW_DRAWER_TRANSFERFLAG_NOT_AFFECTED EnumBasErr = BasMoreErrBegin + 1080 //转账状态未改变
	BASERR_MARKETING_LUCKDRAW_ACTIVITY_HAS_ON                  EnumBasErr = BasMoreErrBegin + 1081 //活动已开始

)

// cobank定义的错误码
const (
	ErrorFailed = 20000 + iota
	ErrorUserInvalid
	ErrorUserFrozen
	ErrorUserNoAuthority
	ErrorParseDataFailed
	ErrorParameter
	ErrorAsssetNotOnline
	ErrorRepeatOrder
	ErrorNotEnoughBalance
	ErrorAccountAbnormal
	ErrorCreateOrder
	ErrorSendDataFailed
	ErrorNotFindPayAddress
)

var EnumBasErr_desc = map[EnumBasErr]string{
	BASERR_SUCCESS: "Success",
	//
	BASERR_INVALID_PARAMETER:            "InvalidParameter",
	BASERR_UNSUPPORTED_METHOD:           "UnsupportedMethod",
	BASERR_URLBODYTOOLARGE_METHOD:       "RequestBodyTooLarge",
	BASERR_INVALID_CONTENTLENGTH_METHOD: "ContentLengthInvalid",
	BASERR_QID_ALREADY_EXISTS:           "QidAlreadyExists",
	//
	BASERR_TOKEN_INVALID:          "TokenInvalidOrNoLongerValid",
	BASERR_TOKEN_EXPIRED:          "TokenExpired",
	BASERR_INCORRECT_SIGNATURE:    "IncorrectSignature",
	BASERR_UNAUTHORIZED_METHOD:    "UnauthorizedMethod",
	BASERR_UNAUTHORIZED_PARAMETER: "UnauthorizedParameter",
	BASERR_ILLEGAL_DATA:           "IllegalData",

	BASERR_INCORRECT_ACCOUNT_PWD: "IncorrectAccountOrPassword",
	BASERR_INCORRECT_GA_PWD:      "IncorrectGoogleAuthenticator",
	BASERR_INCORRECT_PWD:         "IncorrectPassword",

	BASERR_WHITELIST_OUTSIDE: "WhiteListOutside",
	BASERR_BLACKLIST_INSIDE:  "BlackListInside",
	BASERR_FROZEN_ACCOUNT:    "AccountNoActived",
	BASERR_FROZEN_METHOD:     "MethodNoActived",
	BASERR_DATAFROM_INVALID:  "DataFromInvalid",
	BASERR_BLOCK_ACCOUNT:     "AccountBlocked",

	//
	BASERR_OBJECT_NOT_FOUND:      "SpecifiedObjectCannotBeFound",
	BASERR_OBJECT_EXISTS:         "SpecifiedObjectAlreadyExists",
	BASERR_OBJECT_DATA_NOT_FOUND: "SpecifiedObjectDataCannotBeFound",
	BASERR_ACCOUNT_NOT_FOUND:     "AccountNotFound",
	BASERR_OBJECT_DATA_SAME:      "ObjectDataSameAsOld",
	BASERR_OBJECT_DATA_NULL:      "ObjectDataNull",
	BASERR_OBJECT_ZERO:           "ObjectIsZero",
	BASERR_OPERATE_FREQUENT:      "OperateFrequentLimit",
	BASERR_INCORRECT_FORMAT:      "IncorrectFormat",
	BASERR_INCORRECT_PUBKEY:      "IncorrectPublicKey",
	BASERR_INCORRECT_FREQUENT:    "IncorrectFrequentLimit",
	BASERR_NOT_ALLOW_STATE:       "NotAllowState",

	//
	BASERR_SERVICE_UNKNOWN_ERROR:           "ServiceUnknownError",
	BASERR_INTERNAL_SERVICE_ACCESS_ERROR:   "InternalServiceAccessError",
	BASERR_INVALID_OPERATION:               "InvalidOperation",
	BASERR_DATA_PACK_ERROR:                 "DataPackError",
	BASERR_DATA_UNPACK_ERROR:               "DataUnpackError",
	BASERR_DATABASE_ERROR:                  "DatabaseErrorOccurredPleaseTryAgain",
	BASERR_SERVICE_TEMPORARILY_UNAVAILABLE: "ServiceTemporarilyUnavailable",
	BASERR_SEND_TIMEOUT:                    "SendTimeout",
	BASERR_STATE_NOT_ALLOW:                 "StateNotAllow",
	//
	BASERR_UNKNOWN_BUG:           "UnknownBugPleaseContactDeveloper",
	BASERR_INTERNAL_CONFIG_ERROR: "InternalConfigErrorPleaseContactDevops",
	//
	BASERR_SYSTEM_INTERNAL_ERROR: "SystemInternalError",
	//
	BASERR_ADMIN_INVALID_VERIFY_STATUS: "InvalidVerifyStatus",
	BASERR_ADMIN_INCORRECT_VERIFYCODE:  "IncorrectVerifyCode",

	//短信邮件模板程序
	BASERR_BASNOTIFY_AWS_ERR:                        "BasNotifyAwsErr",
	BASERR_BASNOTIFY_LANCHUANG_ERR:                  "BasNotifyLanChuangErr",
	BASERR_BASNOTIFY_TWL_ERR:                        "BasNotifyTwiliErr",
	BASERR_BASNOTIFY_TEMPLATE_DEAD:                  "BasNotifyTemplateDead",
	BASERR_BASNOTIFY_RECIPIENT_EMPTY:                "BasNotifyRecipientEmpty",
	BASERR_BASNOTIFY_TEMPLATE_PARSE_FAIL:            "BasNotifyTemplateParseFail",
	BASERR_BASNOTIFY_Nexmo_ERR:                      "BasNotifyNexmoErr",
	BASERR_BASNOTIFY_RongLianYun_ERR:                "BasNotifyRongLianYunErr",
	BASERR_BASNOTIFY_DingDing_ERR:                   "BasNotifyDingDingErr",
	BASERR_BASNOTIFY_DingDing_QunName_NotSet_ERR:    "BasNotifyDingDingQunNameNotSetErr",
	BASERR_BASNOTIFY_DingDing_QunName_NotConfig_ERR: "BasNotifyDingDingQunNameNotConfigErr",

	//裂变活动
	BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_NO_ON:                   "ActivityFissionShareActivityNotOn",
	BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_HAS_OFF:                 "ActivityFissionShareActivityHasOff",
	BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_RED_ZERO:                "ActivityFissionShareActivityRedZero",
	BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_NOFOUND:                 "ActivityFissionShareActivityNoFound",
	BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_RED_EXISTS:              "ActivityFissionShareActivityRedExists",
	BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_RED_NOFOUND:             "ActivityFissionShareActivityRedNoFound",
	BASERR_ACTIVITY_FISSIONSHARE_RED_ROB_ZERO:                     "ActivityFissionShareRedRobZero",
	BASERR_ACTIVITY_FISSIONSHARE_RED_TIMEOUT:                      "ActivityFissionShareRedTimeout",
	BASERR_ACTIVITY_FISSIONSHARE_RED_ROBBER_EXISTS:                "ActivityFissionShareRedRobberExists",
	BASERR_ACTIVITY_FISSIONSHARE_ROBBER_TIMEOUT:                   "ActivityFissionShareRobberTimeout",
	BASERR_ACTIVITY_FISSIONSHARE_ROBBER_NOFOUND:                   "ActivityFissionShareRobberNoFound",
	BASERR_ACTIVITY_FISSIONSHARE_INCORRECT_PRECISION:              "ActivityFissionShareIncorrectPrecision",
	BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_INVALID:                 "ActivityFissionShareActivityInvalid",
	BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_RED_INVALID:             "ActivityFissionShareActivityRedInvalid",
	BASERR_ACTIVITY_FISSIONSHARE_SMS_INCORRECT_VERIFYCODE:         "ActivityFissionShareSmsIncorrectVerifyCode",
	BASERR_ACTIVITY_FISSIONSHARE_SMS_VERIFYID_INVALID:             "ActivityFissionShareSmsVerifyIdInvalid",
	BASERR_ACTIVITY_FISSIONSHARE_ILLEGAL_APIKEY:                   "ActivityFissionShareIllegalApiKey",
	BASERR_ACTIVITY_FISSIONSHARE_ROBBER_TRANSFERFLAG_NOT_AFFECTED: "ActivityFissionShareRobberTransferFlagNotAffected",
	BASERR_ACTIVITY_FISSIONSHARE_RED_ROB_ZERO_AND_ROBED:           "ActivityFissionShareRedRobZeroAndSelfRobed",
	BASERR_ACTIVITY_FISSIONSHARE_SPONSOR_NOFOUND:                  "ActivityFissionShareSponsorNoFound",
	BASERR_ACTIVITY_FISSIONSHARE_ILLEGAL_IP:                       "ActivityFissionShareIllegalIp",

	//抽奖活动
	BASERR_MARKETING_LUCKDRAW_SPONSOR_NOFOUND:                  "MarketingLuckDrawSponsorNoFound",
	BASERR_MARKETING_LUCKDRAW_SPONSOR_INVALID:                  "MarketingLuckDrawSponsorInvalid",
	BASERR_MARKETING_LUCKDRAW_ILLEGAL_IP:                       "MarketingLuckDrawIllegalIp",
	BASERR_MARKETING_LUCKDRAW_ILLEGAL_APIKEY:                   "MarketingLuckDrawIllegalApikey",
	BASERR_MARKETING_LUCKDRAW_SMS_INCORRECT_VERIFYCODE:         "MarketingLuckDrawSmsIncorrectVerifycode",
	BASERR_MARKETING_LUCKDRAW_SMS_VERIFYID_INVALID:             "MarketingLuckDrawSmsVerifyIdInvalid",
	BASERR_MARKETING_LUCKDRAW_ACTIVITY_NO_ON:                   "MarketingLuckDrawActivityNoOn",
	BASERR_MARKETING_LUCKDRAW_ACTIVITY_HAS_OFF:                 "MarketingLuckDrawActivityHasOff",
	BASERR_MARKETING_LUCKDRAW_ACTIVITY_NOFOUND:                 "MarketingLuckDrawActivityNoFound",
	BASERR_MARKETING_LUCKDRAW_AWARD_EXISTS:                     "MarketingLuckDrawAwardExists",
	BASERR_MARKETING_LUCKDRAW_AWARD_NOFOUND:                    "MarketingLuckDrawAwardNoFound",
	BASERR_MARKETING_LUCKDRAW_ACTIVITY_INVALID:                 "MarketingLuckDrawActivityInValid",
	BASERR_MARKETING_LUCKDRAW_INCORRECT_PRECISION:              "MarketingLuckDrawIncorrectPrecision",
	BASERR_MARKETING_LUCKDRAW_DRAWER_ZERO:                      "MarketingLuckDrawDrawerZero",
	BASERR_MARKETING_LUCKDRAW_DRAWER_EXISTS:                    "MarketingLuckDrawDrawerExists",
	BASERR_MARKETING_LUCKDRAW_DRAWER_TIMEOUT:                   "MarketingLuckDrawDrawerTimeout",
	BASERR_MARKETING_LUCKDRAW_DRAWER_NOFOUND:                   "MarketingLuckDrawDrawerNoFound",
	BASERR_MARKETING_LUCKDRAW_DRAWER_INVALID:                   "MarketingLuckDrawDrawerInvalid",
	BASERR_MARKETING_LUCKDRAW_DRAWER_TRANSFERFLAG_NOT_AFFECTED: "MarketingLuckDrawDrawerTransferflagNotAffected",
	BASERR_MARKETING_LUCKDRAW_AWARD_INVALID:                    "MarketingLuckDrawAwardInvalid",
}

func (x EnumBasErr) Code() int {
	return int(x)
}

func (x EnumBasErr) Desc() string {
	var key EnumBasErr
	switch {
	case x >= BASERR_SERVICE_UNKNOWN_ERROR && x < BASERR_SYSTEM_INTERNAL_ERROR:
		key = BASERR_SERVICE_UNKNOWN_ERROR
	case x >= BASERR_SYSTEM_INTERNAL_ERROR:
		key = BASERR_SYSTEM_INTERNAL_ERROR
	default:
		key = x
	}
	desc, ok := EnumBasErr_desc[key]
	if !ok && key != BASERR_SERVICE_UNKNOWN_ERROR {
		desc, _ = EnumBasErr_desc[BASERR_SERVICE_UNKNOWN_ERROR]
	}
	return desc
}

func (x EnumBasErr) OriginDesc() string {
	desc, ok := EnumBasErr_desc[x]
	if !ok {
		desc, _ = EnumBasErr_desc[BASERR_SERVICE_UNKNOWN_ERROR]
	}
	return desc
}

func (x EnumBasErr) EscapedCode() int {
	var code int
	if x == BASERR_SUCCESS {
		code = int(x)
	} else if x >= BasErrBegin {
		code = int(x - BasErrBegin)
	} else {
		code = int(x)
	}
	return code
}

func (x EnumBasErr) String() string {
	return fmt.Sprintf("{\"code\":%d,\"desc\":\"%s\"}", x.Code(), x.Desc())
}

func (x EnumBasErr) Error() string {
	return x.String()
}

/////////////老错误版本////////////////////////////
// error code and message
type ErrorInfo struct {
	Code   int
	Msg    string
	Groups []string
}

var (
	err_msg map[int]*ErrorInfo
)

func AddErrMsg(errId int, errMsg string, groups []string) {
	if _, ok := err_msg[errId]; ok {
		fmt.Printf("Error code %d exist!", errId)
		os.Exit(1)
	}
	err_msg[errId] = &ErrorInfo{errId, errMsg, groups}
}

func GetErrMsg(errId int) string {
	if msgInfo, ok := err_msg[errId]; ok {
		return msgInfo.Msg
	}
	return "service internal error"
}

func GetGroupErrMsg(group string) map[int]*ErrorInfo {
	if group == "" {
		return err_msg
	}

	errMsgs := make(map[int]*ErrorInfo)
	for _, v := range err_msg {
		for _, g := range v.Groups {
			if g == group {
				errMsgs[v.Code] = v
				break
			}
		}
	}

	return errMsgs
}

func init() {
	err_msg = make(map[int]*ErrorInfo)

	AddErrMsg(NoErr, "NoErr", []string{HttpRouterApi, HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrInternal, "ErrInternal", []string{HttpRouterApi, HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrDataCorrupted, "ErrDataCorrupted", []string{HttpRouterApi, HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrCallFailed, "ErrCallFailed", []string{HttpRouterApi, HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrIllegallyCall, "ErrIllegallyCall", []string{HttpRouterApi, HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrNotFindAuth, "ErrNotFindAuth", []string{HttpRouterApi, HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrNotFindSrv, "ErrNotFindSrv", []string{HttpRouterApi, HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrNotFindFunction, "ErrNotFindFunction", []string{HttpRouterApi, HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrConnectSrvFailed, "ErrConnectSrvFailed", []string{HttpRouterApi, HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrRequestInvalid, "ErrRequestInvalid", []string{HttpRouterApi, HttpRouterUser, HttpRouterAdmin})

	AddErrMsg(ErrAccountSrvNoUser, "ErrAccountSrvNoUser", []string{HttpRouterApi, HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrAccountSrvUpdateProfile, "ErrAccountSrvUpdateProfile", []string{HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrAccountSrvListUsers, "ErrAccountSrvListUsers", []string{HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrAccountSrvListUsersCount, "ErrAccountSrvListUsersCount", []string{HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrAccountPubKeyParse, "ErrAccountPubKeyParse", []string{HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrAccountSrvSetFrozen, "ErrAccountSrvSetFrozen", []string{HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrAccountSrvAudite, "ErrAccountSrvAudite", []string{HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrAccountSrvBindChildUser, "ErrAccountSrvBindChildUser", []string{HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrAccountSrvListUserCustodianInfo, "ErrAccountSrvListUserCustodianInfo", []string{HttpRouterAdmin})
	AddErrMsg(ErrAccountSrvUserAssetCustodianInfoAlreadyExist, "ErrAccountSrvUserAssetCustodianInfoAlreadyExist", []string{HttpRouterAdmin})
	AddErrMsg(ErrAccountSrvUserAssetCustodianInfoNotExist, "ErrAccountSrvUserAssetCustodianInfoNotExist", []string{HttpRouterAdmin})
	AddErrMsg(ErrAccountSrvUserAssetCustodianFeeLimit, "ErrAccountSrvUserAssetCustodianFeeLimit", []string{HttpRouterAdmin})
	AddErrMsg(ErrAccountSrvAddUserAssetCustodianInfo, "ErrAccountSrvAddUserAssetCustodianInfo", []string{HttpRouterAdmin})
	AddErrMsg(ErrAccountSrvUpdateUserAssetCustodianInfo, "ErrAccountSrvUpdateUserAssetCustodianInfo", []string{HttpRouterAdmin})
	AddErrMsg(ErrAccountSrvDelUserAssetCustodianInfo, "ErrAccountSrvDelUserAssetCustodianInfo", []string{HttpRouterAdmin})

	AddErrMsg(ErrAuthSrvNoUserKey, "ErrAuthSrvNoUserKey", []string{HttpRouterApi, HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrAuthSrvNoPublicKey, "ErrAuthSrvNoPublicKey", []string{HttpRouterApi, HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrAuthSrvNoApiLevel, "ErrAuthSrvNoApiLevel", []string{HttpRouterApi, HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrAuthSrvUserFrozen, "ErrAuthSrvUserFrozen", []string{HttpRouterApi, HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrAuthSrvIllegalData, "ErrAuthSrvIllegalData", []string{HttpRouterApi, HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrAuthSrvIllegalDataType, "ErrAuthSrvIllegalDataType", []string{HttpRouterUser, HttpRouterAdmin})
	AddErrMsg(ErrAuthSrvIllegalAudite, "ErrAuthSrvIllegalAudite", []string{HttpRouterApi})
	AddErrMsg(ErrAuthSrvIllegalIp, "ErrAuthSrvIllegalIp", []string{HttpRouterApi, HttpRouterUser, HttpRouterAdmin})

	AddErrMsg(ErrPushSrvPushData, "ErrPushSrvPushData", []string{HttpRouterUser, HttpRouterAdmin})

	// cobank错误码信息
	AddErrMsg(ErrorFailed, "ErrorFailed", []string{HttpRouterApi, HttpRouterApi})
	AddErrMsg(ErrorUserInvalid, "ErrorUserInvalid", []string{HttpRouterApi, HttpRouterUser})
	AddErrMsg(ErrorUserFrozen, "ErrorUserFrozen", []string{HttpRouterApi, HttpRouterUser})
	AddErrMsg(ErrorUserNoAuthority, "ErrorUserNoAuthority", []string{HttpRouterApi, HttpRouterUser})
	AddErrMsg(ErrorParseDataFailed, "ErrorParseDataFailed", []string{HttpRouterApi, HttpRouterUser})
	AddErrMsg(ErrorParameter, "ErrorParameter", []string{HttpRouterApi, HttpRouterUser})
	AddErrMsg(ErrorAsssetNotOnline, "ErrorAsssetNotOnline", []string{HttpRouterApi, HttpRouterUser})
	AddErrMsg(ErrorRepeatOrder, "ErrorRepeatOrder", []string{HttpRouterApi, HttpRouterUser})
	AddErrMsg(ErrorNotEnoughBalance, "ErrorNotEnoughBalance", []string{HttpRouterApi, HttpRouterUser})
	AddErrMsg(ErrorAccountAbnormal, "ErrorAccountAbnormal", []string{HttpRouterApi, HttpRouterUser})
	AddErrMsg(ErrorCreateOrder, "ErrorCreateOrder", []string{HttpRouterApi, HttpRouterUser})
	AddErrMsg(ErrorSendDataFailed, "ErrorSendDataFailed", []string{HttpRouterApi, HttpRouterUser})
	AddErrMsg(ErrorNotFindPayAddress, "ErrorNotFindPayAddress", []string{HttpRouterUser, HttpRouterAdmin})

	for code, desc := range EnumBasErr_desc {
		if code == BASERR_SUCCESS {
			continue
		}
		AddErrMsg(int(code), desc, []string{HttpRouterUser, HttpRouterAdmin})
	}
	for _, errInfo := range err_msg {
		if errInfo.Code == NoErr {
			continue
		}
		EnumBasErr_desc[EnumBasErr(errInfo.Code)] = errInfo.Msg
	}
}

const (
	// /////////////////////////////////////////////////////
	// 0, success
	// /////////////////////////////////////////////////////
	// no error
	NoErr = 0

	// /////////////////////////////////////////////////////
	// 10001-11100 common errors
	// /////////////////////////////////////////////////////
	// internal err
	ErrInternal = 10001

	// data corrupted
	ErrDataCorrupted = 10002

	// call failed
	ErrCallFailed = 10003

	// illegally call
	ErrIllegallyCall = 10004

	// not find auth service
	ErrNotFindAuth = 10005

	// not find service
	ErrNotFindSrv = 10006

	// not find function
	ErrNotFindFunction = 10007

	// connect service failed
	ErrConnectSrvFailed = 10008

	// request invalid
	ErrRequestInvalid = 10009

	// /////////////////////////////////////////////////////
	// 11101-11200 account_srv errors
	// /////////////////////////////////////////////////////
	// no user
	ErrAccountSrvNoUser = 11101

	// updateprofile - failed
	ErrAccountSrvUpdateProfile = 11102

	// listusers - failed
	ErrAccountSrvListUsers = 11103

	// listusers count - failed
	ErrAccountSrvListUsersCount = 11104

	// pub key parse
	ErrAccountPubKeyParse = 11105

	// set user frozen failed
	ErrAccountSrvSetFrozen = 11106

	//
	ErrAccountSrvAudite = 11107

	//
	ErrAccountSrvBindChildUser = 11108

	// custodian info
	ErrAccountSrvListUserCustodianInfo              = 11109
	ErrAccountSrvUserAssetCustodianInfoAlreadyExist = 11110
	ErrAccountSrvUserAssetCustodianInfoNotExist     = 11111
	ErrAccountSrvUserAssetCustodianFeeLimit         = 11112
	ErrAccountSrvAddUserAssetCustodianInfo          = 11113
	ErrAccountSrvUpdateUserAssetCustodianInfo       = 11114
	ErrAccountSrvDelUserAssetCustodianInfo          = 11115

	// /////////////////////////////////////////////////////
	// 11201-11300 auth_srv errors
	// /////////////////////////////////////////////////////
	// no user key
	ErrAuthSrvNoUserKey = 11201

	// no public key
	ErrAuthSrvNoPublicKey = 11202

	// no api level
	ErrAuthSrvNoApiLevel = 11203

	// user frozen
	ErrAuthSrvUserFrozen = 11204

	// illegal data
	ErrAuthSrvIllegalData = 11205

	// illegal data type
	ErrAuthSrvIllegalDataType = 11206

	//Audite_status
	ErrAuthSrvIllegalAudite = 11207

	//illegal SourceIp
	ErrAuthSrvIllegalIp = 11208

	// /////////////////////////////////////////////////////
	// 11301-11400 push_srv errors
	// /////////////////////////////////////////////////////
	// illegal data
	ErrPushSrvPushData = 11301
)
