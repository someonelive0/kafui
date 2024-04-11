<template>
  <v-container fluid class="pa-1 ma-1">
    <v-row no-gutters>
      <v-col>
        <v-breadcrumbs :items="['Brokers', '']" class="bg-grey-lighten-5 pa-1 ma-1">
          <template v-slot:prepend>
            <v-icon icon="mdi-server" size="small"></v-icon>
          </template>
          <template v-slot:divider>
              <v-icon icon="mdi-chevron-right"></v-icon>
          </template>
        </v-breadcrumbs>
      </v-col>
    </v-row>

    <v-card>
      <v-tabs density="compact"
        v-model="selectedTab"
        align-tabs="start"
        color="deep-purple-accent-4"
        >
        <v-tab value="Menmbers">Cluster Menmbers</v-tab>
        <v-tab value="Config">Cluster Config</v-tab>
      </v-tabs>
      
      <v-window v-model="selectedTab">
        <v-window-item key="Menmbers" value="Menmbers">
          <v-container fluid>
            <v-table fixed-header density="compact">
              <thead>
                <tr>
                  <th class="text-left">Broker ID</th>
                  <th class="text-left">Host</th>
                  <th class="text-left">Port</th>
                  <th class="text-left">Rack</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="item in brokers"
                  :key="item.id"
                  @click="rowClicked(item)"
                >
                  <td>{{ item.id }}</td>
                  <td>{{ item.host }}</td>
                  <td>{{ item.port }}</td>
                  <td>{{ item.rack }}</td>
                </tr>
              </tbody>
            </v-table>
          </v-container>
        </v-window-item>

        <v-window-item key="Config" value="Config">
          <v-container fluid>
            <Config title="cluster" :name="1" />
          </v-container>
        </v-window-item>

      </v-window>
    </v-card>

  </v-container>
</template>

<script setup lang="ts">
import { ref, isRef } from 'vue';
import Config from './Config.vue'
import { useRouter } from 'vue-router';


// 使用 stat 传递页面参数
let brokers = window.history.state.brokers;
// console.log('window.history.state.broker ', brokers);
let selectedTab = ref("Menmbers"); // 默认选中 Menmbers 页
const router = useRouter(); 

const rowClicked = (row) => {
  // console.log("Clicked item: ", row)
  router.push({
    name: 'Broker',
    state: { broker: row }
  });
}
</script>
