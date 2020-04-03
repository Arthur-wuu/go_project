package api

type EnumID int32

const (
	EVENT_nonesub = "0"
	EVENT_sub     = "1"
	EVENT_unsub   = "2"
	EVENT_SubMax  = "3"

	EnumID_IDNone             EnumID = 0
	EnumID_IDId               EnumID = 1
	EnumID_IDQuoteKlineSingle EnumID = 20
	EnumID_IDMarket           EnumID = 21
	EnumID_IDJianPanBaoShuChu EnumID = 22
)

func NewMSG(id int32) *MSG {
	return new(MSG).SetId(id)
}

type MSG struct {
	Id                      *int32              `json:"Id,omitempty"`
	ResDataQuoteKlineSingle []*QuoteKlineSingle `json:"RepDataQuoteKlineSingle,omitempty"`
	ResDataMarket           []*Market           `json:"RepDataMarket,omitempty"`
	RepDataJianPanBaoShuChu []*JianPanBaoShuChu `json:"RepDataJianPanBaoShuChu,omitempty"`
}

func (this *MSG) AddMarket(elem ...*Market) *MSG {
	if this.Id == nil {
		this.SetId(int32(EnumID_IDMarket))
	}
	if this.ResDataMarket == nil {
		this.ResDataMarket = make([]*Market, 0, len(elem))
	}
	this.ResDataMarket = append(this.ResDataMarket, elem...)
	return this
}

func (this *MSG) AddQuoteKlineSingle(elem ...*QuoteKlineSingle) *MSG {
	if this.Id == nil {
		this.SetId(int32(EnumID_IDQuoteKlineSingle))
	}
	if this.ResDataQuoteKlineSingle == nil {
		this.ResDataQuoteKlineSingle = make([]*QuoteKlineSingle, 0, len(elem))
	}
	this.ResDataQuoteKlineSingle = append(this.ResDataQuoteKlineSingle, elem...)
	return this
}

func (this *MSG) AddJianPanBaoShuChu(elem ...*JianPanBaoShuChu) *MSG {
	if this.Id == nil {
		this.SetId(int32(EnumID_IDJianPanBaoShuChu))
	}
	if this.RepDataJianPanBaoShuChu == nil {
		this.RepDataJianPanBaoShuChu = make([]*JianPanBaoShuChu, 0, len(elem))
	}
	this.RepDataJianPanBaoShuChu = append(this.RepDataJianPanBaoShuChu, elem...)
	return this
}

func (this *MSG) SetId(id int32) *MSG {
	if this.Id == nil {
		this.Id = new(int32)
	}
	*this.Id = id
	return this
}
