import React, { useState, useContext } from "react";
import Context from "../../Context";
import styles from "./index.module.scss";
import { useNavigate } from "react-router";
import AuthError  from "../Error/autherror"; 
//comment



const Login = () => {
    const [user, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const { dispatch, authError, isAuthenticated, sessionToken, accessToken } = useContext(Context);

    const navigate = useNavigate();
   
    
    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        const formData = new FormData();
        formData.append("user", user);
        formData.append("password", password);
        const response = await fetch(`/api/auth/login`, {
        method: "POST",
        body: formData,
        });
        if (response.ok) {
            const data = await response.json();
            if (data.plaidToken) {
                dispatch({ type: "SET_STATE", state: { user_id: data.user_id, user: data.username, isAuthenticated: true, sessionToken: data.token, accessToken: data.plaidToken, linkSuccess: true }});
            } else {
                dispatch({ type: "SET_STATE", state: { user_id: data.user_id, user: data.username, isAuthenticated: true, sessionToken: data.token }});
            }

            navigate("/");
        } else {
            dispatch({ type: "SET_STATE", state: { authError: { error_code: String(response.status) } } });
        }
    };


    
    return (
        <div className={styles.mainContainer}>
            <div className={styles.container}>
            <form onSubmit={handleSubmit} className={styles.loginForm}>
                <input
                type="text"
                value={user}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="User"
                />
                <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Password"
                />
                <button type="submit" className={styles.formButton}>Login</button>
                <label> Not Registered? </label><a href="#">Register Here</a>
            </form>
            </div>
            
            {
                authError.error_code &&  
                <div><AuthError /></div>
            }
           
        </div>
    );

}

export default Login;
