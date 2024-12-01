<template>
  <v-container fluid class="pa-1 ma-1">
    <v-card flat>
      <v-card-title class="d-flex align-center pe-2">
        <v-icon icon="mdi-list-box-outline"></v-icon> &nbsp;
        Topics {{ topics.length }}
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
        <v-btn icon="mdi-plus" size="small" @click="newDialog = true"></v-btn>
      </v-card-title>

      <v-data-table density="compact"
        :headers="headers"
        :items="topics"
        :search="search"
        :items-per-page="-1"
      >
        <template v-slot:item="{ item }">
          <tr 
            @click="rowClicked(item)">
            <td>{{ item.name }}</td>
          </tr>
        </template>
        <template #bottom>
          <!-- Leave this slot empty to hide pagination controls -->
        </template>
      </v-data-table>
    </v-card>

    <v-dialog v-model="newDialog" width="600">
      <v-card
        max-width="500"
        prepend-icon="mdi-pen-plus"
        text="Input Topic name and number of partitions and replicas."
        title="Create New Topic"
      >
        <v-container fluid>
          <v-row dense class="d-flex align-center">
            <v-col cols="5" md="5" sm="5">Topic Name*:</v-col>
            <v-col cols="7" md="7" sm="7">
                <v-text-field hide-details="auto" v-model="topic_name"></v-text-field>
            </v-col>
          </v-row>
          <v-row dense class="d-flex align-center">
            <v-col cols="5" md="5" sm="5">Number of Partitions*:</v-col>
            <v-col cols="7" md="7" sm="7">
                <v-text-field hide-details="auto" v-model="partitions"></v-text-field>
            </v-col>
          </v-row>
          <v-row dense class="d-flex align-center">
            <v-col cols="5" md="5" sm="5">Number of Replicas*:</v-col>
            <v-col cols="7" md="7" sm="7">
                <v-text-field hide-details="auto" v-model="replicas"></v-text-field>
            </v-col>
          </v-row>
        </v-container>
        <template v-slot:actions>
          <v-spacer></v-spacer>
          <v-btn color="blue-darken-4" rounded="0" variant="outlined" text="Cancel" @click="newDialog = false"></v-btn>
          <v-btn color="blue-darken-4" rounded="0" variant="flat" text="Create" @click="createTopic"></v-btn>
        </template>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snackbar" timeout=4000 color="deep-purple-accent-4" elevation="24">
    {{ snacktext }}
    <template v-slot:actions>
      <v-btn color="pink" variant="text" @click="snackbar = false">Close</v-btn>
    </template>
  </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, reactive } from "vue"
import { useRouter } from 'vue-router';


// 使用 stat 传递页面参数, 把字符串数组转换成对象数组
// let topics = window.history.state.topics;
// console.log('window.history.state.topics ', topics);
let topics: Array<object> = reactive([]);
if (window.history.state.topics) {
  for (var i=0,len=window.history.state.topics.length; i<len; i++) {
    topics[i] = {name: window.history.state.topics[i]}
  }
}

const headers: Array<object> = [
  { title: 'Topic Name', align: 'start', sortable: true, key: 'name' },
];
let search = ref('');
const router = useRouter(); 
let newDialog = ref(false);
let topic_name = ref('');
let partitions = ref('1');
let replicas = ref('1');
let snackbar = ref(false);
let snacktext = '';

const refresh = () => {
  window.go.backend.KafkaTool.ListTopics().then((items: Array<string>) => {
    // console.log('Kafkatool.ListTopics ', items);
    for (var i=0,len=items.length; i<len; i++) {
      topics[i] = {name: items[i]}
    }
    // loading.value = false;
  })
  .catch((err: string) => {
    console.error('Kafkatool.ListTopics ', err);
    // snacktext = 'read message failed: ' + err;
    // snackbar.value = true;
    // loading.value = false;
  });
}

const rowClicked = (row) => {
  // console.log("Clicked item: ", row)
  router.push({
    name: 'Topic',
    query: {
        topic: row.name
    }
  });
}

const valid = () => {
  topic_name.value = topic_name.value.trim()
  partitions.value = partitions.value.trim()
  replicas.value = replicas.value.trim()
  if (topic_name.value.length == 0) {
    snacktext = 'topic name con not be empty';
    snackbar.value = true;
    return false;
  }
  if (isNaN(Number(partitions.value)) || isNaN(Number(replicas.value))) {
    snacktext = 'partitions or replicas is not number';
    snackbar.value = true;
    return false;
  }

  return true
}

const createTopic = () => {
  snackbar.value = false;
  if (!valid()) return;
  console.log(topic_name.value, partitions.value, replicas.value);

  window.go.backend.KafkaTool.CreateTopic(topic_name.value, Number(partitions.value), Number(replicas.value)).then(() => {
    // console.log('Kafkatool.ListTopics ', items);
    snacktext = 'create topic success!';
    snackbar.value = true;
    newDialog.value = false;
    refresh();
  })
  .catch((err: string) => {
    // console.error('Kafkatool.CreateTopic ', err);
    snacktext = 'create topic failed: ' + err;
    snackbar.value = true;
  });
}

</script>
