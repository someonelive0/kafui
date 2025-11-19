<template>
  <v-container fluid class="pa-1 ma-1">
    <v-card flat>
      <v-card-title class="d-flex align-center pe-2">
        <v-icon icon="mdi-account-multiple"></v-icon> &nbsp;
        Consumer groups {{ globalGroupNames.length }}
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
      </v-card-title>

      <v-data-table density="compact"
        :headers="headers"
        :items="globalGroupNames"
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

  </v-container>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from 'vue-router';
import { globalGroupNames } from "../datas/kafka";
import { ListGroups } from "../wailsjs/go/backend/KafkaTool";


const headers: Array<object> = [
  { title: 'Group Name', align: 'start', sortable: true, key: 'name' },
];
let search = ref('');
const router = useRouter(); 

const refresh = () => {
  ListGroups().then((items: Array<string>) => {
    globalGroupNames.value = items;
  }).catch((err: string) => {
    console.error('Kafkatool.ListGroups ', err);
  });
}

const rowClicked = (row: string) => {
  router.push({
    name: 'Group',
    query: { group: row }
  });
}

</script>
