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
      <v-btn color="red-darken-4" @click="setGroupOffset">Change Offset</v-btn>
    </template>
  </v-data-table>

  <v-snackbar v-model="snackbar" timeout=2000 color="deep-purple-accent-4" elevation="24">
    {{ snacktext }}
    <template v-slot:actions>
      <v-btn color="pink" variant="text" @click="snackbar = false" >Close</v-btn>
    </template>
  </v-snackbar>

</template>

<script setup lang="ts">
import { ref, reactive, defineProps } from 'vue';
import {onBeforeMount,onMounted,onBeforeUpdate,onUnmounted} from "vue"
import { backend } from '../wailsjs/go/models';


const { name } = defineProps(['name']) // 可以简写 解构

const headers = [
      { title: 'Topic', align: 'start', sortable: true, key: 'topic' },
      { title: 'Partition', align: 'end', key: 'partition' },
      { title: 'First Offset', align: 'end', key: 'first_offset' },
      { title: 'Last Offset', align: 'end', key: 'last_offset' },
      { title: 'Committed Offset', align: 'end', key: 'committed_offset' },
      { title: 'Metadata', align: 'end', key: 'metadata' },
      { title: 'Operation', align: 'end', key: 'operation' },
    ];

let offsets: Array<backend.GroupOffset> = reactive([]);
let loading =ref(true);
let search = ref('');
const sortBy = [{ key: 'config_name', order: 'asc' }];
let snackbar = ref(false);
let snacktext = '';

onMounted(() => {
  window.go.backend.KafkaTool.GetGroupOffset(name).then((items: Array<backend.GroupOffset>) => {
    console.log('Kafkatool.GetGroupOffset ', items);
    if (items != null) {
      offsets = items;
    }
    loading.value = false;
  })
  .catch((err: string) => {
    console.error('Kafkatool.GetGroupOffset ', err);
    loading.value = false;
  });
})

const setGroupOffset = () => {
  console.log('setGroupOffset ', name);
  window.go.backend.KafkaTool.SetGroupOffset1(name).then(() => {
    snacktext = 'set group ' + name + ' success!';
    snackbar.value = true;
  })
  .catch((err: string) => {
    // console.error('Kafkatool.WriteMsg ', err);
    snacktext = 'set group ' + name + ' failed: ' + err;
    snackbar.value = true;
  });

}
</script>
