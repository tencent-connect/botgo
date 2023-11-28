package token

import (
	"context"
	"fmt"
	"github.com/tencent-connect/botgo/errs"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestManager_GetAccessToken(t *testing.T) {
	type fields struct {
		appID             uint64
		appSecret         string
		Type              Type
		token             *AccessToken
		lock              sync.RWMutex
		forceRefreshToken chan interface{}
		closeCh           chan int
		once              sync.Once
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{fields: fields{
			appID:     123,
			appSecret: "123",
			Type:      TypeQQBot,
			token: &AccessToken{
				Token:     "123",
				ExpiresIn: 123,
			},
			lock: sync.RWMutex{},
		}, want: "123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				appID:             tt.fields.appID,
				appSecret:         tt.fields.appSecret,
				Type:              tt.fields.Type,
				token:             tt.fields.token,
				lock:              tt.fields.lock,
				forceRefreshToken: tt.fields.forceRefreshToken,
				closeCh:           tt.fields.closeCh,
				once:              tt.fields.once,
			}
			if got := m.GetAccessToken(); got.GetToken() != tt.want {
				t.Errorf("GetAccessToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetAppID(t *testing.T) {
	type fields struct {
		appID             uint64
		appSecret         string
		Type              Type
		token             *AccessToken
		lock              sync.RWMutex
		forceRefreshToken chan interface{}
		closeCh           chan int
		once              sync.Once
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{fields: fields{
			appID:     123,
			appSecret: "123",
			Type:      TypeQQBot,
			token: &AccessToken{
				Token:     "123",
				ExpiresIn: 123,
			},
			lock: sync.RWMutex{},
		}, want: 123},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				appID:             tt.fields.appID,
				appSecret:         tt.fields.appSecret,
				Type:              tt.fields.Type,
				token:             tt.fields.token,
				lock:              tt.fields.lock,
				forceRefreshToken: tt.fields.forceRefreshToken,
				closeCh:           tt.fields.closeCh,
				once:              tt.fields.once,
			}
			if got := m.GetAppID(); got != tt.want {
				t.Errorf("GetAppID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_Init(t *testing.T) {
	type fields struct {
		appID             uint64
		appSecret         string
		Type              Type
		token             *AccessToken
		lock              sync.RWMutex
		forceRefreshToken chan interface{}
		closeCh           chan int
		once              sync.Once
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				appID:     123,
				appSecret: "",
				Type:      "",
				token:     &AccessToken{},
			},
			args:    args{context.Background()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				appID:             tt.fields.appID,
				appSecret:         tt.fields.appSecret,
				Type:              tt.fields.Type,
				token:             tt.fields.token,
				lock:              tt.fields.lock,
				forceRefreshToken: tt.fields.forceRefreshToken,
				closeCh:           tt.fields.closeCh,
				once:              tt.fields.once,
			}
			if err := m.Init(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_MarshalJSON(t *testing.T) {
	type fields struct {
		appID             uint64
		appSecret         string
		Type              Type
		token             *AccessToken
		lock              sync.RWMutex
		forceRefreshToken chan interface{}
		closeCh           chan int
		once              sync.Once
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				appID:     1,
				appSecret: "2",
				Type:      TypeQQBot,
				token: &AccessToken{
					Token:      "1",
					ExpiresIn:  2,
					UpdateTime: time.Time{},
				},
			},
			want:    []byte("{\"AppID\":1,\"ClientSecret\":\"2\",\"Type\":\"QQBot\",\"Token\":{\"Token\":\"1\",\"ExpiresIn\":2,\"UpdateTime\":\"0001-01-01T00:00:00Z\"}}"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				appID:             tt.fields.appID,
				appSecret:         tt.fields.appSecret,
				Type:              tt.fields.Type,
				token:             tt.fields.token,
				lock:              tt.fields.lock,
				forceRefreshToken: tt.fields.forceRefreshToken,
				closeCh:           tt.fields.closeCh,
				once:              tt.fields.once,
			}
			got, err := m.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_Stop(t *testing.T) {
	type fields struct {
		appID             uint64
		appSecret         string
		Type              Type
		token             *AccessToken
		lock              sync.RWMutex
		forceRefreshToken chan interface{}
		closeCh           chan int
		once              sync.Once
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "normal",
			fields: fields{
				closeCh: make(chan int, 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				appID:             tt.fields.appID,
				appSecret:         tt.fields.appSecret,
				Type:              tt.fields.Type,
				token:             tt.fields.token,
				lock:              tt.fields.lock,
				forceRefreshToken: tt.fields.forceRefreshToken,
				closeCh:           tt.fields.closeCh,
				once:              tt.fields.once,
			}
			m.Close()
			if !isChanClose(m.closeCh) {
				t.Errorf("manager not stop")
			}
		})
	}
}

func TestManager_TokenStr(t *testing.T) {
	type fields struct {
		appID             uint64
		appSecret         string
		Type              Type
		token             *AccessToken
		lock              sync.RWMutex
		forceRefreshToken chan interface{}
		closeCh           chan int
		once              sync.Once
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "normal",
			fields: fields{

				Type: TypeQQBot,
				token: &AccessToken{
					Token: "123",
				},
			},
			want: "QQBot 123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				appID:             tt.fields.appID,
				appSecret:         tt.fields.appSecret,
				Type:              tt.fields.Type,
				token:             tt.fields.token,
				lock:              tt.fields.lock,
				forceRefreshToken: tt.fields.forceRefreshToken,
				closeCh:           tt.fields.closeCh,
				once:              tt.fields.once,
			}
			if got := m.GetTokenValue(); got != tt.want {
				t.Errorf("GetTokenValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_UnmarshalJSON(t *testing.T) {
	type fields struct {
		appID             uint64
		appSecret         string
		Type              Type
		token             *AccessToken
		lock              sync.RWMutex
		forceRefreshToken chan interface{}
		closeCh           chan int
		once              sync.Once
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				appID:     1,
				appSecret: "2",
				Type:      TypeQQBot,
				token: &AccessToken{
					Token:      "1",
					ExpiresIn:  2,
					UpdateTime: time.Time{},
				},
			},
			args:    args{data: []byte("{\"AppID\":1,\"ClientSecret\":\"2\",\"Type\":\"QQBot\",\"Token\":{\"Token\":\"1\",\"ExpiresIn\":2,\"UpdateTime\":\"0001-01-01T00:00:00Z\"}}")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{}
			if err := m.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if m.appID != tt.fields.appID || m.Type != tt.fields.Type || m.token.Token != tt.fields.token.Token ||
				m.token.UpdateTime != tt.fields.token.UpdateTime || m.token.ExpiresIn != tt.fields.token.ExpiresIn {
				t.Errorf("UnmarshalJSON() error %+v", m)
			}
		})
	}
}

func TestNewBotTokenManager(t *testing.T) {
	type args struct {
		appID  uint64
		secret string
	}
	tests := []struct {
		name string
		args args
		want *Manager
	}{
		{
			name: "normal",
			args: args{
				appID:  1,
				secret: "2",
			},
			want: &Manager{
				appID:     1,
				appSecret: "2",
				Type:      TypeQQBot,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBotTokenManager(tt.args.appID, tt.args.secret); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBotTokenManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewManager(t *testing.T) {
	type args struct {
		tokenType Type
	}
	tests := []struct {
		name string
		args args
		want *Manager
	}{
		{
			name: "normal",
			args: args{
				tokenType: TypeQQBot,
			},
			want: &Manager{Type: TypeQQBot},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewManager(tt.args.tokenType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isChanClose(t *testing.T) {
	type args struct {
		ch chan int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		prepare func(chan int)
	}{
		{
			name: "normal",
			args: args{ch: make(chan int, 1)},
			want: false,
		},
		{
			name: "closed",
			args: args{ch: make(chan int, 1)},
			want: true,
			prepare: func(ch chan int) {
				close(ch)
			},
		},
		{
			name: "ch nil",
			args: args{},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(tt.args.ch)
			}
			if got := isChanClose(tt.args.ch); got != tt.want {
				t.Errorf("isChanClose() = %v, want %v", got, tt.want)
			}
		})
	}
}

var mockTokenExpiresIn int64 = 10

// MockTokenManager mock ITokenManager
type MockTokenManager struct {
	*Manager
}

func (m MockTokenManager) Init(ctx context.Context) (err error) {
	m.forceRefreshToken = make(chan interface{}, 10)
	m.closeCh = make(chan int, 1)
	return nil
}

func (m MockTokenManager) doRefreshToken() error {
	m.token = &AccessToken{
		Token:      "123",
		ExpiresIn:  mockTokenExpiresIn,
		UpdateTime: time.Now(),
	}
	fmt.Printf("%+v refresh token\r\n", time.Now())
	return nil
}

func Test_startRefreshAccessToken(t *testing.T) {
	type args struct {
		ctx context.Context
		m   ITokenManager
	}
	tests := []struct {
		name         string
		args         *args
		wantErr      bool
		prepare      func(args *args)
		checkManager func(manager ITokenManager) error
		cleanUp      func(manager ITokenManager)
	}{
		{
			name: "manager stop",
			args: &args{
				ctx: nil,
				m: &MockTokenManager{&Manager{
					appID:             1,
					appSecret:         "1",
					Type:              TypeQQBot,
					forceRefreshToken: make(chan interface{}, 10),
					closeCh:           make(chan int, 1),
				}},
			},
			wantErr: false,
			prepare: func(args *args) {
				args.m.Close()
			},
			checkManager: func(manager ITokenManager) error {
				if manager.State() != ManagerStateStopped {
					return errs.New(-1, "manager not stopped")
				}
				return nil
			},
		},
		{
			name: "force refresh",
			args: &args{
				ctx: context.TODO(),
				m: &MockTokenManager{&Manager{
					appID:             1,
					appSecret:         "1",
					Type:              TypeQQBot,
					forceRefreshToken: make(chan interface{}, 10),
					closeCh:           make(chan int, 1),
				}},
			},
			wantErr: false,
			checkManager: func(manager ITokenManager) error {
				manager.GetRefreshSigCh() <- "test"
				time.Sleep(time.Duration(100) * time.Millisecond)
				gapSec := manager.GetAccessToken().UpdateTime.Sub(time.Now()).Seconds()
				if gapSec > 1 || gapSec < float64(-1) {
					return errs.New(-1, "force refresh not work")
				}
				return nil
			},
		},
		{
			name: "force refresh",
			args: &args{
				ctx: context.TODO(),
				m: &MockTokenManager{&Manager{
					appID:             1,
					appSecret:         "1",
					Type:              TypeQQBot,
					forceRefreshToken: make(chan interface{}, 10),
					closeCh:           make(chan int, 1),
				}},
			},
			wantErr: false,
			checkManager: func(manager ITokenManager) error {
				manager.GetRefreshSigCh() <- "test"
				time.Sleep(time.Duration(100) * time.Millisecond)
				gapSec := manager.GetAccessToken().UpdateTime.Sub(time.Now()).Seconds()
				if gapSec > 1 || gapSec < float64(-1) {
					return errs.New(-1, "force refresh not work")
				}
				return nil
			},
			cleanUp: func(manager ITokenManager) {
				manager.Close()
			},
		},
		{
			name: "context deadline exceed",
			args: &args{
				m: &MockTokenManager{&Manager{
					appID:             1,
					appSecret:         "1",
					Type:              TypeQQBot,
					forceRefreshToken: make(chan interface{}, 10),
					closeCh:           make(chan int, 1),
				}},
			},
			wantErr: false,
			prepare: func(args *args) {
				ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Millisecond)
				args.ctx = ctx
				cancel()
			},
			checkManager: func(manager ITokenManager) error {
				time.Sleep(time.Duration(100) * time.Millisecond)
				if manager.State() != ManagerStateStopped {
					return errs.New(-1, fmt.Sprintf("ctx deadline exceed,manager not stop,%+v", manager.State()))
				}
				return nil
			},
			cleanUp: func(manager ITokenManager) {
				manager.Close()
			},
		},
		{
			name: "test refresh Process",
			args: &args{
				ctx: context.TODO(),
				m: &MockTokenManager{&Manager{
					appID:             1,
					appSecret:         "1",
					Type:              TypeQQBot,
					forceRefreshToken: make(chan interface{}, 10),
					closeCh:           make(chan int, 1),
				}},
			},
			wantErr: false,
			checkManager: func(manager ITokenManager) error {
				gapSec := manager.GetAccessToken().UpdateTime.Sub(time.Now()).Seconds()
				if gapSec > 1 || gapSec < float64(-1) {
					return errs.New(-1, "token not valued")
				}
				fmt.Printf("gap sec:%+v\r\n", gapSec)
				time.Sleep(time.Duration(mockTokenExpiresIn) * time.Second)
				gapSec = manager.GetAccessToken().UpdateTime.Sub(time.Now()).Seconds()
				//tokenTTL = math.floor(tokenTTL)故差值可能大于1
				if gapSec > 2 || gapSec < float64(-2) {
					return errs.New(-1,
						fmt.Sprintf("ticker refresh not work,gapSec:%+v,token:%+v", gapSec, manager.GetAccessToken()))
				}
				fmt.Printf("gap sec:%+v\r\n", gapSec)
				return nil
			},
			cleanUp: func(manager ITokenManager) {
				manager.Close()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(tt.args)
			}
			if err := startRefreshAccessToken(tt.args.ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("startRefreshAccessToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.checkManager(tt.args.m); err != nil {
				t.Errorf("checkManager() error = %v, ", err)
			}
			if tt.cleanUp != nil {
				tt.cleanUp(tt.args.m)
			}
		})
	}
}
