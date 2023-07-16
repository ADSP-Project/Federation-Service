import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Header from './components/Header';
import ShopList from './components/ShopList';
import MainPage from './components/MainPage';
import GlobalStyles from './globalStyles';
import { useState, useEffect } from 'react';

function App() {
  
  const [shops, setShops] = useState([]);
  console.log(import.meta.env.VITE_FEDERATION_SERVICE)

  const MY_WEBHOOK_URL = `${import.meta.env.VITE_FEDERATION_SERVICE}/webhook`;

  useEffect(() => {
    // Your hardcoded shops
    const hardcodedShops = [
      { Id: "1", Name: "Tech Mart", Description: "Your one stop for all tech gadgets", WebhookURL: "http://localhost:8001/webhook", PublicKey: "", Image: "../../../tech-mart.jpeg" },
      { Id: "2", Name: "Garden Central", Description: "Everything you need for your garden", WebhookURL: "http://localhost:8002/webhook", PublicKey: "", Image: "../../../garden.jpg" },
      { Id: "3", Name: "Sports Gear Galore", Description: "Sports equipment for all ages", WebhookURL: "http://localhost:8003/webhook", PublicKey: "", Image: "../../../sport.jpeg" },
      { Id: "4", Name: "Fashion Boutique", Description: "Latest fashion trends for you", WebhookURL: "http://localhost:8004/webhook", PublicKey: "", Image: "../../../fashion.jpeg" },
      { Id: "5", Name: "Pet Paradise", Description: "Pet food and accessories", WebhookURL: "http://localhost:8005/webhook", PublicKey: "", Image: "../../../pet.jpeg" },
      { Id: "6", Name: "Home Decor Hub", Description: "Make your home beautiful", WebhookURL: "http://localhost:8006/webhook", PublicKey: "", Image: "../../../home.webp" },
      { Id: "7", Name: "Beauty Bliss", Description: "Beauty products for everyone", WebhookURL: "http://localhost:8007/webhook", PublicKey: "", Image: "../../../beauty.jpeg" },
      { Id: "8", Name: "Fitness Fanatics", Description: "Gym equipment and sportswear", WebhookURL: "http://localhost:8008/webhook", PublicKey: "", Image: "../../../fitness.jpeg" },
      { Id: "9", Name: "Kids Kingdom", Description: "Toys and clothes for kids", WebhookURL: "http://localhost:8009/webhook", PublicKey: "", Image: "../../../kids.webp" },
      { Id: "10", Name: "Auto Accessories", Description: "Accessories for your vehicle", WebhookURL: "http://localhost:8010/webhook", PublicKey: "", Image: "../../../auto.webp" },
      { Id: "11", Name: "Healthy Harvest", Description: "Organic produce for your home", WebhookURL: "http://localhost:8011/webhook", PublicKey: "", Image: "../../../health.jpeg" },
      { Id: "12", Name: "Book Barn", Description: "Books for all genres", WebhookURL: "http://localhost:8012/webhook", PublicKey: "", Image: "../../../book.jpeg" },
      { Id: "13", Name: "Music Mania", Description: "Instruments and music gear", WebhookURL: "http://localhost:8013/webhook", PublicKey: "", Image: "../../../music.jpeg" },
      { Id: "14", Name: "Travel Treasures", Description: "Travel gear for your adventures", WebhookURL: "http://localhost:8014/webhook", PublicKey: "", Image: "../../../travel.jpeg" },
      { Id: "15", Name: "Artistic Alley", Description: "Art supplies and crafts", WebhookURL: "http://localhost:8015/webhook", PublicKey: "", Image: "../../../art.jpeg" },
      { Id: "16", Name: "Outdoor Outfitters", Description: "Gear for your outdoor activities", WebhookURL: "http://localhost:8016/webhook", PublicKey: "", Image: "../../../outdoor.jpeg" },
      { Id: "17", Name: "Gourmet Grocers", Description: "Fine foods and ingredients", WebhookURL: "http://localhost:8017/webhook", PublicKey: "", Image: "../../../gourmet.jpeg" },
      { Id: "18", Name: "Stationery Stop", Description: "Office supplies and stationery", WebhookURL: "http://localhost:8018/webhook", PublicKey: "", Image: "../../../stationery.jpeg" },
      { Id: "19", Name: "Tool Town", Description: "Tools for your DIY projects", WebhookURL: "http://localhost:8019/webhook", PublicKey: "", Image: "../../../tool.jpeg" },
      { Id: "20", Name: "Luxury Linens", Description: "High-end linens for your home", WebhookURL: "http://localhost:8020/webhook", PublicKey: "",Image: "../../../lux.jpeg" }
    ];    
  
    fetch(`${import.meta.env.VITE_FEDERATION_SERVICE}/api/v1/shops`)
      .then(response => response.json())
      .then(data => {
        const filteredShops = data.map(shop => ({
          ...shop,
          Image: shop.Image ? shop.Image : "../../../bouqiue.webp" // If the shop data includes an image path, use it; otherwise, use a default image.
        })).filter(shop => shop.WebhookURL !== MY_WEBHOOK_URL);
            // Append the hardcoded shops to the fetched shops
        setShops([...filteredShops, ...hardcodedShops]);
      })
      .catch((error) => {
        console.error('Error:', error);
        // If there is an error fetching, set the hardcoded shops
        setShops(hardcodedShops);
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
