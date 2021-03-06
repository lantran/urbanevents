var React = window.React = require('react'),
    Router = require('react-router'),
    Cities = require("./components/Cities"),
    CityHeader = require("./components/CityHeader"),
    ListEvents = require("./components/ListEvents"),
    EventActions = require("./actions/EventActions"),
    MappedEvents = require("./components/MappedEvents"),
    mountNode = document.getElementById("app");

var DefaultRoute = Router.DefaultRoute;
var Route = Router.Route;
var RouteHandler = Router.RouteHandler;

var ButtonGroup = require('react-bootstrap/lib/buttonGroup');
var Button = require('react-bootstrap/lib/button');

//http://getbootstrap.com/components/#btn-groups-justified
var App = React.createClass({
  render: function () {
    return (
      <div className="holder">
        <RouteHandler/>
      </div>
    );
  }
});

var routes = (
  <Route name="app" handler={App} path="/">
    <DefaultRoute handler={Cities} />
    <Route name="cities" path="/cities" handler={Cities}/>
    <Route name="city" path="/cities/:cityId" handler={CityHeader}>
      <Route name="map" path="map" handler={MappedEvents}/>
      <Route name="events" path="events" handler={ListEvents}/>
    </Route>
  </Route>
);

Router.run(routes, function (Handler) {
  React.render(<Handler/>, mountNode);
});

EventActions.start();
