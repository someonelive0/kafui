<template>
  <v-container fluid class="pa-1 ma-1">
    <v-card flat>
      <v-card-title class="d-flex align-center pe-2">
        <v-icon icon="mdi-list-box-outline"></v-icon> &nbsp;
          Kafka Connection Config items {{ globalKafkaConfigs.length }}
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
        ><v-tooltip activator="parent" location="bottom">Filter</v-tooltip>
        </v-text-field>&nbsp;
        <!-- <v-btn icon="mdi-refresh" size="small" @click="refresh"></v-btn>&nbsp; -->
        <v-tooltip text="New Kafka Connection" location="bottom">
          <template v-slot:activator="{ props }">
            <v-btn v-bind="props" icon="mdi-plus" size="small" @click="addClick"></v-btn>
          </template>
        </v-tooltip>
      </v-card-title>

      <v-data-table density="compact"
        :headers="headers"
        :items="globalKafkaConfigs"
        :search="search"
        :items-per-page="-1"
        hover
      >
        <template v-slot:item.name="{ item }">
          <v-chip :text="item.name" border="thin opacity-25" prepend-icon="mdi-database" label
            @click="rowEdit(item)">
            <template v-slot:prepend>
              <v-icon color="medium-emphasis"></v-icon>
            </template>
          </v-chip>
        </template>

        <template v-slot:item.actions="{ item }">
          <div class="d-flex ga-2 justify-end">
            <v-icon color="medium-emphasis" icon="mdi-pencil" size="small" @click="rowEdit(item)"></v-icon>
            <v-icon color="medium-emphasis" icon="mdi-delete" size="small" @click="rowDelete(item)"></v-icon>
          </div>
        </template>

        <template #bottom>
          <!-- Leave this slot empty to hide pagination controls -->
        </template>
      </v-data-table>
    </v-card>

    <v-dialog v-model="newDialog" width="600">
      <v-card
        max-width="500"
        text="Inuput connection name, can select existed conntion as template. Then change items in end of list in table."
        title="Create new kafka connection"
      >
        <template v-slot:prepend>
          <v-icon color="green" icon="mdi-pen-plus"></v-icon>
        </template>
        <v-container fluid>
          <v-row dense class="d-flex align-center">
            <v-col cols="5" md="5" sm="5">Connectin Name *:</v-col>
            <v-col cols="7" md="7" sm="7">
                <v-text-field hide-details="auto" v-model="newKafkaName" placeholder="mydb"
                  persistent-hint hint="* Self define connectin name"></v-text-field>
            </v-col>
          </v-row>
          <v-row dense class="d-flex align-center">
            <v-col cols="5" md="5" sm="5">Use Template :</v-col>
            <v-col cols="7" md="7" sm="7">
              <v-select density="default" label="Existed conntions"
                :items="tplNames" required v-model="dupName"
                persistent-hint hint="Copy from existed conntions"></v-select>
            </v-col>
          </v-row>
        </v-container>
        <template v-slot:actions>
          <v-spacer></v-spacer>
          <v-btn color="blue-darken-4" rounded="0" variant="outlined" text="Cancel" @click="newDialog = false"></v-btn>
          <v-btn color="blue-darken-4" rounded="0" variant="flat" text="Create" @click="createConn"></v-btn>
        </template>
      </v-card>
    </v-dialog>

    <v-dialog v-model="deleteDialog" max-width="500" >
      <v-card
        text="This kafka connectin will be delete, and lost all configs in this connection."
        :title="'Really delete connectin [' + deleteKafkaName + '] ?'" 
      >
        <template v-slot:prepend>
          <v-icon color="red" icon="mdi-alert"></v-icon>
        </template>
        <template v-slot:actions>
          <v-spacer></v-spacer>
          <v-btn border @click="deleteDialog = false">Cancel</v-btn>&nbsp;
          <v-btn border color="red-darken-4" @click="deleteConn">Delete</v-btn>
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

  <v-dialog v-model="configDialog" width="600">
    <KafkaConfig :kafka-name="selectedKafkaName" />
  </v-dialog>

