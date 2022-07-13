import React from 'react';
import ReactDOM from 'react-dom/client';
import { Provider } from 'react-redux';
import { Routes, Route, BrowserRouter } from 'react-router-dom';

import 'antd/dist/antd.less';

import store from './store';

import Auth from './containers/Auth'
import App from './containers/App'
import Dashboard from './containers/Dashboard'
import DomainList from './containers/DomainList'
import DomainShow from './containers/DomainShow'
import SuffixList from './containers/SuffixList'
import WhoisList from './containers/WhoisList'
import WhoisQuery from './containers/WhoisQuery'
import WhoisShow from './containers/WhoisShow'

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <BrowserRouter>
    <Provider store={store}>
      <Routes>
        <Route path="/auth" element={<Auth />}/>
        <Route path="/" element={<App />}>
          <Route index element={<Dashboard />}/>
          <Route path="/domain" element={<DomainList />}/>
          <Route path="/domain/:domainID" element={<DomainShow />}/>
          <Route path="/suffix" element={<SuffixList />}/>
          <Route path="/whois/query" element={<WhoisQuery />}/>
          <Route path="/whois/:whoisID" element={<WhoisShow />}/>
          <Route path="/whois" element={<WhoisList />}/>
        </Route>
      </Routes>
    </Provider>
  </BrowserRouter>
);
