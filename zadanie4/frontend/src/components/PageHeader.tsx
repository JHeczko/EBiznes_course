import "./pageHeader.css";
import { Link } from "react-router-dom";


function PageHeader() {
    return (
        <header>
            <h1 className="page_title"><Link id="store-name" to="/">Store</Link></h1>
            <div className="container">
                <ul className="horizontal-list">
                    <li><Link className="menu-element" to="/cart">Cart</Link></li>
                    <li><Link className="menu-element" to="/products">Products</Link></li>
                    <li><Link className="menu-element" to="/info">Info</Link></li>
                </ul>
            </div>
        </header>
    )
}

export default PageHeader