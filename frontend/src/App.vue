<template>
  <v-layout class="rounded rounded-md">

    <v-navigation-drawer
        permanent
        v-model="drawer"
        :rail="rail"
        @click="rail = false"
    >

      <v-list density="compact">
        <v-list-item @click="refresh" :title="connection_name" :subtitle="connection_addr">
          <template v-slot:prepend>
            <v-avatar color="grey-lighten-1">
              <v-icon color="white">mdi-refresh</v-icon>
            </v-avatar>
          </template>
        </v-list-item>
      </v-list>

      <v-divider></v-divider>

      <v-list density="compact" nav>
        <v-list-item prepend-icon="mdi-view-dashboard" title="Dashboard" value="inbox" rounded="shaped" size="x-small" class="customPrepend"
        @click="gotoDashboard" ></v-list-item>

        <v-list-group value="Brokers" >
          <template v-slot:activator="{ props }">
            <v-list-item color="success" class="customPrepend"
              v-bind="props"
              prepend-icon="mdi-server"
              title="Brokers"
              @click="gotoBrokers()"
            ></v-list-item>
          </template>

          <v-list-item rounded="shaped" size="x-small" color="warning" class="customPrepend"
            v-for="(broker, i) in brokers"
            :key="i"
            prepend-icon="mdi-fridge"
            :title="broker.host+':'+broker.port"
            :value="broker.id" @click="gotoBroker(broker, i)"
          ></v-list-item>
        </v-list-group>

        <v-list-group value="Topics" >
          <template v-slot:activator="{ props }">
            <v-list-item color="success" class="customPrepend"
              v-bind="props"
              prepend-icon="mdi-list-box-outline"
              title="Topics"
              @click="gotoTopics()"
            ></v-list-item>
          </template>

           <v-list-item rounded="shaped" size="x-small" color="warning" class="customPrepend"
            v-for="(topic, i) in topics"
            :key="i"
            prepend-icon="mdi-book-open-variant-outline"
            :title="topic"
            :value="topic" @click="gotoTopic(topic, i)"
          ></v-list-item>
          
        </v-list-group>

        <v-list-group value="Customers" >
          <template v-slot:activator="{ props }">
            <v-list-item color="success" class="customPrepend"
              v-bind="props"
              prepend-icon="mdi-account-multiple"
              title="Customers"
              @click="gotoGroups()"
            ></v-list-item>
          </template>

          <v-list-item rounded="shaped" size="x-small" color="warning" class="customPrepend"
            v-for="(group, i) in groups"
            :key="i"
            prepend-icon="mdi-account-file-text-outline"
            :title="group"
            :value="group+'_group'" @click="gotoGroup(group, i)"
          ></v-list-item>
          
        </v-list-group>
      </v-list>

      <template v-slot:append>
        <div class="pa-2">
          <v-btn density="compact" size="small" block>
            Copyright @ 2024
          </v-btn>
        </div>
      </template>
    </v-navigation-drawer>

    <v-app-bar :elevation="2" density="compact">
      <template v-slot:prepend>
        <v-app-bar-nav-icon @click.stop="rail = !rail"></v-app-bar-nav-icon>
      </template>

      <v-app-bar-title>Kafui</v-app-bar-title>

      <template v-slot:append>
        <v-btn icon="mdi-cog" @click="setting"></v-btn>
        <v-btn icon="mdi-magnify"></v-btn>
        <!-- <v-btn icon="mdi-dots-vertical"></v-btn> -->
        <v-menu>
            <template v-slot:activator="{ props }">
              <v-btn icon="mdi-dots-vertical" v-bind="props"></v-btn>
            </template>
            <v-list density="compact">
              <v-list-item density="compact" prepend-icon="mdi-information" title="About" @click="about()" />
            </v-list>
          </v-menu>
      </template>
    </v-app-bar>

    <!-- <v-main class="d-flex align-center justify-center" style="min-height: 300px;"> -->
    <v-main style="min-height: 300px;">
      <router-view :key="route.query"/>
    </v-main>

    <!-- <v-footer name="footer" density="compact" app
      class="bg-teal text-center d-flex flex-column"
    >
      <div class="bg-teal d-flex w-100 align-center px-4">
        {{ new Date().getFullYear() }} — <strong>Kafui</strong>

        <v-spacer></v-spacer>
        <v-icon icon="mdi-home" size="x-small" />
        <v-icon icon="mdi-calendar" size="x-small" />
        <v-icon icon="mdi-paperclip" size="x-small" />
      </div>
    </v-footer> -->

  </v-layout>

  <v-dialog v-model="setting_dialog" width="600">
    <Setting :myconfig="myconfig" @settingCancel="settingCancel" @settingSave="settingSave"/>
  </v-dialog>

  <v-dialog v-model="about_dialog" width="auto">
    <About />
  </v-dialog>

  <v-snackbar v-model="snackbar" timeout=2000 color="deep-purple-accent-4" elevation="24">
    {{ snacktext }}
    <template v-slot:actions>
      <v-btn color="pink" variant="text" @click="snackbar = false" >Close</v-btn>
    </template>
  </v-snackbar>

