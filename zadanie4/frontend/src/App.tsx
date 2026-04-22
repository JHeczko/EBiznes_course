import './App.css'
import Router from './router/Router.tsx'
import PageHeader from "./components/PageHeader.tsx";
import {BrowserRouter} from "react-router-dom";

function App() {

    return (
        <>
            <BrowserRouter>
                <PageHeader />
                <Router />
            </BrowserRouter>
        </>
    )
}

export default App
