import { useForm, SubmitHandler } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { CdsCard } from '@cds/react/card';
import { CdsDivider } from '@cds/react/divider';
import { CdsButton } from '@cds/react/button';
import { CdsFormGroup } from '@cds/react/forms';
import { CdsInput } from "@cds/react/input";
import { CdsPassword } from "@cds/react/password";
import { accountAPI } from '../api/account';
import { useAccessToken } from '../api/token';

type Inputs = {
    newPassword: string
};

function UpdatePassword() {
    const { register, handleSubmit, reset, formState: { errors } } = useForm<Inputs>();
    const [accessToken, _] = useAccessToken();

    const navigate = useNavigate();

    const onSubmit: SubmitHandler<Inputs> = async data => {
        try {
            await accountAPI.updatePassword(data.newPassword, accessToken)
        } catch(e) {

        }
        // TODO: avoid code duplication with App.tsx for logout
        sessionStorage.removeItem('token')
        navigate("/")
        navigate(0)
    }

    return (
        <CdsCard>
            <div cds-layout="vertical gap:md">
                <div cds-text="section" cds-layout="p-y:sm">
                    Update Password
                </div>
                <CdsDivider></CdsDivider>
                <form onSubmit = {handleSubmit(onSubmit)}>
                    <CdsFormGroup layout="horizontal">
                        <CdsPassword>
                            <label>New Password</label>
                            <input type="password" {...register("newPassword")} />
                        </CdsPassword>
                        <CdsButton type="submit">Submit</CdsButton>
                    </CdsFormGroup>
                </form>
            </div>
        </CdsCard>
    );
}

export default function Settings() {
    return (
        <main>
            <div cds-layout="vertical gap:lg">
                <p cds-text="title">Settings</p>
                <UpdatePassword />
            </div>
        </main>
    );
}
