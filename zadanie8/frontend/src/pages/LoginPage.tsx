import { useState} from 'react';
import type {ChangeEvent, FormEvent} from 'react';
import {Link, useNavigate} from 'react-router-dom';
import { loginUser } from '../services/api.ts';
import "./LoginPage.css";

function LoginPage() {
    const [email, setEmail] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const navigate = useNavigate();

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        try {
            const data = await loginUser(email, password);
            localStorage.setItem("auth_token", data.token);
            console.log("Zalogowano:", data);
            navigate("/");
        } catch (err) {
            alert((err as Error).message);
        }
    };

    return (
        <main className="auth-page">
            <form className="auth-form" onSubmit={handleSubmit}>
                <h2>Logowanie</h2>
                <input
                    type="email" placeholder="Email" required
                    onChange={(e: ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)}
                />
                <input
                    type="password" placeholder="Hasło" required
                    onChange={(e: ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)}
                />
                <button type="submit" className="cta-button">Zaloguj</button>

                <div style={{textAlign: 'center', margin: '10px 0'}}>lub</div>

                <a href="http://localhost:13000/auth/google/login" className="google-button">
                    Zaloguj przez Google
                </a>

                <Link to="/register" className="google-button">
                    Rejestruj
                </Link>
            </form>
        </main>
    );
}

export default LoginPage;