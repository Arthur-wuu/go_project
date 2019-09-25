package controllers

import (
	"github.com/kataras/iris"
	"BastionPay/pay-user-merchant-api/common"
	"BastionPay/pay-user-merchant-api/models"
	"go.uber.org/zap"
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/pay-user-merchant-api/api"
)


type Merchant struct{
	Controllers
}

func (this *Merchant) Create(ctx iris.Context) {
	userInfo,err := GetAppUserInfo(ctx)
	userId := *userInfo.UserId
	if userId == 0 {
		ZapLog().With(zap.Int64("userId", userId)).Error("userId is 0 err")
		this.ExceptionSerive(ctx, apibackend.BASERR_UNKNOWN_BUG.Code(), apibackend.BASERR_UNKNOWN_BUG.Desc())
		return
	}

	mer,err := new(models.Merchant).GetByPayeeId(userId)
	if err != nil {
		ZapLog().With(zap.Error(err), zap.Int64("userId", userId)).Error("GetByPayeeId err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}
	if mer == nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), apibackend.BASERR_OBJECT_NOT_FOUND.Desc())
		return
	}
	this.Response(ctx, &api.ResMerchant{
		ID: mer.ID,
		MerchantId: mer.MerchantId,
		MerchantName: mer.MerchantName,
		NotifyUrl: mer.NotifyUrl,
		SignType: mer.SignType,
		SignKey:  mer.SignKey,
		LanguageType: mer.LanguageType,
		LegalCurrency: mer.LegalCurrency,
		Contact:       mer.Contact,
		ContactPhone :  mer.ContactPhone,
		ContactEmail :  mer.ContactEmail,
		Country  :      mer.Country,
		CreateTime     : mer.CreateTime,
		LastUpdateTime : mer.LastUpdateTime,
	})
}

func (this *Merchant) Get(ctx iris.Context) {
	userInfo,err := GetAppUserInfo(ctx)
	userId := *userInfo.UserId
	if userId == 0 {
		ZapLog().With(zap.Int64("userId", userId)).Error("userId is 0 err")
		this.ExceptionSerive(ctx, apibackend.BASERR_UNKNOWN_BUG.Code(), apibackend.BASERR_UNKNOWN_BUG.Desc())
		return
	}

	mer,err := new(models.Merchant).GetByPayeeId(userId)
	if err != nil {
		ZapLog().With(zap.Error(err), zap.Int64("userId", userId)).Error("GetByPayeeId err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}
	if mer == nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), apibackend.BASERR_OBJECT_NOT_FOUND.Desc())
		return
	}
	this.Response(ctx, &api.ResMerchant{
		ID: mer.ID,
		MerchantId: mer.MerchantId,
		MerchantName: mer.MerchantName,
		NotifyUrl: mer.NotifyUrl,
		SignType: mer.SignType,
		SignKey:  mer.SignKey,
		LanguageType: mer.LanguageType,
		LegalCurrency: mer.LegalCurrency,
		Contact:       mer.Contact,
		ContactPhone :  mer.ContactPhone,
		ContactEmail :  mer.ContactEmail,
		Country  :      mer.Country,
		CreateTime     : mer.CreateTime,
		LastUpdateTime : mer.LastUpdateTime,
	})
}

func (this *Merchant) Add(ctx iris.Context) {
	userId := int64(common.GetUserIdFromCtx(ctx))
	if userId == 0 {
		ZapLog().With(zap.Int64("userId", userId)).Error("userId is 0 err")
		this.ExceptionSerive(ctx, apibackend.BASERR_UNKNOWN_BUG.Code(), apibackend.BASERR_UNKNOWN_BUG.Desc())
		return
	}
	param := new(api.MerchantAdd)
	if err := Tools.ShouldBindJSON(ctx, param); err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		return
	}
	param.PayeeId = &userId
	modelParam:= new(models.Merchant).ParseAdd(param)

	//商户ID 用随机字符
	genMerchantId := GetRandomString(8)
	modelParam.MerchantId = &genMerchantId

	uniFlag,err := modelParam.Unique()
	if err != nil {
		ZapLog().With(zap.Error(err), zap.Int64("userId", userId)).Error("Unique err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}
	if !uniFlag {
		ZapLog().With(zap.Int64("userId", userId)).Error("Unique err")
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_EXISTS.Code(), apibackend.BASERR_OBJECT_EXISTS.Desc())
		return
	}
	mer,err := modelParam.Add()
	if err != nil {
		ZapLog().With(zap.Error(err), zap.Int64("userId", userId)).Error("Add err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}
	this.Response(ctx, &api.ResMerchant{
		ID: mer.ID,
		MerchantId: mer.MerchantId,
		MerchantName: mer.MerchantName,
		NotifyUrl: mer.NotifyUrl,
		SignType: mer.SignType,
		SignKey:  mer.SignKey,
		LanguageType: mer.LanguageType,
		LegalCurrency: mer.LegalCurrency,
		Contact:       mer.Contact,
		ContactPhone :  mer.ContactPhone,
		ContactEmail :  mer.ContactEmail,
		Country  :      mer.Country,
		CreateTime     : mer.CreateTime,
		LastUpdateTime : mer.LastUpdateTime,
	})
}

func (this *Merchant) Update(ctx iris.Context) {
	userId := int64(common.GetUserIdFromCtx(ctx))
	if userId == 0 {
		ZapLog().With(zap.Int64("userId", userId)).Error("userId is 0 err")
		this.ExceptionSerive(ctx, apibackend.BASERR_UNKNOWN_BUG.Code(), apibackend.BASERR_UNKNOWN_BUG.Desc())
		return
	}
	param := new(api.Merchant)
	if err := Tools.ShouldBindJSON(ctx, param); err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		return
	}
	param.PayeeId = &userId
	modelParam:= new(models.Merchant).Parse(param)
	mer,err := modelParam.Update()
	if err != nil {
		ZapLog().With(zap.Error(err), zap.Int64("userId", userId)).Error("Update err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}
	this.Response(ctx, &api.ResMerchant{
		ID: mer.ID,
		MerchantId: mer.MerchantId,
		MerchantName: mer.MerchantName,
		NotifyUrl: mer.NotifyUrl,
		SignType: mer.SignType,
		SignKey:  mer.SignKey,
		LanguageType: mer.LanguageType,
		LegalCurrency: mer.LegalCurrency,
		Contact:       mer.Contact,
		ContactPhone :  mer.ContactPhone,
		ContactEmail :  mer.ContactEmail,
		Country  :      mer.Country,
		CreateTime     : mer.CreateTime,
		LastUpdateTime : mer.LastUpdateTime,
	})
}