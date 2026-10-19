package entity

type OrderChannel string

const (
	OrderChannelWeb      OrderChannel = "web"
	OrderChannelMobile   OrderChannel = "mobile"
	OrderChannelCall     OrderChannel = "call"
	OrderChannelLine     OrderChannel = "line"
	OrderChannelFacebook OrderChannel = "facebook"
	OrderChannelWalkIn   OrderChannel = "walk_in"
)

func (c OrderChannel) IsValidOrderChannel() bool {
	switch c {
	case OrderChannelWeb,
		OrderChannelMobile,
		OrderChannelCall,
		OrderChannelLine,
		OrderChannelFacebook,
		OrderChannelWalkIn:
		return true
	default:
		return false
	}
}