</template>


<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from 'vue-router';
import { globalAddKafkaConfig, globalDeleteKafkaConfig, globalKafkaConfigs } from "../datas/global";
import { AddKafkaConfig, DeleteKafkaConfig } from "../wailsjs/go/backend/ConfigService";
import { backend } from "../wailsjs/go/models";
import KafkaConfig from './KafkaConfig.vue';


const headers: Array<object> = [
  { title: 'Config Name', align: 'start', sortable: true, key: 'name' },
  { title: 'Brokers', align: 'start', sortable: true, key: 'brokers' },
  { title: 'SASL Mechanism', align: 'start', sortable: true, key: 'sasl_mechanism' },
  { title: 'User', align: 'start', sortable: true, key: 'user' },
  // { title: 'Timeout', align: 'start', sortable: true, key: 'timeout' },
  { title: 'Operation', key: 'actions', align: 'end', sortable: false },
];

let search = ref('');
const router = useRouter(); 
let newDialog = ref(false);
let deleteDialog = ref(false);
var configDialog = ref(false);
let snackbar = ref(false);
let snacktext = '';
let snackcolor = 'deep-purple-darken-4';
let selectedKafkaName = ref('');
let deleteKafkaName = ref('');
let newKafkaName = ref('');
let dupName = ref('empty config'); // when add conn config, as config template
let tplNames = ref(['empty config']); // when add conn config, as config template list


const addClick = () => {
  tplNames.value = globalKafkaConfigs.value.map((item: backend.KafkaConfig) => item.name);
  tplNames.value.unshift('empty config');
  newDialog.value = true;
}

const rowEdit = (row : backend.KafkaConfig) => {
  console.log("rowEdit item: ", row);
  selectedKafkaName.value = row.name;
  configDialog.value = true;
}

const rowDelete = (row : backend.KafkaConfig) => {
  // console.log("rowDelete item: ", row);
  deleteKafkaName.value = row.name;
  deleteDialog.value = true;
}

const deleteConn = () => {
  // console.log("Connctions deleteConn : ", deleteKafkaName.value);
  deleteDialog.value = false;

  DeleteKafkaConfig(deleteKafkaName.value).then(() => {
    showSnackBar('Delete kafka connection success: '+ deleteKafkaName.value, true);
    // update global vars
    globalDeleteKafkaConfig(deleteKafkaName.value);
    newDialog.value = false;
  }).catch((err: string) => {
    showSnackBar('Delete kafka connection failed: '+ err, false);
  });
}

const valid = () => {
  newKafkaName.value = newKafkaName.value.trim()
  if (newKafkaName.value.length == 0) {
    showSnackBar('Kafka connection name is empty, set it', false);
    return false;
  }
  // if (isNaN(Number(partitions.value)) || isNaN(Number(replicas.value))) {
  //   showSnackBar('partitions or replicas is not number', false);
  //   return false;
  // }

  return true
}

const createConn = () => {
  snackbar.value = false;
  if (!valid()) return;
  // console.log('Connctions createConn', newKafkaName.value, dupName.value);

  AddKafkaConfig(newKafkaName.value, dupName.value).then((kafkaconfig : backend.KafkaConfig) => {
    showSnackBar('Create kafka connection success: '+kafkaconfig.name, true);
    // update global vars
    globalAddKafkaConfig(kafkaconfig);
    newDialog.value = false;
  }).catch((err: string) => {
    showSnackBar('Create kafka connection failed: '+err, false);
  });
}

const showSnackBar = (text: string, success: boolean) => {
    snackbar.value = false;
    snacktext = text;
    snackcolor = success ? 'deep-purple-darken-4' : 'deep-orange-darken-3';
    snackbar.value = true;
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
