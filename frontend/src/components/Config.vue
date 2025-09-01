<template>
  <v-card flat>
    <v-card-title class="d-flex align-center pe-2">
      <v-icon icon="mdi-list-box-outline"></v-icon> &nbsp;
        {{ card_title }}
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        label="Filter"
        prepend-inner-icon="mdi-filter-outline"
        variant="outlined"
        hide-details
        single-line
        density="compact"
        clearable
        ><v-tooltip activator="parent" location="bottom">Filter by keyword</v-tooltip>
      </v-text-field>&nbsp;
      <v-btn icon="mdi-refresh" size="small" @click="refresh"></v-btn>
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
      hover
    >
    <template v-slot:item="{ item }">
      <tr :class="getRowClass(item)"
        @click="rowClicked(item)">
        <td>{{ item.config_name }}</td>
        <td>{{ item.config_value }}</td>
        <td>{{ item.readonly }}</td>
        <td>{{ item.is_default }}</td>
        <td>{{ item.is_sensitive }}</td>
        <td class="d-flex ga-2 justify-end">
          <v-icon color="medium-emphasis" icon="mdi-pencil" size="small" @click="edit(item)"></v-icon>
        </td>
      </tr>
    </template>
    <template #bottom>
      <!-- Leave this slot empty to hide pagination controls -->
    </template>
    </v-data-table>

    <v-dialog v-model="dialog" max-width="500">
      <v-card
        :subtitle="`${isEditing ? 'Set' : 'Create'} config item of ${title}`"
        :title="`${isEditing ? 'Edit' : 'Add'} config item of ${title}`"
      >
        <template v-slot:text>
          <v-row>
            <v-col cols="6">
              <v-text-field v-model="formModel.topic" label="Object" color="black" disabled></v-text-field>
            </v-col>

            <v-col cols="6" md="6">
              <v-text-field v-model="formModel.config_name" label="Config Name" disabled></v-text-field>
            </v-col>

            <v-col cols="12">
              <v-text-field v-model="formModel.config_value" label="Config Value" variant="outlined"></v-text-field>
            </v-col>

          </v-row>
        </template>

        <v-divider></v-divider>

        <v-card-actions class="bg-surface-light">
          <v-btn text="Cancel" variant="plain" @click="dialog = false"></v-btn>
          <v-spacer></v-spacer>
          <v-btn text="Save" variant="outlined" 
            prepend-icon="mdi-check-circle" color="blue-darken-4"
            @click="save"></v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>

  <v-snackbar v-model="snackbar" timeout=2000 color="deep-purple-darken-3" elevation="24">
    {{ snacktext }}
    <template v-slot:actions>
      <v-btn color="grey" variant="text" @click="snackbar = false">Close</v-btn>
    </template>
  </v-snackbar>
</template>


<script setup lang="ts">
import { defineProps, onMounted, reactive, ref, shallowRef, toRef } from "vue";
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
const sortBy = ref([{ key: 'config_name', order: 'asc' }]);
let snackbar = ref(false);
let snacktext = '';

// for edit group offset
const formModel = ref([])
const dialog = shallowRef(false)
const isEditing = toRef(() => !!formModel.value.topic)


onMounted(() => {
  refresh();
})

const refresh = () => {
  loading.value = true;
  let getConf = null;
  if (title == 'topic') {
    getConf = window.go.backend.KafkaTool.GetTopicConfig;
  } else if (title == 'broker') {
    getConf = window.go.backend.KafkaTool.GetBrokerConfig;
  } else if (title == 'cluster') {
    getConf = window.go.backend.KafkaTool.GetClusterConfig;
  } else {
    console.error('unknow title ', title);
    return;
  }
  
  // console.log(title, name);
  getConf(name.toString()).then((items: Array<backend.ConfigEntry>) => {
    // console.log('Kafkatool.getConf ', title, items);
    configs = items;
    loading.value = false;
  })
  .catch((err: string) => {
    console.error('Kafkatool.GetTopicConfig ', err);
    snacktext = 'get configs of ' + name.toString() + ' failed: ' + err;
    snackbar.value = true;
    loading.value = false;
  });
}

const edit = (item) => {
  // const found = books.value.find(book => book.id === id)
  // console.log("edit item: " + item.topic);
  formModel.value = item;
  formModel.value.topic = name.toString();
  dialog.value = true
}

const save = () => {
  // console.log("save item: " + formModel.value.topic);
  formModel.value.config_value = formModel.value.config_value.trim();
  if (formModel.value.config_value.length == 0) {
    snacktext = 'Error: New config value can not be empty!';
    snackbar.value = true;
    return;
  }

  let setConf = null;
  if (title == 'topic') {
    setConf = window.go.backend.KafkaTool.SetTopicConfig;
  } else if (title == 'broker') {
    setConf = window.go.backend.KafkaTool.SetBrokerConfig;
  } else if (title == 'cluster') {
    setConf = window.go.backend.KafkaTool.SetClusterConfig;
  } else {
    console.error('unknow title ', title);
    return;
  }

  setConf(formModel.value.topic, formModel.value.config_name, formModel.value.config_value ).then(() => {
    snacktext = 'set config ' + name + ' success!';
    snackbar.value = true;
    dialog.value = false
    refresh();
  })
  .catch((err: string) => {
    // console.error('Kafkatool.setConf ', err);
    snacktext = 'set config ' + name + ' failed: ' + err;
    snackbar.value = true;
  });

}

let selectedRowName = 'compression.type';
const rowClicked = (row: backend.ConfigEntry) => {
  // console.log("Clicked item: ", row.config_name)
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
