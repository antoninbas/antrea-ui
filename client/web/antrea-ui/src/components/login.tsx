import { useForm, SubmitHandler } from "react-hook-form";
import { CdsButton } from '@cds/react/button';
import { CdsFormGroup } from '@cds/react/forms';
import { CdsInput } from "@cds/react/input";
import { CdsPassword } from "@cds/react/password";
import { authAPI } from '../api/auth';

type Inputs = {
    username: string
    password: string
};

export default function Login(props: { setToken: (token: string) => void }) {
    const { register, handleSubmit, reset, formState: { errors } } = useForm<Inputs>();
    const setToken = props.setToken;

    const onSubmit: SubmitHandler<Inputs> = async data => {
        try {
            const token = await authAPI.login(data.username, data.password)
            if (token) setToken(token.accessToken)
        } catch(e) {

        }
    }

    return (
        <form onSubmit = {handleSubmit(onSubmit)}>
            <CdsFormGroup layout="horizontal">
                <CdsInput>
                    <label>Username</label>
                    <input {...register("username")} defaultValue="admin" />
                </CdsInput>
                <CdsPassword>
                    <label>Password</label>
                    <input type="password" {...register("password")} />
                </CdsPassword>
                <CdsButton type="submit">Login</CdsButton>
            </CdsFormGroup>
        </form>
    );
}
