package application

type UseCases struct {
	Registration   Registration
	Login          Login
	UploadOrder    UploadOrder
	GetOrders      GetOrders
	GetBalance     GetBalance
	UploadWithdraw UploadWithdraw
	GetWithdrawals GetWithdrawals
}
