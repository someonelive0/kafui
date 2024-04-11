<template>
  <v-card flat>
    <v-card-title class="d-flex align-center pe-2">
      <v-icon icon="mdi-list-box-outline"></v-icon> &nbsp;
        Topic Name: {{ name }}
      <v-spacer></v-spacer>
        Message number: {{ number }}
    </v-card-title>

    <v-table fixed-header density="compact">
      <thead>
        <tr>
          <th class="text-left">
            Partition ID
          </th>
          <th class="text-left">
            Leader
          </th>
          <th class="text-left">
            First Offset
          </th>
          <th class="text-left">
            Last Offset
          </th>
          <th class="text-left">
            Size
          </th>
          <th class="text-left">
            Replicas
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="item in partitions"
          :key="item.id"
        >
          <td>{{ item.id }}</td>
          <td>{{ item.leader.host }}:{{ item.leader.port }}</td>
          <td>{{ item.first_offset }}</td>
          <td>{{ item.last_offset }}</td>
          <td>{{ item.number }}</td>
          <td><pre>{{ printReplicas(item.replicas) }}</pre></td>
        </tr>
      </tbody>
    </v-table>

  </v-card>
</template>

<script setup lang="ts">
import { ref, defineProps } from "vue"
import {onBeforeMount,onMounted,onBeforeUpdate,onUnmounted} from "vue"


const { name } = defineProps(['name']) // 可以简写 解构
let number = ref(0);
let partitions: Array<object> = ref([]);
let loading = true;

onMounted(() => {
  window.go.backend.KafkaTool.GetTopicPartition(name).then(items => {
    // console.log('Kafkatool.GetTopicPartition ', items);
    partitions.value = items;
    number.value = 0;
    var total = 0
    for(var i = 0, l = items.length; i < l; i++) {
      total = total + items[i].number;
    }
    number.value = total;
  })
  .catch(err => {
    console.error('Kafkatool.GetTopicPartition ', err);
  });
  loading = false;
});

const printReplicas = (replicas: Array<object>) => {
  var s = '';
  for (var i=0,l=replicas.length; i<l; i++) {
    if (i > 0) {
      s += '\n'
    }
    s += 'ID ' + replicas[i].id + ', ' + replicas[i].host + ':' + replicas[i].port;
  }
  return s;
}

</script>
