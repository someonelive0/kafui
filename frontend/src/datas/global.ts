/*
 * src/datas/global.ts
 * 直接用 ref/reactive 在模块作用域创建全局响应式数据，然后导出即可
 */
import { ref } from 'vue';
import { backend } from "../wailsjs/go/models";


// 全局响应式变量
// export const globalKafkaConfigsNum = ref(0)
export const globalKafkaConfigNames = ref<string[]>([]);
export const globalKafkaConfigs = ref<backend.KafkaConfig[]>([]);


export function globalSetKafkaConfigs(kafkaconfigs: backend.KafkaConfig[]) {
  let tmparray = kafkaconfigs.map(kafkaconfig => {
    return kafkaconfig.name;
  });
  globalKafkaConfigNames.value = tmparray;
  globalKafkaConfigs.value = kafkaconfigs;
}

export function globalAddKafkaConfig(kafkaconfig: backend.KafkaConfig) {
  globalKafkaConfigs.value.push(kafkaconfig);
  globalKafkaConfigNames.value.push(kafkaconfig.name);
}

// Just update kafkaconfig, not change name of config
export function globalUpdateKafkaConfig(kafkaconfig: backend.KafkaConfig) {
  let i = -1;
  for (let j = 0; j < globalKafkaConfigs.value.length; j++) {
    if (globalKafkaConfigs.value[j].name == kafkaconfig.name) {
      i = j;
      break;
    }
  }
  if (i >= 0) {
    globalKafkaConfigs.value[i] = kafkaconfig;
  }
}

export function globalDeleteKafkaConfig(kafkaName: string) {
  let i = -1;
  for (let j = 0; j < globalKafkaConfigs.value.length; j++) {
    if (globalKafkaConfigs.value[j].name == kafkaName) {
      i = j;
      break;
    }
  }
  if (i >= 0) {
    globalKafkaConfigNames.value.splice(i, 1);
    globalKafkaConfigs.value.splice(i, 1);
  }
}



/*
  Map<string, ref<string>> of reactive
  set: globalScannerStatusMap.set(connName, ref('none'));
  change:   globalScannerStatusMap.get(connName.value).value = 'running';
  bind: globalScannerStatusMap.get(connName)
*/
// export const globalScannerStatusMap = reactive(new Map());

/*
  Map<string, reactive<ScannerStatus>> of reactive
  set: globalScannerMap.set(kafkaconfig.name, reactive(new ScannerStatus(kafkaconfig.name)));
  change:   globalScannerMap.get(connName.value).logs += 'running...';
  bind: globalScannerMap.get(connName).logs
*/
// export const globalScannerMap = reactive(new Map());
