import React, { useContext, createContext, useState } from 'react'

type AccessTokenContext = [string, React.Dispatch<React.SetStateAction<string>>]

function AccessTokenProvider(props: React.PropsWithChildren) {
    const [accessToken, setAccessToken] = useState<string>("")
    return <AccessToken.Provider value={[accessToken, setAccessToken]} {...props} />
}

const AccessToken = createContext<AccessTokenContext>(["", {} as React.Dispatch<React.SetStateAction<string>>])

const useAccessToken = (): AccessTokenContext => useContext<AccessTokenContext>(AccessToken)

export { AccessTokenProvider, useAccessToken }
