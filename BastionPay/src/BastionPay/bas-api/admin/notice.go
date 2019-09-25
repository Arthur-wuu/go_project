package admin

const (
	STATUS_ALIVE_Online       = 0
	STATUS_ALIVE_PreOnline    = 1
	STATUS_Alive_Offline      = 3 //包含STATUS_ALIVE_PreOnline+STATUS_Alive_AfterOffline
	STATUS_Alive_AfterOffline = 2
)

type NoticeListParam struct {
	TotalLines      int                    `json:"total_lines,omitempty"`
	PageIndex       int                    `json:"page_index,omitempty"`
	MaxDispLines    int                    `json:"max_disp_lines,omitempty"`
	StartCreatedAt  *int64                 `json:"start_created_at,omitempty"`
	EndCreatedAt    *int64                 `json:"end_created_at,omitempty"`
	StartUpdatedAt  *int64                 `json:"start_updated_at,omitempty"`
	EndUpdatedAt    *int64                 `json:"end_updated_at,omitempty"`
	StartOnlinedAt  *int64                 `json:"start_onlined_at,omitempty"`
	EndOnlinedAt    *int64                 `json:"end_onlined_at,omitempty"`
	StartOfflinedAt *int64                 `json:"start_offlined_at,omitempty"`
	EndOfflinedAt   *int64                 `json:"end_offlined_at,omitempty"`
	Language        *string                `json:"language,omitempty"`
	Focus           *bool                  `json:"focus,omitempty"`
	Race            *bool                  `json:"race,omitempty"`
	Title           *string                `json:"title,omitempty"`
	Order           *string                `json:"order,omitempty"`
	Desc            *bool                  `json:"desc,omitempty"`
	Alive           *int                   `json:"alive,omitempty"`
	Id              *uint                  `json:"id,omitempty"`
	Condition       map[string]interface{} `json:"condition,omitempty"`
	Filter          []string               `json:"filter,omitempty"`
}

type NoticeInfo struct {
	Id *uint `json:"id,omitempty"`
	//创建时间
	CreatedAt *int64 `json:"created_at,omitempty"`
	//更新时间
	UpdatedAt *int64 `json:"updated_at,omitempty"`

	OnlinedAt  *int64 `json:"onlined_at,omitempty"`
	OfflinedAt *int64 `json:"offlined_at,omitempty"`
	// 上下线时间
	// 语言
	Language *string `json:"language,omitempty"`
	// 置顶标志
	Focus *bool `json:"focus,omitempty"`
	Race  *bool `json:"race,omitempty"`
	// 标题
	Title  *string `json:"title,omitempty"`
	Author *string `json:"author,omitempty"`
	// 摘要
	Abstract *string `json:"abstract,omitempty"`

	Content *string `json:"content,omitempty"`

	IsRead *bool `json:"readflag,omitempty"`
	Alive  *int  `json:"alive,omitempty"`
}

type ResNoticeList struct {
	Notices []NoticeInfo `json:"notices,omitempty"`

	TotalLines   uint `json:"total_lines, omitempty" `
	PageIndex    uint `json:"page_index, omitempty" `
	MaxDispLines uint `json:"max_disp_lines, omitempty"`
}

type NoticeIdsParam struct {
	Ids []int `json:"id, omitempty"`
}

type CountUserNoticesParam struct {
	ReadFlag        *bool                  `json:"readflag,omitempty"`
	StartCreatedAt  *int64                 `json:"start_created_at,omitempty"`
	EndCreatedAt    *int64                 `json:"end_created_at,omitempty"`
	StartUpdatedAt  *int64                 `json:"start_updated_at,omitempty"`
	EndUpdatedAt    *int64                 `json:"end_updated_at,omitempty"`
	StartOnlinedAt  *int64                 `json:"start_onlined_at,omitempty"`
	EndOnlinedAt    *int64                 `json:"end_onlined_at,omitempty"`
	StartOfflinedAt *int64                 `json:"start_offlined_at,omitempty"`
	EndOfflinedAt   *int64                 `json:"end_offlined_at,omitempty"`
	Language        *string                `json:"language,omitempty"`
	Focus           *bool                  `json:"focus,omitempty"`
	Race            *bool                  `json:"race,omitempty"`
	Alive           *int                   `json:"alive,omitempty"`
	Condition       map[string]interface{} `json:"condition,omitempty"`
}

type ResCountNotices struct {
	AllCount    *uint `json:"allcount,omitempty"`
	ReadCount   *uint `json:"rcount,omitempty"`
	UnReadCount *uint `json:"urcount,omitempty"`
}

//type NoticeInfoParam struct {
//	Id       *uint  `json:"id"`
//	AlivedAt *int64 `json:"alived_at"`
//	// 语言
//	Language *string `json:"language"`
//	//是否下线
//	Alive *bool `json:"alive"`
//	// 置顶标志
//	Focus *bool `json:"focus"`
//	// 标题
//	Title *string `json:"title"`
//	// 摘要
//	Abstract *string `json:"abstract"`
//	// 内容
//	Content *string `json:"content"`
//}

type ResNoticeIdState struct {
	Id     *uint   `json:"id,omitempty"`
	State  *bool   `json:"state,omitempty"`
	ErrMsg *string `json:"errmsg,omitempty"`
}
