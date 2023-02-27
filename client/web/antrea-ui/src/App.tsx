import React from 'react';
import logo from './logo.svg';
import './App.css'
import { Outlet, Link } from "react-router-dom";
import NavTab from './components/nav';

function App() {
    return (
        <div cds-text="body" cds-theme="dark">
            {/* 100vh to fill the whole screen */}
            <div style={{height: "fit-content", minHeight:"100vh"}} cds-layout="vertical gap:md align:top">
                <header cds-layout="horizontal wrap:none gap:md m-t:lg">
                    <Link to="/">
                        <img src={logo} alt="logo" style={{height: "2rem"}}/>
                    </Link>
                    <p cds-text="heading" cds-layout="align:vertical-center">Antrea UI</p>
                </header>
                <div cds-layout="horizontal align:top wrap:none" style={{height: "100%"}}>
                    <NavTab />
                    <div cds-layout="vertical align:stretch p:md gap:md">
                        <Outlet />
                    </div>
                </div>
            </div>
        </div>
    );
}

export default App;
