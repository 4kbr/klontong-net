package domain

// TODO: ErrAccountNotFound, ErrUnbalancedJournal, ErrInsufficientBalance,
//       ErrBelowMinimumPayout, ErrPayoutInProgress, ErrNoPayoutAccount,
//       ErrEarningNotMatured, ErrDisburserUnavailable
//
// ErrUnbalancedJournal tidak boleh pernah muncul di produksi. Kalau muncul,
// itu bug di kode kita, dan transaksinya harus dibatalkan — bukan dipaksakan.
