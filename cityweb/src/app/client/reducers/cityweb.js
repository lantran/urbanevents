import { ActionTypes } from '#app/actions'

const initialState = {
  cities: [],
  across: [],
  cityQuery: null
}

export default function cityweb(state = initialState, action) {
  switch (action.type) {
    case ActionTypes.SET_ACROSS:
      return {
        ...state,
        across: action.cityGeoevents,
        cityQuery: action.q
      }

    case ActionTypes.SET_CITIES:
      return {
        ...state,
        cities: action.cities
      };

    default:
      return state;
  }
}
