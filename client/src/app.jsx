import React from 'react';
import cx from 'classnames';
import {BrowserRouter, NavLink, Route, Routes} from 'react-router-dom';
import Env from './env';
import Mem from './mem';
import Probe from './probe';
import Dns from './dns';
import KeyGen from './keygen';
import Request from './request';
import Disconnected from './disconnected'
import MemQ from './memq'
import ConnErrorContext from './connerror'

// `end` keeps the match exact, so "/" does not stay highlighted on every other page.
function NavItem({to, children}) {
  return (
    <NavLink to={to} end className={({isActive}) => cx("nav-item", {active: isActive})}>
      {children}
    </NavLink>
  );
}

export default class App extends React.Component {
  constructor(props) {
    super(props);
    // Built once rather than per render: a fresh object each time would re-render every consumer.
    this.connError = {
      reportConnError: () => {
        if (this.disconnected) {
          this.disconnected.reportConnError()
        }
      }
    };
  }

  render () {
    let addrs = [];
    for (let a of this.props.page.addrs) {
      addrs.push(<span key={a}>{a}</span>, " ")
    }

    // kuard serves the same app under "", /a, /b and /c; basename makes every route and link
    // relative to whichever one this page was served from. "" is not a valid basename.
    let base = this.props.page.urlBase || "/";

    return (
      <div className="top">
        <div className="title">
          <div className="alert alert-danger" role="alert">
            <svg className="icon icon-notification"><use xlinkHref="#icon-notification"></use></svg> { " " }
            <b>WARNING:</b> This server may expose sensitive and secret information. Be careful.
          </div>
          <Disconnected ref={el => this.disconnected = el}/>
          <h2><samp>{this.props.page.hostname}</samp></h2>
          <div>Demo application version <i>{this.props.page.version}</i></div>
          <div>Serving on {addrs}</div>
        </div>

        <ConnErrorContext.Provider value={this.connError}>
          <BrowserRouter basename={base}>
            <div className="nav-container">
              <div className="nav">
                <NavItem to="/">Request Details</NavItem>
                <NavItem to="/-/env">Server Env</NavItem>
                <NavItem to="/-/mem">Memory</NavItem>
                <NavItem to="/-/liveness">Liveness Probe</NavItem>
                <NavItem to="/-/readiness">Readiness Probe</NavItem>
                <NavItem to="/-/dns">DNS Query</NavItem>
                <NavItem to="/-/keygen">KeyGen Workload</NavItem>
                <NavItem to="/-/memq">MemQ Server</NavItem>
                {/* Served by Go, not React -- must be a real navigation, not a client-side route. */}
                <a className="nav-item" href={this.props.page.urlBase+"/fs/"}>File system browser</a>
              </div>
              <div className="content">
                <Routes>
                  <Route path="/" element={<Request page={this.props.page}/>}/>
                  <Route path="/-/env" element={<Env apiPath={this.props.page.urlBase+"/env/api"}/>}/>
                  <Route path="/-/mem" element={<Mem apiPath={this.props.page.urlBase+"/mem/api"}/>}/>
                  {/* Distinct keys: both routes render Probe, so without them React would reuse the
                      instance and carry one probe's polled history onto the other's page. */}
                  <Route path="/-/liveness" element={<Probe key="liveness" serverPath={this.props.page.urlBase+"/healthy"}/>}/>
                  <Route path="/-/readiness" element={<Probe key="readiness" serverPath={this.props.page.urlBase+"/ready"}/>}/>
                  <Route path="/-/dns" element={<Dns serverPath={this.props.page.urlBase+"/dns"}/>}/>
                  <Route path="/-/keygen" element={<KeyGen serverPath={this.props.page.urlBase+"/keygen"}/>}/>
                  <Route path="/-/memq" element={<MemQ serverPath={this.props.page.urlBase+"/memq"}/>}/>
                </Routes>
              </div>
            </div>
          </BrowserRouter>
        </ConnErrorContext.Provider>
      </div>
    )
  }
}
