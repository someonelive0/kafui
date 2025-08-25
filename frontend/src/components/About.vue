<template>
    <v-card flat
      class="mx-auto"
      prepend-icon="mdi-information"
      width="400"
    >

      <template v-slot:title>
        Kafui
      </template>
      <template v-slot:subtitle>
        Version: 0.0.4 2025-08-25
      </template>

      <v-card-text>
        Kafui is a Kafka GUI client.
        <v-img :src="myLogo" width="330" ></v-img>
      </v-card-text>

      <v-card-actions>
        <v-btn color="secondary" text="GITHUB" @click="visitGithub"></v-btn>
        <v-spacer></v-spacer>
        https://github.com/someonelive0/kafui
      </v-card-actions>

    </v-card>
</template>


<script setup lang="ts">
import myLogo from "../assets/images/appicon.png";
import { Goto } from "../wailsjs/go/main/App";


const visitGithub = () => {
  var url = "https://github.com/someonelive0/kafui";

  Goto(url).then(() => { // window.go.main.App.Goto

  })
  .catch((err: string) => {
    console.error('KafkaTool.ListGroups', err);
  });
}

</script>
