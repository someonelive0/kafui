<template>
  <v-card flat>
    <v-card-title class="d-flex align-center pe-2">
      <v-icon icon="mdi-list-box-outline"></v-icon> &nbsp;
      {{ name }} Messages {{ msgs.length }}
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        label="Search"
        prepend-inner-icon="mdi-magnify"
        variant="outlined"
        hide-details
        single-line
        density="compact"
      ></v-text-field>&nbsp;
      <v-btn icon="mdi-refresh" size="small" @click="refresh"></v-btn>&nbsp;
      <v-btn icon="mdi-plus" size="small" @click="showNewDialog"></v-btn>
    </v-card-title>

    <v-data-table
      :headers="headers"
      :items="msgs"
      :search="search"
      :loading="loading"
      :items-per-page="10"
      density="compact"
      item-value="offset"
    >
      <template v-slot:item="{ item }">
        <tr @click="rowClicked(item)">
          <td>{{ item.time }}</td>
          <td>{{ item.offset }}</td>
          <td>{{ item.partition }}</td>
          <td>{{ item.value }}</td>
        </tr>
      </template>
    </v-data-table>
    
    <v-dialog v-model="newDialog" width="600">
      <v-card
        max-width="600"
        prepend-icon="mdi-pen-plus"
        text="Input message's key and value."
        title="Write Message to Topic"
      >
        <v-container fluid>
          <v-row>
            <v-col cols="12" sm="12">
              <v-textarea v-model="msgkey"
                label="Key"
                row-height="15"
                rows="2"
                variant="outlined"
                auto-grow
              ></v-textarea>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" sm="12">
              <v-textarea v-model="msgvalue"
                label="Value"
                row-height="25"
                rows="8"
                variant="outlined"
                auto-grow
                shaped
              ></v-textarea>
            </v-col>
          </v-row>
        </v-container>
        <template v-slot:actions>
          <v-spacer></v-spacer>
          <v-btn color="blue-darken-4" rounded="0" variant="outlined" text="Close" @click="newDialog = false"></v-btn>
          <v-btn color="blue-darken-4" rounded="0" variant="flat" text="Write" @click="writeMsg"></v-btn>
        </template>
      </v-card>
    </v-dialog>
  </v-card>

  <v-snackbar v-model="snackbar" timeout=2000 color="deep-purple-accent-4" elevation="24">
    {{ snacktext }}
    <template v-slot:actions>
      <v-btn color="pink" variant="text" @click="snackbar = false">Close</v-btn>
    </template>
  </v-snackbar>
</template>

<script setup lang="ts">
import { ref, defineProps } from "vue"
import {onBeforeMount,onMounted,onBeforeUpdate,onUnmounted} from "vue"


const { name } = defineProps(['name']) // 可以简写 解构
let msgs: Array<object> = ref([]);
let loading = ref(true);
let search = ref('');
let newDialog = ref(false);
let msgkey = ref('');
let msgvalue = ref('');
let snackbar = ref(false);
let snacktext = '';

const headers = [
  { title: 'Timestamp', align: 'start', key: 'time' },
  { title: 'Offset', align: 'start', key: 'offset' },
  { title: 'Partition', align: 'start', sortable: true, key: 'partition' },
  // { title: 'Key', align: 'end', key: 'key',
    // value: item => window.atob(item.key) // 如果类型在golang是[]byte，则自动用base64编码，所以要用base64解码
  // },
  { title: 'Value', align: 'start', key: 'value',
    // value: item => window.atob(item.value)
  },
];

onMounted(() => {
  refresh();
})

const refresh = () => {
  // if (msgs.value) {
  //   msgs.value.splice(0);
  // }
  loading.value = true; // why not work?

  // -1 means partition, 3 means timeout
  window.go.backend.KafkaTool.ReadMsgs(name, -1, 3).then(items => {
    console.log('Kafkatool.ReadMsgs ', items);
    msgs.value = items;
    loading.value = false;
  })
  .catch(err => {
    console.error('Kafkatool.ReadMsgs ', err);
    loading.value = false;
  });
}

const showNewDialog = () => {
  newDialog.value = true;
}

const writeMsg = () => {
  // console.log(msgkey.value, msgvalue.value);
  msgkey.value = msgkey.value.trim();
  msgvalue.value = msgvalue.value.trim();
  if (msgvalue.value.length == 0) {
    snacktext = 'value can not be empty!'
    snackbar.value = true;
    return;
  }

  window.go.backend.KafkaTool.WriteMsg(name, msgkey.value, msgvalue.value).then(items => {
    snacktext = 'write message success!';
    snackbar.value = true;
    newDialog.value = false;
    refresh();
  })
  .catch(err => {
    // console.error('Kafkatool.WriteMsg ', err);
    snacktext = 'write message failed: ' + err;
    snackbar.value = true;
  });
}

const rowClicked = (row) => {
  console.log('rowClicked ', typeof row.time);
}

</script>
