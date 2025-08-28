<template>
  <v-card flat>
    <v-card-title class="d-flex align-center pe-2">
      <v-icon icon="mdi-list-box-outline"></v-icon> &nbsp;
      {{ name }} Msgs {{ msgs.length }}
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
      ><v-tooltip activator="parent" location="bottom">Filter keyword</v-tooltip>
      </v-text-field>&nbsp;
      <v-text-field
        v-model="partition"
        label="Partition"
        prepend-inner-icon="mdi-paperclip"
        variant="outlined"
        hide-details
        single-line
        density="compact"
      ><v-tooltip activator="parent" location="bottom">Partition: 0, 1, 2...</v-tooltip>
      </v-text-field>&nbsp;
      <v-text-field
        v-model="limit"
        label="Limit"
        prepend-inner-icon="mdi-page-last"
        variant="outlined"
        hide-details
        single-line
        density="compact"
      ><v-tooltip activator="parent" location="bottom">Limit messages, 0 means no limit, -100 means last 100</v-tooltip>
      </v-text-field>&nbsp;
      <v-btn icon="mdi-refresh" size="small" @click="refresh"></v-btn>&nbsp;
      <v-btn icon="mdi-plus" size="small" @click="showNewMsgDialog"></v-btn>
    </v-card-title>

    <v-data-table
      :headers="headers"
      :items="msgs"
      :search="search"
      :loading="loading"
      :items-per-page="10"
      density="compact"
      item-value="offset"
      hover
    >
      <template v-slot:item="{ item }">
        <tr :class="getRowClass(item)"
          @click="rowClicked(item)">
          <td>{{ item.time }}</td>
          <td>{{ item.offset }}</td>
          <td>{{ item.partition }}</td>
          <td>{{ item.value }}</td>
        </tr>
      </template>
    </v-data-table>
    
    <v-divider :thickness="1" class="border-opacity-75" color="info"></v-divider>

    <v-card variant="outlined" hover>
      <v-card-subtitle class="d-flex align-center pt-4">
        <v-icon icon="mdi-file-document"></v-icon> &nbsp;
        Partition: {{ selectedPartition }}
        <v-spacer></v-spacer>
        Offset: {{ selectedOffset }}
        <v-spacer></v-spacer>
        Timestamp: {{ selectedTimestamp }}
        <v-spacer></v-spacer>
        <v-btn class="text-none" variant="flat" color="#5865f2" 
          prepend-icon="mdi-content-copy" 
          size="small" border
          @click="copyToClipboard">
          Copy
          <v-tooltip activator="parent" location="bottom">Copy to clipboard</v-tooltip>
        </v-btn>&nbsp;
        <v-btn class="text-none" variant="flat" color="#5865f2" 
          prepend-icon="mdi-file-send" 
          size="small" border
          @click="resendMsg">
          Resend
          <v-tooltip activator="parent" location="bottom">Resend this message to topic</v-tooltip>
        </v-btn>
      </v-card-subtitle>

      <json-viewer :value="jsonData"></json-viewer>
    </v-card>
  
    <v-dialog v-model="newMsgDialog" width="600">
      <v-card
        max-width="600" density="compact"
        prepend-icon="mdi-pen-plus"
        text="Input message's key and value."
        title="Write Message to Topic"
      >
        <v-container fluid >
          <v-row>
            <v-col cols="12" sm="12">
              <v-textarea v-model="msgkey"
                label="Key"
                row-height="15"
                rows="2"
                variant="outlined"
                auto-grow
                persistent-hint hint="Message's key can be empty"
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
                persistent-hint hint="Message's value is requried"
              ></v-textarea>
            </v-col>
          </v-row>
        </v-container>
        <template v-slot:actions>
          <v-spacer></v-spacer>
          <v-btn color="blue-darken-4" rounded="0" variant="outlined" text="Close" @click="newMsgDialog = false"></v-btn>
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
import { defineProps, onMounted, reactive, ref } from "vue";
import JsonViewer from 'vue-json-viewer';
import { backend } from "../wailsjs/go/models";


const { name } = defineProps(['name']) // 可以简写 解构
let msgs: Array<backend.Message> = reactive([]);
let loading = ref(true);
let search = ref('');
let limit = ref('1000');
let partition = ref('0');
let newMsgDialog = ref(false);
let msgkey = ref('');
let msgvalue = ref('');
let snackbar = ref(false);
let snacktext = '';
let jsonData = ref('');
let selectedPartition = ref(0);
let selectedOffset = ref(0);
let selectedTimestamp = ref('');
let selectedKey = '';
let selectedValue = '';


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
  // if (msgs) {
  //   msgs.splice(0);
  // }
  loading.value = true; // why not work?

  // -1 means partition, 3 means timeout, limit means records limit
  window.go.backend.KafkaTool.ReadMsgsLimit(name, 
      parseInt(partition.value), parseInt(limit.value), 8)
  .then((items: Array<backend.Message>) => {
    // console.log('Kafkatool.ReadMsgs ', items);
    msgs = items;
    loading.value = false;
  })
  .catch((err: string) => {
    // console.error('Kafkatool.ReadMsgs ', err);
    snacktext = 'read message failed: ' + err;
    snackbar.value = true;
    loading.value = false;
  });
}

const showNewMsgDialog = () => {
  newMsgDialog.value = true;
}

const writeMsg = () => {
  // console.log(msgkey.value, msgvalue.value);
  msgkey.value = msgkey.value.trim();
  msgvalue.value = msgvalue.value.trim();
  if (msgvalue.value.length == 0) {
    snacktext = 'message value can not be empty!'
    snackbar.value = true;
    return;
  }

  window.go.backend.KafkaTool.WriteMsg(name, msgkey.value, msgvalue.value).then(() => {
    snacktext = 'write message success!';
    snackbar.value = true;
    newMsgDialog.value = false;
    refresh();
  })
  .catch((err: string) => {
    // console.error('Kafkatool.WriteMsg ', err);
    snacktext = 'write message failed: ' + err;
    snackbar.value = true;
  });
}

const rowClicked = (row: backend.Message) => {
  // console.log('rowClicked ', row.time, row.offset);
  selectedPartition.value = row.partition;
  selectedOffset.value = row.offset;
  selectedTimestamp.value =  row.time;
  selectedKey = row.key;
  selectedValue = row.value;
  try {
    jsonData.value = JSON.parse(row.value);
  } catch (error) {
    jsonData.value = row.value;
  }
}

const copyToClipboard = () => {
  const text = JSON.stringify(jsonData.value);
  navigator.clipboard.writeText(text).then(() => {
    // console.log('Text copied to clipboard');
  }).catch(err => {
    console.error('Failed to copy: ', err);
  });
}

const resendMsg = () => {
  if (selectedValue.length == 0) {
    snacktext = 'resend message is empty';
    snackbar.value = true;
    return;
  }

  window.go.backend.KafkaTool.WriteMsg(name, selectedKey, selectedValue).then(() => {
    snacktext = 'resend message success!';
    snackbar.value = true;
    newMsgDialog.value = false;
    refresh();
  })
  .catch((err: string) => {
    // console.error('Kafkatool.WriteMsg ', err);
    snacktext = 'resend message failed: ' + err;
    snackbar.value = true;
  });
}

const getRowClass = (row: backend.Message) => {
  if (row.offset === selectedOffset.value) {
    return 'msg-highlight';
  }
  return '';
}

</script>

<style scoped>
.msg-highlight {
  background-color: rgba(173, 223, 252, 0.765);
}
</style>
