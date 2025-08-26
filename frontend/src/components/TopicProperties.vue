<template>
  <v-card flat>
    <v-card-title class="d-flex align-center pe-2">
      <v-icon icon="mdi-list-box-outline"></v-icon> &nbsp;
        Topic: {{ name }}
      <v-spacer></v-spacer>
        Message number: {{ number }}
      <v-spacer></v-spacer>
      <v-btn icon="mdi-refresh" size="small" @click="refresh"></v-btn>&nbsp;&nbsp;
      <v-dialog v-model="dialog" max-width="500" >
        <template v-slot:activator="{ props: activatorProps }">
          <v-btn color="red-darken-4" variant="outlined" v-bind="activatorProps">Delete This Topic</v-btn>
        </template>

        <v-card
          prepend-icon="mdi-map-marker"
          text="This topic will be deleted, and will lost all data."
          :title="'Really delete topic ' + name + ' ?'" 
        >
          <template v-slot:actions>
            <v-spacer></v-spacer>
            <v-btn @click="dialog = false">Cancel</v-btn>&nbsp;
            <v-btn color="red-darken-4" @click="deleteTopic">Delete</v-btn>
          </template>
        </v-card>

      </v-dialog>
    </v-card-title>

    <v-table fixed-header density="compact" hover>
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

    <v-snackbar v-model="snackbar" timeout=2000 color="deep-purple-accent-4" elevation="24">
      {{ snacktext }}
      <template v-slot:actions>
        <v-btn color="pink" variant="text" @click="snackbar = false" >Close</v-btn>
      </template>
    </v-snackbar>
  </v-card>
</template>

<script setup lang="ts">
import { defineProps, onMounted, ref } from "vue";
import { backend } from "../wailsjs/go/models";


const { name } = defineProps(['name']) // 可以简写 解构
let number = ref(0);
let partitions: Array<backend.Partition> = ref([]);
let loading = true;
let dialog = ref(false);
let snackbar = ref(false);
let snacktext = '';

onMounted(() => {
  refresh();
});

const refresh = () => {
  window.go.backend.KafkaTool.GetTopicPartition(name).then((items: Array<backend.Partition>) => {
    // console.log('Kafkatool.GetTopicPartition ', items);
    partitions.value = items;
    number.value = 0;
    var total = 0
    for(var i = 0, l = items.length; i < l; i++) {
      total = total + items[i].number;
    }
    number.value = total;
  })
  .catch((err: string) => {
    console.error('Kafkatool.GetTopicPartition ', err);
  });
  loading = false;
};

const printReplicas = (replicas: Array<backend.Broker>): string => {
  var s = '';
  for (var i=0,l=replicas.length; i<l; i++) {
    if (i > 0) {
      s += '\n'
    }
    s += 'ID ' + replicas[i].id + ', ' + replicas[i].host + ':' + replicas[i].port;
  }
  return s;
}

const deleteTopic = () => {
  // console.log('deleteTopic ', name);
  window.go.backend.KafkaTool.DeleteTopic(name).then(() => {
    snacktext = 'delete topic ' + name + ' success!';
    snackbar.value = true;
    dialog.value = false;
  })
  .catch((err: string) => {
    // console.error('Kafkatool.WriteMsg ', err);
    snacktext = 'delete topic ' + name + ' failed: ' + err;
    snackbar.value = true;
  });
}

</script>
