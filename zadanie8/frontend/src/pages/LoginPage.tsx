import { useState, useEffect } from 'react';
import type { ChangeEvent, FormEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { loginUser } from '../services/api.ts';
import "./LoginPage.css";
import toast from "react-hot-toast";

function LoginPage() {
    const [email, setEmail] = useState<string>('');
    const [password, setPassword] = useState<string>('');

    const [isLoggedIn, setIsLoggedIn] = useState<boolean>(() => {
        return localStorage.getItem("auth_token") !== null;
    });

    const [savedUserId, setSavedUserId] = useState<string | null>(() => {
        return localStorage.getItem("user_id");
    });

    const navigate = useNavigate();

    // 🔥 POPRAWIONY HOOK: Brak bezpośredniego wywoływania setState w efekcie
    useEffect(() => {
        const queryParams = new URLSearchParams(window.location.search);
        const tokenFromUrl = queryParams.get("token");
        const userIdFromUrl = queryParams.get("user_id");

        if (tokenFromUrl && userIdFromUrl) {
            // 1. Zapisujemy bezpiecznie dane do pamięci podręcznej przeglądarki
            localStorage.setItem("auth_token", tokenFromUrl);
            localStorage.setItem("user_id", userIdFromUrl);

            toast.success("Successfully logged in via Google!");

            // 2. Natychmiast uciekamy na stronę główną.
            // Strona główna po zamontowaniu i tak sprawdzi localStorage i rozpozna usera jako zalogowanego.
            navigate("/", { replace: true });
        }
    }, [navigate]);

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        try {
            const data = await loginUser(email, password);

            localStorage.setItem("auth_token", data.token);
            localStorage.setItem("user_id", String(data.user_id));

            setIsLoggedIn(true);
            setSavedUserId(String(data.user_id));

            console.log("Logged in:", data);
            toast.success("Login Successful");
            navigate("/");
        } catch (err) {
            toast.error((err as Error).message);
        }
    };

    const handleLogout = () => {
        localStorage.removeItem("auth_token");
        localStorage.removeItem("user_id");

        setIsLoggedIn(false);
        setSavedUserId(null);

        toast.success("Logged out successfully");
        navigate("/login");
    };

    return (
        <main className="auth-page">
            {isLoggedIn ? (
                /* VIEW FOR LOGGED IN USER */
                <div className="auth-form text-center">
                    <h2>You are already logged in!</h2>
                    <p>
                        Logged in as user ID: <strong>{savedUserId || "Unknown"}</strong>
                    </p>

                    <button onClick={handleLogout} className="google-button logout-button">
                        Log Out
                    </button>

                    <div style={{ marginTop: '20px' }}>
                        <Link to="/" className="home-link">
                            Back to Home Page
                        </Link>
                    </div>
                </div>
            ) : (
                /* LOGIN FORM */
                <form className="auth-form" onSubmit={handleSubmit}>
                    <h2>Login</h2>
                    <input
                        type="email" placeholder="Email" required
                        onChange={(e: ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)}
                    />
                    <input
                        type="password" placeholder="Password" required
                        onChange={(e: ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)}
                    />
                    <button type="submit" className="cta-button google-button">Sign In</button>

                    <div className="text-center">or</div>

                    <a href="http://localhost:13000/auth/google/login" className="google-button">
                        Sign in with Google
                    </a>

                    <Link to="/register" className="google-button">
                        Register
                    </Link>
                </form>
            )}
        </main>
    );
}

export default LoginPage;