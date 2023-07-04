import { useState, useEffect } from 'react';
import { Tile, TileHeader, TileBody, TileFooter, Image, JoinButton } from './ShopTile.styles';
import RightsDropdown from "../RightsDropdown"

const ShopTile = ({ shop }) => {
  const [randomImage, setRandomImage] = useState('');
  const imageArray = ['../../../shoe.jpeg', '../../../vintage.jpeg', '../../../spot.jpeg', '../../../sex.avif'];

  useEffect(() => {
    const randomIndex = Math.floor(Math.random() * imageArray.length);
    setRandomImage(imageArray[randomIndex]);
  }, []);

  const [selectedRights, setSelectedRights] = useState([]);

  const handleJoinButtonClick = async () => {
    const currentShopId = 23139;

    const partnershipRequest = {
      shopId: currentShopId,
      partnerId: shop.id,
      rights: selectedRights
    };

    const res = await fetch('http://localhost:8091/api/v1/partnerships/request', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(partnershipRequest)
    });

    const data = await res.json();
    console.log(data);
  };


  return (
    <Tile>
      <Image src={randomImage} alt={shop.name} />
      <TileHeader>{shop.name}</TileHeader>
      <TileBody>{shop.description}</TileBody>
      <TileFooter>
        <RightsDropdown
          options={["canEarnCommission", "canShareInventory", "canShareData", "canCoPromote", "canSell"]}
          selectedRights={selectedRights}
          onChange={setSelectedRights}
        />
        <JoinButton onClick={handleJoinButtonClick}>Join</JoinButton>
      </TileFooter>
    </Tile>
  );
};

export default ShopTile;
