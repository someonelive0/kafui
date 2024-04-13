<template>
  <v-card flat>
    <v-card-title class="d-flex align-center pe-2">
      <v-icon icon="mdi-list-box-outline"></v-icon> &nbsp;
        {{ card_title }}
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        label="Search"
        prepend-inner-icon="mdi-magnify"
        variant="outlined"
        hide-details
        single-line
        density="compact"
      ></v-text-field>
    </v-card-title>

    <v-data-table
      v-model:sort-by="sortBy"
      :headers="headers"
      :items="configs"
      :search="search"
      :loading="loading"
      :items-per-page="-1"
      item-value="config_name"
      density="compact"
    >
    <template v-slot:item="{ item }">
      <tr :class="getRowClass(item)"
        @click="rowClicked(item)">
        <td>{{ item.config_name }}</td>
        <td>{{ item.config_value }}</td>
        <td>{{ item.readonly }}</td>
        <td>{{ item.is_default }}</td>
        <td>{{ item.is_sensitive }}</td>
      </tr>
    </template>
    <template #bottom>
      <!-- Leave this slot empty to hide pagination controls -->
    </template>
    </v-data-table>
  </v-card>

  <v-snackbar v-model="snackbar" timeout=2000 color="deep-purple-accent-4" elevation="24">
    {{ snacktext }}
    <template v-slot:actions>
      <v-btn color="pink" variant="text" @click="snackbar = false">Close</v-btn>
    </template>
  </v-snackbar>
</template>


<script setup lang="ts">
import { ref, reactive, defineProps } from "vue"
import {onBeforeMount,onMounted,onBeforeUpdate,onUnmounted} from "vue"
import { backend } from "../wailsjs/go/models";


// 调用defineProps方法并获取父组件传递的数据
// const props = defineProps(['title'])
let { title, name } = defineProps(['title', 'name']) // 可以简写 解构
const card_title = 'Configs of ' + title + ' "' + name + '"';

const headers: Array<object> = [
  { title: 'Name', align: 'start', sortable: true, key: 'config_name' },
  { title: 'Value', align: 'end', key: 'config_value' },
  // { title: 'Type', align: 'end', key: 'config_type' },
  // { title: 'Source', align: 'end', key: 'config_source' },
  { title: 'readonly', align: 'end', key: 'readonly' },
  { title: 'is_default', align: 'end', key: 'is_default' },
  { title: 'is_sensitive', align: 'end', key: 'is_sensitive' },
];
let configs: Array<backend.ConfigEntry> = reactive([]);
let loading = ref(true);
let search = ref('');
const sortBy = [{ key: 'config_name', order: 'asc' }];
let snackbar = ref(false);
let snacktext = '';

onMounted(() => {
  let getConf = null;
  if (title == 'topic') {
    getConf = window.go.backend.KafkaTool.GetTopicConfig;
  } else if (title == 'broker') {
    getConf = window.go.backend.KafkaTool.GetBrokerConfig;
    name = name.toString();
  } else if (title == 'cluster') {
    getConf = window.go.backend.KafkaTool.GetClusterConfig;
    name = name.toString();
  } else {
    console.error('unknow title ', title);
    return;
  }
  
  // console.log(title, name);
  getConf(name).then((items: Array<backend.ConfigEntry>) => {
    console.log('Kafkatool.getConf ', title, items);
    configs = items;
    loading.value = false;
  })
  .catch((err: string) => {
    console.error('Kafkatool.GetTopicConfig ', err);
    snacktext = 'get configs of ' + name + ' failed: ' + err;
    snackbar.value = true;
    loading.value = false;
  });
})

let selectedRowName = 'compression.type';
const rowClicked = (row: backend.ConfigEntry) => {
  console.log("Clicked item: ", row.config_name)
  selectedRowName = row.config_name;
}

const getRowClass = (row: backend.ConfigEntry) => {
  // console.log("getRowClass: ", row.config_name)
  // if (selectedRowName == row.config_name) {
  if (row.readonly === false) {
    return 'highlight';
  }
  return '';
}

</script>

<style>
/* 隔列变色 */
.v-table tbody tr td:nth-child(even) {
  background-color: rgba(250, 250, 250, 0.6);
}
.v-table tbody tr td:nth-child(odd) {
  background-color: rgba(244, 245, 245, 0.9);
}

/* 隔行变色 */
/* .v-table tbody tr:nth-child(even) {
  background-color: rgba(250, 250, 250, 0.55);
}

.v-table tbody tr:nth-child(odd) {
  background-color: rgba(244, 245, 245, 0.865);
} */

.highlight {
  background-color: rgba(180, 180, 4, 0.765);
}
</style>
