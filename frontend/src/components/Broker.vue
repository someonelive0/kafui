<template>
  <v-container fluid class="pa-1 ma-1">
    <v-row no-gutters>
      <v-col>
        <v-breadcrumbs :items="breadcrumbs" class="bg-grey-lighten-5 pa-1 ma-1">
          <template v-slot:prepend>
            <v-icon icon="mdi-fridge" size="small"></v-icon>
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
        <v-tab value="Properties">Properties</v-tab>
        <v-tab value="Config">Config</v-tab>
      </v-tabs>
      
      <v-window v-model="selectedTab">
        <v-window-item key="Properties" value="Properties">
          <v-container fluid>
            <v-card class="pa-1 ma-1"
              variant="tonal"
              theme="dark"
              >
              <v-card-item>
                <div>
                  <div class="text-overline mb-1">
                    ID: {{ broker.id }}
                  </div>
                  <div class="text-h6 mb-1">
                    Host: {{ broker.host }}
                  </div>
                  <div class="text-caption">
                    Port: {{ broker.port }}
                  </div>
                  <div class="text-caption">
                    Rack: {{ broker.rack }}
                  </div>
                </div>
              </v-card-item>
              <!-- <v-card-actions>
                <v-btn @click="window.WindowReload()">Click me</v-btn>
              </v-card-actions> -->
            </v-card>
          </v-container>
        </v-window-item>

        <v-window-item key="Config" value="Config">
          <v-container fluid>
            <Config title="broker" :name="broker.id" />
          </v-container>
        </v-window-item>

      </v-window>
    </v-card>

  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import Config from './Config.vue';


// 使用 stat 传递页面参数
let broker = window.history.state.broker;
// console.log('window.history.state.broker ', broker);
var breadcrumbs = ref([
  { title: 'Broker', disabled: false, },
  { title: broker.id, disabled: false, }
]);
let selectedTab = ref("Properties"); // 默认选中 Properties 页

</script>
