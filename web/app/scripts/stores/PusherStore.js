var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('events').EventEmitter;
var AppConstants = require('../constants/AppConstants');
var assign = require('object-assign');

var CHANGE_EVENT = 'change';

var _geoevents = [];

function addEvent(geoevent) {
  _geoevents.unshift(geoevent);
  _geoevents = _geoevents.slice(0, 1000);
}

var PusherStore = assign({}, EventEmitter.prototype, {
  getAll: function() {
    return _geoevents;
  },

  last: function() {
    return _geoevents[0];
  },

  emitChange: function() {
    this.emit(CHANGE_EVENT);
  },

  addChangeListener: function(callback) {
    this.on(CHANGE_EVENT, callback);
  },

  removeChangeListener: function(callback) {
    this.removeListener(CHANGE_EVENT, callback);
  }
});

// Register callback to handle all updates
AppDispatcher.register(function(action) {
  switch(action.actionType) {
    case AppConstants.PUSHER_TWEET:
      addEvent(action.geoevent);
      PusherStore.emitChange(CHANGE_EVENT);
      break;
    case AppConstants.PUSHER_RESET_STORE:
      _geoevents.length = 0;
      break;
  }
});

module.exports = PusherStore;
