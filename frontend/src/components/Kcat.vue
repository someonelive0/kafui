<template>
  <div>
    <v-card class="my-4">
      <div class="text-h3"> Kcat </div>
      <div class="text-caption pa-2 border-t">
        <div class="font-weight-medium">Kcat is a commond line tool of kafka, that from https://github.com/edenhill/kcat</div>
        <div class="font-weight-medium">Recommond kcat, it's another tool of kafka, it's better than Kafui.</div>
      </div>
    </v-card>

    <v-card class="my-4">
      <div class="text-h4"> Install </div>
      <div class="text-caption pa-2 border-t">
        <div class="font-weight-medium">Download or compile kcat on linux, it depends on librdkafka. See https://github.com/edenhill/kcat</div>
      </div>
    </v-card>

    <v-card class="my-4">
      <div class="text-h4"> Usage </div>
      <div class="text-caption pa-2 border-t">
        <div class="opacity-70">see help</div>
        <div class="font-weight-medium">$ kcat -h</div>
      </div>

      <div class="text-h5"> List topics </div>
      <div class="text-caption pa-2 border-t">
        <div class="font-weight-medium">$ kcat -b mybroker -L</div>
      </div>
      <div class="text-caption pa-2 border-t">
        <div class="font-weight-medium">$ kcat -b mybroker -L | grep topicname </div>
      </div>

      <div class="text-h5"> High-level balanced KafkaConsumer </div>
      <div class="text-caption pa-2 border-t">
        <div class="opacity-70">subscribe to topic1 and topic2 (requires broker >=0.9.0 and librdkafka version >=0.9.1)</div>
        <div class="font-weight-medium">$ kcat -b mybroker -G mygroup -o beginning topic1 topic2</div>
        <div class="font-weight-medium">$ kcat -b mybroker -G mygroup topic1 topic2</div>
      </div>

      <div class="text-h5"> Read message </div>
      <div class="text-caption pa-2 border-t">
        <div class="opacity-70">read messages without Consumer Group, and will wait when read EOF</div>
      <div class="font-weight-medium">$ kcat -b mybroker -t topic1 </div>
      </div>
      <div class="text-caption pa-2 border-t">
        <div class="opacity-70">Read the last 2000 messages from 'topic1' topic, partition 0, then exit</div>
      <div class="font-weight-medium">$ kcat -C -b mybroker -t topic1 -p 0 -o -2000 -e </div>
      </div>
      <div class="text-caption pa-2 border-t">
        <div class="opacity-70">Output consumed messages according to format string</div>
      <div class="font-weight-medium">$ kcat -b mybroker -t topic1 -f 'Topic %t[%p], offset: %o, key: %k, payload: %S bytes: %s\n' </div>
      </div>

      <div class="text-h5"> Produce messages </div>
      <div class="text-caption pa-2 border-t">
        <div class="opacity-70">Produce messages from file (one file is one message)</div>
      <div class="font-weight-medium">$ kcat -P -b mybroker -t filedrop -p 0 myfile1.bin /etc/motd thirdfile.tgz </div>
      </div>
      <div class="text-caption pa-2 border-t">
        <div class="opacity-70">Produce messages transactionally (one single transaction for all messages)</div>
      <div class="font-weight-medium">$ kcat -P -b mybroker -t mytopic -X transactional.id=myproducerapp </div>
      </div>

      <div class="text-h5"> SASL authentication </div>
      <div class="text-caption pa-2 border-t">
        <div class="opacity-70">Connect to cluster using SSL and SASL PLAIN authentication</div>
        <div class="font-weight-medium">$ kcat -b mybroker -X security.protocol=SASL_SSL -X sasl.mechanism=PLAIN -X sasl.username=myapikey -X sasl.password=myapisecret ... </div>
      </div>
      <div class="text-caption pa-2 border-t">
        <div class="opacity-70">use config file </div>
        <div class="font-weight-medium">$ kcat -b mybroker -F kcat.conf </div>
      </div>
      <div class="text-caption pa-2 border-t">
        <div class="opacity-70">$ cat kcat.conf</div>
        <pre class="font-weight-medium">

security.protocol=SASL_PLAINTEXT
sasl.mechanisms=PLAIN
sasl.username=USER
sasl.password=PASSWORD

        </pre>
      </div>

      <div class="text-h5"> Query messages </div>
      <div class="text-caption pa-2 border-t">
        <div class="opacity-70">Query offset(s) by timestamp(s)</div>
      <div class="font-weight-medium">$ kcat -b mybroker -Q -t mytopic:3:2389238523 -t mytopic2:0:18921841 </div>
      </div>
      <div class="text-caption pa-2 border-t">
        <div class="opacity-70">Consume messages between two timestamps</div>
      <div class="font-weight-medium">$ kcat -b mybroker -C -t mytopic -o s@1568276612443 -o e@1568276617901 </div>
      </div>

    </v-card>


  </div>
</template>