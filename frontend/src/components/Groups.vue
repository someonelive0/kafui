<template>
  <v-container fluid class="pa-1 ma-1">
    <v-card flat>
      <v-card-title class="d-flex align-center pe-2">
        <v-icon icon="mdi-account-multiple"></v-icon> &nbsp;
        Consumer groups {{ groups.length }}
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

      <v-data-table density="compact"
        :headers="headers"
        :items="groups"
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

  </v-container>
</template>

<script setup lang="ts">
import { ref, defineProps } from "vue"
import { useRouter } from 'vue-router';


// 使用 stat 传递页面参数
// let groups = window.history.state.groups;
// console.log('window.history.state.groups ', groups);
let groups: Array<object> = [];
if (window.history.state.groups) {
  for (var i=0,len=window.history.state.groups.length; i<len; i++) {
    groups[i] = {name: window.history.state.groups[i]}
  }
}

const headers: Array<object> = [
  { title: 'Group Name', align: 'start', sortable: true, key: 'name' },
];
let search = ref('');

const router = useRouter(); 

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
