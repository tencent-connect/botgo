package apitest

import (
	"testing"

	"github.com/tencent-connect/botgo/dto"
)

var (
	scheduleID string
)

func Test_schedule(t *testing.T) {
	t.Run(
		"拉取日程列表", func(t *testing.T) {
			rsp, err := api.ListSchedules(ctx, testChannelID, 0)
			if err != nil {
				t.Error(err)
			}
			for _, v := range rsp {
				t.Logf("%v", v)
				scheduleID = v.ID
			}
		},
	)
	t.Run(
		"拉取单个日程", func(t *testing.T) {
			rsp, err := api.GetSchedule(ctx, testChannelID, scheduleID)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%v", rsp)
		},
	)
	t.Run(
		"创建日程", func(t *testing.T) {
			schedule := &dto.Schedule{
				Name:           "测试创建",
				StartTimestamp: "1639110300000",
				EndTimestamp:   "1639110900000",
				RemindType:     "0",
			}
			rsp, err := api.CreateSchedule(ctx, testChannelID, schedule)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%v", rsp)
			scheduleID = rsp.ID
		},
	)
	t.Run(
		"修改日程", func(t *testing.T) {
			schedule := &dto.Schedule{
				Name:           "测试修改",
				StartTimestamp: "1639110300000",
				EndTimestamp:   "1639110900000",
				RemindType:     "0",
			}
			rsp, err := api.ModifySchedule(ctx, testChannelID, scheduleID, schedule)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%v", rsp)
			scheduleID = rsp.ID
		},
	)
	t.Run(
		"删除日程", func(t *testing.T) {
			err := api.DeleteSchedule(ctx, testChannelID, scheduleID)
			if err != nil {
				t.Error(err)
			}
		},
	)
}
