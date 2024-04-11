<template>
    <v-card prepend-icon="mdi-cog" title="Setting" >
      <v-card-text>
        <v-row dense class="d-flex align-center">
            <v-col cols="4" md="4" sm="4">Connection Name:</v-col>
            <v-col cols="8" md="8" sm="8">
                <v-text-field required :value="name"></v-text-field>
            </v-col>
        </v-row>

        <v-row dense class="d-flex align-center">
            <v-col cols="4" md="4" sm="4">Brokers:</v-col>
            <v-col cols="8" md="8" sm="8">
                <v-text-field :value="brokers" hint="Example: broker1:9092,broker2:9092"></v-text-field>
            </v-col>
        </v-row>

        <v-row dense class="d-flex align-center">
            <v-col cols="4" md="4" sm="4">SASL Mechanism:</v-col>
            <v-col cols="8" md="8" sm="8">
                <v-autocomplete :value="sasl_mechanism"
                    :items="['None', 'SASL_PLAINTEXT']"
                    auto-select-first
                >
                </v-autocomplete>
            </v-col>
        </v-row>

        <v-row dense class="d-flex align-center">
            <v-col cols="4" md="4" sm="4">User:</v-col>
            <v-col cols="4" md="4" sm="4">
                <v-text-field :value="user"></v-text-field>
            </v-col>

            <v-col cols="4" md="4" sm="4">
                <v-text-field type="password" :value="password"></v-text-field>
            </v-col>
        </v-row>

        <small class="text-caption text-medium-emphasis">*indicates required field</small>
      </v-card-text>

        <v-divider></v-divider>

        <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn text="Close" variant="plain" @click="cancel"></v-btn>
            <v-btn color="primary" text="Save" variant="tonal" @click="save"></v-btn>
        </v-card-actions>
    </v-card>
</template>

<script setup lang="ts">
import { ref, defineProps, defineEmits } from "vue"


const { myconfig } = defineProps(['myconfig']) // 可以简写 解构
console.log('setting... ', myconfig);
// 调用defineEmits方法 并接受父组件给绑定的事件
const emit = defineEmits(['settingCancel', 'settingSave'])

let name = ref('');
let brokers = ref('');
let sasl_mechanism = ref('None');
let user = ref('');
let password = ref('');

name.value = myconfig.kafka.name;
for (var i=0; i<myconfig.kafka.brokers.length; i++) {
    brokers.value += myconfig.kafka.brokers + ',';
}
sasl_mechanism.value = myconfig.kafka.sasl_mechanism;
user.value = myconfig.kafka.user;
password.value = myconfig.kafka.password;

const cancel = () => {
    emit("settingCancel")
}

const save = () => {
    const tmpconfig = {
        name: name.value,
        brokers: brokers.value,
        sasl_mechanism:sasl_mechanism.value,
        user: user.value,
        password: password.value,
    }
    emit("settingSave", tmpconfig)
}

</script>
