
没到消息回调时，先去缓存 查接收方的状态，无则主动查询状态，并缓存。若用户在线，则不做任何处理。若不在线，则判断推送时间及推送消息。
用户不在线，也要控制下推送频率，大概2分钟最多推3条.

在线状态回调，实时更新缓存状态。缓存 有超时时间，比如10分钟，防止腾讯丢消息的情况。



1. 回调接收模块https。
2. 用户状态管理模块。
3. 红点推送模块

