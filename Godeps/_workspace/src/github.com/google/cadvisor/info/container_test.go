// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package info

import (
	"reflect"
	"testing"
	"time"
)

func TestStatsStartTime(t *testing.T) {
	N := 10
	stats := make([]*ContainerStats, 0, N)
	ct := time.Now()
	for i := 0; i < N; i++ {
		s := &ContainerStats{
			Timestamp: ct.Add(time.Duration(i) * time.Second),
		}
		stats = append(stats, s)
	}
	cinfo := &ContainerInfo{
		ContainerReference: ContainerReference{
			Name: "/some/container",
		},
		Stats: stats,
	}
	ref := ct.Add(time.Duration(N-1) * time.Second)
	end := cinfo.StatsEndTime()

	if !ref.Equal(end) {
		t.Errorf("end time is %v; should be %v", end, ref)
	}
}

func TestStatsEndTime(t *testing.T) {
	N := 10
	stats := make([]*ContainerStats, 0, N)
	ct := time.Now()
	for i := 0; i < N; i++ {
		s := &ContainerStats{
			Timestamp: ct.Add(time.Duration(i) * time.Second),
		}
		stats = append(stats, s)
	}
	cinfo := &ContainerInfo{
		ContainerReference: ContainerReference{
			Name: "/some/container",
		},
		Stats: stats,
	}
	ref := ct
	start := cinfo.StatsStartTime()

	if !ref.Equal(start) {
		t.Errorf("start time is %v; should be %v", start, ref)
	}
}

func createStats(cpuUsage, memUsage uint64, timestamp time.Time) *ContainerStats {
	stats := &ContainerStats{
		Cpu:    &CpuStats{},
		Memory: &MemoryStats{},
	}
	stats.Cpu.Usage.PerCpu = []uint64{cpuUsage}
	stats.Cpu.Usage.Total = cpuUsage
	stats.Cpu.Usage.System = 0
	stats.Cpu.Usage.User = cpuUsage
	stats.Memory.Usage = memUsage
	stats.Timestamp = timestamp
	return stats
}

func TestContainerStatsCopy(t *testing.T) {
	stats := createStats(100, 101, time.Now())
	shadowStats := stats.Copy(nil)
	if !reflect.DeepEqual(stats, shadowStats) {
		t.Errorf("Copy() returned different object")
	}
	stats.Cpu.Usage.PerCpu[0] = shadowStats.Cpu.Usage.PerCpu[0] + 1
	stats.Cpu.Load = shadowStats.Cpu.Load + 1
	stats.Memory.Usage = shadowStats.Memory.Usage + 1
	if reflect.DeepEqual(stats, shadowStats) {
		t.Errorf("Copy() did not deeply copy the object")
	}
	stats = shadowStats.Copy(stats)
	if !reflect.DeepEqual(stats, shadowStats) {
		t.Errorf("Copy() returned different object")
	}
}
