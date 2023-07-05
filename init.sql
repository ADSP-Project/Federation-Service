CREATE TABLE shops (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) UNIQUE,
  webhookURL VARCHAR(255),
  publicKey VARCHAR(1024)
);

CREATE TABLE partners (
  shopId SERIAL PRIMARY KEY,
  shopName VARCHAR(255),
  canEarnCommission BOOLEAN,
  canShareInventory BOOLEAN,
  canShareData BOOLEAN,
  canCoPromote BOOLEAN,
  canSell BOOLEAN,
  requestStatus VARCHAR(255)
);