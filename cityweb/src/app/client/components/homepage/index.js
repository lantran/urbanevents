import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { getCitiesAsync, getAcrossAsync } from '#app/actions';
import { citytile, citytileGrid, searchBar, hidden } from './styles';
import Geoevent from '#app/components/geoevent';
import classNames from 'classnames';

class Homepage extends Component {
  constructor(props) {
    super(props);

    let { query } = this.props.location
    let q = query && query.q
    this.state = { q: q }
  }

  componentDidMount() {
    const { getCitiesAsync, getAcrossAsync } = this.props;
    getCitiesAsync();

    let { query } = this.props.location
    let q = query && query.q
    if (q) {
      getAcrossAsync(q.trim());
    }
  }

  componentWillReceiveProps(newState) {
    const { getAcrossAsync, cityQuery } = this.props;

    let { query } = newState.location
    let newQ = query && query.q

    // Used to manually manage when a user clicks the back
    // button for queries.
    // Crazy this has to be done by hand with Url Queries.
    // Taken care of when using Url Params in React Router.
    if (newQ && newQ != cityQuery) {
      this.setState({q: newQ});
      getAcrossAsync(newQ.trim());
    }
  }

  handleQueryChange(e) {
    this.setState({q: e.target.value});
  }

  handleSearch(q) {
    if(q) {
      const { history } = this.props;
      const transitionTo = history.pushState.bind(history, null);
      const { getAcrossAsync } = this.props;
      getAcrossAsync(q.trim(), transitionTo);
    }
  }

  handleSubmit(e) {
    e.preventDefault();
    this.handleSearch(this.state.q);
    this.refs.q.blur();
  }

  hideCity(city, e) {
    let cityDOM = this.refs["city-" + city.key]
    cityDOM.className = classNames(citytile, hidden)
  }

  /* Change this landing page to a list of cities?
   * Show a few tiles showing the hearts of the city perhaps as
   * a jpg or a leaflet map?
   */
  render() {
    let { cities, across } = this.props;
    let component = this;

    return <div>
      <Helmet
        title='New Tweet City'
        meta={[
          {
            property: 'og:title',
            content: 'New Tweet City Media Search'
          }
        ]}
      />

      <form onSubmit={this.handleSubmit.bind(this)} className={searchBar + " uk-form"}>
        <input type="search" name="q" ref="q" placeholder="Enter a word"
          tabIndex="1"
          value={this.state.q}
          onChange={this.handleQueryChange.bind(this)}
        />
        <input className="uk-button" type="submit" tabIndex="2" value="Search Across Cities"/>
      </form>

      <div className={citytileGrid}>
        {cities.map(function(city) {
          return <div key={city.key} className={citytile} ref={"city-"+city.key}>
            <a className="uk-close" onClick={(e) => { component.hideCity(city, e)}.bind(component)}></a>
            <h1>{city.display}</h1>
            <div className="uk-flex uk-flex-column uk-flex-middle uk-flex-nowrap">
              {(across[city.key] || []).map(function(geoevent) {
                return <Geoevent geoevent={geoevent} key={geoevent.id}/>
              })}
            </div>
          </div>
        })}
      </div>
    </div>;
  }
}

function mapStateToProps(state) {
  return {
    ...state.cityweb
  }
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({
    getCitiesAsync,
    getAcrossAsync
  }, dispatch);
}

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Homepage);
