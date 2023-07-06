import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Header from './components/Header';
import ShopList from './components/ShopList';
import MainPage from './components/MainPage';
import GlobalStyles from './globalStyles';
import { useState, useEffect } from 'react';

function App() {
  
  const [shops, setShops] = useState([]);

  const MY_WEBHOOK_URL = "http://localhost:8091/webhook";

  useEffect(() => {
    fetch('http://localhost:8091/api/v1/shops')
      .then(response => response.json())
      .then(data => {
        const filteredShops = data.filter(shop => shop.WebhookURL !== MY_WEBHOOK_URL);
        setShops(filteredShops);
      })
      .catch((error) => {
        console.error('Error:', error);
      });
  }, []);

  return (
    <Router>
      <div className="App">
        <GlobalStyles />
        <Header />
        <Layout>
          <Routes>
            <Route path="/" element={<MainPage shops={shops} />} />
            <Route path="/partners" element={<ShopList shops={shops} />}/>
          </Routes>
        </Layout>
      </div>
    </Router>
  );
}

export default App;
