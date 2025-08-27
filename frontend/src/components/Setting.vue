<template>
    <v-card prepend-icon="mdi-cog" title="Setting" >
      <v-card-text>
        <v-row dense class="d-flex align-center">
            <v-col cols="4" md="4" sm="4">Connection Name:</v-col>
            <v-col cols="8" md="8" sm="8">
                <v-text-field :rules="rules" hide-details="auto" v-model="name"></v-text-field>
            </v-col>
        </v-row>

        <v-row dense class="d-flex align-center">
            <v-col cols="4" md="4" sm="4">Brokers:</v-col>
            <v-col cols="8" md="8" sm="8">
                <v-text-field :rules="rules" hide-details="auto" v-model="brokers" hint="Example: broker1:9092,broker2:9092"></v-text-field>
            </v-col>
        </v-row>

        <v-row dense class="d-flex align-center">
            <v-col cols="4" md="4" sm="4">SASL Mechanism:</v-col>
            <v-col cols="8" md="8" sm="8">
                <v-select :items="['None', 'SASL_PLAINTEXT']" required v-model="sasl_mechanism"></v-select>
            </v-col>
        </v-row>

        <v-row dense class="d-flex align-center">
            <v-col cols="4" md="4" sm="4">User:</v-col>
            <v-col cols="4" md="4" sm="4">
                <v-text-field v-model="user"></v-text-field>
            </v-col>

            <v-col cols="4" md="4" sm="4">
                <v-text-field type="password" v-model="password"></v-text-field>
            </v-col>
        </v-row>

        <small class="text-caption text-medium-emphasis">*indicates required field</small>
      </v-card-text>

        <v-divider></v-divider>

        <v-card-actions>
            <v-btn color="secondary" text="Test Connction" variant="tonal" 
              prepend-icon="mdi-connection" @click="test"></v-btn>
            <v-spacer></v-spacer>
            <v-btn text="Close" variant="plain" @click="cancel"></v-btn>
            <v-btn color="primary" text="Save" variant="tonal" 
              prepend-icon="mdi-check-circle" @click="save"></v-btn>
        </v-card-actions>
    </v-card>

    <v-snackbar v-model="snackbar" timeout=4000 color="deep-purple-accent-4" elevation="24">
        {{ snacktext }}
        <template v-slot:actions>
        <v-btn color="pink" variant="text" @click="snackbar = false" >Close</v-btn>
        </template>
    </v-snackbar>
</template>

<script setup lang="ts">
import { defineEmits, defineProps, ref } from "vue";
import { TestKafka } from "../wailsjs/go/main/App";
import { backend } from "../wailsjs/go/models";


const rules = [
    (value: string) => !!value || 'Required.',
    (value: string) => (value && value.length >= 3) || 'Min 3 characters',
];

const { myconfig } = defineProps(['myconfig']); // 可以简写 解构
console.log('setting... ', myconfig);
// 调用defineEmits方法 并接受父组件给绑定的事件
const emit = defineEmits(['settingCancel', 'settingSave'])

let name = ref('');
let brokers = ref('');
let sasl_mechanism = ref('None');
let user = ref('');
let password = ref('');
let snackbar = ref(false);
let snacktext = '';

name.value = myconfig.kafka.name;
for (var i=0; i<myconfig.kafka.brokers.length; i++) {
    if (i>0) brokers.value += ',';
    brokers.value += myconfig.kafka.brokers;
}
sasl_mechanism.value = myconfig.kafka.sasl_mechanism;
user.value = myconfig.kafka.user;
password.value = myconfig.kafka.password; // password in undefined.

const cancel = () => {
    emit("settingCancel")
}

const valid = () => {
    name.value = name.value.trim()
    brokers.value = brokers.value.trim()
    if (name.value.length == 0) {
        snacktext = 'name con not be empty';
        snackbar.value = true;
        return false;
    }
    if (brokers.value.length == 0) {
        snacktext = 'brokers con not be empty';
        snackbar.value = true;
        return false;
    }
    return true
}

const save = () => {
    if (!valid()) return;
    const tmpconfig  = { // backend.Myconfig
        title: myconfig.title,
        license: myconfig.license,
        kafka: {
            name: name.value,
            brokers: brokers.value.split(','),
            sasl_mechanism: sasl_mechanism.value == 'None' ? '' : sasl_mechanism.value,
            user: user.value,
            password: password.value,
        },
        zookeeper: myconfig.zookeeper,
    }
    emit("settingSave", tmpconfig)
}

const test = () => {
    snackbar.value = false;
    if (!valid()) return;

    const kafka = {
        name: name.value,
        brokers: brokers.value.split(','),
        sasl_mechanism: sasl_mechanism.value == 'None' ? '' : sasl_mechanism.value,
        user: user.value,
        password: password.value,
    }

    TestKafka(kafka).then((leader: backend.Broker) => { // window.go.main.App.TestKafka
    snacktext = 'Test connection success! Leader is ' + leader.host + ':' + leader.port;
    snackbar.value = true;
  })
  .catch((err: string) => {
    snacktext = 'Test connection faile: ' + err;
    snackbar.value = true;
  });
}

</script>
