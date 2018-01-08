// Copyright 2017 luoji

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//    http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package server

import (
	"time"

	"github.com/boltmq/common/logger"
	"github.com/boltmq/common/utils/system"
)

// controllerTasks broker控制器的各种任务
// Author: tianyuliang
// Since: 2017/10/11
type controllerTasks struct {
	brokerController            *BrokerController
	deleteTopicTask             *system.Ticker
	brokerStatsRecordTask       *system.Ticker
	persistConsumerOffsetTask   *system.Ticker
	scanUnSubscribedTopicTask   *system.Ticker
	fetchNameServerAddrTask     *system.Ticker
	slaveSynchronizeTask        *system.Ticker
	printMasterAndSlaveDiffTask *system.Ticker
	registerAllBrokerTask       *system.Ticker
}

func newControllerTasks(controller *BrokerController) *controllerTasks {
	controllerTask := &controllerTasks{
		brokerController: controller,
	}
	return controllerTask
}

func (ctasks *controllerTasks) shutdown() bool {
	if ctasks == nil {
		return false
	}
	if ctasks.printMasterAndSlaveDiffTask != nil {
		ctasks.printMasterAndSlaveDiffTask.Stop()
		logger.Info("printMasterAndSlaveDiffTask stop ok")
	}

	if ctasks.deleteTopicTask != nil {
		ctasks.deleteTopicTask.Stop()
		logger.Info("deleteTopicTask stop ok")
	}
	if ctasks.brokerStatsRecordTask != nil {
		ctasks.brokerStatsRecordTask.Stop()
		logger.Info("brokerStatsRecordTask stop ok")
	}
	if ctasks.persistConsumerOffsetTask != nil {
		ctasks.persistConsumerOffsetTask.Stop()
		logger.Info("persistConsumerOffsetTask stop ok")
	}
	if ctasks.scanUnSubscribedTopicTask != nil {
		ctasks.scanUnSubscribedTopicTask.Stop()
		logger.Info("scanUnSubscribedTopicTask stop ok")
	}
	if ctasks.fetchNameServerAddrTask != nil {
		ctasks.fetchNameServerAddrTask.Stop()
		logger.Info("fetchNameServerAddrTask stop ok")
	}
	if ctasks.slaveSynchronizeTask != nil {
		ctasks.slaveSynchronizeTask.Stop()
		logger.Info("slaveSynchronizeTask stop ok")
	}
	if ctasks.registerAllBrokerTask != nil {
		ctasks.registerAllBrokerTask.Stop()
		logger.Info("registerAllBrokerTask stop ok")
	}
	return true
}

// startDeleteTopicTask 清除未使用Topic
// Author: tianyuliang
// Since: 2017/10/10
func (ctasks *controllerTasks) startDeleteTopicTask() {
	if ctasks.deleteTopicTask != nil {
		return
	}

	ctasks.deleteTopicTask = system.NewTicker(false, 5*time.Minute, 10*time.Second, func() {
		topics := ctasks.brokerController.tpConfigManager.tpCfgSerialWrapper.TpConfigTable.Keys()
		removedTopicCount := ctasks.brokerController.messageStore.CleanUnusedTopic(topics)
		if removedTopicCount > 0 {
			logger.Infof("deleteTopicTask removed topic count: %d", removedTopicCount)
		}
	})
	ctasks.deleteTopicTask.Start()
	logger.Infof("deleteTopicTask start ok")
}

// startBrokerStatsRecordTask 定时统计broker各类信息
// Author: tianyuliang
// Since: 2017/10/10
func (ctasks *controllerTasks) startBrokerStatsRecordTask() {
	initialDelay := system.ComputNextMorningTimeMillis() - system.CurrentTimeMillis()
	ctasks.brokerStatsRecordTask = system.NewTicker(false, time.Duration(initialDelay)*time.Millisecond, 24*time.Hour, func() {
		ctasks.brokerController.brokerStatsRelatedStore.Record()
	})
	ctasks.brokerStatsRecordTask.Start()
	logger.Infof("brokerStatsRecordTask start ok")
}

// startPersistConsumerOffsetTask 定时写入ConsumerOffset文件
// Author: tianyuliang
// Since: 2017/10/10
func (ctasks *controllerTasks) startPersistConsumerOffsetTask() {
	period := time.Duration(ctasks.brokerController.cfg.Broker.FlushConsumerOffsetInterval) * time.Millisecond
	ctasks.persistConsumerOffsetTask = system.NewTicker(false, 10*time.Second, period, func() {
		ctasks.brokerController.csmOffsetManager.cfgManagerLoader.persist()
	})
	ctasks.persistConsumerOffsetTask.Start()
	logger.Infof("persistConsumerOffsetTask start ok")
}

// startCcanUnSubscribedTopicTask 扫描被删除Topic，并删除该Topic对应的Offset
// Author: tianyuliang
// Since: 2017/10/10
func (ctasks *controllerTasks) startScanUnSubscribedTopicTask() {
	ctasks.scanUnSubscribedTopicTask = system.NewTicker(false, 10*time.Minute, 1*time.Hour, func() {
		ctasks.brokerController.csmOffsetManager.scanUnsubscribedTopic()
	})
	ctasks.scanUnSubscribedTopicTask.Start()
	logger.Infof("scanUnSubscribedTopicTask start ok")
}

// startFetchNameServerAddrTask 更新Namesrv地址列表
// Author: tianyuliang
// Since: 2017/10/10
func (ctasks *controllerTasks) startFetchNameServerAddrTask() {
	ctasks.fetchNameServerAddrTask = system.NewTicker(false, 10*time.Second, 2*time.Minute, func() {
		ctasks.brokerController.callOuter.FetchNameServerAddr()
	})
	ctasks.fetchNameServerAddrTask.Start()
	logger.Infof("fetchNameServerAddrTask start ok")
}

// startSlaveSynchronizeTask 启动“Slave同步所有数据”任务
// Author: tianyuliang
// Since: 2017/10/10
func (ctasks *controllerTasks) startSlaveSynchronizeTask() {
	ctasks.slaveSynchronizeTask = system.NewTicker(false, 10*time.Second, 1*time.Minute, func() {
		ctasks.brokerController.slaveSync.syncAll()
	})
	ctasks.slaveSynchronizeTask.Start()
	logger.Infof("slaveSynchronizeTask start ok")
}

// startPrintMasterAndSlaveDiffTask 启动“输出主从偏移量差值”任务
// Author: tianyuliang
// Since: 2017/10/10
func (ctasks *controllerTasks) startPrintMasterAndSlaveDiffTask() {
	ctasks.printMasterAndSlaveDiffTask = system.NewTicker(false, 10*time.Second, 1*time.Minute, func() {
		diff := ctasks.brokerController.messageStore.SlaveFallBehindMuch()
		if diff > 0 {
			logger.Infof("slave fall behind master, how much: %d bytes", diff) // warn and notify me
		}
	})
	ctasks.printMasterAndSlaveDiffTask.Start()
	logger.Infof("printMasterAndSlaveDiffTask start ok")
}

// startRegisterAllBrokerTask 注册所有Broker
// Author: tianyuliang
// Since: 2017/10/10
func (ctasks *controllerTasks) startRegisterAllBrokerTask() {
	ctasks.registerAllBrokerTask = system.NewTicker(false, 10*time.Second, 30*time.Second, func() {
		ctasks.brokerController.registerBrokerAll(true, false)
	})
	ctasks.registerAllBrokerTask.Start()
	logger.Infof("registerAllBrokerTask start ok")
}
