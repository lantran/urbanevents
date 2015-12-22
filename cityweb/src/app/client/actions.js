import { pushPath } from 'redux-simple-router'

/* CityWeb Action Types */
var keyMirror = require('keymirror')
var actionTypes = keyMirror({
  SET_ACROSS: null,
  SET_GEOEVENTS: null,
  SET_CITIES: null,
  SET_CURRENT_CITY: null,
  GET_CITY: null
})

export const ActionTypes = actionTypes

/* CityWeb Action Creators */

export function getAcrossAsync(q) {
  return (dispatch, getState) => {
    if (!q || q.length == 0) {
      return (dispatch) => { dispatch(clearAcross()) }
    }

    let url = "/api/v1/across/search?q=" + q;

    return fetch(url).then((result) => {
      return result.json();
    }).then(rawGeoevents => {
      let moldedGeoevents = {}
      rawGeoevents.forEach(cityGeoevent => {
        moldedGeoevents[cityGeoevent.key] = cityGeoevent.geoevents
      })

      dispatch(setAcross(q, moldedGeoevents))
    });
  }
}

export function setAcross(q, cityGeoevents) {
  return { type: actionTypes.SET_ACROSS, q, cityGeoevents };
}

export function clearAcross() {
  return { type: actionTypes.SET_ACROSS, null, cityGeoevents: {} };
}

export function getGeoeventsAsync(cityKey, q) {
  return (dispatch, getState) => {
    if (!q || q.length == 0) {
      return (dispatch) => { dispatch(clearGeoevents()) }
    }

    let url = "/api/v1/cities/" + cityKey + "/search?q=" + q;

    return fetch(url).then((result) => {
      return result.json();
    }).then(geoevents => {
      dispatch(setGeoevents(q, geoevents))
      //dispatch(pushPath('/' + cityKey + '?q=' + q), getState());
    });
  }
}

export function setGeoevents(q, geoevents) {
  return { type: actionTypes.SET_GEOEVENTS, q, geoevents };
}

export function clearGeoevents() {
  return { type: actionTypes.SET_GEOEVENTS, q: null, geoevents: [] };
}

export function getCitiesAsync() {
  return (dispatch) => {
    return fetch('/api/v1/cities').then(function(result) {
      return result.json();
    }).then(cities => { dispatch(setCities(cities)) });
  };
}

export function setCities(cities) {
  return { type: actionTypes.SET_CITIES, cities };
}

export function setCurrentCity(cityKey) {
  return { type: actionTypes.SET_CURRENT_CITY, cityKey };
}

export function getCurrentCity() {
  return { type: actionTypes.GET_CURRENT_CITY };
}
