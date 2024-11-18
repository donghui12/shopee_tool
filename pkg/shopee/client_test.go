package shopee

import (
	"fmt"
    "testing"
)

func TestLogin(t *testing.T) {
    client, config := setupTestClient(t)
	fmt.Println(config)

    tests := []struct {
        name     string
        phone    string
        password string
        wantErr  bool
    }{
        {
            name:     "正常登录",
            phone:    config.Phone,
            password: config.Password,
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := client.Login(tt.phone, tt.password)
            if (err != nil) != tt.wantErr {
                t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestGetProductList(t *testing.T) {
    client, config := setupTestClient(t)
    loginTestClient(t, client, config)

    tests := []struct {
        name     string
        request  *ProductListRequest
        wantErr  bool
    }{
        {
            name: "正常获取商品列表",
            request: &ProductListRequest{
                PageSize: 20,
                PageNo:   1,
                SortBy:   "create_time",
                SortType: 1,
            },
            wantErr: false,
        },
        {
            name: "页码超出范围",
            request: &ProductListRequest{
                PageSize: 20,
                PageNo:   9999,
                SortBy:   "create_time",
                SortType: 1,
            },
            wantErr: false, // 空列表不算错误
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            resp, err := client.GetProductList(tt.request)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetProductList() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if err == nil {
                t.Logf("获取到 %d 个商品", len(resp.Data.Products))
            }
        })
    }
} 