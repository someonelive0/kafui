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
        ></v-text-field>
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

  </v-container>
</template>

<script setup lang="ts">
import { ref, defineProps } from "vue"
import { useRouter } from 'vue-router';


// 使用 stat 传递页面参数, 把字符串数组转换成对象数组
// let topics = window.history.state.topics;
// console.log('window.history.state.topics ', topics);
let topics: Array<object> = [];
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

const rowClicked = (row) => {
  // console.log("Clicked item: ", row)
  router.push({
    name: 'Topic',
    query: {
        topic: row.name
    }
  });
}

</script>
