<template>
  <v-data-table
    :headers="headers"
    :items="offsets"
    :search="search"
    :loading="loading"
    :items-per-page="-1"
    density="compact"
    item-value="topic"
  >
    <template #bottom>
      <!-- Leave this slot empty to hide pagination controls -->
    </template>
  </v-data-table>
</template>

<script setup lang="ts">
import { ref, defineProps } from 'vue';
import {onBeforeMount,onMounted,onBeforeUpdate,onUnmounted} from "vue"


const { name } = defineProps(['name']) // 可以简写 解构

const headers = [
      { title: 'Topic', align: 'start', sortable: true, key: 'topic' },
      { title: 'Partition', align: 'end', key: 'partition' },
      { title: 'First Offset', align: 'end', key: 'first_offset' },
      { title: 'Last Offset', align: 'end', key: 'last_offset' },
      { title: 'Committed Offset', align: 'end', key: 'committed_offset' },
      { title: 'Metadata', align: 'end', key: 'metadata' },
    ];

let offsets: Array<object> = ref([]);
let loading =ref(true);
let search = ref('');
const sortBy = [{ key: 'config_name', order: 'asc' }];

onMounted(() => {
  window.go.backend.KafkaTool.GetGroupOffset(name).then(items => {
    console.log('Kafkatool.GetGroupOffset ', items);
    if (items != null) {
      offsets.value = items;
    }
    loading.value = false;
  })
  .catch(err => {
    console.error('Kafkatool.GetGroupOffset ', err);
    loading.value = false;
  });
})

</script>
