import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import ListView from './ListView';
import searchView from './SearchView';
import {HashRouter as Router, Route} from "react-router-dom"
//import * as serviceWorker from './serviceWorker';

const routing = (
    <Router>
      <React.Fragment>
        <Route exact path="/" component={App} />
        <Route path="/list" component={ListView} />
        <Route path="/search" component={searchView} />
      </React.Fragment>
    </Router>
  )
  ReactDOM.render(routing, document.getElementById('root'))
// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
//serviceWorker.unregister();
