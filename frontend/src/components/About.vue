<template>
    <v-card flat
      class="mx-auto"
      prepend-icon="mdi-information"
      width="400"
    >
      <!-- <v-img
        class="align-end text-white"
        height="200"
        src="@/assests/images/appicon.png"
        cover
      >
        <v-card-title>Kafui</v-card-title>
        <v-card-subtitle>Versin: 0.0.1</v-card-subtitle>
      </v-img> -->

      <template v-slot:title>
        Kafui
      </template>
      <template v-slot:subtitle>
        Versin: 0.0.3 2025-08-22
      </template>  
      <v-card-text>
        Kafui is a Kafka GUI client.
      </v-card-text>
    </v-card>
</template>
