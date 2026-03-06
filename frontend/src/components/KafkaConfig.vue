<template>
    <v-card prepend-icon="mdi-cog" title="Kafka Config" >
      <v-card-text>
        <v-row dense class="d-flex align-center">
            <v-col cols="4" md="4" sm="4">* Connection Name:</v-col>
            <v-col cols="8" md="8" sm="8">
                <v-text-field density="compact" disabled bg-color="blue-grey-darken-4"
                  prepend-inner-icon="mdi-rename"
                  :rules="stringRules" hide-details="auto" v-model="newKafkaConfig.name"
                  placeholder="mykafka"
                  persistent-hint hint="self define connection name"></v-text-field>
            </v-col>
        </v-row>

        <v-row dense class="d-flex align-center">
            <v-col cols="4" md="4" sm="4">* Brokers:</v-col>
            <v-col cols="8" md="8" sm="8">
                <v-text-field density="compact" bg-color="light-green-lighten-4"
                  prepend-inner-icon="mdi-server"
                  hide-details="auto" v-model="newKafkaConfig.brokers"
                  placeholder="localhost:9092"
                  persistent-hint hint="Example: broker1:9092,broker2:9092"></v-text-field>
            </v-col>
        </v-row>

        <v-row dense class="d-flex align-center">
            <v-col cols="4" md="4" sm="4">SASL Mechanism:</v-col>
            <v-col cols="8" md="8" sm="8">
                <v-select density="compact" bg-color="yellow-lighten-4"
                  prepend-inner-icon="mdi-check-bold"
                  :items="['None', 'PLAIN', 'SCRAM-SHA-256', 'SCRAM-SHA-512']"
                  required v-model="newKafkaConfig.sasl_mechanism"
                  persistent-hint hint="None means not use SASL"></v-select>
            </v-col>
        </v-row>

        <v-row dense class="d-flex align-center">
            <v-col cols="4" md="4" sm="4">User & Password:</v-col>
            <v-col cols="4" md="4" sm="4">
                <v-text-field density="compact"
                  bg-color="lime-lighten-4" prepend-inner-icon="mdi-account"
                  label="User" v-model="newKafkaConfig.user"></v-text-field>
            </v-col>

            <v-col cols="4" md="4" sm="4">
                <v-text-field density="compact" 
                    :type="showpwd ? 'text' : 'password'"
                    :append-icon="showpwd ? 'mdi-eye' : 'mdi-eye-off'"
                    @click:append="showpwd = !showpwd"
                    bg-color="lime-lighten-4"
                    label="Password"  v-model="newKafkaConfig.password"></v-text-field>
            </v-col>
        </v-row>

        <v-row dense class="d-flex align-center">
            <v-col cols="4" md="4" sm="4">Timeout:</v-col>
            <v-col cols="4" md="4" sm="4">
                <v-text-field density="compact"
                  bg-color="blue-grey-lighten-5" prepend-inner-icon="mdi-clock-time-five"
                  label="Timeout" v-model="newKafkaConfig.timeout"></v-text-field>
            </v-col>
        </v-row>

        <small class="text-caption text-medium-emphasis">*indicates required field</small>
      </v-card-text>

        <v-divider></v-divider>

        <v-card-actions>
            <v-btn color="secondary" text="Test Connction" variant="tonal" 
              prepend-icon="mdi-connection" @click="test"></v-btn>
            <v-spacer></v-spacer>
            <v-btn text="Cancel" variant="plain" @click="cancel"></v-btn>
            <v-btn color="primary" text="Save" variant="tonal" 
              prepend-icon="mdi-check-circle" @click="save"></v-btn>
        </v-card-actions>
    </v-card>

    <v-snackbar v-model="snackbar" timeout=4000 :color="snackcolor" elevation="24">
        {{ snacktext }}
        <template v-slot:actions>
        <v-btn color="grey" variant="text" @click="snackbar = false">Close</v-btn>
        </template>
    </v-snackbar>
</template>


