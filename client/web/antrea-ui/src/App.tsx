import React, { useState, useEffect, useRef} from 'react';
import logo from './logo.svg';
import './App.css'
import { Outlet, Link, useNavigate } from "react-router-dom";
import NavTab from './components/nav';
import Login from './components/login';
import { useLogout} from './components/logout';
import { CdsButton } from '@cds/react/button';
import { APIErrorProvider, APIErrorNotification } from './components/errors';
import { Provider, useSelector, useDispatch } from 'react-redux'
import type { RootState } from './store'
import { store, setToken } from './store'
import { authAPI } from './api/auth'

function LoginWall(props: React.PropsWithChildren) {
    const token = useSelector((state: RootState) => state.token)
    const dispatch = useDispatch()

    useEffect(() => {
        async function refreshToken() {
            try {
                await authAPI.refreshToken()
            } catch (error) {

            }
        }

        if (token === undefined) {
            // try a refresh
            refreshToken()
        }
    }, [token])

    function doSetToken(token: string) {
        dispatch(setToken(token))
    }

    if (!token) {
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
    const [logoutComplete, logout] = useLogout();

    return (
        <div cds-layout="vertical p:md gap:md">
            <CdsButton type="button" action="outline" onClick={()=> { logout(); }}>Logout</CdsButton>
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
                <Provider store={store}>
                    <div cds-layout="horizontal align:top wrap:none" style={{height: "100%"}}>
                        <NavTab />
                        <LoginWall>
                            <APIErrorProvider>
                                <Outlet />
                                <APIErrorNotification />
                            </APIErrorProvider>
                        </LoginWall>
                    </div>
                    <Logout />
                </Provider>
            </div>
        </div>
    );
}

export default App;
