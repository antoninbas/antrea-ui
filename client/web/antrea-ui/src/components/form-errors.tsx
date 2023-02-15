import React from 'react';
import { CdsAlertGroup, CdsAlert } from "@cds/react/alert";

interface Props {
};

export function ErrorMessageContainer(props: React.PropsWithChildren<Props>) {
    return (
        <CdsAlertGroup type="banner" status="danger">
            <CdsAlert>{props.children}</CdsAlert>
        </CdsAlertGroup>
    );
}