<script setup lang="ts">
import { defineProps, onMounted, ref } from "vue";
import { globalUpdateKafkaConfig } from "../datas/global";
import { GetKafkaConfig, TestKafkaConfig, UpdateKafkaConfig } from "../wailsjs/go/backend/ConfigService";
import { backend } from "../wailsjs/go/models";


// 属性绑定，参考 <KafkaConfig :kafka-name="kafkaName" />
const { kafkaName } = defineProps(['kafkaName']); // 可以简写 解构
// console.log('KafkaConfig kafka-name', kafkaName);

const showpwd = ref(false);
let snackbar = ref(false);
let snacktext = '';
let snackcolor = 'deep-purple-darken-4';
let oldKafkaConfig = {} as backend.KafkaConfig;
let newKafkaConfig = ref({} as backend.KafkaConfig);

const stringRules = [
    (value: string) => !!value || 'Required.',
    (value: string) => (value && value.length >= 2) || 'Min 2 characters',
];

onMounted(() => {
    // console.log("Kafkaconfig init connName: " + connName);
    GetKafkaConfig(kafkaName).then((kafkaconfig : backend.KafkaConfig) => {
        oldKafkaConfig = JSON.parse(JSON.stringify(kafkaconfig)); // deep copy old conn config
        newKafkaConfig.value = kafkaconfig;
    }).catch((err: string) => {
        showSnackBar('ConfigService.GetKafkaConfig faile: '+ err, false);
    });
});

const cancel = () => {
    let tmpKafkaConfig = {} as backend.KafkaConfig;
    tmpKafkaConfig = JSON.parse(JSON.stringify(oldKafkaConfig)); // deep copy old conn config to a tmp object
    newKafkaConfig.value = tmpKafkaConfig;
}

const valid = () => {
    if (newKafkaConfig.value.brokers.length == 0) {
        showSnackBar('brokers con not be empty', false);
        return false;
    }
    // 由于编辑框会把[]string 变成 string，所以需要转换一下
    if (typeof newKafkaConfig.value.brokers == 'string') {
        const tmp = newKafkaConfig.value.brokers as string;
        const brokers = tmp.split(',');
        newKafkaConfig.value.brokers = brokers;
    }
    // 由于编辑框会把 number 变成 string，所以需要转换一下
    if (typeof newKafkaConfig.value.timeout == 'string') {
        let tmp = newKafkaConfig.value.timeout as string;
        tmp = tmp.trim();
        if (tmp.length == 0 || isNaN(Number(tmp)) ) {
            showSnackBar('Error: timeout must be numeric!', false);
            return;
        }
        newKafkaConfig.value.timeout = parseInt(tmp);
    }
    return true
}

const save = () => {
    if (!valid()) return;

    UpdateKafkaConfig(newKafkaConfig.value).then(() => {
        showSnackBar('Save config success', true);
        globalUpdateKafkaConfig(newKafkaConfig.value);
    }).catch((err: string) => {
        showSnackBar('Save config failed: '+ err, false);
    });
}

const test = () => {
    console.log('newKafkaConfig1: ', newKafkaConfig.value);
    if (!valid()) return;
    console.log('newKafkaConfig2: ', newKafkaConfig.value);

    TestKafkaConfig(newKafkaConfig.value).then((leader: backend.Broker) => {
        showSnackBar('Test connection success! Leader is ' + leader.host + ':' + leader.port, true);
    }).catch((err: string) => {
        showSnackBar('Test connection faile: ' + err, false);
    });

    // GetApiVersions(newKafkaConfig.value).then((versions: string) => {
    //     console.log("ApiVersions: "+versions);
    // }).catch((err: string) => {
    //     console.log("GetApiVersions failed: "+err);
    // });
}

const showSnackBar = (text: string, success: boolean) => {
    snackbar.value = false;
    snacktext = text;
    snackcolor = success ? 'deep-purple-darken-4' : 'deep-orange-darken-3';
    snackbar.value = true;
}

</script>
