package shopee

import (
    "testing"
)

// func TestLogin(t *testing.T) {
//     client, config := setupTestClient(t)
// 	fmt.Println(config)

//     tests := []struct {
//         name     string
//         phone    string
//         password string
//         wantErr  bool
//     }{
//         {
//             name:     "正常登录",
//             phone:    config.Phone,
//             password: config.Password,
//             wantErr:  false,
//         },
//     }

//     for _, tt := range tests {
//         t.Run(tt.name, func(t *testing.T) {
//             err := client.Login(tt.phone, tt.password)
//             if (err != nil) != tt.wantErr {
//                 t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
//             }
//         })
//     }
// }

func TestGetProductList(t *testing.T) {
    client := NewClient()

	cookies := "SPC_CNSC_SESSION=c23dc6dcd04bcb38a0e457f84cd66139_2_2375038;"
	shopID := "1350463881"
	region := "sg"

    productIDs, err := client.GetProductList(cookies, shopID, region)
    if err != nil {
        t.Errorf("GetProductList() error = %v", err)
    }
    t.Logf("productIDs: %v", productIDs)
	t.Logf("len(productIDs): %d", len(productIDs))
} 

func TestGetMerchantShopList(t *testing.T) {
	client := NewClient()
	cookies := "SPC_CNSC_SESSION=4ed26470f6a268de0cf09be7e8362d69_2_2375038;"
	merchantShopList, err := client.GetMerchantShopList(cookies)
	if err != nil {
		t.Errorf("GetMerchantShopList() error = %v", err)
	}
	t.Logf("merchantShopList: %v", merchantShopList)
}

func TestUpdateProductInfo(t *testing.T) {
	client := NewClient()
	cookies := "SPC_CNSC_SESSION=c9ad3caf0d1d2d15d25d6e752a6c5723_2_2375038;"
	productID := int64(28760843741)
	day := 10
	shopID := "1350463893"
	region := "my"
	err := client.UpdateProductInfo(productID, day, cookies, shopID, region)
	if err != nil {
		t.Errorf("UpdateProductInfo() error = %v", err)
	}
}

