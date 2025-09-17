import React, { useState, useContext } from "react";
import Context from "../../Context";
import styles from "./index.module.scss";

 

const handleSubmit = async(e: React.FormEvent) => {
    e.preventDefault();


}

const Register = () => {

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [passworConf, setPasswordConf] = useState("");
    const { dispatch } = useContext(Context);

    return (
        <div className={styles.container}>
        <form onSubmit={handleSubmit} className={styles.loginForm}>
            <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Email"
            />
            <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Password"
            />
            <button type="submit" className={styles.formButton}>Register</button>
            <label> Already Registered? </label><a href="#">Login Here</a>
        </form>
        </div>
    );

}

export default Register;