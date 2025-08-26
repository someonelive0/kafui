<template>
  <v-container fluid class="pa-1 ma-1">
    <v-card flat>
      <v-card-title class="d-flex align-center pe-2">
        <v-icon icon="mdi-account-multiple"></v-icon> &nbsp;
        Consumer groups {{ groups.length }}
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
        :items="groups"
        :search="search"
        :items-per-page="-1"
        hover
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

  </v-container>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRoute, useRouter } from 'vue-router';


const { query, params } = useRoute();
// console.log('{ query, params } = useRoute() ', query, params);
const param_groups = ref(query.groups).value;

let groups: Array<object> = [];
// console.log('param_groups ', param_groups);
if (param_groups != null && param_groups != undefined && Array.isArray(param_groups)) {
  for (var i=0,len=param_groups.length; i<len; i++) {
    groups[i] = {name: param_groups[i]}
  }
}

const headers: Array<object> = [
  { title: 'Group Name', align: 'start', sortable: true, key: 'name' },
];
let search = ref('');

const router = useRouter(); 

const refresh = () => {
  window.go.backend.KafkaTool.ListGroups().then((items: Array<string>) => {
    // console.log('Kafkatool.ListGroups ', items);
    for (var i=0,len=items.length; i<len; i++) {
      groups[i] = {name: items[i]}
    }
    // loading.value = false;
  })
  .catch((err: string) => {
    console.error('Kafkatool.ListGroups ', err);
    // snacktext = 'read message failed: ' + err;
    // snackbar.value = true;
    // loading.value = false;
  });
}

const rowClicked = (row) => {
  // console.log("Clicked item: ", row)
  router.push({
    name: 'Group',
    query: {
        group: row.name
    }
  });
}

</script>