</template>

<script setup lang="ts">
import { ref } from 'vue';
import {onBeforeMount,onMounted,onBeforeUpdate,onUnmounted} from "vue"
import { useRouter, useRoute } from 'vue-router';
import Setting from './components/Setting.vue'
import About from './components/About.vue'


const router = useRouter(); 
const route = useRoute(); 
const drawer = ref(true);
const rail = ref(false);
var brokers = ref([]);
var topics = ref([]);
var groups = ref([]);
var setting_dialog = ref(false);
var about_dialog = ref(false);
var myconfig = ref(null);
var connection_name = ref('');
var connection_addr = ref('');
let snackbar = ref(false);
let snacktext = '';

onMounted(() => {
  getMyconfig();
});

const refresh = () => {
  getBrokers();
  getTopics();
  getGroups();
}

const getMyconfig = () => {
  window.go.main.App.GetMyconfig().then(item => {
    console.log('App.GetMyconfig ', item);
    myconfig.value = item;
    connection_name.value = item.kafka.name;
    connection_addr.value = item.kafka.brokers[0];
  })
  .catch(err => {
    console.error('KafkaTool.ListTopics', err);
  });
}

const getBrokers = () => {
  // window.go.backend.ZkTool.ListBrokers(zk_hosts).then(items => {
  //   console.log('ZkTool.ListBrokers ', items);
  //   brokers = items
  // })
  // .catch(err => {
  //   console.error('ZkTool.ListBrokers ', err);
  // });

  window.go.backend.KafkaTool.ListBrokers().then(items => {
    console.log('Kafkatool.ListBrokers ', items);
    brokers = items;
    snacktext = 'get brokers success!';
    snackbar.value = true;
  })
  .catch(err => {
    console.error('Kafkatool.ListBrokers ', err);
    snacktext = 'get brokers failed: ' + err;
    snackbar.value = true;
  });
}

const getTopics = () => {
  // window.go.backend.ZkTool.ListTopics(zk_hosts).then(items => {
  //   console.log('ZkTool.ListTopics ', items);
  //   topics = items
  // })
  // .catch(err => {
  //   console.error('ZkTool.ListTopics', err);
  // });

  window.go.backend.KafkaTool.ListTopics().then(items => {
    console.log('KafkaTool.ListTopics ', items);
    topics = items
  })
  .catch(err => {
    console.error('KafkaTool.ListTopics', err);
  });
}

const getGroups = () => {
  window.go.backend.KafkaTool.ListGroups().then(items => {
    console.log('KafkaTool.ListGroups ', items);
    groups = items
  })
  .catch(err => {
    console.error('KafkaTool.ListGroups', err);
  });
}

const gotoDashboard = () => {
  router.push({
    name:'Dashboard',
    query: {
        num_brokers: brokers.length,
        num_topics: topics.length,
        num_groups: groups.length
    }
  });
}

const gotoBroker = (broker, i) => {
  // console.log('选择 broker ', broker, i);
  router.push({
    name: 'Broker',
    state: { broker: broker }
  });
}

const gotoBrokers = () => {
  router.push({
    name: 'Brokers',
    state: { brokers: brokers }
  });
}

const gotoTopic = (topic: string, i) => {
  // console.log('选择 topic ', topic, i);
  router.push({
    name: 'Topic',
    query: {
        id: i,
        topic: topic
    }
  });
}

const gotoTopics = () => {
  router.push({
    name: 'Topics',
    state: { topics: topics }
  });
}

const gotoGroup = (group: string, i) => {
  // console.log('选择 group ', group, i);
  router.push({
    name: 'Group',
    query: {
        id: i,
        group: group
    }
  });
}

const gotoGroups = () => {
  router.push({
    name: 'Groups',
    state: { groups: groups }
  });
}

const setting = () => {
  setting_dialog.value = true;
}
const settingCancel = () => {
  setting_dialog.value = false;
}
const settingSave = (item: any) => {
  console.log(item);
  setting_dialog.value = false;
}

const about = () => {
  about_dialog.value = true;
}
</script>

<style scoped>
.v-list-item--density-compact.v-list-item--one-line {
  min-height: 30px;
  font-size: 12px;
}
.v-list-group {
  --list-indent-size: 8px;
  --prepend-width: 0px;
}
.v-list-item__spacer {
  width: 8px;
}
list-item__prepend>.v-icon~.v-list-item__spacer, .v-list-item__prepend>.v-tooltip~.v-list-item__spacer {
  width: 8px;
}

.customPrepend :deep(.v-list-item__prepend .v-list-item__spacer) {
  width: 8px;
}
</style>