// Package keyboard 消息按钮
package keyboard

// ActionType 按钮操作类型
type ActionType uint32

// PermissionType 按钮的权限类型
type PermissionType uint32

const (
	// ActionTypeURL http 或 小程序 客户端识别 schema, data字段为链接
	ActionTypeURL ActionType = 0
	// ActionTypeCallback 回调互动回调地址, data 传给互动回调地址
	ActionTypeCallback ActionType = 1
	// ActionTypeAtBot at机器人, 根据 at_bot_show_channel_list 决定在当前频道或用户选择频道,自动在输入框 @bot data
	ActionTypeAtBot ActionType = 2
	// ActionTypeMQQAPI 客户端native跳转链接
	ActionTypeMQQAPI ActionType = 3
	// ActionTypeSubscribe 订阅按钮
	ActionTypeSubscribe ActionType = 4

	// PermissionTypeSpecifyUserIDs 仅指定这条消息的人可操作
	PermissionTypeSpecifyUserIDs PermissionType = 0
	// PermissionTypManager  仅频道管理者可操作
	PermissionTypManager PermissionType = 1
	// PermissionTypAll  所有人可操作
	PermissionTypAll PermissionType = 2
	// PermissionTypSpecifyRoleIDs 指定身份组可操作
	PermissionTypSpecifyRoleIDs PermissionType = 3
)

// MessageKeyboard 消息按钮组件
type MessageKeyboard struct {
	ID      string          `json:"id,omitempty"`      // 消息按钮组件模板 ID
	Content *CustomKeyboard `json:"content,omitempty"` // 消息按钮组件自定义内容
}

// CustomKeyboard 自定义 Keyboard
type CustomKeyboard struct {
	Rows  []*Row         `json:"rows,omitempty"`  // 行数组
	Style *KeyboardStyle `json:"style,omitempty"` // 按钮样式
}

// KeyboardStyle 键盘样式
type KeyboardStyle struct {
	FontSize string `json:"font_size,omitempty"` // 字体大小
}

// Row 每行结构
type Row struct {
	Buttons []*Button `json:"buttons,omitempty"` // 每行按钮
}

// Button 单个按纽
type Button struct {
	ID         string      `json:"id,omitempty"`          // 按钮 ID
	RenderData *RenderData `json:"render_data,omitempty"` // 渲染展示字段
	Action     *Action     `json:"action,omitempty"`      // 该按纽操作相关字段
	GroupID    string      `json:"group_id,omitempty"`    // 分组ID, 同一分组内有一个按钮操作后, 其它按钮则变灰不可点击 注意:只有当action.type = 1 时才有效
}

// RenderData  按纽渲染展示
type RenderData struct {
	Label        string `json:"label,omitempty"`         // 按纽上的文字
	VisitedLabel string `json:"visited_label,omitempty"` // 点击后按纽上文字
	Style        int    `json:"style,omitempty"`         // 按钮样式，0：灰色线框，1：蓝色线框 3: 白色背景+红色字体, 4:蓝色背景+白色字体
}

// Action 按纽点击操作
type Action struct {
	Type                 ActionType    `json:"type,omitempty"`        // 操作类型
	Permission           *Permission   `json:"permission,omitempty"`  // 可操作
	ClickLimit           uint32        `json:"click_limit,omitempty"` // 可点击的次数, 默认不限
	Data                 string        `json:"data,omitempty"`        // 操作相关数据
	Enter                bool          `json:"enter"`
	AtBotShowChannelList bool          `json:"at_bot_show_channel_list,omitempty"` // false:当前 true:弹出展示子频道选择器
	SubscribeData        SubscribeData `json:"subscribe_data,omitempty"`           // 订阅按钮数据，type=ActionTypeSubscribe时使用
	Modal                *Modal        `json:"modal,omitempty"`                    // 用户点击二次确认操作
}

// Modal 二次确认数据
type Modal struct {
	Content     string `json:"content,omitempty"`      // 二次确认的提示文本,如果不为空则会进行二次确认. 注意:最多40个字符, 不能有URL
	ConfirmText string `json:"confirm_text,omitempty"` // 二次确认提示确认按钮中展示的文字,可以为空,  默认为"确认" 注意:最多4个字符
	CancelText  string `json:"cancel_text,omitempty"`  // 二次确认提示取消按钮中的文字,可以为空,默认为"取消" 注意:最多4个字符
}

// Permission 按纽操作权限
type Permission struct {
	// Type 操作权限类型
	Type PermissionType `json:"type,omitempty"`
	// SpecifyRoleIDs 身份组
	SpecifyRoleIDs []string `json:"specify_role_ids,omitempty"`
	// SpecifyUserIDs 指定 UserID
	SpecifyUserIDs []string `json:"specify_user_ids,omitempty"`
}

// TemplateID 对模板id的封装，兼容官方模板和自定义模板
type TemplateID struct {
	// 这两个字段互斥，只填入一个
	TemplateID       uint32 `json:"template_id,omitempty"`        // 官方提供的模板id
	CustomTemplateID string `json:"custom_template_id,omitempty"` // 自定义模板
}

// SubscribeData 订阅按钮数据
type SubscribeData struct {
	TemplateIDs []*TemplateID `json:"template_ids,omitempty"` // 订阅按钮对应的模板id列表
}
