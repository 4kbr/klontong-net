package domain

// TODO:
//   type Notification struct { ID; RecipientUserID; Kind; Payload; Channel;
//                              ReadAt, SentAt *time.Time; FailedReason string;
//                              CreatedAt }
//   type Preference struct { UserID; EmailEnabled, PushEnabled, WhatsAppEnabled bool;
//                            OrderUpdates, Promotions, SellerUpdates bool }
//   func (p Preference) ShouldSend(kind, channel string) bool
