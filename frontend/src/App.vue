<template>
  <v-layout class="rounded rounded-md">

    <v-navigation-drawer
        permanent
        v-model="drawer"
        :rail="rail"
        @click="rail = false"
    >

      <v-list density="compact" style="padding: 0px;">
        <v-list-item @click="connect()" :title="connection_name" :subtitle="connection_addr"
          :class="{ 'active-connection': kafkaConnected === 1 }">
          <template v-slot:prepend>
            <v-avatar :color="iconColor">
              <v-icon color="white">mdi-connection</v-icon>
            </v-avatar>
          </template>
        </v-list-item>
      </v-list>

      <!-- <v-divider></v-divider> -->

      <v-toolbar flat density="compact" height="50"
        rounded="shaped" border
        color="grey-lighten-1" class="pl-4 ma-0">
        <v-tooltip text="Refresh config" location="bottom">
          <template v-slot:activator="{ props }">
            <v-btn v-bind="props" density="compact" size="small" icon="mdi-apache-kafka"
              @click="refresh()"></v-btn>
          </template>
        </v-tooltip>
        <!-- <v-tooltip text="导出" location="bottom">
          <template v-slot:activator="{ props }">
            <v-btn v-bind="props" density="compact" size="small" icon="mdi-export" 
              @click="gotoRoute('Connections')"></v-btn>
          </template>
        </v-tooltip> -->
        <v-select v-model="currentKafkaName"
          variant="underlined"
          class="ma-2 pa-2"
          density="default"
          :items="globalKafkaConfigNames"
        ></v-select>
      </v-toolbar>

      <v-list density="compact" nav>
        <v-list-item prepend-icon="mdi-view-dashboard" title="Dashboard" value="inbox" rounded="shaped" size="x-small" class="customPrepend"
          @click="gotoRoute('Dashboard')" ></v-list-item>

        <v-list-group value="Brokers" >
          <template v-slot:activator="{ props }">
            <v-list-item color="success" class="customPrepend"
              v-bind="props"
              prepend-icon="mdi-server"
              title="Brokers"
              @click="gotoRoute('Brokers')"
            ></v-list-item>
          </template>

          <v-list-item rounded="shaped" size="x-small" color="warning" class="customPrepend"
            v-for="(broker, i) in globalBrokers"
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
              @click="gotoRoute('Topics')"
            ></v-list-item>
          </template>

           <v-list-item rounded="shaped" size="x-small" color="warning" class="customPrepend"
            v-for="(topic, i) in globalTopicNames"
            :key="i"
            prepend-icon="mdi-book-open-variant-outline"
            :title="topic"
            :value="topic" @click="gotoTopic(topic, i)"
          ></v-list-item>
          
        </v-list-group>

        <v-list-group value="Consumer Groups" >
          <template v-slot:activator="{ props }">
            <v-list-item color="success" class="customPrepend"
              v-bind="props"
              prepend-icon="mdi-account-multiple"
              title="Consumer Groups"
              @click="gotoRoute('Groups')"
            ></v-list-item>
          </template>

          <v-list-item rounded="shaped" size="x-small" color="warning" class="customPrepend"
            v-for="(group, i) in globalGroupNames"
            :key="i"
            prepend-icon="mdi-account-file-text-outline"
            :title="group"
            :value="group+'_group'" @click="gotoGroup(group, i)"
          ></v-list-item>
          
        </v-list-group>

        <!-- <v-list-item prepend-icon="mdi-network-outline" title="ZooKepper" value="inbox" rounded="shaped" size="x-small" class="customPrepend"
        @click="gotoZooKeeper" ></v-list-item> -->

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
        <v-btn icon="mdi-cog" @click="gotoRoute('Connections')"></v-btn>
        <!-- <v-btn icon="mdi-magnify"></v-btn> -->
        <v-menu>
            <template v-slot:activator="{ props }">
              <v-btn icon="mdi-dots-vertical" v-bind="props"></v-btn>
            </template>
            <v-list density="compact">
              <v-list-item density="compact" prepend-icon="mdi-information" title="About" @click="about()" />
              <v-list-item density="compact" prepend-icon="mdi-hammer-screwdriver" title="Kcat" @click="router.push({ name: 'Kcat', })" />
            </v-list>
          </v-menu>
      </template>
    </v-app-bar>

    <v-main style="min-height: 300px;">
      <router-view :key="$route.fullPath"/>
    </v-main>

  </v-layout>

  <v-dialog v-model="about_dialog" width="auto">
    <About />
  </v-dialog>

  <v-snackbar v-model="snackbar" timeout=4000 :color="snackcolor" elevation="24">
    {{ snacktext }}
    <template v-slot:actions>
      <v-btn color="grey" variant="text" @click="snackbar = false">Close</v-btn>
    </template>
  </v-snackbar>

