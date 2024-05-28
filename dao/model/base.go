package model

type DataStatus int32

const (
	DataStatusUnknown DataStatus = iota //无效值，未知错误
	DataStatusEnable                    //有效的，默认
	DataStatusDisable                   //已失效，下架
	DataStatusDelete                    //已删除（逻辑删除）
)
