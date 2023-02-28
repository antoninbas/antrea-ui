import React, { useState, useEffect, useRef} from 'react';
import logo from './logo.svg';
import './App.css'
import { Outlet, Link, useNavigate } from "react-router-dom";
import NavTab from './components/nav';
import Login from './components/login';
import { AccessTokenProvider, useAccessToken } from './api/token';
import { CdsButton } from '@cds/react/button';

function saveToken(token: string) {
    sessionStorage.setItem('token', token);
}

function retrieveToken(): string | undefined {
    const token = sessionStorage.getItem('token');
    return token ? token : undefined;
}

function removeToken() {
    sessionStorage.removeItem('token')
}

function LoginWall(props: React.PropsWithChildren) {
    const [accessToken, setAccessToken] = useAccessToken();
    const savedToken = retrieveToken()

    function doSetToken(token: string) {
        saveToken(token)
        setAccessToken(token)
    }

    useEffect(() => {
        if (savedToken && !accessToken) {
            setAccessToken(savedToken)
        }
    }, [])

    if (!accessToken) {
        return (
            <div cds-layout="vertical p:md gap:md">
                <p cds-text="section" >Please log in</p>
                <Login setToken={doSetToken} />
            </div>
        )
    }

    return (
        <div cds-layout="vertical align:stretch p:md gap:md">
            {props.children}
        </div>
    );
}

function Logout() {
    const navigate = useNavigate();

    return (
        <div cds-layout="vertical p:md gap:md">
            <CdsButton type="button" action="outline" onClick={()=> { removeToken(); navigate(0); }}>Logout</CdsButton>
        </div>
    )
}

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
                <AccessTokenProvider>
                    <div cds-layout="horizontal align:top wrap:none" style={{height: "100%"}}>
                        <NavTab />
                        <LoginWall>
                            <Outlet />
                        </LoginWall>
                    </div>
                    <Logout />
                </AccessTokenProvider>
            </div>
        </div>
    );
}

export default App;
