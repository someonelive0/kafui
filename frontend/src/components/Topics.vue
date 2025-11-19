<template>
  <v-container fluid class="pa-1 ma-1">
    <v-card flat>
      <v-card-title class="d-flex align-center pe-2">
        <v-icon icon="mdi-list-box-outline"></v-icon> &nbsp;
        Topics {{ globalTopicNames.length }}
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
        <v-btn icon="mdi-refresh" size="small" @click="refresh"></v-btn>&nbsp;
        <v-tooltip text="Create new topic" location="bottom">
          <template v-slot:activator="{ props }">
            <v-btn v-bind="props" icon="mdi-plus" size="small" @click="newDialog = true"></v-btn>
          </template>
        </v-tooltip>
      </v-card-title>

      <v-data-table density="compact"
        :headers="headers"
        :items="globalTopicNames"
        :search="search"
        :items-per-page="-1"
        hover
      >
        <template v-slot:item="{ item }">
          <tr 
            @click="rowClicked(item)">
            <td>{{ item }}</td>
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
        text="Input Topic name and number of partitions and replicas."
        title="Create New Topic"
      >
        <template v-slot:prepend>
          <v-icon color="green" icon="mdi-pen-plus"></v-icon>
        </template>
        <v-container fluid>
          <v-row dense class="d-flex align-center">
            <v-col cols="5" md="5" sm="5">Topic Name*:</v-col>
            <v-col cols="7" md="7" sm="7">
                <v-text-field hide-details="auto" v-model="topic_name" placeholder="mytopic"></v-text-field>
            </v-col>
          </v-row>
          <v-row dense class="d-flex align-center">
            <v-col cols="5" md="5" sm="5">Number of Partitions*:</v-col>
            <v-col cols="7" md="7" sm="7">
                <v-text-field hide-details="auto" v-model="partitions" placeholder="1"></v-text-field>
            </v-col>
          </v-row>
          <v-row dense class="d-flex align-center">
            <v-col cols="5" md="5" sm="5">Number of Replicas*:</v-col>
            <v-col cols="7" md="7" sm="7">
                <v-text-field hide-details="auto" v-model="replicas" placeholder="1"></v-text-field>
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

    <v-snackbar v-model="snackbar" timeout=4000 :color="snackcolor" elevation="24">
      {{ snacktext }}
      <template v-slot:actions>
        <v-btn color="grey" variant="text" @click="snackbar = false">Close</v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>


<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from 'vue-router';
import { globalTopicNames } from "../datas/kafka";
import { CreateTopic, ListTopics } from "../wailsjs/go/backend/KafkaTool";


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
let snackcolor = 'deep-purple-darken-4';

const refresh = () => {
  ListTopics().then((items: Array<string>) => {
    globalTopicNames.value = items;
  }).catch((err: string) => {
    showSnackBar('KafkaTool.ListTopics failed: '+err, false);
  });
}

const rowClicked = (row: string) => {
  router.push({
    name: 'Topic',
    query: { topic: row }
  });
}

const valid = () => {
  topic_name.value = topic_name.value.trim()
  partitions.value = partitions.value.trim()
  replicas.value = replicas.value.trim()
  if (topic_name.value.length == 0) {
    showSnackBar('topic name con not be empty', false);
    return false;
  }
  if (isNaN(Number(partitions.value)) || isNaN(Number(replicas.value))) {
    showSnackBar('partitions or replicas is not number', false);
    return false;
  }

  return true
}

const createTopic = () => {
  snackbar.value = false;
  if (!valid()) return;
  // console.log('createTopic', topic_name.value, partitions.value, replicas.value);

  CreateTopic(topic_name.value, Number(partitions.value), Number(replicas.value)).then(() => {
    showSnackBar('create topic success!', true);
    refresh();
  }).catch((err: string) => {
    showSnackBar('create topic failed: ' + err, false);
  });
}

const showSnackBar = (text: string, success: boolean) => {
    snackbar.value = false;
    snacktext = text;
    snackcolor = success ? 'deep-purple-darken-4' : 'deep-orange-darken-3';
    snackbar.value = true;
}
</script>
