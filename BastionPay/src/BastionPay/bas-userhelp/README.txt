1. 需求设计
   a. 用户的求助信息增加，查询
   b. app的版本信息的设置，查询



2. 接口设计

    用户的求助信息增加，查询
    /v1/user-help/message/add
    /v1/user-help-bk/message/get
    /v1/user-help-bk/message/update
    /v1/user-help-bk/message/del
    /v1/user-help-bk/message/list



    app的版本信息的查询
    /v1/app-version/get

    app的版本信息的设置，查询，更新
    /v1/app-version-bk/add
    /v1/app-version-bk/list
    /v1/app-version-bk/update



3. 注意点

4. 运行
    ./bas-userhelp -conf_path=config.yaml -log_path=zap.conf

