import Vue from 'vue'
import Router from 'vue-router'
import Resource from 'vue-resource'
import App from './App'

import TournamentList from './components/TournamentList.vue'
import Tournament from './components/Tournament.vue'
import New from './components/New.vue'
import Join from './components/Join.vue'
import Match from './components/Match.vue'

// install router
Vue.use(Router)
Vue.use(Resource)

// routing
var router = new Router({
  'hashbang': false,
  'history': true
})

router.map({
  '/towerfall/': {
    component: TournamentList
  },
  '/towerfall/new/': {
    component: New
  },
  '/towerfall/:tournament/': {
    component: Tournament
  },
  '/towerfall/:tournament/join/': {
    component: Join
  },
  '/towerfall/:tournament/:kind/:match/': {
    component: Match
  }
})

router.beforeEach(function () {
  window.scrollTo(0, 0)

  router.app.populate()
  router.app.connect()
})

// As long as we only have Drunken TowerFall on drunkenfall.com, we should
// always redirect to the towerfall app right away.
router.redirect({
  '/': '/towerfall/'
})

router.start(App, 'app')
