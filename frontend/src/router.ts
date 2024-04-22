import { createRouter, createWebHashHistory } from 'vue-router'
import Home from './components/Home.vue'
import Dashboard from './components/Dashboard.vue'
import Brokers from './components/Brokers.vue'
import Broker from './components/Broker.vue'
import Topic from './components/Topic.vue'
import Topics from './components/Topics.vue'
import Group from './components/Group.vue'
import Groups from './components/Groups.vue'
import ZooKeeper from './components/ZooKeeper.vue'


const routes = [
    { path: '/', name: 'Dashboard', component: Dashboard },
    { path: '/brokers', name: 'Brokers', component: Brokers,
      props: true, mete: { title: 'brokers' }
    },
    { path: '/broker', name: 'Broker', component: Broker,
      props: true, mete: { title: 'broker' }
    },
    { path: '/topics', name: 'Topics', component: Topics,
      props: true, mete: { title: 'topics' }
    },
    { path: '/topic', name: 'Topic', component: Topic,
      props: true, mete: { title: 'topic' }
    },
    { path: '/group', name: 'Group', component: Group,
      props: true, mete: { title: 'group' }
    },
    { path: '/groups', name: 'Groups', component: Groups,
      props: true, mete: { title: 'groups' }
    },
    { path: '/zk', name: 'ZooKeeper', component: ZooKeeper },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router
