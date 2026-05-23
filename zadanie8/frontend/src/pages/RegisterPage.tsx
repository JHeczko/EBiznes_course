import { useState} from 'react';
import type {ChangeEvent, FormEvent} from 'react';
import { useNavigate } from 'react-router-dom';
import { registerUser } from '../services/api.ts';
import "./LoginPage.css";
import toast from "react-hot-toast";

function RegisterPage() {
    const [email, setEmail] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const navigate = useNavigate();

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();

        // BOMBKA 1: Podstawowa walidacja długości hasła
        if (password.length < 6) {
            toast.error("Hasło musi mieć minimum 6 znaków");
            return;
        }

        try {
            await registerUser(email, password);
            navigate('/login'); // Po sukcesie wracamy do logowania
            toast.success("Register successful");
        } catch (err) {
            toast.error((err as Error).message);
        }
    };

    return (
        <main className="auth-page">
            <form className="auth-form" onSubmit={handleSubmit}>
                <h2>Rejestracja</h2>
                <input
                    type="email" placeholder="Email" required
                    onChange={(e: ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)}
                />
                <input
                    type="password" placeholder="Hasło" required
                    onChange={(e: ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)}
                />
                <button type="submit" className="cta-button">Zarejestruj się</button>
            </form>
        </main>
    );
}

export default RegisterPage;