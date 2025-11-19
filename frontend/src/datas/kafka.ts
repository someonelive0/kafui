/*
 * src/datas/kafka.ts
 * 直接用 ref/reactive 在模块作用域创建全局响应式数据，然后导出即可
 */
import { ref } from 'vue';
import { backend } from "../wailsjs/go/models";


// 全局响应式变量
// export const globalKafkaConfigsNum = ref(0)
export const globalBrokerNames = ref<string[]>([]);
export const globalBrokers = ref<backend.Broker[]>([]);

export const globalTopicNames = ref<string[]>([]);
export const globalGroupNames = ref<string[]>([]);


export function globalSetBrokers(brokers: backend.Broker[]) {
  let tmparray = brokers.map(broker => {
    return broker.host+broker.port;
  });
  globalBrokerNames.value = tmparray;
  globalBrokers.value = brokers;
}

export function globalGetBroker(id: number) {
    for (let j = 0; j < globalBrokers.value.length; j++) {
        if (globalBrokers.value[j].id == id) {
            return globalBrokers.value[j];
        }
    }
    return null;
}
