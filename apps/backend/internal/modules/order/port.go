package order

// TODO:
//   type OrderInfo struct { ID uuid.UUID; Number string; BuyerUserID uuid.UUID;
//                           Status string; GrandTotal money.Amount }
//   type SuborderInfo struct { ID, OrderID, SellerID, OutletID uuid.UUID;
//                              Number, Status, FulfillmentMethod string
//                              TotalAmount, CommissionAmount, SellerEarningAmount money.Amount }
//
//   type Port interface {
//       GetOrder(ctx, orderID uuid.UUID) (OrderInfo, error)
//       GetSuborder(ctx, suborderID uuid.UUID) (SuborderInfo, error)
//       HasPurchased(ctx, userID, variantID uuid.UUID) (orderItemID uuid.UUID, ok bool, err error)
//   }
//
// HasPurchased dipakai `review`: hanya yang benar-benar membeli dan menerima
// yang boleh menulis ulasan, dan hanya sekali per barang yang dibeli.