</template>


<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import About from './components/About.vue';
import { globalKafkaConfigNames, globalSetKafkaConfigs } from "./datas/global";
import { globalBrokers, globalGroupNames, globalResetKafkaDatas, globalSetBrokers, globalTopicNames } from "./datas/kafka";
import { GetKafkaConfig, GetKafkaConfigs } from "./wailsjs/go/backend/ConfigService";
import { Init, ListBrokers, ListGroups, ListTopics } from "./wailsjs/go/backend/KafkaTool";
import { backend } from "./wailsjs/go/models";


let currentKafkaName = ref('');
let kafkaConnected = 0;
let iconColor = "grey";
const router = useRouter(); 
const route = useRoute(); 
const drawer = ref(true);
const rail = ref(false);
var about_dialog = ref(false);
var connection_name = ref('Select kafka name');
var connection_addr = ref('none');
let snackbar = ref(false);
let snacktext = '';
let snackcolor = 'deep-purple-darken-4';

onMounted(() => {
  refresh();
});

const refresh = () => {
  getMyconfig();
}

const connect = () => {
  if (currentKafkaName.value.length == 0) {
    showSnackBar('Please select kafka connection name', false);
    return;
  }
  globalResetKafkaDatas();

  GetKafkaConfig(currentKafkaName.value).then((kafkaconfig: backend.KafkaConfig) => {
    // console.log('ConfigService.GetKafkaConfig', kafkaconfig);
    connection_name.value = kafkaconfig.name;
    connection_addr.value = kafkaconfig.brokers.join();

    // 这里重新配置kafkatool的连接信息
    Init(kafkaconfig).then(() => {
      getBrokers();
      getTopics();
      getGroups();
      gotoRoute('Dashboard');
    }).catch((err: string) => {
      showSnackBar('Kafkatool.Init failed: '+ err, false);
      return;
    });
  }).catch((err: string) => {
    showSnackBar('ConfigService.GetKafkaConfig failed: '+ err, false);
    return;
  });
}

const getMyconfig = () => {
  GetKafkaConfigs().then((kafkaconfigs : backend.KafkaConfig[]) => {
    // console.log('ConfigService.GetKafkaConfigs', kafkaconfigs);
    globalSetKafkaConfigs(kafkaconfigs);
    if (kafkaconfigs.length > 0) currentKafkaName.value = kafkaconfigs[0].name;
    showSnackBar('GetKafkaConfigs success!', true);
  }).catch((err: string) => {
    showSnackBar('ConfigService.GetKafkaConfigs failed: '+ err, false);
  });
}

const getBrokers = () => {
  ListBrokers().then((items: backend.Broker[]) => {
    // console.log('Kafkatool.ListBrokers ', items);
    globalSetBrokers(items);
    kafkaConnected = 1;
    iconColor = "blue-darken-1";
    showSnackBar('get brokers success!', true);
  }).catch((err: string) => {
    // console.error('Kafkatool.ListBrokers ', err);
    showSnackBar('get brokers failed: ' + err, false);
    kafkaConnected = 0;
    iconColor = "grey";
  });
}

const getTopics = () => {
  ListTopics().then((items: Array<string>) => {
    // console.log('KafkaTool.ListTopics ', items);
    globalTopicNames.value = items;
  }).catch((err: string) => {
    // console.error('KafkaTool.ListTopics', err);
    showSnackBar('get topics failed: ' + err, false);
  });
}

const getGroups = () => {
  ListGroups().then((items: Array<string>) => {
    // console.log('KafkaTool.ListGroups ', items);
    globalGroupNames.value = items;
  }).catch((err: string) => {
    // console.error('KafkaTool.ListGroups', err);
    showSnackBar('get groups failed: ' + err, false);
  });
}

const gotoRoute = (routeName: string) => {
  router.push({
    name: routeName,
    query: { }
  });
}

const gotoBroker = (broker: backend.Broker, i: number) => {
  router.push({
    name: 'Broker',
    query: { broker_id: broker.id }
  });
}

const gotoTopic = (topic: string, i: number) => {
  router.push({
    name: 'Topic',
    query: { id: i, topic: topic }
  });
}

const gotoGroup = (group: string, i: number) => {
  router.push({
    name: 'Group',
    query: { id: i, group: group }
  });
}

const about = () => {
  about_dialog.value = true;
}

const showSnackBar = (text: string, success: boolean) => {
    snackbar.value = false;
    snacktext = text;
    snackcolor = success ? 'deep-purple-darken-4' : 'deep-orange-darken-3';
    snackbar.value = true;
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

/* change connection color by myself */
.active-connection {
  background-color:hwb(200 70% 5%) !important;
  /* Change this to the color you want */
  color: #110000 !important;
  /* Change the text color to match the background color */
}

</style>
