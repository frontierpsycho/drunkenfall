import moment from 'moment'
import _ from 'lodash'

export default class Match {
  static fromObject (obj) {
    let m = new Match()
    Object.assign(m, obj)

    m.started = moment(m.started)
    m.ended = moment(m.ended)

    switch (m.kind) {
      case 'tryout':
      case 'semi':
        m.end = 10
        break
      case 'final':
        m.end = 20
        break
    }

    return m
  }

  get end () { return this._end }
  set end (value) { this._end = value }

  get isStarted () {
    // match is started if 'started' is defined and NOT equal to go's zero date
    return !isGoZeroDateOrFalsy(this.started)
  }

  get isEnded () {
    // match is ended if 'ended' is defined and NOT equal to go's zero date
    return !isGoZeroDateOrFalsy(this.ended)
  }

  get canStart () {
    return !this.isStarted
  }

  get canEnd () {
    // can't end if already ended
    if (this.isEnded) {
      return false
    }

    // can end if at least one player has enough kills (ie >= end)
    return _.some(this.players, (player) => { return player.kills >= this.end })
  }

  get isRunning () {
    return this.isStarted && !this.isEnded
  }
};

function isGoZeroDateOrFalsy (m) {
  if (m) {
    return m.isSame("0001-01-01T00:00:00Z")
  } else {
    return true
  }
}
