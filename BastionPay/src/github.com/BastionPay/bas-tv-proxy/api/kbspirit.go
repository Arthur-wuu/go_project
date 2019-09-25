package api

type JPBLeiXing int32

const (
	JPBLeiXing_TYPE_OBJ        JPBLeiXing = 0
)

var JPBLeiXing_name = map[int32]string{
	0: "TYPE_OBJ",
}
var JPBLeiXing_value = map[string]int32{
	"TYPE_OBJ":        0,
}

func (x JPBLeiXing) Enum() *JPBLeiXing {
	p := new(JPBLeiXing)
	*p = x
	return p
}
func (x JPBLeiXing) String() string {
	return JPBLeiXing_name[int32(x)]
}

type JPBShuJu struct {
	DaiMa            *string `json:"DaiMa,omitempty"`
	MingCheng        *string `json:"MingCheng,omitempty"`
	ShuXing          *string `json:"ShuXing,omitempty"`
	KuoZhan          *string `json:"KuoZhan,omitempty"`
}

func (m *JPBShuJu) Reset()         { *m = JPBShuJu{} }

func (m *JPBShuJu) GetDaiMa() string {
	if m != nil && m.DaiMa != nil {
		return *m.DaiMa
	}
	return ""
}

func (m *JPBShuJu) GetMingCheng() string {
	if m != nil && m.MingCheng != nil {
		return *m.MingCheng
	}
	return ""
}

func (m *JPBShuJu) GetShuXing() string {
	if m != nil && m.ShuXing != nil {
		return *m.ShuXing
	}
	return ""
}

func (m *JPBShuJu) GetKuoZhan() string {
	if m != nil && m.KuoZhan != nil {
		return *m.KuoZhan
	}
	return ""
}

func NewJPBShuChu(leixing JPBLeiXing, shuJu []*JPBShuJu, shiChang, shiChangQuanMing string, updateTime int64) *JPBShuChu {
	return &JPBShuChu{
		LeiXing:&leixing ,
		ShuJu:shuJu,
		ShiChang: &shiChang,
		ShiChangQuanMing:&shiChangQuanMing,
		UpdateTime:&updateTime,
	}
}

type JPBShuChu struct {
	LeiXing          *JPBLeiXing `json:"LeiXing,omitempty"`
	ShiChang         *string     `json:"ShiChang,omitempty"`
	ShiChangQuanMing *string     `json:"ShiChangQuanMing,omitempty"`
	UpdateTime       *int64      `json:"UpdateTime,omitempty"`
	ShuJu            []*JPBShuJu `json:"ShuJu,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *JPBShuChu) Reset()         { *m = JPBShuChu{} }
func (*JPBShuChu) ProtoMessage()    {}

func (m *JPBShuChu) GetLeiXing() JPBLeiXing {
	if m != nil && m.LeiXing != nil {
		return *m.LeiXing
	}
	return JPBLeiXing_TYPE_OBJ
}

func (m *JPBShuChu) GetShuJu() []*JPBShuJu {
	if m != nil {
		return m.ShuJu
	}
	return nil
}

func NewJianPanBaoShuChu(guanjianzi string,JieGuo []*JPBShuChu) *JianPanBaoShuChu {
	return &JianPanBaoShuChu{
		GuanJianZi: &guanjianzi,
		JieGuo: JieGuo,
	}
}
//
type JianPanBaoShuChu struct {
	GuanJianZi       *string      `protobuf:"bytes,1,req" json:"GuanJianZi,omitempty"`
	JieGuo           []*JPBShuChu `protobuf:"bytes,2,rep" json:"JieGuo,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *JianPanBaoShuChu) Reset()         { *m = JianPanBaoShuChu{} }
func (*JianPanBaoShuChu) ProtoMessage()    {}

func (m *JianPanBaoShuChu) GetGuanJianZi() string {
	if m != nil && m.GuanJianZi != nil {
		return *m.GuanJianZi
	}
	return ""
}

func (m *JianPanBaoShuChu) GetJieGuo() []*JPBShuChu {
	if m != nil {
		return m.JieGuo
	}
	return nil
}

