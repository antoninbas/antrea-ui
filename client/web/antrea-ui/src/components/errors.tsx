import React, { useState, useEffect, useCallback, useContext } from 'react';
import { CdsAlertGroup, CdsAlert } from "@cds/react/alert";
import { APIError } from '../api/common';

interface APIErrorContextType {
    error: APIError | null
    addError: (error: APIError) => void
    removeError: () => void
}

export const APIErrorContext = React.createContext<APIErrorContextType>({
    error: null,
    addError: (error: APIError) => {},
    removeError: () => {}
});

export function APIErrorProvider(props: React.PropsWithChildren) {
    const [error, setError] = useState<APIError | null>(null);

    const removeError = () => setError(null);

    const addError = (error: APIError) => setError(error);

    const contextValue = {
        error,
        addError: useCallback((error: APIError) => addError(error), []),
        removeError: useCallback(() => removeError(), [])
    };

    return (
        <APIErrorContext.Provider value={contextValue}>
            {props.children}
        </APIErrorContext.Provider>
    );
}

export function useAPIError() {
  const { error, addError, removeError } = useContext(APIErrorContext);
  return { error, addError, removeError };
}

export function APIErrorNotification() {
    const { error, removeError } = useAPIError();

    const handleClose = () => {
        removeError();
    };

    if (!error) {
        return null;
    }

    return (
        <CdsAlertGroup type="banner" status="danger">
            <CdsAlert closable onCloseChange={()=>handleClose()}>{error.message}, {error.code}, {error.status}</CdsAlert>
        </CdsAlertGroup>
    );
}
