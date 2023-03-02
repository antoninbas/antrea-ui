import React, { useState, useEffect, useRef} from 'react';
import { useNavigate } from "react-router-dom";
import { AccessTokenProvider, useAccessToken } from '../api/token';
import { authAPI } from '../api/auth';

export function useLogout(): [boolean, (() => Promise<void>)] {
    const navigate = useNavigate();
    const [accessToken, setAccessToken] = useAccessToken();
    const [logoutComplete, setLogoutComplete] = useState<boolean>(false);

    async function logout() {
        await authAPI.logout()
        setAccessToken("")
        setLogoutComplete(true)
        navigate("/")
        // navigate(0)
    }

    return [logoutComplete, logout];
}
